package handlers_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/douglas/file-server/internal/adapters/handlers"
	"github.com/douglas/file-server/internal/core/domain"
	"github.com/douglas/file-server/internal/core/services"
	"github.com/douglas/file-server/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestMux(t *testing.T, rootDir string) (*http.ServeMux, portsOrService) {
	t.Helper()
	fileSvc, err := services.NewFileService(rootDir)
	require.NoError(t, err)

	healthSvc := &testutils.MockHealthService{}
	handler, err := handlers.NewHandler(fileSvc, healthSvc)
	require.NoError(t, err)

	mux := http.NewServeMux()
	err = handler.RegisterRoutes(mux)
	require.NoError(t, err)

	return mux, portsOrService{fileSvc: fileSvc, healthSvc: healthSvc}
}

type portsOrService struct {
	fileSvc   *services.LocalFileService
	healthSvc *testutils.MockHealthService
}

func setupTestEnvironment(t *testing.T) (string, *http.ServeMux) {
	t.Helper()
	tempDir := t.TempDir()

	require.NoError(t, os.MkdirAll(filepath.Join(tempDir, "documents"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "sample.txt"), []byte("sample file content for download"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "documents", "guide.pdf"), []byte("guide content"), 0644))

	mux, _ := setupTestMux(t, tempDir)
	return tempDir, mux
}

func TestFileBrowserHandler(t *testing.T) {
	_, mux := setupTestEnvironment(t)

	t.Run("GET / renderiza explorador na raiz", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Header().Get("Content-Type"), "text/html")
		assert.Contains(t, rec.Body.String(), "sample.txt")
		assert.Contains(t, rec.Body.String(), "documents")
		assert.Contains(t, rec.Body.String(), `href="/view/sample.txt"`)
		assert.Contains(t, rec.Body.String(), `href="/download/sample.txt"`)
		assert.Contains(t, rec.Body.String(), `href="/zip/documents"`)
	})

	t.Run("GET /files/documents renderiza subdiretório", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/files/documents", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "guide.pdf")
	})

	t.Run("GET /files em caminho que é arquivo redireciona para download", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/files/sample.txt", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
		assert.Equal(t, "/download/sample.txt", rec.Header().Get("Location"))
	})

	t.Run("GET /files com path traversal retorna 403 Forbidden", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/files/%2e%2e/%2e%2e/etc/passwd", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("GET /files em diretório inexistente retorna 404 Not Found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/files/fantasma", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("POST /files retorna 405 Method Not Allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/files", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})
}

