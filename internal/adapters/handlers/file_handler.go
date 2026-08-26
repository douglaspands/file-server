package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/douglas/file-server/internal/core/services"
	"github.com/douglas/file-server/internal/version"
)

// extractSubpath extrai o caminho relativo após um prefixo específico na URL.
func extractSubpath(r *http.Request, prefix string) string {
	rawPath := r.URL.Path
	trimmed := strings.TrimPrefix(rawPath, prefix)
	trimmed = strings.TrimPrefix(trimmed, "/")
	unescaped, err := url.PathUnescape(trimmed)
	if err != nil {
		return trimmed
	}
	return unescaped
}

// FileBrowserHandler renderiza a interface web do explorador de arquivos.
func (h *Handler) FileBrowserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var subpath string
	if r.URL.Path == "/" {
		subpath = ""
	} else {
		subpath = extractSubpath(r, "/files")
	}

	listing, err := h.fileService.ListDirectory(r.Context(), subpath)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrPathTraversal):
			http.Error(w, "Acesso negado: fora do diretório raiz", http.StatusForbidden)
		case errors.Is(err, services.ErrNotFound):
			http.Error(w, "Diretório não encontrado", http.StatusNotFound)
		case errors.Is(err, services.ErrNotADirectory):
			// Redireciona para download caso o usuário tenha clicado em um arquivo
			http.Redirect(w, r, "/download/"+subpath, http.StatusTemporaryRedirect)
		default:
			http.Error(w, fmt.Sprintf("Erro ao listar diretório: %v", err), http.StatusInternalServerError)
		}
		return
	}

	data := map[string]interface{}{
		"Listing": listing,
		"Version": version.Get().Version,
		"IsDev":   version.Get().IsDev,
		"RootDir": h.fileService.GetRootDir(),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, ok := h.templates["explorer.html"]
	if !ok {
		http.Error(w, "Template explorer.html não encontrado", http.StatusInternalServerError)
		return
	}
	if err := tmpl.ExecuteTemplate(w, "explorer.html", data); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao renderizar template: %v", err), http.StatusInternalServerError)
	}
}

// FilesAPIHandler retorna a listagem de arquivos e metadados em formato JSON.
func (h *Handler) FilesAPIHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	subpath := extractSubpath(r, "/api/files")
	listing, err := h.fileService.ListDirectory(r.Context(), subpath)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case errors.Is(err, services.ErrPathTraversal):
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "acesso negado: path traversal"})
		case errors.Is(err, services.ErrNotFound):
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "diretório não encontrado"})
		case errors.Is(err, services.ErrNotADirectory):
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "o caminho informado não é um diretório"})
		default:
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(listing)
}

// DownloadFileHandler realiza o streaming de download de arquivo com suporte a Range requests (HTTP 206).
func (h *Handler) DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	subpath := extractSubpath(r, "/download")
	file, info, err := h.fileService.GetFile(r.Context(), subpath)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrPathTraversal):
			http.Error(w, "Acesso negado: fora do diretório raiz", http.StatusForbidden)
		case errors.Is(err, services.ErrNotFound):
			http.Error(w, "Arquivo não encontrado", http.StatusNotFound)
		case errors.Is(err, services.ErrIsDirectory):
			http.Error(w, "O caminho informado é um diretório", http.StatusBadRequest)
		default:
			http.Error(w, fmt.Sprintf("Erro ao abrir arquivo: %v", err), http.StatusInternalServerError)
		}
		return
	}
	defer file.Close()

	// Header para forçar o download com o nome correto
	filename := filepath.Base(info.Name())
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	// ServeContent cuida automaticamente de Content-Type, Content-Length, HTTP Range (206) e Last-Modified
	http.ServeContent(w, r, info.Name(), info.ModTime(), file)
}

// DownloadZipHandler compacta um diretório em ZIP sob demanda via streaming direto para a resposta.
func (h *Handler) DownloadZipHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	subpath := extractSubpath(r, "/zip")

	// Determina o nome do arquivo ZIP baixado
	var zipName string
	if subpath == "" || subpath == "." {
		rootBase := filepath.Base(h.fileService.GetRootDir())
		if rootBase == "." || rootBase == "/" || rootBase == "" {
			zipName = "files.zip"
		} else {
			zipName = rootBase + ".zip"
		}
	} else {
		zipName = filepath.Base(subpath) + ".zip"
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, zipName))
	w.Header().Set("X-Content-Type-Options", "nosniff")

	err := h.fileService.StreamZip(r.Context(), subpath, w)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrPathTraversal):
			http.Error(w, "Acesso negado", http.StatusForbidden)
		case errors.Is(err, services.ErrNotFound):
			http.Error(w, "Diretório não encontrado", http.StatusNotFound)
		case errors.Is(err, services.ErrNotADirectory):
			http.Error(w, "O caminho não é um diretório", http.StatusBadRequest)
		default:
			// Se o streaming já tiver iniciado os bytes, o erro pode ser registrado
			return
		}
	}
}

// UploadFilesHandler processa o envio multipart de arquivos para o diretório indicado.
func (h *Handler) UploadFilesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	subpath := extractSubpath(r, "/upload")

	// Limite de 128 MB em memória antes de usar temporários
	if err := r.ParseMultipartForm(128 * 1024 * 1024); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao processar formulário multipart: %v", err), http.StatusBadRequest)
		return
	}

	var fileHeaders []*multipart.FileHeader
	if r.MultipartForm != nil && r.MultipartForm.File != nil {
		if files, ok := r.MultipartForm.File["files"]; ok {
			fileHeaders = append(fileHeaders, files...)
		}
		if file, ok := r.MultipartForm.File["file"]; ok {
			fileHeaders = append(fileHeaders, file...)
		}
	}

	if len(fileHeaders) == 0 {
		http.Error(w, "Nenhum arquivo enviado", http.StatusBadRequest)
		return
	}

	results, err := h.fileService.SaveUploadedFiles(r.Context(), subpath, fileHeaders)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrPathTraversal):
			http.Error(w, "Acesso negado: fora do diretório raiz", http.StatusForbidden)
		case errors.Is(err, services.ErrNotFound):
			http.Error(w, "Diretório de destino não encontrado", http.StatusNotFound)
		case errors.Is(err, services.ErrNotADirectory):
			http.Error(w, "O destino não é um diretório", http.StatusBadRequest)
		default:
			http.Error(w, fmt.Sprintf("Erro ao salvar arquivos: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Resposta em JSON se requisitada
	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(results)
		return
	}

	// Redireciona o navegador de volta para a visualização da pasta
	redirectPath := "/files/" + subpath
	if subpath == "" {
		redirectPath = "/"
	}
	http.Redirect(w, r, redirectPath, http.StatusSeeOther)
}
