package sftp

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
)

var (
	// ErrReadOnly indica tentativa de gravação quando o servidor está em modo somente leitura.
	ErrReadOnly = sftp.ErrSSHFxPermissionDenied

	// ErrPathTraversal indica tentativa de acessar arquivos fora do diretório raiz.
	ErrPathTraversal = sftp.ErrSSHFxPermissionDenied

	// ErrNotFound indica arquivo ou pasta não encontrada.
	ErrNotFound = sftp.ErrSSHFxNoSuchFile
)

// FSHandler gerencia as operações do sistema de arquivos para o servidor SFTP com sandbox e controle de permissões.
type FSHandler struct {
	rootDir  string
	readOnly bool
}

// NewFSHandler inicializa um novo FSHandler para o diretório raiz informado.
func NewFSHandler(rootDir string, readOnly bool) (*FSHandler, error) {
	if rootDir == "" {
		rootDir = "."
	}

	absRoot, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, fmt.Errorf("diretório raiz inválido: %w", err)
	}

	info, err := os.Stat(absRoot)
	if err != nil {
		return nil, fmt.Errorf("erro ao acessar diretório raiz '%s': %w", absRoot, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("o caminho raiz '%s' não é um diretório", absRoot)
	}

	canonicalRoot, err := filepath.EvalSymlinks(absRoot)
	if err != nil {
		canonicalRoot = absRoot
	}

	return &FSHandler{
		rootDir:  canonicalRoot,
		readOnly: readOnly,
	}, nil
}

// ResolvePath sanitiza o caminho virtual vindo do cliente SFTP e garante confinamento estrito no diretório raiz.
func (h *FSHandler) ResolvePath(virtualPath string) (string, error) {
	if strings.Contains(virtualPath, "\x00") {
		return "", ErrPathTraversal
	}

	normalized := filepath.ToSlash(virtualPath)

	// Se o caminho for absoluto no SO hospedeiro e não estiver dentro do diretório raiz
	if filepath.IsAbs(virtualPath) && !strings.HasPrefix(filepath.Clean(virtualPath), h.rootDir) {
		// Se for um caminho virtual SFTP (começando com /), mapeia relativo à raiz
		cleanVirtual := strings.TrimPrefix(filepath.Clean("/"+normalized), "/")
		normalized = cleanVirtual
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
				return "", ErrPathTraversal
			}
		} else {
			depth++
		}
	}

	targetPath := filepath.Clean(filepath.Join(h.rootDir, filepath.FromSlash(normalized)))

	rel, err := filepath.Rel(h.rootDir, targetPath)
	if err != nil || strings.HasPrefix(rel, "..") || rel == ".." || (filepath.IsAbs(rel) && !strings.HasPrefix(rel, h.rootDir)) {
		return "", ErrPathTraversal
	}

	// Se o caminho já existe e for um symlink, valida se o destino escapa do sandbox
	if targetInfo, err := os.Lstat(targetPath); err == nil {
		if targetInfo.Mode()&os.ModeSymlink != 0 {
			resolved, err := filepath.EvalSymlinks(targetPath)
			if err == nil {
				relEval, relEvalErr := filepath.Rel(h.rootDir, resolved)
				if relEvalErr != nil || strings.HasPrefix(relEval, "..") || relEval == ".." {
					return "", ErrPathTraversal
				}
			}
		}
	}

	return targetPath, nil
}

// Fileread implementa sftp.FileReader para leitura e download de arquivos.
func (h *FSHandler) Fileread(r *sftp.Request) (io.ReaderAt, error) {
	targetPath, err := h.ResolvePath(r.Filepath)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		return nil, sftp.ErrSSHFxFailure
	}

	if info.IsDir() {
		return nil, sftp.ErrSSHFxFailure
	}

	f, err := os.Open(targetPath)
	if err != nil {
		return nil, sftp.ErrSSHFxPermissionDenied
	}

	return f, nil
}

// Filewrite implementa sftp.FileWriter para gravação e upload de arquivos.
func (h *FSHandler) Filewrite(r *sftp.Request) (io.WriterAt, error) {
	if h.readOnly {
		return nil, ErrReadOnly
	}

	targetPath, err := h.ResolvePath(r.Filepath)
	if err != nil {
		return nil, err
	}

	// Valida se o diretório pai está dentro do sandbox
	parentDir := filepath.Dir(targetPath)
	if parentResolved, err := filepath.EvalSymlinks(parentDir); err == nil {
		relParent, err := filepath.Rel(h.rootDir, parentResolved)
		if err != nil || strings.HasPrefix(relParent, "..") || relParent == ".." {
			return nil, ErrPathTraversal
		}
	}

	var openFlags int
	pflags := r.Pflags()
	if pflags.Write {
		openFlags |= os.O_WRONLY
	}
	if pflags.Append {
		openFlags |= os.O_APPEND
	}
	if pflags.Creat {
		openFlags |= os.O_CREATE
	}
	if pflags.Trunc {
		openFlags |= os.O_TRUNC
	}
	if pflags.Excl {
		openFlags |= os.O_EXCL
	}
	if openFlags == 0 {
		openFlags = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	}

	f, err := os.OpenFile(targetPath, openFlags, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		return nil, sftp.ErrSSHFxPermissionDenied
	}

	return f, nil
}

