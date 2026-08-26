package ftp_test

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	adapterFtp "github.com/douglas/file-server/internal/adapters/ftp"
	"github.com/douglas/file-server/internal/core/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsePassivePorts(t *testing.T) {
	t.Run("faixa válida", func(t *testing.T) {
		pr, err := adapterFtp.ParsePassivePorts("50000-50100")
		require.NoError(t, err)
		require.NotNil(t, pr)
		assert.Equal(t, 50000, pr.Start)
		assert.Equal(t, 50100, pr.End)
	})

	t.Run("string vazia ou zero retorna nil", func(t *testing.T) {
		pr, err := adapterFtp.ParsePassivePorts("")
		require.NoError(t, err)
		assert.Nil(t, pr)

		pr, err = adapterFtp.ParsePassivePorts("0")
		require.NoError(t, err)
		assert.Nil(t, pr)
	})

	t.Run("faixa inválida", func(t *testing.T) {
		_, err := adapterFtp.ParsePassivePorts("invalido")
		assert.Error(t, err)

		_, err = adapterFtp.ParsePassivePorts("50000-40000")
		assert.Error(t, err)

		_, err = adapterFtp.ParsePassivePorts("0-100")
		assert.Error(t, err)

		_, err = adapterFtp.ParsePassivePorts("50000-70000")
		assert.Error(t, err)
	})
}

func TestDriver_Unit(t *testing.T) {
	tempDir := t.TempDir()

	testFile := filepath.Join(tempDir, "arquivo.txt")
	err := os.WriteFile(testFile, []byte("conteudo teste"), 0644)
	require.NoError(t, err)

	t.Run("criação com diretório inválido", func(t *testing.T) {
		_, err := adapterFtp.NewDriver(adapterFtp.ServerOptions{TargetDir: filepath.Join(tempDir, "inexistente")})
		assert.Error(t, err)

		_, err = adapterFtp.NewDriver(adapterFtp.ServerOptions{TargetDir: testFile})
		assert.Error(t, err)
	})

	t.Run("criação com opções de TLS", func(t *testing.T) {
		// TLS autoassinado
		d, err := adapterFtp.NewDriver(adapterFtp.ServerOptions{
			TargetDir: tempDir,
			UseTLS:    true,
			Host:      "127.0.0.1",
			Port:      2121,
		})
		require.NoError(t, err)
		tlsCfg, err := d.GetTLSConfig()
		require.NoError(t, err)
		assert.NotNil(t, tlsCfg)

		// TLS com certificados em disco
		cert, err := services.GenerateSelfSignedCertificate("127.0.0.1")
		require.NoError(t, err)
		certPEM, keyPEM, err := services.CertificateToPEM(cert)
		require.NoError(t, err)

		certPath := filepath.Join(tempDir, "cert.pem")
		keyPath := filepath.Join(tempDir, "key.pem")
		require.NoError(t, os.WriteFile(certPath, certPEM, 0600))
		require.NoError(t, os.WriteFile(keyPath, keyPEM, 0600))

		dCustom, err := adapterFtp.NewDriver(adapterFtp.ServerOptions{
			TargetDir: tempDir,
			TLSCert:   certPath,
			TLSKey:    keyPath,
		})
		require.NoError(t, err)
		customCfg, err := dCustom.GetTLSConfig()
		require.NoError(t, err)
		assert.NotNil(t, customCfg)

		// Falha com certificado incompleto
		_, err = adapterFtp.NewDriver(adapterFtp.ServerOptions{
			TargetDir: tempDir,
			TLSCert:   certPath,
		})
		assert.Error(t, err)
	})

	t.Run("GetSettings, ClientConnected e ClientDisconnected", func(t *testing.T) {
		d, err := adapterFtp.NewDriver(adapterFtp.ServerOptions{
			TargetDir:    tempDir,
			Host:         "127.0.0.1",
			Port:         2121,
			PassivePorts: "50000-50050",
		})
		require.NoError(t, err)

		settings, err := d.GetSettings()
		require.NoError(t, err)
		assert.Equal(t, "127.0.0.1:2121", settings.ListenAddr)

		welcome, err := d.ClientConnected(nil)
		require.NoError(t, err)
		assert.NotEmpty(t, welcome)

		d.ClientDisconnected(nil)
	})

	t.Run("AuthUser com credenciais corretas e incorretas", func(t *testing.T) {
		d, err := adapterFtp.NewDriver(adapterFtp.ServerOptions{
			TargetDir: tempDir,
			User:      "ftpuser",
			Pass:      "ftppass",
			ReadOnly:  false,
		})
		require.NoError(t, err)

		// Sucesso
		fs, err := d.AuthUser(nil, "ftpuser", "ftppass")
		require.NoError(t, err)
		assert.NotNil(t, fs)

		// Falha
		_, err = d.AuthUser(nil, "ftpuser", "senha_errada")
		assert.Error(t, err)
		assert.ErrorIs(t, err, adapterFtp.ErrInvalidCredentials)
	})

	t.Run("AuthUser em modo ReadOnly", func(t *testing.T) {
		d, err := adapterFtp.NewDriver(adapterFtp.ServerOptions{
			TargetDir: tempDir,
			ReadOnly:  true,
		})
		require.NoError(t, err)

		fs, err := d.AuthUser(nil, "", "")
		require.NoError(t, err)
		assert.NotNil(t, fs)

		// Tentativa de escrita deve falhar
		f, err := fs.Create("teste_ro.txt")
		if err == nil {
			_, err = f.Write([]byte("teste"))
			_ = f.Close()
		}
		assert.Error(t, err)
	})
}

func TestFTPServer_LifecycleAndBanner(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("inicialização e exibição de banner", func(t *testing.T) {
		opts := adapterFtp.ServerOptions{
			Host:         "0.0.0.0",
			Port:         2121,
			TargetDir:    tempDir,
			User:         "admin",
			Pass:         "senha123",
			ReadOnly:     true,
			PassivePorts: "50000-50100",
			UseTLS:       true,
		}

		server, err := adapterFtp.NewServer(opts)
		require.NoError(t, err)
		assert.NotNil(t, server)

		adapterFtp.LogStartupBanner(opts, true)
		adapterFtp.LogStartupBanner(adapterFtp.ServerOptions{
			Host:      "127.0.0.1",
			Port:      2121,
			TargetDir: tempDir,
			User:      "user",
		}, false)
	})

	t.Run("execução graciosa via Run com cancelamento por contexto", func(t *testing.T) {
		// Procura uma porta livre para testar o listener
		l, err := net.Listen("tcp", "127.0.0.1:0")
		require.NoError(t, err)
		freePort := l.Addr().(*net.TCPAddr).Port
		_ = l.Close()

		opts := adapterFtp.ServerOptions{
			Host:      "127.0.0.1",
			Port:      freePort,
			TargetDir: tempDir,
			User:      "user",
			Pass:      "pass",
		}

		server, err := adapterFtp.NewServer(opts)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		doneChan := make(chan error, 1)

		go func() {
			doneChan <- server.Run(ctx)
		}()

		time.Sleep(100 * time.Millisecond)
		assert.NotEmpty(t, server.Addr())

		cancel()

		select {
		case err := <-doneChan:
			assert.NoError(t, err)
		case <-time.After(4 * time.Second):
			t.Fatal("tempo limite excedido ao encerrar servidor FTP")
		}
	})
}
