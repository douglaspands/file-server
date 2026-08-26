package gui

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/douglas/file-server/internal/version"
	"github.com/douglas/file-server/web"
)

// Server representa o servidor web dedicado para renderizar e controlar a GUI do desktop.
type Server struct {
	controller   *Controller
	folderPicker FolderPickerInterface
	httpServer   *http.Server
	listener     net.Listener
	tmpl         *template.Template
}

// NewServer cria uma nova instância do servidor HTTP da GUI.
func NewServer(controller *Controller, folderPicker FolderPickerInterface) (*Server, error) {
	if folderPicker == nil {
		folderPicker = NewNativeFolderPicker()
	}

	tmpl, err := template.ParseFS(web.GetTemplatesFS(), "templates/pages/gui_launcher.html")
	if err != nil {
		// Fallback se template ainda não estiver carregado
		tmpl = template.New("gui_launcher")
	}

	s := &Server{
		controller:   controller,
		folderPicker: folderPicker,
		tmpl:         tmpl,
	}

	return s, nil
}

// RegisterRoutes conecta todas as rotas e endpoints REST/SSE da GUI ao roteador.
func (s *Server) RegisterRoutes(mux *http.ServeMux) error {
	// Arquivos estáticos
	staticFS, err := web.GetStaticFileSystem()
	if err != nil {
		return fmt.Errorf("falha ao carregar sistema de arquivos estático: %w", err)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticFS)))

	// Rotas da Interface
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/api/status", s.handleStatus)
	mux.HandleFunc("/api/interfaces", s.handleInterfaces)
	mux.HandleFunc("/api/share-message", s.handleShareMessage)
	mux.HandleFunc("/api/server/start", s.handleStartServer)
	mux.HandleFunc("/api/server/stop", s.handleStopServer)
	mux.HandleFunc("/api/picker/folder", s.handlePickFolder)
	mux.HandleFunc("/api/logs/stream", s.handleLogStream)
	mux.HandleFunc("/api/app/open-browser", s.handleOpenBrowser)
	mux.HandleFunc("/api/app/close", s.handleCloseApp)

	return nil
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Recarrega template caso necessário
	tmpl, err := template.ParseFS(web.GetTemplatesFS(), "templates/pages/gui_launcher.html")
	if err == nil {
		s.tmpl = tmpl
	}

	data := map[string]interface{}{
		"Version":    version.Get().Version,
		"Status":     s.controller.GetStatus(),
		"Config":     s.controller.GetConfig(),
		"InitialDir": s.controller.GetConfig().TargetDir,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.tmpl.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao renderizar interface gráfica: %v", err), http.StatusInternalServerError)
	}
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(s.controller.GetStatus())
}

func (s *Server) handleInterfaces(w http.ResponseWriter, r *http.Request) {
	portStr := r.URL.Query().Get("port")
	port, _ := strconv.Atoi(portStr)
	if port <= 0 {
		port = 8080
	}
	protocol := r.URL.Query().Get("protocol")
	if protocol == "" {
		protocol = "web"
	}
	isTLS := r.URL.Query().Get("tls") == "true"

	ifaces := DetectNetworkInterfaces(port, protocol, isTLS)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ifaces)
}

func (s *Server) handleShareMessage(w http.ResponseWriter, r *http.Request) {
	status := s.controller.GetStatus()
	cfg := s.controller.GetConfig()
	ifaces := DetectNetworkInterfaces(cfg.Port, string(cfg.Protocol), cfg.UseTLS)

	msg := FormatShareMessage(cfg, ifaces)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": msg,
		"status":  status,
	})
}

func (s *Server) handleStartServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var cfg ServerConfig
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		http.Error(w, fmt.Sprintf("Corpo de requisição inválido: %v", err), http.StatusBadRequest)
		return
	}

	if err := s.controller.StartServer(cfg); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(s.controller.GetStatus())
}

func (s *Server) handleStopServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := s.controller.StopServer(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(s.controller.GetStatus())
}

func (s *Server) handlePickFolder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		CurrentDir string `json:"currentDir"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	selected, err := s.folderPicker.PickFolder(r.Context(), req.CurrentDir)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
			"path":    req.CurrentDir,
		})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"path":    selected,
	})
}

func (s *Server) handleLogStream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming não suportado pelo cliente", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	broadcaster := s.controller.GetBroadcaster()
	history, ch := broadcaster.Subscribe()
	defer broadcaster.Unsubscribe(ch)

	// Envia histórico inicial acumulado
	for _, line := range history {
		data, _ := json.Marshal(map[string]string{"message": line})
		fmt.Fprintf(w, "data: %s\n\n", data)
	}
	flusher.Flush()

	// Envia novas entradas em tempo real
	for {
		select {
		case <-r.Context().Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			data, _ := json.Marshal(map[string]string{"message": msg})
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}

func (s *Server) handleOpenBrowser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		status := s.controller.GetStatus()
		req.URL = status.LocalURL
	}

	go func() {
		_ = OpenURLInBrowser(req.URL)
	}()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (s *Server) handleCloseApp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]bool{"closing": true})

	go func() {
		time.Sleep(200 * time.Millisecond)
		_ = s.controller.StopServer()
		if s.httpServer != nil {
			_ = s.httpServer.Close()
		}
		os.Exit(0)
	}()
}

// Start inicializa o listener e serve as requisições HTTP da interface gráfica.
func (s *Server) Start(host string, port int) (string, error) {
	mux := http.NewServeMux()
	if err := s.RegisterRoutes(mux); err != nil {
		return "", err
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		// Fallback para porta aleatória disponível caso a porta solicitada esteja ocupada
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:0", host))
		if err != nil {
			return "", fmt.Errorf("falha ao abrir porta para a GUI: %w", err)
		}
	}

	s.listener = listener
	actualPort := listener.Addr().(*net.TCPAddr).Port
	s.controller.SetGUIAddr(host, actualPort)

	s.httpServer = &http.Server{
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	guiURL := fmt.Sprintf("http://127.0.0.1:%d", actualPort)

	go func() {
		if err := s.httpServer.Serve(listener); err != nil && err != http.ErrServerClosed {
			fmt.Printf("⚠️ Erro no servidor da GUI: %v\n", err)
		}
	}()

	return guiURL, nil
}

// Stop finaliza o servidor HTTP da GUI.
func (s *Server) Stop() error {
	if s.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}
