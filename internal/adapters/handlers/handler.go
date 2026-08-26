package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/douglas/file-server/internal/core/ports"
	"github.com/douglas/file-server/internal/version"
	"github.com/douglas/file-server/web"
)

// Handler gerencia o roteamento e a execução de requisições HTTP.
type Handler struct {
	fileService   ports.FileService
	healthService ports.HealthService
	templates     map[string]*template.Template
}

// NewHandler cria uma nova instância de Handler com templates parseados modularmente por página.
func NewHandler(fileService ports.FileService, healthService ports.HealthService) (*Handler, error) {
	tmplFS := web.GetTemplatesFS()
	pages := []string{"explorer.html", "index.html"}
	templates := make(map[string]*template.Template)

	for _, page := range pages {
		tmpl, err := template.ParseFS(tmplFS,
			"templates/layouts/*.html",
			"templates/partials/*.html",
			"templates/pages/"+page,
		)
		if err != nil {
			return nil, fmt.Errorf("falha ao parsear template para a página %s: %w", page, err)
		}
		templates[page] = tmpl
	}

	// Template para renderização isolada de partials
	partialsTmpl, err := template.ParseFS(tmplFS, "templates/partials/*.html")
	if err != nil {
		return nil, fmt.Errorf("falha ao parsear templates parciais: %w", err)
	}
	templates["partials"] = partialsTmpl

	return &Handler{
		fileService:   fileService,
		healthService: healthService,
		templates:     templates,
	}, nil
}

// RegisterRoutes registra todas as rotas da aplicação no mux HTTP fornecido.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) error {
	// Arquivos estáticos embutidos
	staticFS, err := web.GetStaticFileSystem()
	if err != nil {
		return fmt.Errorf("falha ao obter sistema de arquivos estáticos: %w", err)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticFS)))

	// Rotas do Explorador de Arquivos e Páginas Web
	mux.HandleFunc("/", h.FileBrowserHandler)
	mux.HandleFunc("/files/", h.FileBrowserHandler)
	mux.HandleFunc("/files", h.FileBrowserHandler)
	mux.HandleFunc("/status", h.Home)
	mux.HandleFunc("/partials/health", h.HealthPartial)
	mux.HandleFunc("/_live_reload", LiveReloadHandler)

	// Rotas de Transferência e Operações de Arquivos
	mux.HandleFunc("/download/", h.DownloadFileHandler)
	mux.HandleFunc("/download", h.DownloadFileHandler)
	mux.HandleFunc("/zip/", h.DownloadZipHandler)
	mux.HandleFunc("/zip", h.DownloadZipHandler)
	mux.HandleFunc("/upload/", h.UploadFilesHandler)
	mux.HandleFunc("/upload", h.UploadFilesHandler)

	// Rotas de API JSON
	mux.HandleFunc("/api/files/", h.FilesAPIHandler)
	mux.HandleFunc("/api/files", h.FilesAPIHandler)
	mux.HandleFunc("/api/health", h.HealthAPI)
	mux.HandleFunc("/api/version", h.VersionAPI)

	return nil
}

// Home renderiza a página de status do sistema.
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	health, err := h.healthService.Check(r.Context())
	if err != nil {
		http.Error(w, "Erro ao obter status do sistema", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Version": version.Get().Version,
		"Health":  health,
		"IsDev":   version.Get().IsDev,
		"RootDir": h.fileService.GetRootDir(),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, ok := h.templates["index.html"]
	if !ok {
		http.Error(w, "Template não encontrado", http.StatusInternalServerError)
		return
	}
	if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao renderizar template: %v", err), http.StatusInternalServerError)
	}
}

// HealthPartial renderiza apenas o fragmento HTML do card de saúde para requisições HTMX.
func (h *Handler) HealthPartial(w http.ResponseWriter, r *http.Request) {
	health, err := h.healthService.Check(r.Context())
	if err != nil {
		http.Error(w, "Erro ao obter status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, ok := h.templates["partials"]
	if !ok {
		http.Error(w, "Template não encontrado", http.StatusInternalServerError)
		return
	}
	if err := tmpl.ExecuteTemplate(w, "health_card.html", health); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao renderizar partial: %v", err), http.StatusInternalServerError)
	}
}

// HealthAPI retorna o status de saúde em JSON.
func (h *Handler) HealthAPI(w http.ResponseWriter, r *http.Request) {
	health, err := h.healthService.Check(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(health)
}

// VersionAPI retorna os detalhes de versão em JSON.
func (h *Handler) VersionAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(version.Get())
}
