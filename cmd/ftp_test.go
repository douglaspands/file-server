package cmd_test

import (
	"context"
	"net"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/douglas/file-server/cmd"
	adapterftp "github.com/douglas/file-server/internal/adapters/ftp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFTPCommand(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("execução com flag help", func(t *testing.T) {
		output, err := executeCommand("ftp", "--help")
		require.NoError(t, err)
		assert.Contains(t, output, "Inicia o servidor FTP")
		assert.Contains(t, output, "--tls")
		assert.Contains(t, output, "--tls-cert")
		assert.Contains(t, output, "--tls-key")
		assert.Contains(t, output, "--passive-ports")
		assert.Contains(t, output, "--read-only")
	})

	t.Run("falha com diretório inexistente", func(t *testing.T) {
		_, err := executeCommand("ftp", filepath.Join(tempDir, "inexistente"))
		assert.Error(t, err)
	})

	t.Run("falha com muitos argumentos", func(t *testing.T) {
		_, err := executeCommand("ftp", tempDir, "argumento_extra")
		assert.Error(t, err)
	})

	t.Run("execução com RunFTPServerWithOptions e cancelamento gracioso", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:0")
		require.NoError(t, err)
		freePort := l.Addr().(*net.TCPAddr).Port
		_ = l.Close()

		opts := adapterftp.ServerOptions{
			Host:      "127.0.0.1",
			Port:      freePort,
			TargetDir: tempDir,
			User:      "admin",
			Pass:      "admin123",
		}

		ctx, cancel := context.WithCancel(context.Background())
		doneChan := make(chan error, 1)

		go func() {
			doneChan <- cmd.RunFTPServerWithOptions(ctx, opts)
		}()

		time.Sleep(100 * time.Millisecond)
		cancel()

		select {
		case err := <-doneChan:
			assert.NoError(t, err)
		case <-time.After(3 * time.Second):
			t.Fatal("timeout ao aguardar encerramento do FTP")
		}
	})

	t.Run("RunFTPServerWithOptions com opções inválidas", func(t *testing.T) {
		opts := adapterftp.ServerOptions{
			TargetDir: filepath.Join(tempDir, "inexistente"),
		}
		err := cmd.RunFTPServerWithOptions(context.Background(), opts)
		assert.Error(t, err)
	})

	t.Run("execução de ftp via executeCommand com interrupção por sinal", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:0")
		require.NoError(t, err)
		freePort := l.Addr().(*net.TCPAddr).Port
		_ = l.Close()

		doneChan := make(chan error, 1)
		go func() {
			_, err := executeCommand("ftp", "--port", "19121", "--tls", tempDir)
			doneChan <- err
		}()

		time.Sleep(100 * time.Millisecond)
		_ = freePort
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)

		select {
		case err := <-doneChan:
			assert.NoError(t, err)
		case <-time.After(3 * time.Second):
			t.Fatal("timeout aguardando término de ftp")
		}
	})
}
