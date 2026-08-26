package gui

import (
	"context"
	"fmt"
	"net/http"
	"time"

	adapterftp "github.com/douglas/file-server/internal/adapters/ftp"
	"github.com/douglas/file-server/internal/adapters/handlers"
	adaptersftp "github.com/douglas/file-server/internal/adapters/sftp"
	"github.com/douglas/file-server/internal/core/services"
)

// DefaultRunner implementa a execução real dos serviços HTTP, FTP e SFTP.
type DefaultRunner struct{}

// NewDefaultRunner cria uma nova instância padrão do executor de servidores.
func NewDefaultRunner() *DefaultRunner {
	return &DefaultRunner{}
}

// RunWeb inicializa e roda o servidor Web HTTP/HTTPS até o cancelamento do contexto.
func (r *DefaultRunner) RunWeb(ctx context.Context, host string, port int, targetDir string, useTLS bool, certFile, keyFile string) error {
	fileSvc, err := services.NewFileService(targetDir)
	if err != nil {
		return fmt.Errorf("erro ao inicializar serviço de arquivos: %w", err)
	}

	healthSvc := services.NewHealthService()
	h, err := handlers.NewHandler(fileSvc, healthSvc)
	if err != nil {
		return fmt.Errorf("erro ao inicializar handlers: %w", err)
	}

	mux := http.NewServeMux()
	if err := h.RegisterRoutes(mux); err != nil {
		return fmt.Errorf("erro ao registrar rotas: %w", err)
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if certFile != "" || keyFile != "" {
		if certFile == "" || keyFile == "" {
			return fmt.Errorf("ambos certificado e chave privada devem ser fornecidos")
		}
		tlsCfg, err := services.LoadTLSConfig(certFile, keyFile)
		if err != nil {
			return fmt.Errorf("erro ao carregar certificados TLS: %w", err)
		}
		server.TLSConfig = tlsCfg
	} else if useTLS {
		tlsCfg, err := services.CreateSelfSignedTLSConfig(host, addr, "localhost", "127.0.0.1")
		if err != nil {
			return fmt.Errorf("erro ao gerar certificado TLS autoassinado: %w", err)
		}
		server.TLSConfig = tlsCfg
	}

	errChan := make(chan error, 1)
	go func() {
		if server.TLSConfig != nil {
			if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
				errChan <- err
			}
		} else {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				errChan <- err
			}
		}
	}()

	select {
	case err := <-errChan:
		return fmt.Errorf("erro no servidor web: %w", err)
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return server.Shutdown(shutdownCtx)
}

// RunFTP inicializa e roda o servidor FTP/FTPS.
func (r *DefaultRunner) RunFTP(ctx context.Context, opts adapterftp.ServerOptions) error {
	server, err := adapterftp.NewServer(opts)
	if err != nil {
		return fmt.Errorf("erro ao inicializar servidor FTP: %w", err)
	}
	return server.Run(ctx)
}

// RunSFTP inicializa e roda o servidor SFTP sobre SSH.
func (r *DefaultRunner) RunSFTP(ctx context.Context, opts adaptersftp.ServerOptions) error {
	server, err := adaptersftp.NewServer(opts)
	if err != nil {
		return fmt.Errorf("erro ao inicializar servidor SFTP: %w", err)
	}
	return server.Run(ctx)
}
