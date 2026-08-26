package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/douglas/file-server/internal/adapters/handlers"
	"github.com/douglas/file-server/internal/core/domain"
	"github.com/douglas/file-server/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerRoutes(t *testing.T) {
	mockHealth := &testutils.MockHealthService{
		Status: domain.HealthStatus{
			Status:    "healthy",
			Timestamp: time.Now().UTC(),
			Version:   "v1.0.0",
			Uptime:    "5m",
			Details:   map[string]string{"env": "test"},
		},
	}
	mockFile := &testutils.MockFileService{
		RootDir: "/tmp/mock",
	}

	handler, err := handlers.NewHandler(mockFile, mockHealth)
	require.NoError(t, err)

	mux := http.NewServeMux()
	err = handler.RegisterRoutes(mux)
	require.NoError(t, err)

	t.Run("Given GET request to root (/) When path is valid Then renders file explorer successfully", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Header().Get("Content-Type"), "text/html")
		assert.Contains(t, rec.Body.String(), "File Server")
	})

	t.Run("Given GET request to /status When path is valid Then renders system status HTML", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/status", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Header().Get("Content-Type"), "text/html")
		assert.Contains(t, rec.Body.String(), "File Server - Status do Sistema")
	})

	t.Run("Given GET request to /partials/health When called via HTMX Then renders health card partial", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/partials/health", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Header().Get("Content-Type"), "text/html")
		assert.Contains(t, rec.Body.String(), "health-status-container")
		assert.Contains(t, rec.Body.String(), "healthy")
	})

	t.Run("Given GET request to /api/health When called Then returns JSON health status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var resp domain.HealthStatus
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "healthy", resp.Status)
		assert.Equal(t, "v1.0.0", resp.Version)
	})

	t.Run("Given GET request to /api/version When called Then returns JSON version info", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/version", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var resp map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Contains(t, resp, "version")
		assert.Contains(t, resp, "go_version")
	})

	t.Run("Given health service failure When requesting /api/health Then returns 500 error", func(t *testing.T) {
		failingMock := &testutils.MockHealthService{
			Err: errors.New("database connection failed"),
		}
		failingHandler, err := handlers.NewHandler(mockFile, failingMock)
		require.NoError(t, err)

		failingMux := http.NewServeMux()
		_ = failingHandler.RegisterRoutes(failingMux)

		req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
		rec := httptest.NewRecorder()

		failingMux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "database connection failed")
	})
}

func TestHandlerErrors(t *testing.T) {
	mockFile := &testutils.MockFileService{}

	t.Run("Given failing health service When calling /partials/health Then returns 500", func(t *testing.T) {
		failingMock := &testutils.MockHealthService{Err: errors.New("health check failure")}
		h, err := handlers.NewHandler(mockFile, failingMock)
		require.NoError(t, err)

		mux := http.NewServeMux()
		_ = h.RegisterRoutes(mux)

		req := httptest.NewRequest(http.MethodGet, "/partials/health", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("Given failing health service When calling /status Then returns 500", func(t *testing.T) {
		failingMock := &testutils.MockHealthService{Err: errors.New("health check failure")}
		h, err := handlers.NewHandler(mockFile, failingMock)
		require.NoError(t, err)

		mux := http.NewServeMux()
		_ = h.RegisterRoutes(mux)

		req := httptest.NewRequest(http.MethodGet, "/status", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestLiveReloadHandler(t *testing.T) {
	t.Run("Given GET request to /_live_reload When client connects Then streams SSE headers and connected event", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/_live_reload", nil)
		ctx, cancel := req.Context(), func() {}
		req = req.WithContext(ctx)

		rec := httptest.NewRecorder()
		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()

		go handlers.LiveReloadHandler(rec, req)
		time.Sleep(100 * time.Millisecond)

		assert.Equal(t, "text/event-stream", rec.Header().Get("Content-Type"))
		assert.Contains(t, rec.Body.String(), "data: connected")
	})
}
