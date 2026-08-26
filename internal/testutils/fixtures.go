package testutils

import (
	"context"
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
