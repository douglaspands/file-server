package services

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/douglas/file-server/internal/core/domain"
	"github.com/douglas/file-server/internal/core/ports"
)

var (
	// ErrPathTraversal indica tentativa de acessar caminhos fora da raiz configurada.
	ErrPathTraversal = errors.New("tentativa de path traversal detectada: acesso negado fora do diretório raiz")

	// ErrNotFound indica que o arquivo ou pasta solicitada não existe.
	ErrNotFound = errors.New("arquivo ou diretório não encontrado")

	// ErrNotADirectory indica que o caminho esperado como diretório é um arquivo.
	ErrNotADirectory = errors.New("o caminho solicitado não é um diretório")

	// ErrIsDirectory indica que o caminho esperado como arquivo é um diretório.
	ErrIsDirectory = errors.New("o caminho solicitado é um diretório, não um arquivo")

	// ErrPermissionDenied indica falha de permissão no sistema de arquivos.
	ErrPermissionDenied = errors.New("permissão negada para acessar o caminho")
)

// LocalFileService implementa ports.FileService operando diretamente no sistema de arquivos local com sandbox.
type LocalFileService struct {
	rootDir string
}

// Garante que LocalFileService implementa ports.FileService.
var _ ports.FileService = (*LocalFileService)(nil)

// NewFileService cria uma nova instância de LocalFileService validando e normalizando o diretório raiz.
func NewFileService(rootDir string) (*LocalFileService, error) {
	if rootDir == "" {
		rootDir = "."
	}

	absPath, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, fmt.Errorf("caminho de diretório inválido: %w", err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("diretório raiz não existe: %s", absPath)
		}
		return nil, fmt.Errorf("erro ao acessar diretório raiz: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("o caminho informado não é um diretório: %s", absPath)
	}

	canonicalRoot, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		canonicalRoot = absPath
	}

	return &LocalFileService{
		rootDir: canonicalRoot,
	}, nil
}

// GetRootDir retorna o caminho absoluto canônico do diretório raiz.
func (s *LocalFileService) GetRootDir() string {
	return s.rootDir
}

// ResolveAndValidatePath resolve o caminho relativo contra o diretório raiz e valida contra path traversal.
func (s *LocalFileService) ResolveAndValidatePath(relPath string) (string, string, error) {
	// Se contiver caracteres nulos, bloqueia imediatamente
	if strings.Contains(relPath, "\x00") {
		return "", "", ErrPathTraversal
	}

	normalized := filepath.ToSlash(relPath)

	// Se o caminho for absoluto no SO e não estiver contido na raiz
	if filepath.IsAbs(relPath) && !strings.HasPrefix(filepath.Clean(relPath), s.rootDir) {
		return "", "", ErrPathTraversal
	}

	// Validação de contagem de profundidade de segmentos ..
	parts := strings.Split(normalized, "/")
	depth := 0
	for _, p := range parts {
		if p == "" || p == "." {
			continue
		}
		if p == ".." {
			depth--
			if depth < 0 {
				return "", "", ErrPathTraversal
			}
		} else {
			depth++
		}
	}

	targetPath := filepath.Clean(filepath.Join(s.rootDir, filepath.FromSlash(normalized)))

	// Validação de fronteira 1: checagem do caminho relativo computado
	rel, err := filepath.Rel(s.rootDir, targetPath)
	if err != nil || strings.HasPrefix(rel, "..") || rel == ".." || filepath.IsAbs(rel) {
		return "", "", ErrPathTraversal
	}

	// Validação de fronteira 2: checagem de symlinks se o alvo existir
	if info, statErr := os.Lstat(targetPath); statErr == nil {
		if info.Mode()&os.ModeSymlink != 0 {
			evalTarget, evalErr := filepath.EvalSymlinks(targetPath)
			if evalErr != nil {
				return "", "", ErrNotFound
			}
			relEval, relEvalErr := filepath.Rel(s.rootDir, evalTarget)
			if relEvalErr != nil || strings.HasPrefix(relEval, "..") || relEval == ".." || filepath.IsAbs(relEval) {
				return "", "", ErrPathTraversal
			}
			targetPath = evalTarget
			rel, _ = filepath.Rel(s.rootDir, targetPath)
		}
	}

	cleanRel := rel
	if cleanRel == "." {
		cleanRel = ""
	}

	return targetPath, filepath.ToSlash(cleanRel), nil
}

