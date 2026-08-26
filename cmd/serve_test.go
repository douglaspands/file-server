package cmd_test

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/douglas/file-server/cmd"
	"github.com/douglas/file-server/internal/core/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpandHomeDir(t *testing.T) {
	home, err := os.UserHomeDir()
	require.NoError(t, err)
	require.NotEmpty(t, home)

	t.Run("expande til sozinho (~)", func(t *testing.T) {
		expanded, err := cmd.ExpandHomeDir("~")
		require.NoError(t, err)
		assert.Equal(t, home, expanded)
	})

	t.Run("expande caminho com prefixo de til (~/pasta)", func(t *testing.T) {
		expanded, err := cmd.ExpandHomeDir("~/documentos/fotos")
		require.NoError(t, err)
		assert.Equal(t, filepath.Join(home, "documentos", "fotos"), expanded)
	})

	t.Run("mantém caminho relativo comum inalterado", func(t *testing.T) {
		expanded, err := cmd.ExpandHomeDir("./dados/teste")
		require.NoError(t, err)
		assert.Equal(t, "./dados/teste", expanded)
	})

	t.Run("mantém caminho absoluto inalterado", func(t *testing.T) {
		expanded, err := cmd.ExpandHomeDir("/var/log")
		require.NoError(t, err)
		assert.Equal(t, "/var/log", expanded)
	})
}

func TestResolveDirectory(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("fallback para diretório atual quando vazio", func(t *testing.T) {
		dir, err := cmd.ResolveDirectory(nil, "")
		require.NoError(t, err)
		assert.NotEmpty(t, dir)
	})

	t.Run("argumento posicional válido", func(t *testing.T) {
		dir, err := cmd.ResolveDirectory([]string{tempDir}, "")
		require.NoError(t, err)
		assert.Equal(t, tempDir, dir)
	})

	t.Run("flag explícita válida", func(t *testing.T) {
		dir, err := cmd.ResolveDirectory(nil, tempDir)
		require.NoError(t, err)
		assert.Equal(t, tempDir, dir)
	})

	t.Run("caminho relativo com ponto e barra (./)", func(t *testing.T) {
		dir, err := cmd.ResolveDirectory([]string{"./"}, "")
		require.NoError(t, err)
		cwd, err := os.Getwd()
		require.NoError(t, err)
		assert.Equal(t, cwd, dir)
	})

	t.Run("caminho relativo de diretório pai (../)", func(t *testing.T) {
		subDir := filepath.Join(tempDir, "subnivel1", "subnivel2")
		require.NoError(t, os.MkdirAll(subDir, 0755))

		// Resolve ../ a partir de subDir
		originalWd, err := os.Getwd()
		require.NoError(t, err)
		defer func() { _ = os.Chdir(originalWd) }()

		require.NoError(t, os.Chdir(subDir))
		dir, err := cmd.ResolveDirectory([]string{"../"}, "")
		require.NoError(t, err)
		expectedParent := filepath.Join(tempDir, "subnivel1")
		assert.Equal(t, expectedParent, dir)
	})

	t.Run("expansão de til (~) para home directory existente", func(t *testing.T) {
		home, err := os.UserHomeDir()
		require.NoError(t, err)

		dir, err := cmd.ResolveDirectory([]string{"~"}, "")
		require.NoError(t, err)
		assert.Equal(t, home, dir)
	})

	t.Run("prioridade de argumento posicional sobre flag", func(t *testing.T) {
		otherDir := t.TempDir()
		dir, err := cmd.ResolveDirectory([]string{tempDir}, otherDir)
		require.NoError(t, err)
		assert.Equal(t, tempDir, dir)
	})

	t.Run("erro com diretório inexistente", func(t *testing.T) {
		_, err := cmd.ResolveDirectory([]string{filepath.Join(tempDir, "fantasma")}, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "não existe")
	})

	t.Run("erro quando o caminho for um arquivo", func(t *testing.T) {
		tempFile := filepath.Join(tempDir, "arquivo.txt")
		require.NoError(t, os.WriteFile(tempFile, []byte("data"), 0644))

		_, err := cmd.ResolveDirectory([]string{tempFile}, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "não é um diretório")
	})
}

func TestSetupServer(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("sucesso ao inicializar servidor com diretório válido", func(t *testing.T) {
		server, err := cmd.SetupServer("127.0.0.1", 9091, tempDir)
		require.NoError(t, err)
		assert.Equal(t, "127.0.0.1:9091", server.Addr)
		assert.Nil(t, server.TLSConfig)
	})

	t.Run("erro ao inicializar servidor com diretório inexistente", func(t *testing.T) {
		_, err := cmd.SetupServer("127.0.0.1", 9091, filepath.Join(tempDir, "nao_existe"))
		assert.Error(t, err)
	})
}

