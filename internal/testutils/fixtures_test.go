package testutils_test

import (
	"context"
	"errors"
	"testing"

	"github.com/douglas/file-server/internal/core/domain"
	"github.com/douglas/file-server/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockHealthService(t *testing.T) {
	t.Run("Given default mock service When checking Then returns default status", func(t *testing.T) {
		mock := &testutils.MockHealthService{}
		status, err := mock.Check(context.Background())
		require.NoError(t, err)
		assert.Equal(t, "healthy", status.Status)
	})

	t.Run("Given mock service with custom status When checking Then returns configured status", func(t *testing.T) {
		mock := &testutils.MockHealthService{
			Status: domain.HealthStatus{Status: "custom"},
		}
		status, err := mock.Check(context.Background())
		require.NoError(t, err)
		assert.Equal(t, "custom", status.Status)
	})

	t.Run("Given mock service with error When checking Then returns error", func(t *testing.T) {
		mock := &testutils.MockHealthService{
			Err: errors.New("mock error"),
		}
		_, err := mock.Check(context.Background())
		require.Error(t, err)
	})
}
