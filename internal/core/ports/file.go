package ports

import (
	"context"
	"io"
	"mime/multipart"
	"os"

	"github.com/douglas/file-server/internal/core/domain"
)

// FileService define os contratos de serviço para navegação, download, streaming e manipulação de arquivos com sandboxing estrito.
type FileService interface {
	// GetRootDir retorna o caminho absoluto configurado como raiz do compartilhamento.
	GetRootDir() string

	// ResolveAndValidatePath resolve o caminho relativo contra o diretório raiz e valida contra path traversal.
	// Retorna o caminho absoluto canônico seguro e o caminho relativo limpo.
	ResolveAndValidatePath(relPath string) (absPath string, cleanRel string, err error)

	// ListDirectory lista os conteúdos de um diretório com metadados e breadcrumbs.
	ListDirectory(ctx context.Context, relPath string) (*domain.DirectoryListing, error)

	// GetFile recupera um arquivo aberto e seu FileInfo para download (compatível com http.ServeContent e Range Requests).
	// O chamador é responsável por fechar o *os.File retornado quando err == nil.
	GetFile(ctx context.Context, relPath string) (*os.File, os.FileInfo, error)

	// StreamZip compacta o diretório em streaming diretamente para o io.Writer sem resíduos em disco.
	StreamZip(ctx context.Context, relPath string, w io.Writer) error

	// SaveUploadedFiles grava arquivos multipart enviados no diretório de destino validado.
	SaveUploadedFiles(ctx context.Context, relDir string, files []*multipart.FileHeader) ([]domain.UploadResult, error)
}