// Filecmd implementa sftp.FileCmder para criação, remoção, renomeação e modificação de atributos.
func (h *FSHandler) Filecmd(r *sftp.Request) error {
	targetPath, err := h.ResolvePath(r.Filepath)
	if err != nil {
		return err
	}

	switch r.Method {
	case "Setstat":
		if h.readOnly {
			return ErrReadOnly
		}
		return nil

	case "Rename":
		if h.readOnly {
			return ErrReadOnly
		}
		dstPath, err := h.ResolvePath(r.Target)
		if err != nil {
			return err
		}
		if err := os.Rename(targetPath, dstPath); err != nil {
			if os.IsNotExist(err) {
				return ErrNotFound
			}
			return sftp.ErrSSHFxFailure
		}
		return nil

	case "Rmdir":
		if h.readOnly {
			return ErrReadOnly
		}
		info, err := os.Stat(targetPath)
		if err != nil {
			if os.IsNotExist(err) {
				return ErrNotFound
			}
			return sftp.ErrSSHFxFailure
		}
		if !info.IsDir() {
			return sftp.ErrSSHFxFailure
		}
		if err := os.Remove(targetPath); err != nil {
			return sftp.ErrSSHFxFailure
		}
		return nil

	case "Mkdir":
		if h.readOnly {
			return ErrReadOnly
		}
		if err := os.Mkdir(targetPath, 0755); err != nil {
			return sftp.ErrSSHFxFailure
		}
		return nil

	case "Remove":
		if h.readOnly {
			return ErrReadOnly
		}
		if err := os.Remove(targetPath); err != nil {
			if os.IsNotExist(err) {
				return ErrNotFound
			}
			return sftp.ErrSSHFxFailure
		}
		return nil

	case "Symlink":
		if h.readOnly {
			return ErrReadOnly
		}
		linkTarget, err := h.ResolvePath(r.Target)
		if err != nil {
			return err
		}
		if err := os.Symlink(linkTarget, targetPath); err != nil {
			return sftp.ErrSSHFxFailure
		}
		return nil

	default:
		return sftp.ErrSSHFxOpUnsupported
	}
}

// Filelist implementa sftp.FileLister para listagem de pastas e consulta de atributos (Stat/Readlink).
func (h *FSHandler) Filelist(r *sftp.Request) (sftp.ListerAt, error) {
	targetPath, err := h.ResolvePath(r.Filepath)
	if err != nil {
		return nil, err
	}

	switch r.Method {
	case "List":
		entries, err := os.ReadDir(targetPath)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, ErrNotFound
			}
			return nil, sftp.ErrSSHFxPermissionDenied
		}

		var fileInfos []os.FileInfo
		for _, entry := range entries {
			info, err := entry.Info()
			if err == nil {
				fileInfos = append(fileInfos, info)
			}
		}
		return listerAt(fileInfos), nil

	case "Stat":
		info, err := os.Stat(targetPath)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, ErrNotFound
			}
			return nil, sftp.ErrSSHFxPermissionDenied
		}
		return listerAt([]os.FileInfo{info}), nil

	case "Readlink":
		_, err := os.Readlink(targetPath)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, ErrNotFound
			}
			return nil, sftp.ErrSSHFxFailure
		}
		info, err := os.Stat(targetPath)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, ErrNotFound
			}
			return nil, sftp.ErrSSHFxFailure
		}
		return listerAt([]os.FileInfo{info}), nil

	default:
		return nil, sftp.ErrSSHFxOpUnsupported
	}
}

// ToHandlers converte o FSHandler para a estrutura sftp.Handlers pronta para uso pelo RequestServer.
func (h *FSHandler) ToHandlers() sftp.Handlers {
	return sftp.Handlers{
		FileGet:  h,
		FilePut:  h,
		FileCmd:  h,
		FileList: h,
	}
}

// listerAt implementa sftp.ListerAt para fatiamento de lista de arquivos.
type listerAt []os.FileInfo

func (l listerAt) ListAt(f []os.FileInfo, offset int64) (int, error) {
	if offset >= int64(len(l)) {
		return 0, io.EOF
	}
	n := copy(f, l[offset:])
	if n < len(f) {
		return n, io.EOF
	}
	return n, nil
}