func TestSetupServerWithOptions(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("inicialização com TLS autoassinado", func(t *testing.T) {
		server, err := cmd.SetupServerWithOptions(cmd.ServerOptions{
			Host:      "127.0.0.1",
			Port:      9092,
			TargetDir: tempDir,
			UseTLS:    true,
		})
		require.NoError(t, err)
		assert.Equal(t, "127.0.0.1:9092", server.Addr)
		require.NotNil(t, server.TLSConfig)
		assert.Len(t, server.TLSConfig.Certificates, 1)
	})

	t.Run("inicialização com certificados customizados válidos", func(t *testing.T) {
		cert, err := services.GenerateSelfSignedCertificate("custom.local")
		require.NoError(t, err)

		certPEM, keyPEM, err := services.CertificateToPEM(cert)
		require.NoError(t, err)

		certPath := filepath.Join(tempDir, "custom_cert.pem")
		keyPath := filepath.Join(tempDir, "custom_key.pem")

		require.NoError(t, os.WriteFile(certPath, certPEM, 0600))
		require.NoError(t, os.WriteFile(keyPath, keyPEM, 0600))

		server, err := cmd.SetupServerWithOptions(cmd.ServerOptions{
			Host:      "127.0.0.1",
			Port:      9093,
			TargetDir: tempDir,
			TLSCert:   certPath,
			TLSKey:    keyPath,
		})
		require.NoError(t, err)
		require.NotNil(t, server.TLSConfig)
		assert.Len(t, server.TLSConfig.Certificates, 1)
	})

	t.Run("erro quando apenas --tls-cert é fornecido", func(t *testing.T) {
		_, err := cmd.SetupServerWithOptions(cmd.ServerOptions{
			Host:      "127.0.0.1",
			Port:      9094,
			TargetDir: tempDir,
			TLSCert:   "cert.pem",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ambos --tls-cert e --tls-key devem ser fornecidos")
	})

	t.Run("erro quando apenas --tls-key é fornecido", func(t *testing.T) {
		_, err := cmd.SetupServerWithOptions(cmd.ServerOptions{
			Host:      "127.0.0.1",
			Port:      9095,
			TargetDir: tempDir,
			TLSKey:    "key.pem",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ambos --tls-cert e --tls-key devem ser fornecidos")
	})

	t.Run("erro quando certificados customizados não existem", func(t *testing.T) {
		_, err := cmd.SetupServerWithOptions(cmd.ServerOptions{
			Host:      "127.0.0.1",
			Port:      9096,
			TargetDir: tempDir,
			TLSCert:   filepath.Join(tempDir, "inexistente.pem"),
			TLSKey:    filepath.Join(tempDir, "inexistente.key"),
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao carregar certificados TLS")
	})
}

func TestRunServerWithContext(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("encerramento gracioso ao cancelar contexto em HTTP", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()

		err := cmd.RunServerWithContext(ctx, "127.0.0.1", 19090, tempDir)
		require.NoError(t, err)
	})

	t.Run("encerramento gracioso ao cancelar contexto com TLS ativo", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		serverErr := make(chan error, 1)
		go func() {
			serverErr <- cmd.RunServerWithOptions(ctx, cmd.ServerOptions{
				Host:      "127.0.0.1",
				Port:      19092,
				TargetDir: tempDir,
				UseTLS:    true,
			})
		}()

		// Aguarda o servidor iniciar e faz uma requisição HTTPS
		time.Sleep(100 * time.Millisecond)

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr, Timeout: 2 * time.Second}

		resp, err := client.Get("https://127.0.0.1:19092/api/health")
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			_ = resp.Body.Close()
		}

		cancel()
		require.NoError(t, <-serverErr)
	})

	t.Run("erro ao tentar iniciar com diretório inválido", func(t *testing.T) {
		ctx := context.Background()
		err := cmd.RunServerWithContext(ctx, "127.0.0.1", 19091, filepath.Join(tempDir, "invalid"))
		assert.Error(t, err)
	})
}

func TestFormatAccessURLs(t *testing.T) {
	t.Run("host 0.0.0.0 em HTTP", func(t *testing.T) {
		localURL, lanURLs := cmd.FormatAccessURLs("0.0.0.0", 8080, false)
		assert.Equal(t, "http://127.0.0.1:8080", localURL)
		for _, u := range lanURLs {
			assert.Contains(t, u, "http://")
			assert.Contains(t, u, ":8080")
		}
	})

	t.Run("host 0.0.0.0 em HTTPS", func(t *testing.T) {
		localURL, lanURLs := cmd.FormatAccessURLs("0.0.0.0", 8443, true)
		assert.Equal(t, "https://127.0.0.1:8443", localURL)
		for _, u := range lanURLs {
			assert.Contains(t, u, "https://")
			assert.Contains(t, u, ":8443")
		}
	})

	t.Run("host específico (127.0.0.1)", func(t *testing.T) {
		localURL, lanURLs := cmd.FormatAccessURLs("127.0.0.1", 9000, false)
		assert.Equal(t, "http://127.0.0.1:9000", localURL)
		assert.Empty(t, lanURLs)
	})
}

func TestLogStartupBanner(t *testing.T) {
	t.Run("executa sem erros em HTTP", func(t *testing.T) {
		assert.NotPanics(t, func() {
			cmd.LogStartupBanner(cmd.ServerOptions{
				Host:      "0.0.0.0",
				Port:      8080,
				TargetDir: "/tmp",
			}, false)
		})
	})

	t.Run("executa sem erros em HTTPS", func(t *testing.T) {
		assert.NotPanics(t, func() {
			cmd.LogStartupBanner(cmd.ServerOptions{
				Host:      "0.0.0.0",
				Port:      8443,
				TargetDir: "/tmp",
				UseTLS:    true,
			}, true)
		})
	})
}
