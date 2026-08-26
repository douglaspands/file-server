package cmd_test

import (
	"context"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/douglas/file-server/cmd"
	adaptersftp "github.com/douglas/file-server/internal/adapters/sftp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSFTPCommand(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("execução com flag help", func(t *testing.T) {
		output, err := executeCommand("sftp", "--help")
		require.NoError(t, err)
		assert.Contains(t, output, "Inicia o servidor de transferência de arquivos SFTP")
		assert.Contains(t, output, "--auth-key")
		assert.Contains(t, output, "--host-key")
		assert.Contains(t, output, "--read-only")
	})

	t.Run("falha com diretório inexistente", func(t *testing.T) {
		_, err := executeCommand("sftp", filepath.Join(tempDir, "inexistente"))
		assert.Error(t, err)
	})

	t.Run("falha com muitos argumentos", func(t *testing.T) {
		_, err := executeCommand("sftp", tempDir, "argumento_extra")
		assert.Error(t, err)
	})

	t.Run("execução com RunSFTPServerWithOptions e cancelamento gracioso", func(t *testing.T) {
		opts := adaptersftp.ServerOptions{
			Host:      "127.0.0.1",
			Port:      0,
			TargetDir: tempDir,
			User:      "admin",
			Pass:      "admin123",
		}

		ctx, cancel := context.WithCancel(context.Background())
		doneChan := make(chan error, 1)

		go func() {
			doneChan <- cmd.RunSFTPServerWithOptions(ctx, opts)
		}()

		time.Sleep(100 * time.Millisecond)
		cancel()

		select {
		case err := <-doneChan:
			assert.NoError(t, err)
		case <-time.After(3 * time.Second):
			t.Fatal("timeout ao aguardar encerramento do SFTP")
		}
	})

	t.Run("RunSFTPServerWithOptions com opções inválidas", func(t *testing.T) {
		opts := adaptersftp.ServerOptions{
			TargetDir: filepath.Join(tempDir, "inexistente"),
		}
		err := cmd.RunSFTPServerWithOptions(context.Background(), opts)
		assert.Error(t, err)
	})

	t.Run("execução de sftp via executeCommand com interrupção por sinal", func(t *testing.T) {
		doneChan := make(chan error, 1)
		go func() {
			_, err := executeCommand("sftp", "--port", "0", tempDir)
			doneChan <- err
		}()

		time.Sleep(100 * time.Millisecond)
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)

		select {
		case err := <-doneChan:
			assert.NoError(t, err)
		case <-time.After(3 * time.Second):
			t.Fatal("timeout aguardando término de sftp")
		}
	})
}
