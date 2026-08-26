package ftp

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/douglas/file-server/internal/core/services"
	"github.com/douglas/file-server/internal/version"
	ftpserver "github.com/fclairamb/ftpserverlib"
)

// ServerOptions define as opções de inicialização do servidor FTP/FTPS.
type ServerOptions struct {
	Host         string
	Port         int
	TargetDir    string
	User         string
	Pass         string
	UseTLS       bool
	TLSCert      string
	TLSKey       string
	PassivePorts string
	ReadOnly     bool
}

// Server encapsula o servidor FTP e suas dependências.
type Server struct {
	opts   ServerOptions
	driver *Driver
	server *ftpserver.FtpServer
}

// NewServer inicializa um novo Server com driver configurado.
func NewServer(opts ServerOptions) (*Server, error) {
	driver, err := NewDriver(opts)
	if err != nil {
		return nil, fmt.Errorf("erro ao configurar driver FTP: %w", err)
	}

	server := ftpserver.NewFtpServer(driver)

	return &Server{
		opts:   opts,
		driver: driver,
		server: server,
	}, nil
}

// Addr retorna o endereço em que o servidor está escutando.
func (s *Server) Addr() string {
	if s.server != nil {
		return s.server.Addr()
	}
	return ""
}

// LogStartupBanner exibe as informações de inicialização do servidor FTP no terminal.
func LogStartupBanner(opts ServerOptions, isTLS bool) {
	scheme := "ftp"
	protocol := "FTP (Texto Plano)"
	if isTLS {
		scheme = "ftps"
		protocol = "FTPS (TLS 1.3 / 1.2 Criptografado)"
	}

	var localAccess string
	var lanAccess []string

	if opts.Host == "0.0.0.0" || opts.Host == "" || opts.Host == "::" {
		localAccess = fmt.Sprintf("%s://%s@127.0.0.1:%d", scheme, opts.User, opts.Port)
		seen := make(map[string]bool)
		for _, ip := range services.GetLANIPAddresses() {
			if ip.IsLoopback() {
				continue
			}
			var url string
			if ip.To4() != nil {
				url = fmt.Sprintf("%s://%s@%s:%d", scheme, opts.User, ip.String(), opts.Port)
			} else {
				url = fmt.Sprintf("%s://%s@[%s]:%d", scheme, opts.User, ip.String(), opts.Port)
			}
			if !seen[url] {
				seen[url] = true
				lanAccess = append(lanAccess, url)
			}
		}
	} else {
		localAccess = fmt.Sprintf("%s://%s@%s:%d", scheme, opts.User, opts.Host, opts.Port)
	}

	log.Printf("⚡ FTP Server %s inicializado com sucesso!", version.Get().Version)
	log.Printf("📁 Diretório compartilhado: %s", opts.TargetDir)
	log.Printf("🔒 Protocolo: %s", protocol)
	log.Printf("👤 Usuário: %s", opts.User)
	if opts.Pass != "" {
		log.Printf("🔑 Senha: %s", opts.Pass)
	}
	if opts.ReadOnly {
		log.Printf("🛡️ Modo: Somente Leitura (Read-Only)")
	}
	if opts.PassivePorts != "" {
		log.Printf("📡 Faixa de Portas Passivas: %s", opts.PassivePorts)
	}
	log.Printf("👉 Acesso Local:     %s", localAccess)
	if len(lanAccess) > 0 {
		for _, url := range lanAccess {
			log.Printf("🌐 Acesso Rede (LAN): %s", url)
		}
	}
}

// Run executa o servidor FTP aguardando sinal de cancelamento pelo contexto.
func (s *Server) Run(ctx context.Context) error {
	if err := s.server.Listen(); err != nil {
		return fmt.Errorf("falha ao iniciar listener FTP: %w", err)
	}

	isTLS := s.driver.tlsConfig != nil
	LogStartupBanner(s.opts, isTLS)

	errChan := make(chan error, 1)
	go func() {
		if err := s.server.Serve(); err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return fmt.Errorf("erro fatal no servidor FTP: %w", err)
	case <-ctx.Done():
		log.Println("🛑 Encerrando servidor FTP graciosamente...")
	}

	shutdownDone := make(chan error, 1)
	go func() {
		shutdownDone <- s.server.Stop()
	}()

	select {
	case err := <-shutdownDone:
		if err != nil {
			return fmt.Errorf("erro ao encerrar servidor FTP: %w", err)
		}
	case <-time.After(5 * time.Second):
		log.Println("⚠️ Tempo limite de encerramento do servidor FTP excedido.")
	}

	log.Println("✅ Servidor FTP finalizado com sucesso.")
	return nil
}
