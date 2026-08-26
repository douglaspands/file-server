package domain_test

import (
	"testing"

	"github.com/douglas/file-server/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{-1, "-"},
		{0, "0 B"},
		{512, "512 B"},
		{1023, "1023 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
		{1099511627776, "1.0 TB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, domain.FormatFileSize(tt.bytes))
		})
	}
}

func TestDetectCategory(t *testing.T) {
	tests := []struct {
		name     string
		isDir    bool
		expected domain.FileCategory
	}{
		{"pasta", true, domain.CategoryFolder},
		{"imagem.png", false, domain.CategoryImage},
		{"foto.JPG", false, domain.CategoryImage},
		{"video.mp4", false, domain.CategoryVideo},
		{"musica.mp3", false, domain.CategoryAudio},
		{"arquivo.zip", false, domain.CategoryArchive},
		{"doc.pdf", false, domain.CategoryDocument},
		{"main.go", false, domain.CategoryCode},
		{"desconhecido.xyz123", false, domain.CategoryOther},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, domain.DetectCategory(tt.name, tt.isDir))
		})
	}
}
