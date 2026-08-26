package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/douglas/file-server/internal/adapters/handlers"
	"github.com/douglas/file-server/internal/core/services"
	"github.com/douglas/file-server/internal/version"
	"github.com/spf13/cobra"
)

var (
	port    int
	host    string
	dirFlag string
	useTLS  bool
	tlsCert string
	tlsKey  string
)

// ServerOptions define as opções de inicialização do servidor HTTP/HTTPS.
type ServerOptions struct {
	Host      string
	Port      int
	TargetDir string
	UseTLS    bool
	TLSCert   string
	TLSKey    string
}

// ExpandHomeDir expande o prefixo de til (~) para o diretório home do usuário.
func ExpandHomeDir(path string) (string, error) {
	if path == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			if envHome := os.Getenv("HOME"); envHome != "" {
				return envHome, nil
			}
			return "", fmt.Errorf("não foi possível determinar o diretório home do usuário: %w", err)
		}
		return home, nil
	}
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "~\\") {
		home, err := os.UserHomeDir()
		if err != nil {
			if envHome := os.Getenv("HOME"); envHome != "" {
				return filepath.Join(envHome, path[2:]), nil
			}
			return "", fmt.Errorf("não foi possível determinar o diretório home do usuário: %w", err)
		}
		return filepath.Join(home, path[2:]), nil
	}
	return path, nil
}

// ResolveDirectory determina e valida o diretório a ser servido a partir dos argumentos ou flags,
// suportando caminhos absolutos, relativos (./, ../) e expansão de til (~).
func ResolveDirectory(args []string, flagValue string) (string, error) {
	target := "."
	if len(args) > 0 && strings.TrimSpace(args[0]) != "" {
		target = strings.TrimSpace(args[0])
	} else if strings.TrimSpace(flagValue) != "" {
		target = strings.TrimSpace(flagValue)
	}

	expanded, err := ExpandHomeDir(target)
	if err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(filepath.Clean(expanded))
	if err != nil {
		return "", fmt.Errorf("caminho de diretório inválido '%s': %w", target, err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("diretório informado não existe: %s", absPath)
		}
		return "", fmt.Errorf("erro ao acessar diretório informado '%s': %w", absPath, err)
	}

	if !info.IsDir() {
		return "", fmt.Errorf("o caminho informado não é um diretório: %s", absPath)
	}

	return absPath, nil
}

