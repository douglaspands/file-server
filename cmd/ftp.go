package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	adapterftp "github.com/douglas/file-server/internal/adapters/ftp"
	"github.com/douglas/file-server/internal/core/services"
	"github.com/spf13/cobra"
)

var (
	ftpPort         int
	ftpHost         string
	ftpDirFlag      string
	ftpUser         string
	ftpPass         string
	ftpUseTLS       bool
	ftpTLSCert      string
	ftpTLSKey       string
	ftpPassivePorts string
	ftpReadOnly     bool
)

// ftpCmd inicia o servidor FTP/FTPS.
var ftpCmd = &cobra.Command{
	Use:   "ftp [diretório]",
	Short: "Inicia o servidor de arquivos FTP/FTPS",
	Long:  `Inicia o servidor FTP com suporte a conexões criptografadas FTPS (TLS 1.3 / 1.2), autenticação segura e modo somente leitura.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetDir, err := ResolveDirectory(args, ftpDirFlag)
		if err != nil {
			return err
		}

		user := ftpUser
		if user == "" {
			user = services.DefaultUsername
		}

		pass := ftpPass
		if pass == "" {
			generatedPass, err := services.GenerateRandomPassword(12)
			if err != nil {
				return fmt.Errorf("erro ao gerar senha temporária para FTP: %w", err)
			}
			pass = generatedPass
		}

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		defer stop()

		return RunFTPServerWithOptions(ctx, adapterftp.ServerOptions{
			Host:         ftpHost,
			Port:         ftpPort,
			TargetDir:    targetDir,
			User:         user,
			Pass:         pass,
			UseTLS:       ftpUseTLS,
			TLSCert:      ftpTLSCert,
			TLSKey:       ftpTLSKey,
			PassivePorts: ftpPassivePorts,
			ReadOnly:     ftpReadOnly,
		})
	},
}

// RunFTPServerWithOptions executa o servidor FTP aguardando sinal de cancelamento pelo contexto.
func RunFTPServerWithOptions(ctx context.Context, opts adapterftp.ServerOptions) error {
	server, err := adapterftp.NewServer(opts)
	if err != nil {
		return fmt.Errorf("erro ao inicializar servidor FTP: %w", err)
	}

	return server.Run(ctx)
}

func init() {
	ftpCmd.Flags().IntVarP(&ftpPort, "port", "p", 2121, "porta na qual o servidor FTP irá escutar")
	ftpCmd.Flags().StringVar(&ftpHost, "host", "0.0.0.0", "endereço do host para escuta")
	ftpCmd.Flags().StringVarP(&ftpDirFlag, "dir", "d", "", "caminho do diretório raiz a ser compartilhado")
	ftpCmd.Flags().StringVarP(&ftpUser, "user", "u", "", "usuário de autenticação FTP")
	ftpCmd.Flags().StringVarP(&ftpPass, "pass", "P", "", "senha de autenticação FTP")
	ftpCmd.Flags().BoolVarP(&ftpUseTLS, "tls", "s", false, "habilita FTPS com TLS autoassinado ou customizado")
	ftpCmd.Flags().StringVar(&ftpTLSCert, "tls-cert", "", "caminho do certificado público TLS (PEM)")
	ftpCmd.Flags().StringVar(&ftpTLSKey, "tls-key", "", "caminho da chave privada TLS (PEM)")
	ftpCmd.Flags().StringVar(&ftpPassivePorts, "passive-ports", "", "faixa de portas para modo passivo (ex: 50000-50100)")
	ftpCmd.Flags().BoolVarP(&ftpReadOnly, "read-only", "r", false, "bloqueia uploads, renomeações e exclusões no FTP")
	RootCmd.AddCommand(ftpCmd)
}
