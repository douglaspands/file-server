package gui

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	adapterftp "github.com/douglas/file-server/internal/adapters/ftp"
	adaptersftp "github.com/douglas/file-server/internal/adapters/sftp"
	"github.com/douglas/file-server/internal/core/services"
)

var (
	ErrServerAlreadyRunning = errors.New("o servidor já está em execução")
	ErrServerNotRunning     = errors.New("nenhum servidor está em execução")
	ErrInvalidTargetDir     = errors.New("diretório informado é inválido ou inacessível")
)

// NewController cria e inicializa o gerenciador de estado da GUI.
func NewController(initialDir string, runner ServerRunner, broadcaster *LogBroadcaster) *Controller {
	if runner == nil {
		runner = NewDefaultRunner()
	}
	if broadcaster == nil {
		broadcaster = NewLogBroadcaster(500)
	}

	targetDir := initialDir
	if targetDir == "" {
		cwd, err := os.Getwd()
		if err == nil {
			targetDir = cwd
		} else {
			targetDir = "."
		}
	}
	absTargetDir, err := filepath.Abs(targetDir)
	if err == nil {
		targetDir = absTargetDir
	}

	defaultCfg := ServerConfig{
		Protocol:  ProtocolWeb,
		Host:      "0.0.0.0",
		Port:      8080,
		TargetDir: targetDir,
		UseTLS:    false,
		User:      services.DefaultUsername,
		ReadOnly:  false,
	}

	return &Controller{
		status: ServerStatus{
			IsRunning: false,
			Protocol:  ProtocolWeb,
			Host:      "0.0.0.0",
			Port:      8080,
			TargetDir: targetDir,
		},
		config:      defaultCfg,
		broadcaster: broadcaster,
		runner:      runner,
	}
}

// SetGUIAddr define o endereço do servidor web da própria interface gráfica.
func (c *Controller) SetGUIAddr(host string, port int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.guiHost = host
	c.guiPort = port
}

// GetGUIAddr retorna o host e porta da interface gráfica.
func (c *Controller) GetGUIAddr() (string, int) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.guiHost, c.guiPort
}

// GetBroadcaster retorna a instância de broadcast de logs.
func (c *Controller) GetBroadcaster() *LogBroadcaster {
	return c.broadcaster
}

// GetStatus retorna uma cópia segura do status atual do servidor.
func (c *Controller) GetStatus() ServerStatus {
	c.mu.RLock()
	defer c.mu.RUnlock()

	status := c.status
	if status.IsRunning {
		status.Interfaces = DetectNetworkInterfaces(status.Port, string(status.Protocol), status.UseTLS)
		status.LocalURL, status.LanURLs = BuildAccessURLs(status.Host, status.Port, string(status.Protocol), status.UseTLS)
	}
	return status
}

// GetConfig retorna a configuração corrente.
func (c *Controller) GetConfig() ServerConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config
}

// IsRunning retorna se há algum servidor ativo.
func (c *Controller) IsRunning() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.status.IsRunning
}

