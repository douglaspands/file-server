package ports

import (
	"context"

	"github.com/douglas/file-server/internal/core/domain"
)

// HealthService define a porta de entrada para verificação de integridade do sistema.
type HealthService interface {
	Check(ctx context.Context) (domain.HealthStatus, error)
}

// SystemInfoService define a porta de entrada para obter informações do sistema.
type SystemInfoService interface {
	GetSystemInfo(ctx context.Context) (domain.SystemInfo, error)
}
