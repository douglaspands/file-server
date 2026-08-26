package services_test

import (
	"context"
	"testing"

	"github.com/douglas/file-server/internal/core/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthService(t *testing.T) {
	t.Run("Given initialized health service When checking health Then returns healthy status with details", func(t *testing.T) {
		svc := services.NewHealthService()
		ctx := context.Background()

		status, err := svc.Check(ctx)

		require.NoError(t, err)
		assert.Equal(t, "healthy", status.Status)
		assert.NotEmpty(t, status.Version)
		assert.NotEmpty(t, status.Uptime)
		assert.False(t, status.Timestamp.IsZero())
		assert.Contains(t, status.Details, "go_version")
		assert.Contains(t, status.Details, "platform")
		assert.Contains(t, status.Details, "num_cpu")
	})
}