func TestFilesAPIHandler(t *testing.T) {
	_, mux := setupTestEnvironment(t)

	t.Run("GET /api/files retorna listagem em JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/files", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var listing domain.DirectoryListing
		err := json.Unmarshal(rec.Body.Bytes(), &listing)
		require.NoError(t, err)
		assert.False(t, listing.IsEmpty)
		assert.Equal(t, 2, listing.TotalItems)
	})

	t.Run("GET /api/files/documents retorna subdiretório em JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/files/documents", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var listing domain.DirectoryListing
		err := json.Unmarshal(rec.Body.Bytes(), &listing)
		require.NoError(t, err)
		assert.Equal(t, "documents", listing.CurrentPath)
		assert.Equal(t, 1, listing.TotalFiles)
	})

	t.Run("GET /api/files com path traversal retorna 403 JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/files/%2e%2e/%2e%2e", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("GET /api/files para arquivo retorna 400 JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/files/sample.txt", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestDownloadFileHandler(t *testing.T) {
	_, mux := setupTestEnvironment(t)

	t.Run("GET /download/sample.txt faz download com cabeçalho de anexo", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/download/sample.txt", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Header().Get("Content-Disposition"), `attachment; filename="sample.txt"`)
		assert.Equal(t, "sample file content for download", rec.Body.String())
	})

	t.Run("HEAD /download/sample.txt retorna cabeçalhos de download", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodHead, "/download/sample.txt", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Header().Get("Content-Disposition"), `attachment; filename="sample.txt"`)
		assert.Empty(t, rec.Body.String())
	})

	t.Run("GET /download/sample.txt com header Range retorna HTTP 206 Partial Content", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/download/sample.txt", nil)
		req.Header.Set("Range", "bytes=0-5")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusPartialContent, rec.Code)
		assert.Equal(t, "sample", rec.Body.String())
		assert.Contains(t, rec.Header().Get("Content-Range"), "bytes 0-5/")
	})

	t.Run("GET /download/documents (pasta) retorna 400 Bad Request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/download/documents", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("GET /download/inexistente retorna 404 Not Found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/download/inexistente.zip", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("GET /download com path traversal retorna 403 Forbidden", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/download/%2e%2e/%2e%2e/etc/shadow", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}

func TestViewFileHandler(t *testing.T) {
	_, mux := setupTestEnvironment(t)

	t.Run("GET /view/sample.txt exibe arquivo inline com cabeçalho inline", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/view/sample.txt", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Header().Get("Content-Disposition"), `inline; filename="sample.txt"`)
		assert.Contains(t, rec.Header().Get("Content-Type"), "text/plain")
		assert.Equal(t, "sample file content for download", rec.Body.String())
	})

	t.Run("GET /view/documents/guide.pdf exibe PDF inline com Content-Type application/pdf", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/view/documents/guide.pdf", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Header().Get("Content-Disposition"), `inline; filename="guide.pdf"`)
		assert.Equal(t, "application/pdf", rec.Header().Get("Content-Type"))
		assert.Equal(t, "guide content", rec.Body.String())
	})

	t.Run("HEAD /view/sample.txt retorna cabeçalhos inline sem corpo", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodHead, "/view/sample.txt", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Header().Get("Content-Disposition"), `inline; filename="sample.txt"`)
		assert.Empty(t, rec.Body.String())
	})

	t.Run("GET /view/sample.txt com Range retorna HTTP 206 Partial Content", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/view/sample.txt", nil)
		req.Header.Set("Range", "bytes=0-5")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusPartialContent, rec.Code)
		assert.Equal(t, "sample", rec.Body.String())
		assert.Contains(t, rec.Header().Get("Content-Range"), "bytes 0-5/")
		assert.Contains(t, rec.Header().Get("Content-Disposition"), `inline; filename="sample.txt"`)
	})

	t.Run("GET /view/documents (pasta) retorna 400 Bad Request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/view/documents", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("GET /view/inexistente retorna 404 Not Found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/view/inexistente.pdf", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("GET /view com path traversal retorna 403 Forbidden", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/view/%2e%2e/%2e%2e/etc/passwd", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("POST /view retorna 405 Method Not Allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/view/sample.txt", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})
}

func TestDownloadZipHandler(t *testing.T) {
	_, mux := setupTestEnvironment(t)

	t.Run("GET /zip/documents retorna arquivo ZIP em streaming", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/zip/documents", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/zip", rec.Header().Get("Content-Type"))
		assert.Contains(t, rec.Header().Get("Content-Disposition"), "documents.zip")
		assert.Greater(t, rec.Body.Len(), 0)
	})

	t.Run("GET /zip/ (raiz) retorna ZIP de toda a raiz", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/zip/", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/zip", rec.Header().Get("Content-Type"))
		assert.Greater(t, rec.Body.Len(), 0)
	})

	t.Run("GET /zip/inexistente retorna 404 Not Found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/zip/pasta_inexistente", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestUploadFilesHandler(t *testing.T) {
	tempDir, mux := setupTestEnvironment(t)

	createMultipartRequest := func(urlPath string, files map[string]string, acceptHeader string) *http.Request {
		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		for filename, content := range files {
			part, err := writer.CreateFormFile("files", filename)
			if err == nil {
				_, _ = part.Write([]byte(content))
			}
		}
		_ = writer.Close()

		req := httptest.NewRequest(http.MethodPost, urlPath, &body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		if acceptHeader != "" {
			req.Header.Set("Accept", acceptHeader)
		}
		return req
	}

	t.Run("POST /upload/documents com múltiplos arquivos e Accept JSON", func(t *testing.T) {
		files := map[string]string{
			"nota1.txt": "primeira nota",
			"nota2.txt": "segunda nota",
		}
		req := createMultipartRequest("/upload/documents", files, "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var results []domain.UploadResult
		err := json.Unmarshal(rec.Body.Bytes(), &results)
		require.NoError(t, err)
		assert.Len(t, results, 2)

		assert.FileExists(t, filepath.Join(tempDir, "documents", "nota1.txt"))
		assert.FileExists(t, filepath.Join(tempDir, "documents", "nota2.txt"))
	})

	t.Run("POST /upload/documents com redirecionamento de navegador", func(t *testing.T) {
		files := map[string]string{
			"web_upload.txt": "conteudo web",
		}
		req := createMultipartRequest("/upload/documents", files, "")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusSeeOther, rec.Code)
		assert.Equal(t, "/files/documents", rec.Header().Get("Location"))
	})

	t.Run("POST /upload sem arquivos retorna 400 Bad Request", func(t *testing.T) {
		var body bytes.Buffer
		writer := multipart.NewWriter(&body)
		_ = writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/upload", &body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("POST /upload/documents com campo file unico", func(t *testing.T) {
		var body bytes.Buffer
		writer := multipart.NewWriter(&body)
		part, err := writer.CreateFormFile("file", "single.txt")
		require.NoError(t, err)
		_, _ = part.Write([]byte("single file data"))
		_ = writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/upload/documents", &body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusSeeOther, rec.Code)
		assert.FileExists(t, filepath.Join(tempDir, "documents", "single.txt"))
	})

	t.Run("GET /upload retorna 405 Method Not Allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/upload", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("POST /upload com path traversal retorna 403 Forbidden", func(t *testing.T) {
		files := map[string]string{"evil.txt": "evil"}
		req := createMultipartRequest("/upload/%2e%2e/%2e%2e/etc", files, "")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("POST /upload em arquivo ao inves de pasta retorna 400 Bad Request", func(t *testing.T) {
		files := map[string]string{"test.txt": "test"}
		req := createMultipartRequest("/upload/sample.txt", files, "")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("POST /download retorna 405 Method Not Allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/download/sample.txt", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("POST /zip retorna 405 Method Not Allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/zip/documents", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("GET /zip em arquivo ao inves de pasta retorna 400 Bad Request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/zip/sample.txt", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
