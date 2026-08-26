package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	adaptersftp "github.com/douglas/file-server/internal/adapters/sftp"
	"github.com/douglas/file-server/internal/core/services"
	"github.com/spf13/cobra"
)

var (
	sftpPort     int
	sftpHost     string
	sftpDirFlag  string
	sftpUser     string
	sftpPass     string
	sftpAuthKey  string
	sftpHostKey  string
	sftpReadOnly bool
)

// sftpCmd inicia o servidor SFTP sobre SSH.
var sftpCmd = &cobra.Command{
	Use:   "sftp [diretório]",
	Short: "Inicia o servidor SFTP seguro sobre SSH",
	Long:  `Inicia o servidor de transferência de arquivos SFTP seguro com autenticação, sandbox estrito e criptografia de canal SSHv2 para redes locais.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetDir, err := ResolveDirectory(args, sftpDirFlag)
		if err != nil {
			return err
		}

		user := sftpUser
		if user == "" {
			user = services.DefaultUsername
		}

		pass := sftpPass
		if pass == "" && sftpAuthKey == "" {
			generatedPass, err := services.GenerateRandomPassword(12)
			if err != nil {
				return fmt.Errorf("erro ao gerar senha temporária para SFTP: %w", err)
			}
			pass = generatedPass
		}

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		defer stop()

		return RunSFTPServerWithOptions(ctx, adaptersftp.ServerOptions{
			Host:      sftpHost,
			Port:      sftpPort,
			TargetDir: targetDir,
			User:      user,
			Pass:      pass,
			AuthKey:   sftpAuthKey,
			HostKey:   sftpHostKey,
			ReadOnly:  sftpReadOnly,
		})
	},
}

// RunSFTPServerWithOptions executa o servidor SFTP aguardando sinal de cancelamento pelo contexto.
func RunSFTPServerWithOptions(ctx context.Context, opts adaptersftp.ServerOptions) error {
	server, err := adaptersftp.NewServer(opts)
	if err != nil {
		return fmt.Errorf("erro ao inicializar servidor SFTP: %w", err)
	}

	return server.Run(ctx)
}

func init() {
	sftpCmd.Flags().IntVarP(&sftpPort, "port", "p", 2222, "porta na qual o servidor SFTP irá escutar")
	sftpCmd.Flags().StringVar(&sftpHost, "host", "0.0.0.0", "endereço do host para escuta")
	sftpCmd.Flags().StringVarP(&sftpDirFlag, "dir", "d", "", "caminho do diretório raiz a ser compartilhado")
	sftpCmd.Flags().StringVarP(&sftpUser, "user", "u", "", "usuário de autenticação SFTP")
	sftpCmd.Flags().StringVarP(&sftpPass, "pass", "P", "", "senha de autenticação SFTP")
	sftpCmd.Flags().StringVarP(&sftpAuthKey, "auth-key", "k", "", "caminho para chave pública SSH autorizada")
	sftpCmd.Flags().StringVar(&sftpHostKey, "host-key", "", "caminho para chave privada de host SSH (PEM)")
	sftpCmd.Flags().BoolVarP(&sftpReadOnly, "read-only", "r", false, "bloqueia uploads, renomeações e exclusões no SFTP")
	RootCmd.AddCommand(sftpCmd)
}
