package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	verbose bool
)

// RootCmd representa o comando base executado sem subcomandos.
var RootCmd = &cobra.Command{
	Use:   "file-server",
	Short: "File Server com interface web moderna",
	Long:  `File Server é uma aplicação de gerenciamento e servidor com suporte a CLI e interface web.`,
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
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "arquivo de configuração (default é $HOME/.file-server.yaml)")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "habilita logs detalhados")
}
