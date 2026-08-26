package testutils_test

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
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

func TestMockFileService(t *testing.T) {
	t.Run("Given default mock file service When calling methods Then returns defaults", func(t *testing.T) {
		mock := &testutils.MockFileService{}
		assert.Equal(t, "/mock/root", mock.GetRootDir())

		abs, rel, err := mock.ResolveAndValidatePath("sub")
		require.NoError(t, err)
		assert.Empty(t, abs)
		assert.Empty(t, rel)

		listing, err := mock.ListDirectory(context.Background(), "")
		require.NoError(t, err)
		assert.True(t, listing.IsEmpty)

		_, _, err = mock.GetFile(context.Background(), "file.txt")
		require.NoError(t, err)

		var buf bytes.Buffer
		err = mock.StreamZip(context.Background(), "", &buf)
		require.NoError(t, err)
		assert.NotEmpty(t, buf.Bytes())

		results, err := mock.SaveUploadedFiles(context.Background(), "", []*multipart.FileHeader{
			{Filename: "upload.txt", Size: 10},
		})
		require.NoError(t, err)
		require.Len(t, results, 1)
		assert.True(t, results[0].Success)
	})

	t.Run("Given mock file service with errors When calling methods Then returns errors", func(t *testing.T) {
		testErr := errors.New("file service error")
		mock := &testutils.MockFileService{
			RootDir:       "/custom/root",
			ListErr:       testErr,
			FileErr:       testErr,
			ZipErr:        testErr,
			UploadErr:     testErr,
			ResolveErr:    testErr,
			UploadResults: []domain.UploadResult{{Filename: "custom.txt", Success: true}},
		}

		assert.Equal(t, "/custom/root", mock.GetRootDir())
		_, _, err := mock.ResolveAndValidatePath("test")
		assert.ErrorIs(t, err, testErr)

		_, err = mock.ListDirectory(context.Background(), "")
		assert.ErrorIs(t, err, testErr)

		_, _, err = mock.GetFile(context.Background(), "test")
		assert.ErrorIs(t, err, testErr)

		err = mock.StreamZip(context.Background(), "", &bytes.Buffer{})
		assert.ErrorIs(t, err, testErr)

		_, err = mock.SaveUploadedFiles(context.Background(), "", nil)
		assert.ErrorIs(t, err, testErr)
	})
}
