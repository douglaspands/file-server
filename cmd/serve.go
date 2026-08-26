package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/douglas/file-server/internal/adapters/handlers"
	"github.com/douglas/file-server/internal/core/services"
	"github.com/douglas/file-server/internal/version"
	"github.com/spf13/cobra"
)

var (
	port int
	host string
)

// serveCmd inicia o servidor HTTP web.
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Inicia o servidor web e API",
	Long:  `Inicia o servidor HTTP da aplicação na porta configurada com suporte a graceful shutdown.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		defer stop()
		return RunServerWithContext(ctx, host, port)
	},
}

// SetupServer inicializa o servidor HTTP e suas dependências.
func SetupServer(host string, port int) (*http.Server, error) {
	healthSvc := services.NewHealthService()
	handler, err := handlers.NewHandler(healthSvc)
	if err != nil {
		return nil, fmt.Errorf("erro ao inicializar handlers: %w", err)
	}

	mux := http.NewServeMux()
	if err := handler.RegisterRoutes(mux); err != nil {
		return nil, fmt.Errorf("erro ao registrar rotas: %w", err)
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	return &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}, nil
}

// RunServerWithContext executa o servidor aguardando cancelamento pelo contexto.
func RunServerWithContext(ctx context.Context, host string, port int) error {
	server, err := SetupServer(host, port)
	if err != nil {
		return err
	}

	errChan := make(chan error, 1)
	go func() {
		log.Printf("🚀 Servidor iniciado em http://%s (versão: %s)", server.Addr, version.Get().Version)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return fmt.Errorf("erro fatal no servidor HTTP: %w", err)
	case <-ctx.Done():
		log.Println("🛑 Encerrando servidor graciosamente...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("erro durante encerramento gracioso: %w", err)
	}

	log.Println("✅ Servidor finalizado com sucesso.")
	return nil
}

func init() {
	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "porta na qual o servidor irá escutar")
	serveCmd.Flags().StringVar(&host, "host", "0.0.0.0", "endereço do host para escuta")
	RootCmd.AddCommand(serveCmd)
}