// ListDirectory lista os conteúdos de um diretório com metadados e breadcrumbs.
func (s *LocalFileService) ListDirectory(ctx context.Context, relPath string) (*domain.DirectoryListing, error) {
	absPath, cleanRel, err := s.ResolveAndValidatePath(relPath)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		if os.IsPermission(err) {
			return nil, ErrPermissionDenied
		}
		return nil, err
	}

	if !info.IsDir() {
		return nil, ErrNotADirectory
	}

	entries, err := os.ReadDir(absPath)
	if err != nil {
		if os.IsPermission(err) {
			return nil, ErrPermissionDenied
		}
		return nil, err
	}

	breadcrumbs := buildBreadcrumbs(cleanRel)
	var items []domain.FileItem
	var totalDirs, totalFiles int
	var totalSize int64

	for _, entry := range entries {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		entryInfo, err := entry.Info()
		if err != nil {
			continue
		}

		isDir := entry.IsDir()
		// Caso seja um symlink, inspeciona se aponta para diretório seguro
		if entry.Type()&os.ModeSymlink != 0 {
			symPath := filepath.Join(absPath, entry.Name())
			target, evalErr := filepath.EvalSymlinks(symPath)
			if evalErr != nil {
				continue
			}
			relEval, relEvalErr := filepath.Rel(s.rootDir, target)
			if relEvalErr != nil || strings.HasPrefix(relEval, "..") || relEval == ".." {
				// Symlink aponta para fora do sandbox, oculta por segurança
				continue
			}
			if targetStat, statErr := os.Stat(target); statErr == nil {
				isDir = targetStat.IsDir()
				entryInfo = targetStat
			}
		}

		var itemSize int64
		if isDir {
			totalDirs++
		} else {
			totalFiles++
			itemSize = entryInfo.Size()
			totalSize += itemSize
		}

		itemRelPath := filepath.ToSlash(filepath.Join(cleanRel, entry.Name()))

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		category := domain.DetectCategory(entry.Name(), isDir)

		item := domain.FileItem{
			Name:             entry.Name(),
			RelativePath:     itemRelPath,
			IsDir:            isDir,
			Size:             itemSize,
			FormattedSize:    domain.FormatFileSize(itemSize),
			ModTime:          entryInfo.ModTime(),
			FormattedModTime: entryInfo.ModTime().Format("02/01/2006 15:04"),
			Extension:        ext,
			Category:         category,
			IsViewable:       domain.IsViewableFormat(category, ext),
		}

		items = append(items, item)
	}

	// Ordenação: pastas primeiro (alfabético), depois arquivos (alfabético)
	sort.Slice(items, func(i, j int) bool {
		if items[i].IsDir != items[j].IsDir {
			return items[i].IsDir
		}
		return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
	})

	return &domain.DirectoryListing{
		CurrentPath:        cleanRel,
		Breadcrumbs:        breadcrumbs,
		Items:              items,
		TotalItems:         len(items),
		TotalDirs:          totalDirs,
		TotalFiles:         totalFiles,
		TotalSize:          totalSize,
		FormattedTotalSize: domain.FormatFileSize(totalSize),
		IsEmpty:            len(items) == 0,
		CanUpload:          true,
	}, nil
}

// GetFile recupera um arquivo aberto e seu FileInfo para download.
func (s *LocalFileService) GetFile(ctx context.Context, relPath string) (*os.File, os.FileInfo, error) {
	absPath, _, err := s.ResolveAndValidatePath(relPath)
	if err != nil {
		return nil, nil, err
	}

	file, err := os.Open(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, ErrNotFound
		}
		if os.IsPermission(err) {
			return nil, nil, ErrPermissionDenied
		}
		return nil, nil, err
	}

	info, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return nil, nil, err
	}

	if info.IsDir() {
		_ = file.Close()
		return nil, nil, ErrIsDirectory
	}

	return file, info, nil
}

