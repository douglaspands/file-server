package services_test

import (
	"archive/zip"
	"bytes"
	"context"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
	"testing"

	"github.com/douglas/file-server/internal/core/domain"
	"github.com/douglas/file-server/internal/core/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDir(t *testing.T) string {
	t.Helper()
	tempDir := t.TempDir()

	// Cria estrutura de teste:
	// tempDir/
	//   ├─ docs/
	//   │    └─ manual.pdf
	//   ├─ empty/
	//   ├─ image.png
	//   └─ notes.txt
	require.NoError(t, os.MkdirAll(filepath.Join(tempDir, "docs"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(tempDir, "empty"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "notes.txt"), []byte("hello world"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "image.png"), []byte("fake png content"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "docs", "manual.pdf"), []byte("manual content"), 0644))

	return tempDir
}

func TestNewFileService(t *testing.T) {
	t.Run("sucesso com diretório válido", func(t *testing.T) {
		tempDir := t.TempDir()
		svc, err := services.NewFileService(tempDir)
		require.NoError(t, err)
		assert.NotEmpty(t, svc.GetRootDir())
	})

	t.Run("fallback automático quando rootDir for vazio", func(t *testing.T) {
		svc, err := services.NewFileService("")
		require.NoError(t, err)
		assert.NotEmpty(t, svc.GetRootDir())
	})

	t.Run("erro com diretório inexistente", func(t *testing.T) {
		_, err := services.NewFileService(filepath.Join(t.TempDir(), "nonexistent"))
		assert.Error(t, err)
	})

	t.Run("erro quando o caminho for um arquivo", func(t *testing.T) {
		tempFile := filepath.Join(t.TempDir(), "file.txt")
		require.NoError(t, os.WriteFile(tempFile, []byte("data"), 0644))
		_, err := services.NewFileService(tempFile)
		assert.Error(t, err)
	})
}

func TestLocalFileService_ResolveAndValidatePath_Security(t *testing.T) {
	tempDir := setupTestDir(t)
	outsideDir := t.TempDir()
	outsideFile := filepath.Join(outsideDir, "secret.txt")
	require.NoError(t, os.WriteFile(outsideFile, []byte("secret"), 0644))

	// Cria symlink apontando para fora do sandbox
	symlinkOutside := filepath.Join(tempDir, "sym_outside")
	_ = os.Symlink(outsideDir, symlinkOutside)

	// Cria symlink seguro apontando para dentro do sandbox
	symlinkInside := filepath.Join(tempDir, "sym_inside")
	_ = os.Symlink(filepath.Join(tempDir, "docs"), symlinkInside)

	svc, err := services.NewFileService(tempDir)
	require.NoError(t, err)

	traversalCases := []struct {
		name    string
		relPath string
	}{
		{"dois pontos simples", ".."},
		{"dois pontos com caminho", "../outside"},
		{"múltiplos dois pontos", "../../etc/passwd"},
		{"caminho absoluto /etc", "/etc/passwd"},
		{"caminho com escape embutido", "docs/../../secret"},
		{"symlink para fora", "sym_outside"},
	}

	for _, tc := range traversalCases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := svc.ResolveAndValidatePath(tc.relPath)
			assert.Error(t, err, "deve bloquear: %s", tc.relPath)
		})
	}

	t.Run("caminho válido na raiz", func(t *testing.T) {
		absPath, cleanRel, err := svc.ResolveAndValidatePath("")
		require.NoError(t, err)
		assert.Equal(t, svc.GetRootDir(), absPath)
		assert.Equal(t, "", cleanRel)
	})

	t.Run("caminho válido em subdiretório", func(t *testing.T) {
		absPath, cleanRel, err := svc.ResolveAndValidatePath("docs/manual.pdf")
		require.NoError(t, err)
		assert.Equal(t, filepath.Join(svc.GetRootDir(), "docs", "manual.pdf"), absPath)
		assert.Equal(t, "docs/manual.pdf", cleanRel)
	})
}

