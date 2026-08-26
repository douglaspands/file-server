package cmd_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/douglas/file-server/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommands(t *testing.T) {
	t.Run("Given version command When executed with buffer Then outputs version string", func(t *testing.T) {
		buf := new(bytes.Buffer)
		cmd.RootCmd.SetOut(buf)
		cmd.RootCmd.SetErr(buf)
		cmd.RootCmd.SetArgs([]string{"version"})

		err := cmd.RootCmd.Execute()
		require.NoError(t, err)
	})

	t.Run("Given version command with --json flag When executed Then succeeds", func(t *testing.T) {
		buf := new(bytes.Buffer)
		cmd.RootCmd.SetOut(buf)
		cmd.RootCmd.SetErr(buf)
		cmd.RootCmd.SetArgs([]string{"version", "--json"})

		err := cmd.RootCmd.Execute()
		require.NoError(t, err)
	})

	t.Run("Given root command with help flag When executed Then shows help", func(t *testing.T) {
		buf := new(bytes.Buffer)
		cmd.RootCmd.SetOut(buf)
		cmd.RootCmd.SetErr(buf)
		cmd.RootCmd.SetArgs([]string{"--help"})

		err := cmd.RootCmd.Execute()
		require.NoError(t, err)
		assert.Contains(t, buf.String(), "file-server")
	})

	t.Run("Given serve command with help flag When executed Then shows serve help", func(t *testing.T) {
		buf := new(bytes.Buffer)
		cmd.RootCmd.SetOut(buf)
		cmd.RootCmd.SetErr(buf)
		cmd.RootCmd.SetArgs([]string{"serve", "--help"})

		err := cmd.RootCmd.Execute()
		require.NoError(t, err)
		assert.Contains(t, buf.String(), "Inicia o servidor HTTP")
	})

	t.Run("Given server setup When configuring host and port Then initializes successfully", func(t *testing.T) {
		server, err := cmd.SetupServer("127.0.0.1", 9090)
		require.NoError(t, err)
		assert.Equal(t, "127.0.0.1:9090", server.Addr)
	})

	t.Run("Given running server with context When context cancelled Then shuts down gracefully", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()

		err := cmd.RunServerWithContext(ctx, "127.0.0.1", 18080)
		require.NoError(t, err)
	})
}

func TestExecuteRoot(t *testing.T) {
	t.Run("Given valid args When executing root Then runs without error", func(t *testing.T) {
		err := cmd.ExecuteRoot([]string{"--help"})
		require.NoError(t, err)
	})
}
