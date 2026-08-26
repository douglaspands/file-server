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

func TestIsViewableFormat(t *testing.T) {
	tests := []struct {
		name     string
		category domain.FileCategory
		ext      string
		expected bool
	}{
		{"pasta", domain.CategoryFolder, "", false},
		{"zip archive", domain.CategoryArchive, ".zip", false},
		{"tar gz archive", domain.CategoryArchive, ".tar.gz", false},
		{"imagem png", domain.CategoryImage, ".png", true},
		{"imagem jpg", domain.CategoryImage, "jpg", true},
		{"video mp4", domain.CategoryVideo, ".mp4", true},
		{"video webm", domain.CategoryVideo, ".webm", true},
		{"audio mp3", domain.CategoryAudio, ".mp3", true},
		{"audio wav", domain.CategoryAudio, ".wav", true},
		{"codigo go", domain.CategoryCode, ".go", true},
		{"codigo js", domain.CategoryCode, ".js", true},
		{"documento pdf", domain.CategoryDocument, ".pdf", true},
		{"documento txt", domain.CategoryDocument, ".txt", true},
		{"documento md", domain.CategoryDocument, ".md", true},
		{"documento csv", domain.CategoryDocument, ".csv", true},
		{"documento docx", domain.CategoryDocument, ".docx", false},
		{"documento xlsx", domain.CategoryDocument, ".xlsx", false},
		{"outro json", domain.CategoryOther, ".json", true},
		{"outro log", domain.CategoryOther, ".log", true},
		{"outro binario", domain.CategoryOther, ".bin", false},
		{"outro exe", domain.CategoryOther, ".exe", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, domain.IsViewableFormat(tt.category, tt.ext))
		})
	}
}
