package ftp

import (
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/douglas/file-server/internal/core/services"
	ftpserver "github.com/fclairamb/ftpserverlib"
	"github.com/spf13/afero"
)

var (
	// ErrInvalidCredentials indica credenciais inválidas para o usuário FTP.
	ErrInvalidCredentials = errors.New("credenciais de autenticação FTP inválidas")

	// ErrInvalidPortRange indica formato inválido para faixa de portas de modo passivo.
	ErrInvalidPortRange = errors.New("formato inválido para faixa de portas passivas (esperado inicio-fim, ex: 50000-50100)")
)

// Driver implementa ftpserver.MainDriver gerenciando autenticação, TLS, configurações e sistema de arquivos com sandbox.
type Driver struct {
	opts      ServerOptions
	rootDir   string
	tlsConfig *tls.Config
}

// NewDriver cria e inicializa um novo Driver validando o diretório raiz e opções de TLS.
func NewDriver(opts ServerOptions) (*Driver, error) {
	if opts.TargetDir == "" {
		opts.TargetDir = "."
	}

	absRoot, err := filepath.Abs(opts.TargetDir)
	if err != nil {
		return nil, fmt.Errorf("caminho de diretório raiz inválido: %w", err)
	}

	info, err := os.Stat(absRoot)
	if err != nil {
		return nil, fmt.Errorf("erro ao acessar diretório raiz '%s': %w", absRoot, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("o caminho raiz '%s' não é um diretório", absRoot)
	}

	canonicalRoot, err := filepath.EvalSymlinks(absRoot)
	if err != nil {
		canonicalRoot = absRoot
	}

	var tlsCfg *tls.Config
	if opts.TLSCert != "" || opts.TLSKey != "" {
		if opts.TLSCert == "" || opts.TLSKey == "" {
			return nil, fmt.Errorf("ambos --tls-cert e --tls-key devem ser fornecidos para usar certificado customizado")
		}
		cfg, err := services.LoadTLSConfig(opts.TLSCert, opts.TLSKey)
		if err != nil {
			return nil, fmt.Errorf("erro ao carregar certificados TLS para FTP: %w", err)
		}
		tlsCfg = cfg
	} else if opts.UseTLS {
		addr := fmt.Sprintf("%s:%d", opts.Host, opts.Port)
		cfg, err := services.CreateSelfSignedTLSConfig(opts.Host, addr, "localhost", "127.0.0.1")
		if err != nil {
			return nil, fmt.Errorf("erro ao gerar certificado TLS autoassinado para FTP: %w", err)
		}
		tlsCfg = cfg
	}

	return &Driver{
		opts:      opts,
		rootDir:   canonicalRoot,
		tlsConfig: tlsCfg,
	}, nil
}

// ParsePassivePorts analisa uma string de intervalo de portas (ex: "50000-50100") e retorna o PortRange.
func ParsePassivePorts(rangeStr string) (*ftpserver.PortRange, error) {
	rangeStr = strings.TrimSpace(rangeStr)
	if rangeStr == "" || rangeStr == "0" {
		return nil, nil
	}

	parts := strings.Split(rangeStr, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("%w: '%s'", ErrInvalidPortRange, rangeStr)
	}

	start, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	end, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err1 != nil || err2 != nil || start <= 0 || end <= 0 || start > end || end > 65535 {
		return nil, fmt.Errorf("%w: '%s'", ErrInvalidPortRange, rangeStr)
	}

	return &ftpserver.PortRange{
		Start: start,
		End:   end,
	}, nil
}

// GetSettings define as configurações operacionais do servidor FTP.
func (d *Driver) GetSettings() (*ftpserver.Settings, error) {
	listenAddr := fmt.Sprintf("%s:%d", d.opts.Host, d.opts.Port)

	var pasvRange *ftpserver.PortRange
	if d.opts.PassivePorts != "" {
		pr, err := ParsePassivePorts(d.opts.PassivePorts)
		if err != nil {
			return nil, err
		}
		pasvRange = pr
	}

	tlsReq := ftpserver.ClearOrEncrypted
	if d.tlsConfig != nil {
		tlsReq = ftpserver.MandatoryEncryption
	}

	var pasvGetter ftpserver.PasvPortGetter
	if pasvRange != nil {
		pasvGetter = *pasvRange
	}

	settings := &ftpserver.Settings{
		ListenAddr:               listenAddr,
		Banner:                   "⚡ File Server FTP Service",
		TLSRequired:              tlsReq,
		PassiveTransferPortRange: pasvGetter,
		DisableActiveMode:        false,
		IdleTimeout:              300,
	}

	return settings, nil
}

// ClientConnected envia a mensagem de boas-vindas ao cliente conectado.
func (d *Driver) ClientConnected(cc ftpserver.ClientContext) (string, error) {
	return "Bem-vindo ao File Server FTP!", nil
}

// ClientDisconnected é chamado quando um cliente é desconectado.
func (d *Driver) ClientDisconnected(cc ftpserver.ClientContext) {
}

// AuthUser autentica as credenciais e retorna o sistema de arquivos isolado em sandbox.
func (d *Driver) AuthUser(cc ftpserver.ClientContext, user, pass string) (ftpserver.ClientDriver, error) {
	if d.opts.User != "" && d.opts.Pass != "" {
		if !services.ValidateCredentials(user, pass, d.opts.User, d.opts.Pass) {
			return nil, ErrInvalidCredentials
		}
	}

	baseFs := afero.NewBasePathFs(afero.NewOsFs(), d.rootDir)
	if d.opts.ReadOnly {
		return afero.NewReadOnlyFs(baseFs), nil
	}

	return baseFs, nil
}

// GetTLSConfig retorna a configuração TLS para conexões FTPS.
func (d *Driver) GetTLSConfig() (*tls.Config, error) {
	return d.tlsConfig, nil
}
