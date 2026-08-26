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
	healthService ports.HealthService
	tmpl          *template.Template
}

// NewHandler cria uma nova instância de Handler com templates parseados.
func NewHandler(healthService ports.HealthService) (*Handler, error) {
	tmplFS := web.GetTemplatesFS()
	tmpl, err := template.ParseFS(tmplFS,
		"templates/layouts/*.html",
		"templates/pages/*.html",
		"templates/partials/*.html",
	)
	if err != nil {
		return nil, fmt.Errorf("falha ao parsear templates: %w", err)
	}

	return &Handler{
		healthService: healthService,
		tmpl:          tmpl,
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

	// Rotas de Páginas e HTML
	mux.HandleFunc("/", h.Home)
	mux.HandleFunc("/partials/health", h.HealthPartial)
	mux.HandleFunc("/_live_reload", LiveReloadHandler)

	// Rotas de API JSON
	mux.HandleFunc("/api/health", h.HealthAPI)
	mux.HandleFunc("/api/version", h.VersionAPI)

	return nil
}

// Home renderiza a página inicial.
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	health, err := h.healthService.Check(r.Context())
	if err != nil {
		http.Error(w, "Erro ao obter status do sistema", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Version": version.Get().Version,
		"Health":  health,
		"IsDev":   version.Get().IsDev,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
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
	if err := h.tmpl.ExecuteTemplate(w, "health_card.html", health); err != nil {
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
