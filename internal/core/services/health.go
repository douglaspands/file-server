package services

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/douglas/file-server/internal/core/domain"
	"github.com/douglas/file-server/internal/core/ports"
	"github.com/douglas/file-server/internal/version"
)

type healthService struct {
	startTime time.Time
}

// NewHealthService instancia o serviço de saúde do sistema.
func NewHealthService() ports.HealthService {
	return &healthService{
		startTime: time.Now(),
	}
}

// Check avalia o estado atual da aplicação.
func (s *healthService) Check(ctx context.Context) (domain.HealthStatus, error) {
	vInfo := version.Get()
	uptime := time.Since(s.startTime).Truncate(time.Second)

	return domain.HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now().UTC(),
		Version:   vInfo.Version,
		Uptime:    uptime.String(),
		Details: map[string]string{
			"go_version": vInfo.GoVersion,
			"platform":   vInfo.Platform,
			"num_cpu":    fmt.Sprintf("%d", runtime.NumCPU()),
		},
	}, nil
}
