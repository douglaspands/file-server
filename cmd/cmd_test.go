package cmd_test

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/douglas/file-server/cmd"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func executeCommand(args ...string) (string, error) {
	buf := new(bytes.Buffer)
	cmd.RootCmd.SetOut(buf)
	cmd.RootCmd.SetErr(buf)
	cmd.RootCmd.SetArgs(args)

	// Reseta flags de help para isolamento entre testes
	if f := cmd.RootCmd.Flags().Lookup("help"); f != nil {
		_ = f.Value.Set("false")
	}
	for _, c := range cmd.RootCmd.Commands() {
		if f := c.Flags().Lookup("help"); f != nil {
			_ = f.Value.Set("false")
		}
	}

	err := cmd.RootCmd.Execute()
	return buf.String(), err
}

func TestCommands(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("Given version command When executed with buffer Then outputs version string", func(t *testing.T) {
		output, err := executeCommand("version")
		require.NoError(t, err)
		assert.Contains(t, output, "File Server version")
	})

	t.Run("Given version command with --json flag When executed Then succeeds", func(t *testing.T) {
		output, err := executeCommand("version", "--json")
		require.NoError(t, err)
		assert.Contains(t, output, `"version":`)
	})

	t.Run("Given serve command with invalid directory When executed Then returns error", func(t *testing.T) {
		_, err := executeCommand("serve", filepath.Join(tempDir, "inexistente"))
		assert.Error(t, err)
	})

	t.Run("Given root command with help flag When executed Then shows help and TLS flags", func(t *testing.T) {
		output, err := executeCommand("--help")
		require.NoError(t, err)
		assert.Contains(t, output, "file-server")
		assert.Contains(t, output, "--tls")
		assert.Contains(t, output, "--tls-cert")
		assert.Contains(t, output, "--tls-key")
	})

	t.Run("Given serve command with help flag When executed Then shows serve help and TLS flags", func(t *testing.T) {
		output, err := executeCommand("serve", "--help")
		require.NoError(t, err)
		assert.Contains(t, output, "Inicia o servidor")
		assert.Contains(t, output, "--tls")
		assert.Contains(t, output, "--tls-cert")
		assert.Contains(t, output, "--tls-key")
	})

	t.Run("Given serve command with only --tls-cert flag When executed Then returns validation error", func(t *testing.T) {
		_, err := executeCommand("serve", "--tls-cert", "cert.pem", tempDir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ambos --tls-cert e --tls-key devem ser fornecidos")
	})
}

func TestExecuteRoot(t *testing.T) {
	t.Run("Given valid args When executing root Then runs without error", func(t *testing.T) {
		err := cmd.ExecuteRoot([]string{"--help"})
		require.NoError(t, err)
	})
}

func TestMousetrapDisabled(t *testing.T) {
	t.Run("Given root command initialization When checking mousetrap Then MousetrapHelpText is empty", func(t *testing.T) {
		assert.Empty(t, cobra.MousetrapHelpText, "Cobra MousetrapHelpText deve estar vazio para permitir clique duplo no Windows Explorer")
	})
}
