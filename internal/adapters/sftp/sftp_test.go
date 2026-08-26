package sftp_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"io"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	adapterSftp "github.com/douglas/file-server/internal/adapters/sftp"
	"github.com/douglas/file-server/internal/core/services"
	"github.com/pkg/sftp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

func TestFSHandler_Unit(t *testing.T) {
	tempDir := t.TempDir()

	// Cria estrutura de testes
	subDir := filepath.Join(tempDir, "pasta1")
	err := os.Mkdir(subDir, 0755)
	require.NoError(t, err)

	testFile := filepath.Join(tempDir, "arquivo.txt")
	err = os.WriteFile(testFile, []byte("conteudo de teste"), 0644)
	require.NoError(t, err)

	t.Run("criação de handler com diretório padrão ou inválido", func(t *testing.T) {
		h, err := adapterSftp.NewFSHandler("", false)
		require.NoError(t, err)
		assert.NotNil(t, h)

		h, err = adapterSftp.NewFSHandler(filepath.Join(tempDir, "inexistente"), false)
		assert.Error(t, err)
		assert.Nil(t, h)

		h, err = adapterSftp.NewFSHandler(testFile, false)
		assert.Error(t, err)
		assert.Nil(t, h)
	})

	t.Run("resolução de caminhos e proteção de sandbox", func(t *testing.T) {
		h, err := adapterSftp.NewFSHandler(tempDir, false)
		require.NoError(t, err)

		// Caminho válido
		resolved, err := h.ResolvePath("/arquivo.txt")
		require.NoError(t, err)
		assert.Equal(t, testFile, resolved)

		resolvedSub, err := h.ResolvePath("pasta1")
		require.NoError(t, err)
		assert.Equal(t, subDir, resolvedSub)

		// Path Traversal
		_, err = h.ResolvePath("../../etc/passwd")
		assert.Error(t, err)

		_, err = h.ResolvePath("pasta1/../../..")
		assert.Error(t, err)

		_, err = h.ResolvePath("arquivo.txt\x00malicioso")
		assert.Error(t, err)

		// Symlink externo
		outsideDir := t.TempDir()
		outsideFile := filepath.Join(outsideDir, "secret.txt")
		err = os.WriteFile(outsideFile, []byte("secret"), 0644)
		require.NoError(t, err)

		symlinkPath := filepath.Join(tempDir, "symlink_ext")
		err = os.Symlink(outsideFile, symlinkPath)
		require.NoError(t, err)

		_, err = h.ResolvePath("symlink_ext")
		assert.Error(t, err)
	})

	t.Run("operações em modo somente leitura (read-only)", func(t *testing.T) {
		h, err := adapterSftp.NewFSHandler(tempDir, true)
		require.NoError(t, err)

		// Filewrite deve falhar
		_, err = h.Filewrite(&sftp.Request{Filepath: "novo.txt", Method: "Put"})
		assert.Error(t, err)

		// Filecmd deve falhar para todas operações de modificação
		assert.Error(t, h.Filecmd(&sftp.Request{Filepath: "pasta_nova", Method: "Mkdir"}))
		assert.Error(t, h.Filecmd(&sftp.Request{Filepath: "arquivo.txt", Target: "arquivo_renomeado.txt", Method: "Rename"}))
		assert.Error(t, h.Filecmd(&sftp.Request{Filepath: "arquivo.txt", Method: "Remove"}))
		assert.Error(t, h.Filecmd(&sftp.Request{Filepath: "pasta1", Method: "Rmdir"}))
		assert.Error(t, h.Filecmd(&sftp.Request{Filepath: "link.txt", Target: "arquivo.txt", Method: "Symlink"}))
		assert.Error(t, h.Filecmd(&sftp.Request{Filepath: "arquivo.txt", Method: "Setstat"}))
	})

	t.Run("tratamento de erros em Fileread e Filewrite", func(t *testing.T) {
		h, err := adapterSftp.NewFSHandler(tempDir, false)
		require.NoError(t, err)

		// Fileread em arquivo inexistente
		_, err = h.Fileread(&sftp.Request{Filepath: "nao_existe.txt", Method: "Get"})
		assert.Error(t, err)

		// Fileread em diretório
		_, err = h.Fileread(&sftp.Request{Filepath: "pasta1", Method: "Get"})
		assert.Error(t, err)

		// Filewrite em pasta inexistente no pai
		_, err = h.Filewrite(&sftp.Request{Filepath: "pasta_inexistente/sub/arquivo.txt", Method: "Put"})
		assert.Error(t, err)

		// Filewrite com sucesso
		w, err := h.Filewrite(&sftp.Request{Filepath: "gravado_unit.txt", Method: "Put"})
		require.NoError(t, err)
		assert.NotNil(t, w)
		if closer, ok := w.(io.Closer); ok {
			_ = closer.Close()
		}
	})

	t.Run("tratamento de Filecmd com métodos variados", func(t *testing.T) {
		h, err := adapterSftp.NewFSHandler(tempDir, false)
		require.NoError(t, err)

		// Setstat com sucesso
		err = h.Filecmd(&sftp.Request{Filepath: "arquivo.txt", Method: "Setstat"})
		assert.NoError(t, err)

		// Mkdir
		err = h.Filecmd(&sftp.Request{Filepath: "mkdir_unit", Method: "Mkdir"})
		assert.NoError(t, err)

		// Symlink interno
		err = h.Filecmd(&sftp.Request{Filepath: "symlink_unit", Target: "arquivo.txt", Method: "Symlink"})
		assert.NoError(t, err)

		// Rename
		err = h.Filecmd(&sftp.Request{Filepath: "mkdir_unit", Target: "mkdir_renomeado", Method: "Rename"})
		assert.NoError(t, err)

		// Rmdir
		err = h.Filecmd(&sftp.Request{Filepath: "mkdir_renomeado", Method: "Rmdir"})
		assert.NoError(t, err)

		// Remove arquivo
		err = h.Filecmd(&sftp.Request{Filepath: "symlink_unit", Method: "Remove"})
		assert.NoError(t, err)

		// Método não suportado
		err = h.Filecmd(&sftp.Request{Filepath: "arquivo.txt", Method: "MetodoInvalido"})
		assert.Error(t, err)
	})

	t.Run("tratamento de Filelist com métodos variados", func(t *testing.T) {
		h, err := adapterSftp.NewFSHandler(tempDir, false)
		require.NoError(t, err)

		// List
		lister, err := h.Filelist(&sftp.Request{Filepath: ".", Method: "List"})
		require.NoError(t, err)
		assert.NotNil(t, lister)

		// Stat
		statLister, err := h.Filelist(&sftp.Request{Filepath: "arquivo.txt", Method: "Stat"})
		require.NoError(t, err)
		assert.NotNil(t, statLister)

		// Readlink em link existente
		err = os.Symlink(testFile, filepath.Join(tempDir, "link_para_stat"))
		require.NoError(t, err)
		rlLister, err := h.Filelist(&sftp.Request{Filepath: "link_para_stat", Method: "Readlink"})
		require.NoError(t, err)
		assert.NotNil(t, rlLister)

		// Método não suportado
		_, err = h.Filelist(&sftp.Request{Filepath: "arquivo.txt", Method: "MetodoInvalido"})
		assert.Error(t, err)

		// Teste de paginação do ListerAt
		buf := make([]os.FileInfo, 1)
		n, err := lister.ListAt(buf, 0)
		assert.True(t, n > 0 || err == io.EOF)
		n, err = lister.ListAt(buf, 9999)
		assert.Equal(t, 0, n)
		assert.Equal(t, io.EOF, err)
	})
}