func TestLocalFileService_ListDirectory(t *testing.T) {
	tempDir := setupTestDir(t)
	svc, err := services.NewFileService(tempDir)
	require.NoError(t, err)
	ctx := context.Background()

	t.Run("listagem do diretório raiz", func(t *testing.T) {
		listing, err := svc.ListDirectory(ctx, "")
		require.NoError(t, err)
		assert.False(t, listing.IsEmpty)
		assert.True(t, listing.CanUpload)
		assert.Equal(t, "", listing.CurrentPath)
		assert.Equal(t, 2, listing.TotalDirs)  // docs, empty
		assert.Equal(t, 2, listing.TotalFiles) // image.png, notes.txt
		assert.Equal(t, 4, listing.TotalItems)
		assert.Greater(t, listing.TotalSize, int64(0))

		// Verifica que as pastas vêm primeiro e IsViewable
		require.GreaterOrEqual(t, len(listing.Items), 4)
		assert.True(t, listing.Items[0].IsDir)
		assert.False(t, listing.Items[0].IsViewable)
		assert.True(t, listing.Items[1].IsDir)
		assert.False(t, listing.Items[1].IsViewable)
		assert.False(t, listing.Items[2].IsDir)
		assert.True(t, listing.Items[2].IsViewable) // image.png ou notes.txt

		// Verifica breadcrumbs
		require.Len(t, listing.Breadcrumbs, 1)
		assert.Equal(t, "Início", listing.Breadcrumbs[0].Name)
		assert.True(t, listing.Breadcrumbs[0].IsCurrent)
	})

	t.Run("listagem de subdiretório com múltiplos níveis e breadcrumbs", func(t *testing.T) {
		subDir := filepath.Join(tempDir, "docs", "sublevel")
		require.NoError(t, os.MkdirAll(subDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(subDir, "code.go"), []byte("package main"), 0644))

		listing, err := svc.ListDirectory(ctx, "docs/sublevel")
		require.NoError(t, err)
		assert.Equal(t, "docs/sublevel", listing.CurrentPath)
		assert.Equal(t, 1, listing.TotalFiles)
		assert.Equal(t, domain.CategoryCode, listing.Items[0].Category)

		// Verifica breadcrumbs
		require.Len(t, listing.Breadcrumbs, 3)
		assert.Equal(t, "Início", listing.Breadcrumbs[0].Name)
		assert.False(t, listing.Breadcrumbs[0].IsCurrent)
		assert.Equal(t, "docs", listing.Breadcrumbs[1].Name)
		assert.False(t, listing.Breadcrumbs[1].IsCurrent)
		assert.Equal(t, "sublevel", listing.Breadcrumbs[2].Name)
		assert.True(t, listing.Breadcrumbs[2].IsCurrent)
	})

	t.Run("listagem de diretório vazio", func(t *testing.T) {
		listing, err := svc.ListDirectory(ctx, "empty")
		require.NoError(t, err)
		assert.True(t, listing.IsEmpty)
		assert.Equal(t, 0, listing.TotalItems)
		assert.Equal(t, 0, listing.TotalFiles)
		assert.Equal(t, 0, listing.TotalDirs)
	})

	t.Run("erro ao listar caminho que é arquivo", func(t *testing.T) {
		_, err := svc.ListDirectory(ctx, "notes.txt")
		assert.ErrorIs(t, err, services.ErrNotADirectory)
	})

	t.Run("erro ao listar caminho inexistente", func(t *testing.T) {
		_, err := svc.ListDirectory(ctx, "nao_existe")
		assert.ErrorIs(t, err, services.ErrNotFound)
	})

	t.Run("cancelamento com contexto", func(t *testing.T) {
		cancelCtx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := svc.ListDirectory(cancelCtx, "")
		assert.Error(t, err)
	})
}

