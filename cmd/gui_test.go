package cmd_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/douglas/file-server/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGUICommand(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("Given gui command with --help When executed Then displays help information", func(t *testing.T) {
		output, err := executeCommand("gui", "--help")
		require.NoError(t, err)
		assert.Contains(t, output, "interface gráfica desktop")
		assert.Contains(t, output, "--no-open")
	})

	t.Run("Given gui command with invalid directory When executed Then returns error", func(t *testing.T) {
		_, err := executeCommand("gui", filepath.Join(tempDir, "nao_existe_dir_xyz"))
		assert.Error(t, err)
	})

	t.Run("Given RunGUIWithOptions with cancelled context When executed Then starts and stops cleanly", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		// Cancela após pequeno intervalo para permitir inicialização e shutdown gracioso
		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()

		err := cmd.RunGUIWithOptions(ctx, cmd.GUIOptions{
			Host:       "127.0.0.1",
			Port:       0,
			InitialDir: tempDir,
			NoOpen:     true,
		})

		assert.NoError(t, err)
	})
}
