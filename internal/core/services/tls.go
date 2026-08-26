package services

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
	"strings"
	"time"
)

var (
	// ErrInvalidCertificate indica que o certificado fornecido é inválido.
	ErrInvalidCertificate = errors.New("certificado TLS inválido")

	// ErrCertificateFileNotFound indica que os arquivos de certificado ou chave não foram encontrados.
	ErrCertificateFileNotFound = errors.New("arquivo de certificado ou chave privada TLS não encontrado")
)

// GetLANIPAddresses descobre os endereços IP (v4 e v6) locais das interfaces de rede ativas.
func GetLANIPAddresses() []net.IP {
	var ips []net.IP

	// Adiciona sempre loopbacks padrão
	ips = append(ips, net.ParseIP("127.0.0.1"), net.ParseIP("::1"))

	ifaces, err := net.Interfaces()
	if err != nil {
		return ips
	}

	for _, iface := range ifaces {
		// Ignora interfaces desativadas ou de loopback (já adicionadas manualmente)
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip != nil && !ip.IsLoopback() {
				ips = append(ips, ip)
			}
		}
	}

	return ips
}

// GenerateSelfSignedCertificate gera em memória um certificado X.509 autoassinado e sua chave privada ECDSA P-256.
func GenerateSelfSignedCertificate(hosts ...string) (tls.Certificate, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("falha ao gerar chave privada ECDSA: %w", err)
	}

	// Gera número serial aleatório de 128 bits
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("falha ao gerar serial do certificado: %w", err)
	}

	notBefore := time.Now().Add(-1 * time.Hour)
	notAfter := notBefore.Add(365 * 24 * time.Hour) // Válido por 1 ano

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization:  []string{"File Server Local"},
			CommonName:    "File Server Local",
			Country:       []string{"BR"},
			Province:      []string{"Local"},
			Locality:      []string{"Local"},
			StreetAddress: []string{"LAN"},
			PostalCode:    []string{"00000"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// DNS Names e IPs
	dnsMap := make(map[string]bool)
	ipMap := make(map[string]bool)

	// Adiciona padrões locais
	dnsMap["localhost"] = true
	dnsMap["*.localhost"] = true

	// Adiciona IPs locais detectados
	for _, ip := range GetLANIPAddresses() {
		if ip != nil {
			ipStr := ip.String()
			if !ipMap[ipStr] {
				ipMap[ipStr] = true
				template.IPAddresses = append(template.IPAddresses, ip)
			}
		}
	}

	// Adiciona hosts e IPs informados pelo chamador
	for _, h := range hosts {
		h = strings.TrimSpace(h)
		if h == "" {
			continue
		}
		// Trata casos de host:porta
		if hostOnly, _, err := net.SplitHostPort(h); err == nil {
			h = hostOnly
		}

		if ip := net.ParseIP(h); ip != nil {
			ipStr := ip.String()
			if !ipMap[ipStr] {
				ipMap[ipStr] = true
				template.IPAddresses = append(template.IPAddresses, ip)
			}
		} else {
			if !dnsMap[h] {
				dnsMap[h] = true
			}
		}
	}

	for dns := range dnsMap {
		template.DNSNames = append(template.DNSNames, dns)
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("falha ao criar certificado X.509: %w", err)
	}

	cert := tls.Certificate{
		Certificate: [][]byte{derBytes},
		PrivateKey:  priv,
	}

	return cert, nil
}

// BuildTLSConfig monta uma configuração *tls.Config de alta performance e segurança com o certificado fornecido.
func BuildTLSConfig(cert tls.Certificate) *tls.Config {
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		NextProtos:   []string{"h2", "http/1.1"},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
		},
	}
}

// CreateSelfSignedTLSConfig gera um certificado autoassinado em memória e retorna a configuração TLS pronta para uso.
func CreateSelfSignedTLSConfig(hosts ...string) (*tls.Config, error) {
	cert, err := GenerateSelfSignedCertificate(hosts...)
	if err != nil {
		return nil, err
	}
	return BuildTLSConfig(cert), nil
}

// LoadTLSConfig carrega um certificado e chave privada a partir de arquivos PEM em disco e monta o *tls.Config.
func LoadTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	if _, err := os.Stat(certFile); err != nil {
		return nil, fmt.Errorf("%w: arquivo de certificado '%s': %v", ErrCertificateFileNotFound, certFile, err)
	}
	if _, err := os.Stat(keyFile); err != nil {
		return nil, fmt.Errorf("%w: arquivo de chave privada '%s': %v", ErrCertificateFileNotFound, keyFile, err)
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("%w: falha ao carregar par de chaves TLS: %v", ErrInvalidCertificate, err)
	}

	return BuildTLSConfig(cert), nil
}

// CertificateToPEM converte um tls.Certificate para formato codificado PEM (bytes de certificado e chave).
func CertificateToPEM(cert tls.Certificate) (certPEM []byte, keyPEM []byte, err error) {
	if len(cert.Certificate) == 0 {
		return nil, nil, errors.New("certificado vazio")
	}

	var certBuf []byte
	for _, b := range cert.Certificate {
		block := &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: b,
		}
		certBuf = append(certBuf, pem.EncodeToMemory(block)...)
	}

	privKeyBytes, err := x509.MarshalPKCS8PrivateKey(cert.PrivateKey)
	if err != nil {
		// Tenta MarshalECPrivateKey caso MarshalPKCS8PrivateKey falhe
		if ecKey, ok := cert.PrivateKey.(*ecdsa.PrivateKey); ok {
			var ecErr error
			privKeyBytes, ecErr = x509.MarshalECPrivateKey(ecKey)
			if ecErr != nil {
				return nil, nil, fmt.Errorf("falha ao codificar chave privada ECDSA: %w", ecErr)
			}
		} else {
			return nil, nil, fmt.Errorf("falha ao serializar chave privada: %w", err)
		}
	}

	keyBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privKeyBytes,
	}
	keyBuf := pem.EncodeToMemory(keyBlock)

	return certBuf, keyBuf, nil
}