func TestLocalFileService_GetFile(t *testing.T) {
	tempDir := setupTestDir(t)
	svc, err := services.NewFileService(tempDir)
	require.NoError(t, err)
	ctx := context.Background()

	t.Run("sucesso ao obter arquivo existente", func(t *testing.T) {
		file, info, err := svc.GetFile(ctx, "notes.txt")
		require.NoError(t, err)
		defer file.Close()

		assert.Equal(t, "notes.txt", info.Name())
		assert.Equal(t, int64(11), info.Size())

		buf := make([]byte, 5)
		n, readErr := file.Read(buf)
		require.NoError(t, readErr)
		assert.Equal(t, 5, n)
		assert.Equal(t, "hello", string(buf))
	})

	t.Run("erro ao tentar obter diretório como arquivo", func(t *testing.T) {
		_, _, err := svc.GetFile(ctx, "docs")
		assert.ErrorIs(t, err, services.ErrIsDirectory)
	})

	t.Run("erro ao obter arquivo inexistente", func(t *testing.T) {
		_, _, err := svc.GetFile(ctx, "fantasma.txt")
		assert.ErrorIs(t, err, services.ErrNotFound)
	})
}

func TestLocalFileService_StreamZip(t *testing.T) {
	tempDir := setupTestDir(t)
	svc, err := services.NewFileService(tempDir)
	require.NoError(t, err)
	ctx := context.Background()

	t.Run("sucesso ao gerar ZIP de diretório com arquivos e subpastas", func(t *testing.T) {
		var buf bytes.Buffer
		err := svc.StreamZip(ctx, "", &buf)
		require.NoError(t, err)

		zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
		require.NoError(t, err)

		fileNames := make(map[string]bool)
		for _, f := range zipReader.File {
			fileNames[f.Name] = true
		}

		assert.True(t, fileNames["notes.txt"])
		assert.True(t, fileNames["image.png"])
		assert.True(t, fileNames["docs/manual.pdf"] || fileNames["docs/"])
	})

	t.Run("sucesso ao gerar ZIP de pasta vazia", func(t *testing.T) {
		var buf bytes.Buffer
		err := svc.StreamZip(ctx, "empty", &buf)
		require.NoError(t, err)

		zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
		require.NoError(t, err)
		assert.Empty(t, zipReader.File)
	})

	t.Run("erro ao tentar ZIP de arquivo", func(t *testing.T) {
		var buf bytes.Buffer
		err := svc.StreamZip(ctx, "notes.txt", &buf)
		assert.ErrorIs(t, err, services.ErrNotADirectory)
	})

	t.Run("interrupção de streaming com contexto cancelado", func(t *testing.T) {
		cancelCtx, cancel := context.WithCancel(context.Background())
		cancel()

		var buf bytes.Buffer
		err := svc.StreamZip(cancelCtx, "", &buf)
		assert.Error(t, err)
	})
}

