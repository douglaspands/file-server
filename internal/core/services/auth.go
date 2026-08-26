package services

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/subtle"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

var (
	// ErrInvalidSSHKey indica que a chave SSH fornecida é inválida ou não pôde ser parseada.
	ErrInvalidSSHKey = errors.New("chave SSH inválida")

	// ErrSSHKeyFileNotFound indica que o arquivo de chave SSH não foi encontrado.
	ErrSSHKeyFileNotFound = errors.New("arquivo de chave SSH não encontrado")
)

const (
	// DefaultUsername define o nome de usuário padrão gerado automaticamente caso não seja informado.
	DefaultUsername = "fileserver"

	// charset utilizado para geração de senhas aleatórias seguras e legíveis.
	passwordCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%&*"
)

// GenerateRandomPassword gera uma senha aleatória segura com o comprimento especificado utilizando crypto/rand.
func GenerateRandomPassword(length int) (string, error) {
	if length <= 0 {
		length = 12
	}

	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(passwordCharset)))

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", fmt.Errorf("falha ao gerar caractere aleatório: %w", err)
		}
		result[i] = passwordCharset[num.Int64()]
	}

	return string(result), nil
}

// GenerateSSHHostKey gera uma chave privada Ed25519 em memória e retorna o ssh.Signer correspondente e os bytes PEM.
func GenerateSSHHostKey() (ssh.Signer, []byte, error) {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("falha ao gerar chave Ed25519: %w", err)
	}

	signer, err := ssh.NewSignerFromKey(priv)
	if err != nil {
		return nil, nil, fmt.Errorf("falha ao criar ssh.Signer a partir da chave: %w", err)
	}

	pkcs8Bytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, nil, fmt.Errorf("falha ao serializar chave privada Ed25519: %w", err)
	}

	pemBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pkcs8Bytes,
	}

	return signer, pem.EncodeToMemory(pemBlock), nil
}

// LoadSSHHostKey carrega e parseia uma chave privada de host SSH a partir de um arquivo PEM em disco.
func LoadSSHHostKey(keyPath string) (ssh.Signer, error) {
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: '%s'", ErrSSHKeyFileNotFound, keyPath)
		}
		return nil, fmt.Errorf("erro ao ler arquivo de chave SSH '%s': %w", keyPath, err)
	}

	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("%w: falha ao parsear chave privada '%s': %v", ErrInvalidSSHKey, keyPath, err)
	}

	return signer, nil
}

// LoadSSHPublicKey carrega e parseia uma chave pública autorizada (formato OpenSSH authorized_keys) a partir de um arquivo em disco.
func LoadSSHPublicKey(pubKeyPath string) (ssh.PublicKey, error) {
	keyBytes, err := os.ReadFile(pubKeyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: '%s'", ErrSSHKeyFileNotFound, pubKeyPath)
		}
		return nil, fmt.Errorf("erro ao ler arquivo de chave pública SSH '%s': %w", pubKeyPath, err)
	}

	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(bytes.TrimSpace(keyBytes))
	if err != nil {
		return nil, fmt.Errorf("%w: falha ao parsear chave pública '%s': %v", ErrInvalidSSHKey, pubKeyPath, err)
	}

	return pubKey, nil
}

// ValidateCredentials valida usuário e senha de forma segura em tempo constante para mitigar ataques de temporização.
func ValidateCredentials(inputUser, inputPass, expectedUser, expectedPass string) bool {
	if expectedUser == "" && expectedPass == "" {
		return false
	}

	userMatch := subtle.ConstantTimeCompare([]byte(strings.TrimSpace(inputUser)), []byte(strings.TrimSpace(expectedUser))) == 1
	passMatch := subtle.ConstantTimeCompare([]byte(inputPass), []byte(expectedPass)) == 1

	return userMatch && passMatch
}

// ValidateSSHPublicKey compara se uma chave pública apresentada pelo cliente corresponde à chave pública autorizada.
func ValidateSSHPublicKey(presentedKey ssh.PublicKey, authorizedKey ssh.PublicKey) bool {
	if presentedKey == nil || authorizedKey == nil {
		return false
	}
	return bytes.Equal(presentedKey.Marshal(), authorizedKey.Marshal())
}