// StartServer valida as opções recebidas e inicia o servidor correspondente em background.
func (c *Controller) StartServer(cfg ServerConfig) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.status.IsRunning {
		return ErrServerAlreadyRunning
	}

	// Validação e resolução do diretório compartilhado
	cleanedDir := strings.TrimSpace(cfg.TargetDir)
	if cleanedDir == "" {
		cwd, _ := os.Getwd()
		cleanedDir = cwd
	}
	absDir, err := filepath.Abs(cleanedDir)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidTargetDir, err)
	}
	info, err := os.Stat(absDir)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("%w: '%s'", ErrInvalidTargetDir, absDir)
	}
	cfg.TargetDir = absDir

	if cfg.Port <= 0 || cfg.Port > 65535 {
		switch cfg.Protocol {
		case ProtocolFTP:
			cfg.Port = 2121
		case ProtocolSFTP:
			cfg.Port = 2222
		default:
			cfg.Port = 8080
		}
	}
	if cfg.Host == "" {
		cfg.Host = "0.0.0.0"
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.activeCancel = cancel
	c.config = cfg

	now := time.Now()
	localURL, lanURLs := BuildAccessURLs(cfg.Host, cfg.Port, string(cfg.Protocol), cfg.UseTLS)
	ifaces := DetectNetworkInterfaces(cfg.Port, string(cfg.Protocol), cfg.UseTLS)

	c.status = ServerStatus{
		IsRunning:  true,
		Protocol:   cfg.Protocol,
		Host:       cfg.Host,
		Port:       cfg.Port,
		TargetDir:  cfg.TargetDir,
		UseTLS:     cfg.UseTLS,
		LocalURL:   localURL,
		LanURLs:    lanURLs,
		Interfaces: ifaces,
		User:       cfg.User,
		ReadOnly:   cfg.ReadOnly,
		StartedAt:  &now,
	}

	c.broadcaster.Broadcast(fmt.Sprintf("🚀 [GUI] Iniciando servidor %s na porta %d (%s)...", strings.ToUpper(string(cfg.Protocol)), cfg.Port, cfg.TargetDir))

	go func() {
		var runErr error
		switch cfg.Protocol {
		case ProtocolFTP:
			user := cfg.User
			if user == "" {
				user = services.DefaultUsername
			}
			pass := cfg.Pass
			if pass == "" {
				generated, genErr := services.GenerateRandomPassword(12)
				if genErr == nil {
					pass = generated
				} else {
					pass = "fileserver"
				}
			}
			runErr = c.runner.RunFTP(ctx, adapterftp.ServerOptions{
				Host:         cfg.Host,
				Port:         cfg.Port,
				TargetDir:    cfg.TargetDir,
				User:         user,
				Pass:         pass,
				UseTLS:       cfg.UseTLS,
				TLSCert:      cfg.TLSCert,
				TLSKey:       cfg.TLSKey,
				PassivePorts: cfg.PassivePorts,
				ReadOnly:     cfg.ReadOnly,
			})
		case ProtocolSFTP:
			user := cfg.User
			if user == "" {
				user = services.DefaultUsername
			}
			pass := cfg.Pass
			if pass == "" && cfg.AuthKey == "" {
				generated, genErr := services.GenerateRandomPassword(12)
				if genErr == nil {
					pass = generated
				} else {
					pass = "fileserver"
				}
			}
			runErr = c.runner.RunSFTP(ctx, adaptersftp.ServerOptions{
				Host:      cfg.Host,
				Port:      cfg.Port,
				TargetDir: cfg.TargetDir,
				User:      user,
				Pass:      pass,
				AuthKey:   cfg.AuthKey,
				HostKey:   cfg.HostKey,
				ReadOnly:  cfg.ReadOnly,
			})
		default: // Web / HTTP
			runErr = c.runner.RunWeb(ctx, cfg.Host, cfg.Port, cfg.TargetDir, cfg.UseTLS, cfg.TLSCert, cfg.TLSKey)
		}

		if runErr != nil && !errors.Is(runErr, context.Canceled) {
			c.broadcaster.Broadcast(fmt.Sprintf("❌ [GUI] Erro na execução do servidor: %v", runErr))
		}

		c.mu.Lock()
		c.status.IsRunning = false
		c.status.StartedAt = nil
		c.activeCancel = nil
		c.mu.Unlock()
	}()

	return nil
}

// StopServer interrompe graciosamente o servidor ativo.
func (c *Controller) StopServer() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.status.IsRunning {
		return ErrServerNotRunning
	}

	c.broadcaster.Broadcast("🛑 [GUI] Solicitando interrupção graciosa do servidor...")
	if c.activeCancel != nil {
		c.activeCancel()
		c.activeCancel = nil
	}

	c.status.IsRunning = false
	c.status.StartedAt = nil
	c.broadcaster.Broadcast("✅ [GUI] Servidor finalizado.")
	return nil
}
