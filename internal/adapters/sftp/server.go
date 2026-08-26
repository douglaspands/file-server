package sftp

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/douglas/file-server/internal/core/services"
	"github.com/douglas/file-server/internal/version"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// ServerOptions define as opções de configuração do servidor SFTP.
type ServerOptions struct {
	Host      string
	Port      int
	TargetDir string
	User      string
	Pass      string
	AuthKey   string
	HostKey   string
	ReadOnly  bool
}

// Server gerencia o listener SSH e o subsistema SFTP.
type Server struct {
	opts      ServerOptions
	sshConfig *ssh.ServerConfig
	listener  net.Listener
	fsHandler *FSHandler
	mu        sync.Mutex
	conns     map[net.Conn]struct{}
	closed    bool
}

// NewServer inicializa e configura o servidor SFTP validando chaves e credenciais.
func NewServer(opts ServerOptions) (*Server, error) {
	fsHandler, err := NewFSHandler(opts.TargetDir, opts.ReadOnly)
	if err != nil {
		return nil, fmt.Errorf("erro ao configurar handler de arquivos SFTP: %w", err)
	}

	var authPubKey ssh.PublicKey
	if opts.AuthKey != "" {
		pk, err := services.LoadSSHPublicKey(opts.AuthKey)
		if err != nil {
			return nil, fmt.Errorf("erro ao carregar chave pública autorizada: %w", err)
		}
		authPubKey = pk
	}

	var hostSigner ssh.Signer
	if opts.HostKey != "" {
		signer, err := services.LoadSSHHostKey(opts.HostKey)
		if err != nil {
			return nil, fmt.Errorf("erro ao carregar chave de host SSH: %w", err)
		}
		hostSigner = signer
	} else {
		signer, _, err := services.GenerateSSHHostKey()
		if err != nil {
			return nil, fmt.Errorf("erro ao gerar chave de host SSH em memória: %w", err)
		}
		hostSigner = signer
	}

	sshConfig := &ssh.ServerConfig{
		NoClientAuth: false,
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			if opts.User != "" && opts.Pass != "" {
				if services.ValidateCredentials(c.User(), string(pass), opts.User, opts.Pass) {
					return nil, nil
				}
			}
			return nil, fmt.Errorf("autenticação por senha rejeitada para usuário '%s'", c.User())
		},
		PublicKeyCallback: func(c ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			if authPubKey != nil {
				if services.ValidateSSHPublicKey(key, authPubKey) {
					return nil, nil
				}
			}
			return nil, fmt.Errorf("autenticação por chave pública rejeitada para usuário '%s'", c.User())
		},
	}

	sshConfig.AddHostKey(hostSigner)

	return &Server{
		opts:      opts,
		sshConfig: sshConfig,
		fsHandler: fsHandler,
		conns:     make(map[net.Conn]struct{}),
	}, nil
}

// Addr retorna o endereço de escuta do servidor (útil quando escutando em porta 0 em testes).
func (s *Server) Addr() net.Addr {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.listener != nil {
		return s.listener.Addr()
	}
	return nil
}

