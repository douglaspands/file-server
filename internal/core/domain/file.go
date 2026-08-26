package domain

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// FileCategory define as categorias visuais e semânticas de arquivos.
type FileCategory string

const (
	CategoryFolder   FileCategory = "folder"
	CategoryDocument FileCategory = "document"
	CategoryImage    FileCategory = "image"
	CategoryVideo    FileCategory = "video"
	CategoryAudio    FileCategory = "audio"
	CategoryArchive  FileCategory = "archive"
	CategoryCode     FileCategory = "code"
	CategoryOther    FileCategory = "other"
)

// FileItem representa um arquivo ou diretório contido no sistema servido.
type FileItem struct {
	Name             string       `json:"name"`
	RelativePath     string       `json:"relativePath"`
	IsDir            bool         `json:"isDir"`
	Size             int64        `json:"size"`
	FormattedSize    string       `json:"formattedSize"`
	ModTime          time.Time    `json:"modTime"`
	FormattedModTime string       `json:"formattedModTime"`
	Extension        string       `json:"extension"`
	Category         FileCategory `json:"category"`
	IsViewable       bool         `json:"isViewable"`
}

// Breadcrumb representa um elemento no caminho hierárquico de navegação.
type Breadcrumb struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	IsCurrent bool   `json:"isCurrent"`
}

// DirectoryListing representa os dados consolidados da visualização de uma pasta.
type DirectoryListing struct {
	CurrentPath        string       `json:"currentPath"`
	Breadcrumbs        []Breadcrumb `json:"breadcrumbs"`
	Items              []FileItem   `json:"items"`
	TotalItems         int          `json:"totalItems"`
	TotalDirs          int          `json:"totalDirs"`
	TotalFiles         int          `json:"totalFiles"`
	TotalSize          int64        `json:"totalSize"`
	FormattedTotalSize string       `json:"formattedTotalSize"`
	IsEmpty            bool         `json:"isEmpty"`
	CanUpload          bool         `json:"canUpload"`
}

// UploadResult representa o resultado do processamento de upload de um arquivo.
type UploadResult struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	Success  bool   `json:"success"`
	Message  string `json:"message,omitempty"`
}

// FormatFileSize formata o tamanho em bytes para representação legível (B, KB, MB, GB, TB).
func FormatFileSize(bytes int64) string {
	if bytes < 0 {
		return "-"
	}
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	units := []string{"KB", "MB", "GB", "TB", "PB"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}

// DetectCategory classifica um arquivo pela extensão e se é diretório.
func DetectCategory(name string, isDir bool) FileCategory {
	if isDir {
		return CategoryFolder
	}

	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".svg", ".webp", ".bmp", ".ico", ".tiff":
		return CategoryImage
	case ".mp4", ".mkv", ".avi", ".mov", ".wmv", ".flv", ".webm", ".m4v":
		return CategoryVideo
	case ".mp3", ".wav", ".flac", ".aac", ".ogg", ".wma", ".m4a":
		return CategoryAudio
	case ".zip", ".tar", ".gz", ".tgz", ".rar", ".7z", ".bz2", ".xz", ".iso":
		return CategoryArchive
	case ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt", ".rtf", ".odt", ".ods", ".odp", ".csv", ".tsv":
		return CategoryDocument
	case ".go", ".js", ".ts", ".jsx", ".tsx", ".html", ".css", ".scss", ".json", ".xml", ".yaml", ".yml", ".py", ".rs", ".java", ".c", ".cpp", ".h", ".hpp", ".sh", ".bash", ".sql", ".md", ".env":
		return CategoryCode
	default:
		return CategoryOther
	}
}

// IsViewableFormat determina se um arquivo com a categoria e extensão especificadas pode ser visualizado nativamente no navegador.
func IsViewableFormat(category FileCategory, ext string) bool {
	if category == CategoryFolder || category == CategoryArchive {
		return false
	}
	cleanExt := strings.ToLower(ext)
	if !strings.HasPrefix(cleanExt, ".") && cleanExt != "" {
		cleanExt = "." + cleanExt
	}
	switch category {
	case CategoryImage, CategoryVideo, CategoryAudio, CategoryCode:
		return true
	case CategoryDocument:
		switch cleanExt {
		case ".pdf", ".txt", ".csv", ".tsv", ".md", ".log":
			return true
		default:
			return false
		}
	default:
		switch cleanExt {
		case ".txt", ".log", ".cfg", ".conf", ".ini", ".env", ".json", ".xml", ".yaml", ".yml":
			return true
		default:
			return false
		}
	}
}
