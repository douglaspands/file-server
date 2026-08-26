package services_test

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/douglas/file-server/internal/core/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLANIPAddresses(t *testing.T) {
	ips := services.GetLANIPAddresses()
	assert.NotEmpty(t, ips, "deve retornar pelo menos os IPs de loopback")

	var foundLoopbackV4, foundLoopbackV6 bool
	for _, ip := range ips {
		if ip.String() == "127.0.0.1" {
			foundLoopbackV4 = true
		}
		if ip.String() == "::1" {
			foundLoopbackV6 = true
		}
	}
	assert.True(t, foundLoopbackV4, "deve conter 127.0.0.1")
	assert.True(t, foundLoopbackV6, "deve conter ::1")
}

func TestGenerateSelfSignedCertificate(t *testing.T) {
	t.Run("Gera certificado autoassinado com hosts padrão e customizados", func(t *testing.T) {
		cert, err := services.GenerateSelfSignedCertificate("192.168.1.100", "meuservidor.lan", "10.0.0.5:8443", "")
		require.NoError(t, err)
		require.NotEmpty(t, cert.Certificate)
		require.NotNil(t, cert.PrivateKey)

		x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
		require.NoError(t, err)

		assert.Equal(t, "File Server Local", x509Cert.Subject.CommonName)
		assert.Contains(t, x509Cert.Subject.Organization, "File Server Local")
		assert.Contains(t, x509Cert.Subject.Country, "BR")

		// Valida validade
		now := time.Now()
		assert.True(t, now.After(x509Cert.NotBefore))
		assert.True(t, now.Before(x509Cert.NotAfter))

		// Valida DNS Names
		assert.Contains(t, x509Cert.DNSNames, "localhost")
		assert.Contains(t, x509Cert.DNSNames, "*.localhost")
		assert.Contains(t, x509Cert.DNSNames, "meuservidor.lan")

		// Valida IPs
		var ipStrings []string
		for _, ip := range x509Cert.IPAddresses {
			ipStrings = append(ipStrings, ip.String())
		}
		assert.Contains(t, ipStrings, "127.0.0.1")
		assert.Contains(t, ipStrings, "::1")
		assert.Contains(t, ipStrings, "192.168.1.100")
		assert.Contains(t, ipStrings, "10.0.0.5")

		// Valida KeyUsage e ExtKeyUsage
		assert.Equal(t, x509.KeyUsageKeyEncipherment|x509.KeyUsageDigitalSignature, x509Cert.KeyUsage)
		assert.Contains(t, x509Cert.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
	})
}

func TestBuildTLSConfig(t *testing.T) {
	cert, err := services.GenerateSelfSignedCertificate()
	require.NoError(t, err)

	cfg := services.BuildTLSConfig(cert)
	require.NotNil(t, cfg)

	assert.Equal(t, uint16(tls.VersionTLS12), cfg.MinVersion)
	assert.Contains(t, cfg.NextProtos, "h2")
	assert.Contains(t, cfg.NextProtos, "http/1.1")
	assert.NotEmpty(t, cfg.CipherSuites)
	assert.Len(t, cfg.Certificates, 1)
}

func TestCreateSelfSignedTLSConfig(t *testing.T) {
	cfg, err := services.CreateSelfSignedTLSConfig("localhost", "127.0.0.1")
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Len(t, cfg.Certificates, 1)
}

func TestCertificateToPEMAndLoadTLSConfig(t *testing.T) {
	t.Run("Exporta para PEM e carrega com sucesso via LoadTLSConfig", func(t *testing.T) {
		cert, err := services.GenerateSelfSignedCertificate("teste.local")
		require.NoError(t, err)

		certPEM, keyPEM, err := services.CertificateToPEM(cert)
		require.NoError(t, err)
		assert.NotEmpty(t, certPEM)
		assert.NotEmpty(t, keyPEM)

		tempDir := t.TempDir()
		certPath := filepath.Join(tempDir, "cert.pem")
		keyPath := filepath.Join(tempDir, "key.pem")

		err = os.WriteFile(certPath, certPEM, 0600)
		require.NoError(t, err)
		err = os.WriteFile(keyPath, keyPEM, 0600)
		require.NoError(t, err)

		loadedCfg, err := services.LoadTLSConfig(certPath, keyPath)
		require.NoError(t, err)
		require.NotNil(t, loadedCfg)
		assert.Len(t, loadedCfg.Certificates, 1)
	})

	t.Run("Erro ao exportar certificado vazio", func(t *testing.T) {
		_, _, err := services.CertificateToPEM(tls.Certificate{})
		assert.Error(t, err)
	})

	t.Run("Erro quando arquivo de certificado não existe", func(t *testing.T) {
		tempDir := t.TempDir()
		keyPath := filepath.Join(tempDir, "key.pem")
		err := os.WriteFile(keyPath, []byte("dummy"), 0600)
		require.NoError(t, err)

		_, err = services.LoadTLSConfig(filepath.Join(tempDir, "inexistente.pem"), keyPath)
		assert.ErrorIs(t, err, services.ErrCertificateFileNotFound)
	})

	t.Run("Erro quando arquivo de chave não existe", func(t *testing.T) {
		tempDir := t.TempDir()
		certPath := filepath.Join(tempDir, "cert.pem")
		err := os.WriteFile(certPath, []byte("dummy"), 0600)
		require.NoError(t, err)

		_, err = services.LoadTLSConfig(certPath, filepath.Join(tempDir, "inexistente.key"))
		assert.ErrorIs(t, err, services.ErrCertificateFileNotFound)
	})

	t.Run("Erro quando arquivos contêm conteúdo inválido", func(t *testing.T) {
		tempDir := t.TempDir()
		certPath := filepath.Join(tempDir, "cert.pem")
		keyPath := filepath.Join(tempDir, "key.pem")

		err := os.WriteFile(certPath, []byte("INVALID_CERT_CONTENT"), 0600)
		require.NoError(t, err)
		err = os.WriteFile(keyPath, []byte("INVALID_KEY_CONTENT"), 0600)
		require.NoError(t, err)

		_, err = services.LoadTLSConfig(certPath, keyPath)
		assert.ErrorIs(t, err, services.ErrInvalidCertificate)
	})
}

func TestTLSHandshakeAndHTTPSRequest(t *testing.T) {
	cert, err := services.GenerateSelfSignedCertificate("127.0.0.1", "localhost")
	require.NoError(t, err)

	tlsCfg := services.BuildTLSConfig(cert)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.True(t, r.TLS != nil, "requisição deve ser criptografada com TLS")
		assert.GreaterOrEqual(t, r.TLS.Version, uint16(tls.VersionTLS12), "TLS version >= TLS 1.2")
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("HTTPS Seguro OK"))
	})

	server := httptest.NewUnstartedServer(handler)
	server.TLS = tlsCfg
	server.StartTLS()
	defer server.Close()

	// Cliente com certificado confiado
	certPool := x509.NewCertPool()
	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	require.NoError(t, err)
	certPool.AddCert(x509Cert)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool,
			},
		},
	}

	resp, err := client.Get(server.URL)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "HTTPS Seguro OK", string(body))
}