func TestLocalFileService_SaveUploadedFiles(t *testing.T) {
	tempDir := setupTestDir(t)
	svc, err := services.NewFileService(tempDir)
	require.NoError(t, err)
	ctx := context.Background()

	createHeader := func(filename, content string) *multipart.FileHeader {
		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", `form-data; name="files"; filename="`+filename+`"`)
		h.Set("Content-Type", "text/plain")
		part, _ := w.CreatePart(h)
		_, _ = part.Write([]byte(content))
		_ = w.Close()

		r := multipart.NewReader(&b, w.Boundary())
		form, _ := r.ReadForm(1024 * 1024)
		return form.File["files"][0]
	}

	t.Run("upload de arquivos múltiplos com sucesso", func(t *testing.T) {
		h1 := createHeader("upload1.txt", "conteudo 1")
		h2 := createHeader("upload2.txt", "conteudo 2")

		results, err := svc.SaveUploadedFiles(ctx, "docs", []*multipart.FileHeader{h1, h2})
		require.NoError(t, err)
		require.Len(t, results, 2)
		assert.True(t, results[0].Success)
		assert.True(t, results[1].Success)

		saved1, err := os.ReadFile(filepath.Join(tempDir, "docs", "upload1.txt"))
		require.NoError(t, err)
		assert.Equal(t, "conteudo 1", string(saved1))
	})

	t.Run("upload com nome malicioso sanitizado", func(t *testing.T) {
		h := createHeader("../../malicious.sh", "echo evil")

		results, err := svc.SaveUploadedFiles(ctx, "docs", []*multipart.FileHeader{h})
		require.NoError(t, err)
		require.Len(t, results, 1)
		assert.True(t, results[0].Success)
		assert.Equal(t, "malicious.sh", results[0].Filename)

		// Deve salvar dentro de docs/malicious.sh, nunca fora
		assert.FileExists(t, filepath.Join(tempDir, "docs", "malicious.sh"))
		assert.NoFileExists(t, filepath.Join(tempDir, "malicious.sh"))
	})

	t.Run("erro em diretório de destino inválido ou inexistente", func(t *testing.T) {
		h := createHeader("test.txt", "test")
		_, err := svc.SaveUploadedFiles(ctx, "diretorio_fantasma", []*multipart.FileHeader{h})
		assert.ErrorIs(t, err, services.ErrNotFound)
	})

	t.Run("erro ao tentar salvar upload em caminho que é arquivo", func(t *testing.T) {
		h := createHeader("test.txt", "test")
		_, err := svc.SaveUploadedFiles(ctx, "notes.txt", []*multipart.FileHeader{h})
		assert.ErrorIs(t, err, services.ErrNotADirectory)
	})

	t.Run("upload com nome vazio ou somente pontos", func(t *testing.T) {
		h := createHeader(".", "test")
		results, err := svc.SaveUploadedFiles(ctx, "docs", []*multipart.FileHeader{h})
		require.NoError(t, err)
		require.Len(t, results, 1)
		assert.False(t, results[0].Success)
	})

	t.Run("cancelamento com contexto no upload", func(t *testing.T) {
		cancelCtx, cancel := context.WithCancel(context.Background())
		cancel()
		h := createHeader("cancel.txt", "test")
		_, err := svc.SaveUploadedFiles(cancelCtx, "docs", []*multipart.FileHeader{h})
		assert.Error(t, err)
	})
}

func TestLocalFileService_SymlinksAndEdgeCases(t *testing.T) {
	tempDir := setupTestDir(t)
	outsideDir := t.TempDir()
	outsideFile := filepath.Join(outsideDir, "outside.txt")
	require.NoError(t, os.WriteFile(outsideFile, []byte("outside content"), 0644))

	// Symlink válido para arquivo
	validSymlinkFile := filepath.Join(tempDir, "sym_file.txt")
	_ = os.Symlink(filepath.Join(tempDir, "notes.txt"), validSymlinkFile)

	// Symlink válido para pasta
	validSymlinkDir := filepath.Join(tempDir, "sym_docs")
	_ = os.Symlink(filepath.Join(tempDir, "docs"), validSymlinkDir)

	// Symlink inválido (quebrado)
	brokenSymlink := filepath.Join(tempDir, "broken_sym")
	_ = os.Symlink(filepath.Join(tempDir, "non_existent_target"), brokenSymlink)

	// Symlink externo
	extSymlink := filepath.Join(tempDir, "ext_sym")
	_ = os.Symlink(outsideFile, extSymlink)

	svc, err := services.NewFileService(tempDir)
	require.NoError(t, err)
	ctx := context.Background()

	t.Run("listagem com symlinks variados", func(t *testing.T) {
		listing, err := svc.ListDirectory(ctx, "")
		require.NoError(t, err)
		assert.NotEmpty(t, listing.Items)
	})

	t.Run("zip streaming com symlinks internos", func(t *testing.T) {
		var buf bytes.Buffer
		err := svc.StreamZip(ctx, "", &buf)
		require.NoError(t, err)
		assert.Greater(t, buf.Len(), 0)
	})

	t.Run("tentativa com caractere nulo", func(t *testing.T) {
		_, _, err := svc.ResolveAndValidatePath("notes\x00.txt")
		assert.ErrorIs(t, err, services.ErrPathTraversal)
	})

	t.Run("zip de caminho inexistente", func(t *testing.T) {
		var buf bytes.Buffer
		err := svc.StreamZip(ctx, "caminho_inexistente", &buf)
		assert.ErrorIs(t, err, services.ErrNotFound)
	})
}