// LogStartupBanner imprime as informações de inicialização do servidor SFTP no terminal.
func LogStartupBanner(opts ServerOptions) {
	scheme := "sftp"
	var localAccess string
	var lanAccess []string

	if opts.Host == "0.0.0.0" || opts.Host == "" || opts.Host == "::" {
		localAccess = fmt.Sprintf("%s://%s@127.0.0.1:%d", scheme, opts.User, opts.Port)
		seen := make(map[string]bool)
		for _, ip := range services.GetLANIPAddresses() {
			if ip.IsLoopback() {
				continue
			}
			var url string
			if ip.To4() != nil {
				url = fmt.Sprintf("%s://%s@%s:%d", scheme, opts.User, ip.String(), opts.Port)
			} else {
				url = fmt.Sprintf("%s://%s@[%s]:%d", scheme, opts.User, ip.String(), opts.Port)
			}
			if !seen[url] {
				seen[url] = true
				lanAccess = append(lanAccess, url)
			}
		}
	} else {
		localAccess = fmt.Sprintf("%s://%s@%s:%d", scheme, opts.User, opts.Host, opts.Port)
	}

	log.Printf("⚡ SFTP Server v%s inicializado com sucesso!", version.Get().Version)
	log.Printf("📁 Diretório compartilhado: %s", opts.TargetDir)
	log.Printf("🔒 Protocolo: SFTP (SSHv2 / Criptografia Forte)")
	log.Printf("👤 Usuário: %s", opts.User)
	if opts.Pass != "" {
		log.Printf("🔑 Senha: %s", opts.Pass)
	}
	if opts.AuthKey != "" {
		log.Printf("🔑 Chave Pública Autorizada: %s", opts.AuthKey)
	}
	if opts.ReadOnly {
		log.Printf("🛡️ Modo: Somente Leitura (Read-Only)")
	}
	log.Printf("👉 Acesso Local:     %s", localAccess)
	if len(lanAccess) > 0 {
		for _, url := range lanAccess {
			log.Printf("🌐 Acesso Rede (LAN): %s", url)
		}
	}
}

// Serve escuta conexões no listener informado até ser cancelado ou fechado.
func (s *Server) Serve(listener net.Listener) error {
	s.mu.Lock()
	s.listener = listener
	s.mu.Unlock()

	for {
		nConn, err := listener.Accept()
		if err != nil {
			s.mu.Lock()
			closed := s.closed
			s.mu.Unlock()
			if closed {
				return nil
			}
			return err
		}

		s.trackConn(nConn, true)
		go func(c net.Conn) {
			defer s.trackConn(c, false)
			s.handleConn(c)
		}(nConn)
	}
}

func (s *Server) trackConn(c net.Conn, add bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if add {
		s.conns[c] = struct{}{}
	} else {
		delete(s.conns, c)
		_ = c.Close()
	}
}

func (s *Server) handleConn(nConn net.Conn) {
	sshConn, chans, reqs, err := ssh.NewServerConn(nConn, s.sshConfig)
	if err != nil {
		return
	}
	defer sshConn.Close()

	go ssh.DiscardRequests(reqs)

	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			_ = newChannel.Reject(ssh.UnknownChannelType, "tipo de canal desconhecido")
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			continue
		}

		go func(in <-chan *ssh.Request) {
			for req := range in {
				ok := false
				if req.Type == "subsystem" && len(req.Payload) >= 4 {
					subsystemName := string(req.Payload[4:])
					if subsystemName == "sftp" {
						ok = true
					}
				}
				_ = req.Reply(ok, nil)

				if ok {
					handlers := s.fsHandler.ToHandlers()
					server := sftp.NewRequestServer(channel, handlers)
					_ = server.Serve()
					_ = channel.Close()
					return
				}
			}
		}(requests)
	}
}

// Shutdown encerra o listener e todas as conexões ativas do servidor SFTP.
func (s *Server) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	s.closed = true
	var err error
	if s.listener != nil {
		err = s.listener.Close()
	}
	for conn := range s.conns {
		_ = conn.Close()
	}
	s.mu.Unlock()

	return err
}

// Run inicializa o listener TCP e executa o servidor SFTP aguardando cancelamento pelo contexto.
func (s *Server) Run(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", s.opts.Host, s.opts.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("falha ao iniciar listener SFTP em %s: %w", addr, err)
	}

	// Atualiza porta efetiva em opts caso porta 0 tenha sido usada
	if tcpAddr, ok := listener.Addr().(*net.TCPAddr); ok {
		s.opts.Port = tcpAddr.Port
	}

	LogStartupBanner(s.opts)

	errChan := make(chan error, 1)
	go func() {
		if err := s.Serve(listener); err != nil && !s.closed {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return fmt.Errorf("erro fatal no servidor SFTP: %w", err)
	case <-ctx.Done():
		log.Println("🛑 Encerrando servidor SFTP graciosamente...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.Shutdown(shutdownCtx)
}
