package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/douglas/file-server/internal/adapters/gui"
	"github.com/spf13/cobra"
)

var (
	cfgFile     string
	verbose     bool
	rootDirFlag string
	rootPort    int
	rootHost    string
	rootUseTLS  bool
	rootTLSCert string
	rootTLSKey  string
	rootUseGUI  bool
)

// RootCmd representa o comando base executado sem subcomandos.
var RootCmd = &cobra.Command{
	Use:   "file-server [diretório]",
	Short: "Servidor web de arquivos e diretórios de alta performance",
	Long:  `File Server é uma aplicação portátil de compartilhamento e navegação de arquivos com interface gráfica desktop, streaming de ZIP, uploads, downloads e TLS/HTTPS.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetDir, err := ResolveDirectory(args, rootDirFlag)
		if err != nil {
			return err
		}

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		defer stop()

		// Se a flag --gui for explícita ou estiver em ambiente desktop sem flags de servidor CLI
		shouldLaunchGUI := rootUseGUI || (gui.IsDesktopEnvironment() && !cmd.Flags().Changed("port") && !cmd.Flags().Changed("host") && !cmd.Flags().Changed("tls") && !cmd.Flags().Changed("tls-cert") && !cmd.Flags().Changed("tls-key"))

		if shouldLaunchGUI {
			return RunGUIWithOptions(ctx, GUIOptions{
				Host:       "127.0.0.1",
				Port:       0,
				InitialDir: targetDir,
				NoOpen:     false,
			})
		}

		// Modo CLI Headless tradicional
		return RunServerWithOptions(ctx, ServerOptions{
			Host:      rootHost,
			Port:      rootPort,
			TargetDir: targetDir,
			UseTLS:    rootUseTLS,
			TLSCert:   rootTLSCert,
			TLSKey:    rootTLSKey,
		})
	},
}

// ExecuteRoot executa o comando raiz com argumentos fornecidos.
func ExecuteRoot(args []string) error {
	RootCmd.SetArgs(args)
	return RootCmd.Execute()
}

// Execute executa o comando raiz e trata eventuais erros de inicialização.
func Execute() {
	if err := ExecuteRoot(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Desabilita o mousetrap do Cobra para permitir execução por duplo clique no Windows Explorer
	cobra.MousetrapHelpText = ""

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "arquivo de configuração (default é $HOME/.file-server.yaml)")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "habilita logs detalhados")
	RootCmd.Flags().StringVarP(&rootDirFlag, "dir", "d", "", "caminho do diretório raiz a ser compartilhado")
	RootCmd.Flags().IntVarP(&rootPort, "port", "p", 8080, "porta na qual o servidor irá escutar")
	RootCmd.Flags().StringVar(&rootHost, "host", "0.0.0.0", "endereço do host para escuta")
	RootCmd.Flags().BoolVarP(&rootUseTLS, "tls", "s", false, "habilita HTTPS com certificado autoassinado ou customizado")
	RootCmd.Flags().StringVar(&rootTLSCert, "tls-cert", "", "caminho do arquivo PEM contendo o certificado público TLS")
	RootCmd.Flags().StringVar(&rootTLSKey, "tls-key", "", "caminho do arquivo PEM contendo a chave privada TLS")
	RootCmd.Flags().BoolVar(&rootUseGUI, "gui", false, "inicializa a interface gráfica desktop (GUI)")
}