func TestSFTPServer_EndToEnd(t *testing.T) {
	tempDir := t.TempDir()

	// Arquivo de teste inicial
	testFile := filepath.Join(tempDir, "hello.txt")
	err := os.WriteFile(testFile, []byte("ola mundo sftp"), 0644)
	require.NoError(t, err)

	// Gera par de chaves para autenticação por chave pública
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	sshPub, err := ssh.NewPublicKey(pub)
	require.NoError(t, err)
	signer, err := ssh.NewSignerFromKey(priv)
	require.NoError(t, err)

	pubKeyFile := filepath.Join(tempDir, "authorized_key.pub")
	err = os.WriteFile(pubKeyFile, ssh.MarshalAuthorizedKey(sshPub), 0644)
	require.NoError(t, err)

	// Cria servidor SFTP
	opts := adapterSftp.ServerOptions{
		Host:      "127.0.0.1",
		Port:      0, // Porta randômica
		TargetDir: tempDir,
		User:      "testuser",
		Pass:      "secretpass",
		AuthKey:   pubKeyFile,
		ReadOnly:  false,
	}

	server, err := adapterSftp.NewServer(opts)
	require.NoError(t, err)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	errChan := make(chan error, 1)
	go func() {
		errChan <- server.Serve(listener)
	}()

	time.Sleep(20 * time.Millisecond)
	serverAddr := listener.Addr().String()
	assert.NotNil(t, server.Addr())

	t.Run("autenticação por senha incorreta falha", func(t *testing.T) {
		clientConfig := &ssh.ClientConfig{
			User: "testuser",
			Auth: []ssh.AuthMethod{
				ssh.Password("senha_errada"),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         2 * time.Second,
		}
		_, err := ssh.Dial("tcp", serverAddr, clientConfig)
		assert.Error(t, err)
	})

	t.Run("autenticação por senha correta e operações completas no cliente SFTP", func(t *testing.T) {
		clientConfig := &ssh.ClientConfig{
			User: "testuser",
			Auth: []ssh.AuthMethod{
				ssh.Password("secretpass"),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         2 * time.Second,
		}

		sshConn, err := ssh.Dial("tcp", serverAddr, clientConfig)
		require.NoError(t, err)
		defer sshConn.Close()

		sftpClient, err := sftp.NewClient(sshConn)
		require.NoError(t, err)
		defer sftpClient.Close()

		// 1. Leitura de arquivo existente
		f, err := sftpClient.Open("hello.txt")
		require.NoError(t, err)
		content, err := io.ReadAll(f)
		require.NoError(t, err)
		assert.Equal(t, "ola mundo sftp", string(content))
		_ = f.Close()

		// 2. Listagem de diretório
		entries, err := sftpClient.ReadDir(".")
		require.NoError(t, err)
		assert.NotEmpty(t, entries)

		// 3. Criação e escrita de novo arquivo
		newF, err := sftpClient.Create("upload.txt")
		require.NoError(t, err)
		_, err = newF.Write([]byte("upload feito via sftp"))
		require.NoError(t, err)
		_ = newF.Close()

		// Verifica se foi salvo no disco
		diskContent, err := os.ReadFile(filepath.Join(tempDir, "upload.txt"))
		require.NoError(t, err)
		assert.Equal(t, "upload feito via sftp", string(diskContent))

		// 4. Criação de subdiretório
		err = sftpClient.Mkdir("nova_pasta")
		require.NoError(t, err)

		// 5. Renomeação
		err = sftpClient.Rename("upload.txt", "upload_renomeado.txt")
		require.NoError(t, err)

		// 6. Remoção de arquivo e pasta
		err = sftpClient.Remove("upload_renomeado.txt")
		require.NoError(t, err)
		err = sftpClient.RemoveDirectory("nova_pasta")
		require.NoError(t, err)

		// 7. Stat e Readlink
		stat, err := sftpClient.Stat("hello.txt")
		require.NoError(t, err)
		assert.Equal(t, "hello.txt", stat.Name())

		// 8. Path Traversal bloqueado
		_, err = sftpClient.Open("../../../etc/passwd")
		assert.Error(t, err)
	})

	t.Run("autenticação por chave pública autorizada", func(t *testing.T) {
		clientConfig := &ssh.ClientConfig{
			User: "testuser",
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         2 * time.Second,
		}

		sshConn, err := ssh.Dial("tcp", serverAddr, clientConfig)
		require.NoError(t, err)
		defer sshConn.Close()

		sftpClient, err := sftp.NewClient(sshConn)
		require.NoError(t, err)
		defer sftpClient.Close()

		entries, err := sftpClient.ReadDir(".")
		require.NoError(t, err)
		assert.NotEmpty(t, entries)
	})

	// Encerra servidor
	err = server.Shutdown(context.Background())
	require.NoError(t, err)
}

func TestSFTPServer_LifecycleAndBanner(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("inicialização com host key customizada existente e banner", func(t *testing.T) {
		_, pemBytes, err := services.GenerateSSHHostKey()
		require.NoError(t, err)
		hostKeyPath := filepath.Join(tempDir, "my_host_key.pem")
		err = os.WriteFile(hostKeyPath, pemBytes, 0600)
		require.NoError(t, err)

		opts := adapterSftp.ServerOptions{
			Host:      "0.0.0.0",
			Port:      2222,
			TargetDir: tempDir,
			User:      "admin",
			Pass:      "secret",
			HostKey:   hostKeyPath,
			ReadOnly:  true,
		}

		server, err := adapterSftp.NewServer(opts)
		require.NoError(t, err)
		assert.NotNil(t, server)

		// Executa banner para cobertura
		adapterSftp.LogStartupBanner(opts)
	})

	t.Run("falha com chave pública autorizada inexistente", func(t *testing.T) {
		opts := adapterSftp.ServerOptions{
			Host:      "127.0.0.1",
			Port:      2222,
			TargetDir: tempDir,
			AuthKey:   filepath.Join(tempDir, "inexistente.pub"),
		}
		_, err := adapterSftp.NewServer(opts)
		assert.Error(t, err)
	})

	t.Run("falha com host key inexistente", func(t *testing.T) {
		opts := adapterSftp.ServerOptions{
			Host:      "127.0.0.1",
			Port:      2222,
			TargetDir: tempDir,
			HostKey:   filepath.Join(tempDir, "inexistente.pem"),
		}
		_, err := adapterSftp.NewServer(opts)
		assert.Error(t, err)
	})

	t.Run("execução graciosa via Run com cancelamento por contexto", func(t *testing.T) {
		opts := adapterSftp.ServerOptions{
			Host:      "127.0.0.1",
			Port:      0,
			TargetDir: tempDir,
			User:      "user",
			Pass:      "pass",
		}

		server, err := adapterSftp.NewServer(opts)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		doneChan := make(chan error, 1)

		go func() {
			doneChan <- server.Run(ctx)
		}()

		time.Sleep(100 * time.Millisecond)
		cancel()

		select {
		case err := <-doneChan:
			assert.NoError(t, err)
		case <-time.After(3 * time.Second):
			t.Fatal("tempo limite excedido ao encerrar servidor")
		}
	})
}
