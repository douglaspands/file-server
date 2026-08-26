package version_test

import (
	"strings"
	"testing"

	"github.com/douglas/file-server/internal/version"
	"github.com/stretchr/testify/assert"
)

func TestVersionInfo(t *testing.T) {
	t.Run("Given development environment When getting version Then returns dev status", func(t *testing.T) {
		version.Version = "dev"
		version.Commit = "abcdef"
		version.Date = "2026-08-25T00:00:00Z"

		info := version.Get()

		assert.Equal(t, "dev", info.Version)
		assert.Equal(t, "abcdef", info.Commit)
		assert.Equal(t, "2026-08-25T00:00:00Z", info.BuildDate)
		assert.True(t, info.IsDev)
		assert.Contains(t, info.String(), "development build")
	})

	t.Run("Given release version When getting version Then returns release status", func(t *testing.T) {
		version.Version = "v1.2.3"
		version.Commit = "123456"
		version.Date = "2026-08-25T00:00:00Z"

		info := version.Get()

		assert.Equal(t, "v1.2.3", info.Version)
		assert.False(t, info.IsDev)
		assert.False(t, strings.Contains(info.String(), "development build"))
		assert.Contains(t, info.String(), "File Server version: v1.2.3")
	})

	t.Run("Given empty version string When getting version Then falls back to dev", func(t *testing.T) {
		version.Version = ""
		info := version.Get()

		assert.Equal(t, "dev", info.Version)
		assert.True(t, info.IsDev)
	})
}
