package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/douglas/file-server/internal/version"
	"github.com/spf13/cobra"
)

var jsonOutput bool

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Exibe a versão da aplicação e informações de compilação",
	Long:  `Exibe a versão semântica, hash do commit git, data de compilação e plataforma.`,
	Run: func(cmd *cobra.Command, args []string) {
		info := version.Get()
		if jsonOutput {
			data, err := json.MarshalIndent(info, "", "  ")
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "erro ao serializar versão: %v\n", err)
				return
			}
			fmt.Fprintln(cmd.OutOrStdout(), string(data))
			return
		}
		fmt.Fprintln(cmd.OutOrStdout(), info.String())
	},
}

func init() {
	versionCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "exibe a versão em formato JSON")
	RootCmd.AddCommand(versionCmd)
}
