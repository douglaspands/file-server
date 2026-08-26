package services_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/douglas/file-server/internal/core/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

func TestGenerateRandomPassword(t *testing.T) {
	t.Run("gera senha com tamanho customizado", func(t *testing.T) {
		pass, err := services.GenerateRandomPassword(16)
		require.NoError(t, err)
		assert.Len(t, pass, 16)
	})

	t.Run("gera senha com tamanho padrão quando comprimento <= 0", func(t *testing.T) {
		pass, err := services.GenerateRandomPassword(0)
		require.NoError(t, err)
		assert.Len(t, pass, 12)

		passNeg, err := services.GenerateRandomPassword(-5)
		require.NoError(t, err)
		assert.Len(t, passNeg, 12)
	})

	t.Run("duas chamadas consecutivas geram senhas distintas", func(t *testing.T) {
		p1, err := services.GenerateRandomPassword(12)
		require.NoError(t, err)
		p2, err := services.GenerateRandomPassword(12)
		require.NoError(t, err)
		assert.NotEqual(t, p1, p2)
	})
}

func TestGenerateSSHHostKey(t *testing.T) {
	signer, pemBytes, err := services.GenerateSSHHostKey()
	require.NoError(t, err)
	assert.NotNil(t, signer)
	assert.NotEmpty(t, pemBytes)
	assert.Contains(t, string(pemBytes), "BEGIN PRIVATE KEY")
	assert.Equal(t, ssh.KeyAlgoED25519, signer.PublicKey().Type())
}

func TestLoadSSHHostKey(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("carrega chave privada de host existente com sucesso", func(t *testing.T) {
		_, pemBytes, err := services.GenerateSSHHostKey()
		require.NoError(t, err)

		keyPath := filepath.Join(tempDir, "host_key.pem")
		err = os.WriteFile(keyPath, pemBytes, 0600)
		require.NoError(t, err)

		signer, err := services.LoadSSHHostKey(keyPath)
		require.NoError(t, err)
		assert.NotNil(t, signer)
		assert.Equal(t, ssh.KeyAlgoED25519, signer.PublicKey().Type())
	})

	t.Run("retorna erro quando o arquivo não existe", func(t *testing.T) {
		signer, err := services.LoadSSHHostKey(filepath.Join(tempDir, "nao_existe.pem"))
		assert.Error(t, err)
		assert.Nil(t, signer)
		assert.ErrorIs(t, err, services.ErrSSHKeyFileNotFound)
	})

	t.Run("retorna erro quando o arquivo contém conteúdo inválido", func(t *testing.T) {
		keyPath := filepath.Join(tempDir, "invalido.pem")
		err := os.WriteFile(keyPath, []byte("isto nao e uma chave ssh"), 0600)
		require.NoError(t, err)

		signer, err := services.LoadSSHHostKey(keyPath)
		assert.Error(t, err)
		assert.Nil(t, signer)
		assert.ErrorIs(t, err, services.ErrInvalidSSHKey)
	})
}

func TestLoadSSHPublicKey(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("carrega chave pública SSH autorizada com sucesso", func(t *testing.T) {
		pub, _, err := ed25519.GenerateKey(rand.Reader)
		require.NoError(t, err)

		sshPub, err := ssh.NewPublicKey(pub)
		require.NoError(t, err)

		authorizedKeyBytes := ssh.MarshalAuthorizedKey(sshPub)
		pubKeyPath := filepath.Join(tempDir, "id_ed25519.pub")
		err = os.WriteFile(pubKeyPath, authorizedKeyBytes, 0644)
		require.NoError(t, err)

		loadedPub, err := services.LoadSSHPublicKey(pubKeyPath)
		require.NoError(t, err)
		assert.NotNil(t, loadedPub)
		assert.Equal(t, sshPub.Marshal(), loadedPub.Marshal())
	})

	t.Run("retorna erro quando o arquivo não existe", func(t *testing.T) {
		loadedPub, err := services.LoadSSHPublicKey(filepath.Join(tempDir, "inexistente.pub"))
		assert.Error(t, err)
		assert.Nil(t, loadedPub)
		assert.ErrorIs(t, err, services.ErrSSHKeyFileNotFound)
	})

	t.Run("retorna erro quando o arquivo contém formato inválido", func(t *testing.T) {
		pubKeyPath := filepath.Join(tempDir, "invalido.pub")
		err := os.WriteFile(pubKeyPath, []byte("conteudo-invalido-de-chave"), 0644)
		require.NoError(t, err)

		loadedPub, err := services.LoadSSHPublicKey(pubKeyPath)
		assert.Error(t, err)
		assert.Nil(t, loadedPub)
		assert.ErrorIs(t, err, services.ErrInvalidSSHKey)
	})
}

func TestValidateCredentials(t *testing.T) {
	t.Run("credenciais corretas retornam true", func(t *testing.T) {
		assert.True(t, services.ValidateCredentials("admin", "senha123", "admin", "senha123"))
		assert.True(t, services.ValidateCredentials("  admin  ", "senha123", "admin", "senha123"))
	})

	t.Run("usuário ou senha incorretos retornam false", func(t *testing.T) {
		assert.False(t, services.ValidateCredentials("admin", "senhaErrada", "admin", "senha123"))
		assert.False(t, services.ValidateCredentials("userErrado", "senha123", "admin", "senha123"))
	})

	t.Run("credenciais esperadas vazias retornam false", func(t *testing.T) {
		assert.False(t, services.ValidateCredentials("", "", "", ""))
	})
}

func TestValidateSSHPublicKey(t *testing.T) {
	pub1, _, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	sshPub1, err := ssh.NewPublicKey(pub1)
	require.NoError(t, err)

	pub2, _, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	sshPub2, err := ssh.NewPublicKey(pub2)
	require.NoError(t, err)

	t.Run("chaves iguais retornam true", func(t *testing.T) {
		assert.True(t, services.ValidateSSHPublicKey(sshPub1, sshPub1))
	})

	t.Run("chaves distintas retornam false", func(t *testing.T) {
		assert.False(t, services.ValidateSSHPublicKey(sshPub1, sshPub2))
	})

	t.Run("chave nula retorna false", func(t *testing.T) {
		assert.False(t, services.ValidateSSHPublicKey(nil, sshPub1))
		assert.False(t, services.ValidateSSHPublicKey(sshPub1, nil))
		assert.False(t, services.ValidateSSHPublicKey(nil, nil))
	})
}