// StreamZip compacta o diretório em streaming diretamente para o io.Writer sem resíduos em disco.
func (s *LocalFileService) StreamZip(ctx context.Context, relPath string, w io.Writer) error {
	absPath, _, err := s.ResolveAndValidatePath(relPath)
	if err != nil {
		return err
	}

	stat, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrNotFound
		}
		if os.IsPermission(err) {
			return ErrPermissionDenied
		}
		return err
	}

	if !stat.IsDir() {
		return ErrNotADirectory
	}

	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	visitedSymlinks := make(map[string]bool)

	err = filepath.WalkDir(absPath, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if path == absPath {
			return nil
		}

		relInArchive, relErr := filepath.Rel(absPath, path)
		if relErr != nil {
			return relErr
		}
		archivePath := filepath.ToSlash(relInArchive)

		info, infoErr := d.Info()
		if infoErr != nil {
			return infoErr
		}

		// Validação estrita de symlink para prevenir loops e escapes
		if d.Type()&os.ModeSymlink != 0 {
			target, evalErr := filepath.EvalSymlinks(path)
			if evalErr != nil {
				return nil
			}
			relEval, relEvalErr := filepath.Rel(s.rootDir, target)
			if relEvalErr != nil || strings.HasPrefix(relEval, "..") || relEval == ".." {
				return nil // Pula symlink externo
			}
			if visitedSymlinks[target] {
				return nil // Previne ciclo
			}
			visitedSymlinks[target] = true
			if targetStat, statErr := os.Stat(target); statErr == nil {
				info = targetStat
			}
		}

		header, headerErr := zip.FileInfoHeader(info)
		if headerErr != nil {
			return headerErr
		}

		if info.IsDir() {
			header.Name = archivePath + "/"
			_, createErr := zipWriter.CreateHeader(header)
			return createErr
		}

		header.Name = archivePath
		header.Method = zip.Deflate

		writer, createErr := zipWriter.CreateHeader(header)
		if createErr != nil {
			return createErr
		}

		file, openErr := os.Open(path)
		if openErr != nil {
			return openErr
		}
		defer file.Close()

		_, copyErr := io.Copy(writer, file)
		return copyErr
	})

	if err != nil {
		return err
	}

	return zipWriter.Flush()
}

// SaveUploadedFiles grava arquivos multipart enviados no diretório de destino validado.
func (s *LocalFileService) SaveUploadedFiles(ctx context.Context, relDir string, files []*multipart.FileHeader) ([]domain.UploadResult, error) {
	absDir, _, err := s.ResolveAndValidatePath(relDir)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(absDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		if os.IsPermission(err) {
			return nil, ErrPermissionDenied
		}
		return nil, err
	}

	if !info.IsDir() {
		return nil, ErrNotADirectory
	}

	var results []domain.UploadResult

	for _, fileHeader := range files {
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		default:
		}

		// Sanitiza nome de arquivo removendo caminhos maliciosos
		cleanFilename := filepath.Base(fileHeader.Filename)
		if cleanFilename == "" || cleanFilename == "." || cleanFilename == ".." {
			results = append(results, domain.UploadResult{
				Filename: fileHeader.Filename,
				Success:  false,
				Message:  "Nome de arquivo inválido",
			})
			continue
		}

		targetPath := filepath.Join(absDir, cleanFilename)

		// Valida se o targetPath continua estritamente dentro da raiz
		rel, relErr := filepath.Rel(s.rootDir, targetPath)
		if relErr != nil || strings.HasPrefix(rel, "..") || rel == ".." {
			results = append(results, domain.UploadResult{
				Filename: cleanFilename,
				Success:  false,
				Message:  "Tentativa de path traversal no nome do arquivo",
			})
			continue
		}

		src, openErr := fileHeader.Open()
		if openErr != nil {
			results = append(results, domain.UploadResult{
				Filename: cleanFilename,
				Success:  false,
				Message:  fmt.Sprintf("Falha ao abrir arquivo temporário: %v", openErr),
			})
			continue
		}

		dst, createErr := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if createErr != nil {
			_ = src.Close()
			results = append(results, domain.UploadResult{
				Filename: cleanFilename,
				Success:  false,
				Message:  fmt.Sprintf("Falha ao salvar no disco: %v", createErr),
			})
			continue
		}

		written, copyErr := io.Copy(dst, src)
		_ = src.Close()
		_ = dst.Close()

		if copyErr != nil {
			results = append(results, domain.UploadResult{
				Filename: cleanFilename,
				Success:  false,
				Message:  fmt.Sprintf("Falha na gravação dos dados: %v", copyErr),
			})
			continue
		}

		results = append(results, domain.UploadResult{
			Filename: cleanFilename,
			Size:     written,
			Success:  true,
		})
	}

	return results, nil
}

// buildBreadcrumbs constrói a trilha de navegação a partir do caminho relativo limpo.
func buildBreadcrumbs(cleanRel string) []domain.Breadcrumb {
	crumbs := []domain.Breadcrumb{
		{
			Name:      "Início",
			Path:      "",
			IsCurrent: cleanRel == "",
		},
	}

	if cleanRel == "" {
		return crumbs
	}

	segments := strings.Split(cleanRel, "/")
	var currentPath string

	for i, segment := range segments {
		if segment == "" {
			continue
		}
		if currentPath == "" {
			currentPath = segment
		} else {
			currentPath = currentPath + "/" + segment
		}

		isLast := i == len(segments)-1
		crumbs = append(crumbs, domain.Breadcrumb{
			Name:      segment,
			Path:      currentPath,
			IsCurrent: isLast,
		})
	}

	return crumbs
}