// serveCmd inicia o servidor HTTP web.
var serveCmd = &cobra.Command{
	Use:   "serve [diretório]",
	Short: "Inicia o servidor web de arquivos e API",
	Long:  `Inicia o servidor HTTP/HTTPS da aplicação servindo o diretório informado (ou a pasta atual) com interface web, downloads, uploads e criptografia em trânsito.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetDir, err := ResolveDirectory(args, dirFlag)
		if err != nil {
			return err
		}

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		defer stop()
		return RunServerWithOptions(ctx, ServerOptions{
			Host:      host,
			Port:      port,
			TargetDir: targetDir,
			UseTLS:    useTLS,
			TLSCert:   tlsCert,
			TLSKey:    tlsKey,
		})
	},
}

// SetupServerWithOptions inicializa o servidor HTTP/HTTPS e suas dependências a partir de ServerOptions.
func SetupServerWithOptions(opts ServerOptions) (*http.Server, error) {
	fileSvc, err := services.NewFileService(opts.TargetDir)
	if err != nil {
		return nil, fmt.Errorf("erro ao inicializar serviço de arquivos: %w", err)
	}

	healthSvc := services.NewHealthService()
	handler, err := handlers.NewHandler(fileSvc, healthSvc)
	if err != nil {
		return nil, fmt.Errorf("erro ao inicializar handlers: %w", err)
	}

	mux := http.NewServeMux()
	if err := handler.RegisterRoutes(mux); err != nil {
		return nil, fmt.Errorf("erro ao registrar rotas: %w", err)
	}

	addr := fmt.Sprintf("%s:%d", opts.Host, opts.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Configuração de TLS / HTTPS se solicitado
	if opts.TLSCert != "" || opts.TLSKey != "" {
		if opts.TLSCert == "" || opts.TLSKey == "" {
			return nil, fmt.Errorf("ambos --tls-cert e --tls-key devem ser fornecidos para usar certificado customizado")
		}
		tlsCfg, err := services.LoadTLSConfig(opts.TLSCert, opts.TLSKey)
		if err != nil {
			return nil, fmt.Errorf("erro ao carregar certificados TLS: %w", err)
		}
		server.TLSConfig = tlsCfg
	} else if opts.UseTLS {
		tlsCfg, err := services.CreateSelfSignedTLSConfig(opts.Host, addr, "localhost", "127.0.0.1")
		if err != nil {
			return nil, fmt.Errorf("erro ao gerar certificado TLS autoassinado: %w", err)
		}
		server.TLSConfig = tlsCfg
	}

	return server, nil
}

// SetupServer inicializa o servidor HTTP e suas dependências com opções padrão (compatibilidade).
func SetupServer(host string, port int, targetDir string) (*http.Server, error) {
	return SetupServerWithOptions(ServerOptions{
		Host:      host,
		Port:      port,
		TargetDir: targetDir,
	})
}

// FormatAccessURLs constrói as URLs de acesso local e de rede (LAN) formatadas com o protocolo apropriado.
func FormatAccessURLs(host string, port int, isTLS bool) (localURL string, lanURLs []string) {
	scheme := "http"
	if isTLS {
		scheme = "https"
	}

	if host == "0.0.0.0" || host == "" || host == "::" {
		localURL = fmt.Sprintf("%s://127.0.0.1:%d", scheme, port)
		seen := make(map[string]bool)
		for _, ip := range services.GetLANIPAddresses() {
			if ip.IsLoopback() {
				continue
			}
			ipStr := ip.String()
			var url string
			if ip.To4() != nil {
				url = fmt.Sprintf("%s://%s:%d", scheme, ipStr, port)
			} else {
				url = fmt.Sprintf("%s://[%s]:%d", scheme, ipStr, port)
			}
			if !seen[url] {
				seen[url] = true
				lanURLs = append(lanURLs, url)
			}
		}
	} else {
		localURL = fmt.Sprintf("%s://%s:%d", scheme, host, port)
	}

	return localURL, lanURLs
}

// LogStartupBanner imprime as informações de inicialização e as URLs de acesso disponíveis.
func LogStartupBanner(opts ServerOptions, isTLS bool) {
	localURL, lanURLs := FormatAccessURLs(opts.Host, opts.Port, isTLS)
	protocol := "HTTP"
	if isTLS {
		protocol = "HTTPS (TLS 1.3 / HTTP/2)"
	}

	log.Printf("⚡ File Server %s inicializado com sucesso!", version.Get().Version)
	log.Printf("📁 Diretório compartilhado: %s", opts.TargetDir)
	log.Printf("🔒 Protocolo: %s", protocol)
	log.Printf("👉 Acesso Local:     %s", localURL)
	if len(lanURLs) > 0 {
		for _, url := range lanURLs {
			log.Printf("🌐 Acesso Rede (LAN): %s", url)
		}
	} else {
		log.Printf("🌐 Acesso Rede (LAN): nenhum IP adicional de rede local detectado")
	}
}

// RunServerWithOptions executa o servidor aguardando cancelamento pelo contexto com suporte a TLS.
func RunServerWithOptions(ctx context.Context, opts ServerOptions) error {
	server, err := SetupServerWithOptions(opts)
	if err != nil {
		return err
	}

	isTLS := server.TLSConfig != nil
	LogStartupBanner(opts, isTLS)

	errChan := make(chan error, 1)
	go func() {
		if isTLS {
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
		return fmt.Errorf("erro fatal no servidor web: %w", err)
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

// RunServerWithContext executa o servidor aguardando cancelamento pelo contexto (compatibilidade).
func RunServerWithContext(ctx context.Context, host string, port int, targetDir string) error {
	return RunServerWithOptions(ctx, ServerOptions{
		Host:      host,
		Port:      port,
		TargetDir: targetDir,
	})
}

func init() {
	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "porta na qual o servidor irá escutar")
	serveCmd.Flags().StringVar(&host, "host", "0.0.0.0", "endereço do host para escuta")
	serveCmd.Flags().StringVarP(&dirFlag, "dir", "d", "", "caminho do diretório raiz a ser compartilhado")
	serveCmd.Flags().BoolVarP(&useTLS, "tls", "s", false, "habilita HTTPS com certificado autoassinado ou customizado")
	serveCmd.Flags().StringVar(&tlsCert, "tls-cert", "", "caminho do arquivo PEM contendo o certificado público TLS")
	serveCmd.Flags().StringVar(&tlsKey, "tls-key", "", "caminho do arquivo PEM contendo a chave privada TLS")
	RootCmd.AddCommand(serveCmd)
}
