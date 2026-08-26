package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/douglas/file-server/internal/adapters/gui"
	"github.com/douglas/file-server/internal/version"
	"github.com/spf13/cobra"
)

var (
	guiHostFlag   string
	guiPortFlag   int
	guiDirFlag    string
	guiNoOpenFlag bool
)

// GUIOptions define as opções de execução da interface gráfica desktop.
type GUIOptions struct {
	Host       string
	Port       int
	InitialDir string
	NoOpen     bool
}

// guiCmd representa o subcomando 'gui' do File Server.
var guiCmd = &cobra.Command{
	Use:   "gui [diretório]",
	Short: "Inicia a interface gráfica desktop (GUI)",
	Long:  `Inicia a interface gráfica desktop moderna e nativa com controles visuais para serviços Web, FTP e SFTP, visualização escalável de múltiplos IPs e compartilhamento via QR Code.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetDir, err := ResolveDirectory(args, guiDirFlag)
		if err != nil {
			return err
		}

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		defer stop()

		return RunGUIWithOptions(ctx, GUIOptions{
			Host:       guiHostFlag,
			Port:       guiPortFlag,
			InitialDir: targetDir,
			NoOpen:     guiNoOpenFlag,
		})
	},
}

// RunGUIWithOptions inicializa e orquestra o servidor web da interface desktop e a janela nativa.
func RunGUIWithOptions(ctx context.Context, opts GUIOptions) error {
	if opts.Host == "" {
		opts.Host = "127.0.0.1"
	}

	broadcaster := gui.NewLogBroadcaster(500)
	// Configura log padrão para espelhar mensagens no terminal e no console da GUI
	originalLogOutput := log.Writer()
	log.SetOutput(io.MultiWriter(originalLogOutput, broadcaster))
	defer log.SetOutput(originalLogOutput)

	ctrl := gui.NewController(opts.InitialDir, gui.NewDefaultRunner(), broadcaster)
	picker := gui.NewNativeFolderPicker()

	server, err := gui.NewServer(ctrl, picker)
	if err != nil {
		return fmt.Errorf("erro ao inicializar servidor da GUI: %w", err)
	}

	guiURL, err := server.Start(opts.Host, opts.Port)
	if err != nil {
		return fmt.Errorf("erro ao iniciar listener da GUI: %w", err)
	}

	log.Printf("⚡ File Server Launcher Desktop v%s inicializado!", version.Get().Version)
	log.Printf("👉 Painel Gráfico Desktop: %s", guiURL)
	log.Printf("📁 Diretório Inicial: %s", opts.InitialDir)
	log.Printf("🛑 Pressione Ctrl+C para encerrar a aplicação.")

	// Abre a janela desktop dedicada se estiver em ambiente gráfico e auto-open habilitado
	if !opts.NoOpen && gui.IsDesktopEnvironment() {
		go func() {
			_, _ = gui.LaunchDesktopWindow(ctx, guiURL)
		}()
	}

	<-ctx.Done()
	log.Println("🛑 Encerrando aplicação e interface gráfica desktop...")

	_ = ctrl.StopServer()
	_ = server.Stop()
	log.Println("✅ Interface desktop encerrada com sucesso.")

	return nil
}

func init() {
	guiCmd.Flags().StringVar(&guiHostFlag, "host", "127.0.0.1", "endereço do host para o servidor da GUI")
	guiCmd.Flags().IntVarP(&guiPortFlag, "port", "p", 0, "porta do servidor da GUI (0 para porta aleatória disponível)")
	guiCmd.Flags().StringVarP(&guiDirFlag, "dir", "d", "", "caminho do diretório raiz inicial")
	guiCmd.Flags().BoolVar(&guiNoOpenFlag, "no-open", false, "não abre a janela desktop/navegador automaticamente")
	RootCmd.AddCommand(guiCmd)
}
