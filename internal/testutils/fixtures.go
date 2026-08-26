package testutils

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"time"

	"github.com/douglas/file-server/internal/core/domain"
	"github.com/douglas/file-server/internal/core/ports"
)

// MockHealthService implementa ports.HealthService para testes.
type MockHealthService struct {
	Status domain.HealthStatus
	Err    error
}

var _ ports.HealthService = (*MockHealthService)(nil)

func (m *MockHealthService) Check(ctx context.Context) (domain.HealthStatus, error) {
	if m.Err != nil {
		return domain.HealthStatus{}, m.Err
	}
	if m.Status.Status == "" {
		return domain.HealthStatus{
			Status:    "healthy",
			Timestamp: time.Now().UTC(),
			Version:   "test-version",
			Uptime:    "10s",
			Details: map[string]string{
				"mock": "true",
			},
		}, nil
	}
	return m.Status, nil
}

// MockFileService implementa ports.FileService para testes de adaptadores HTTP.
type MockFileService struct {
	RootDir         string
	Listing         *domain.DirectoryListing
	ListErr         error
	File            *os.File
	FileInfo        os.FileInfo
	FileErr         error
	ZipErr          error
	UploadResults   []domain.UploadResult
	UploadErr       error
	ResolveAbsPath  string
	ResolveCleanRel string
	ResolveErr      error
}

var _ ports.FileService = (*MockFileService)(nil)

func (m *MockFileService) GetRootDir() string {
	if m.RootDir == "" {
		return "/mock/root"
	}
	return m.RootDir
}

func (m *MockFileService) ResolveAndValidatePath(relPath string) (string, string, error) {
	if m.ResolveErr != nil {
		return "", "", m.ResolveErr
	}
	return m.ResolveAbsPath, m.ResolveCleanRel, nil
}

func (m *MockFileService) ListDirectory(ctx context.Context, relPath string) (*domain.DirectoryListing, error) {
	if m.ListErr != nil {
		return nil, m.ListErr
	}
	if m.Listing != nil {
		return m.Listing, nil
	}
	return &domain.DirectoryListing{
		CurrentPath: relPath,
		Items:       []domain.FileItem{},
		Breadcrumbs: []domain.Breadcrumb{
			{Name: "Início", Path: "", IsCurrent: true},
		},
		CanUpload: true,
		IsEmpty:   true,
	}, nil
}

func (m *MockFileService) GetFile(ctx context.Context, relPath string) (*os.File, os.FileInfo, error) {
	if m.FileErr != nil {
		return nil, nil, m.FileErr
	}
	return m.File, m.FileInfo, nil
}

func (m *MockFileService) StreamZip(ctx context.Context, relPath string, w io.Writer) error {
	if m.ZipErr != nil {
		return m.ZipErr
	}
	_, err := w.Write([]byte("PK\x05\x06" + string(make([]byte, 18)))) // mock zip end header
	return err
}

func (m *MockFileService) SaveUploadedFiles(ctx context.Context, relDir string, files []*multipart.FileHeader) ([]domain.UploadResult, error) {
	if m.UploadErr != nil {
		return nil, m.UploadErr
	}
	if m.UploadResults != nil {
		return m.UploadResults, nil
	}
	var results []domain.UploadResult
	for _, f := range files {
		results = append(results, domain.UploadResult{
			Filename: f.Filename,
			Size:     f.Size,
			Success:  true,
		})
	}
	return results, nil
}
