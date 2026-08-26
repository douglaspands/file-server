package gui

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	ErrFolderPickerCancelled = errors.New("seleção de pasta cancelada ou fechada")
	ErrFolderPickerNotFound  = errors.New("nenhuma ferramenta nativa de seleção de pastas disponível no sistema")
)

// FolderPickerInterface define a interface para seleção de pastas nativas.
type FolderPickerInterface interface {
	PickFolder(ctx context.Context, initialDir string) (string, error)
}

// NativeFolderPicker implementa o diálogo nativo do sistema operacional.
type NativeFolderPicker struct {
	lookPath func(file string) (string, error)
	execCmd  func(ctx context.Context, name string, args ...string) ([]byte, error)
}

// NewNativeFolderPicker cria uma instância do selecionador nativo de diretórios.
func NewNativeFolderPicker() *NativeFolderPicker {
	return &NativeFolderPicker{
		lookPath: exec.LookPath,
		execCmd: func(ctx context.Context, name string, args ...string) ([]byte, error) {
			cmd := exec.CommandContext(ctx, name, args...)
			return cmd.Output()
		},
	}
}

// PickFolder abre a janela nativa de seleção de pastas do sistema operacional e retorna o caminho absoluto escolhido.
func (p *NativeFolderPicker) PickFolder(ctx context.Context, initialDir string) (string, error) {
	if initialDir == "" {
		cwd, err := os.Getwd()
		if err == nil {
			initialDir = cwd
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	var selected string
	var err error

	switch runtime.GOOS {
	case osWindows:
		selected, err = p.pickWindows(ctx, initialDir)
	case osDarwin:
		selected, err = p.pickDarwin(ctx, initialDir)
	default: // Linux, FreeBSD, OpenBSD, etc.
		selected, err = p.pickLinux(ctx, initialDir)
	}

	if err != nil {
		return "", err
	}

	selected = strings.TrimSpace(selected)
	if selected == "" {
		return "", ErrFolderPickerCancelled
	}

	absPath, err := filepath.Abs(selected)
	if err != nil {
		return selected, nil
	}

	info, err := os.Stat(absPath)
	if err != nil || !info.IsDir() {
		return "", errors.New("o caminho selecionado não é um diretório válido")
	}

	return absPath, nil
}

func (p *NativeFolderPicker) pickLinux(ctx context.Context, initialDir string) (string, error) {
	// 1. Tenta Zenity (GNOME / Padrão Linux)
	if _, err := p.lookPath("zenity"); err == nil {
		args := []string{"--file-selection", "--directory", "--title=Selecione a pasta para compartilhar"}
		if initialDir != "" {
			args = append(args, "--filename="+initialDir+"/")
		}
		out, err := p.execCmd(ctx, "zenity", args...)
		if err != nil || len(bytes.TrimSpace(out)) == 0 {
			return "", ErrFolderPickerCancelled
		}
		return string(bytes.TrimSpace(out)), nil
	}

	// 2. Tenta KDialog (KDE Plasma)
	if _, err := p.lookPath("kdialog"); err == nil {
		args := []string{"--getexistingdirectory"}
		if initialDir != "" {
			args = append(args, initialDir)
		}
		out, err := p.execCmd(ctx, "kdialog", args...)
		if err != nil || len(bytes.TrimSpace(out)) == 0 {
			return "", ErrFolderPickerCancelled
		}
		return string(bytes.TrimSpace(out)), nil
	}

	// 3. Tenta Qarma
	if _, err := p.lookPath("qarma"); err == nil {
		args := []string{"--file-selection", "--directory", "--title=Selecione a pasta para compartilhar"}
		out, err := p.execCmd(ctx, "qarma", args...)
		if err != nil || len(bytes.TrimSpace(out)) == 0 {
			return "", ErrFolderPickerCancelled
		}
		return string(bytes.TrimSpace(out)), nil
	}

	return "", ErrFolderPickerNotFound
}

func (p *NativeFolderPicker) pickWindows(ctx context.Context, initialDir string) (string, error) {
	script := `
Add-Type -AssemblyName System.Windows.Forms
$f = New-Object System.Windows.Forms.FolderBrowserDialog
$f.Description = 'Selecione a pasta para compartilhar no File Server'
$f.ShowNewFolderButton = $true
if ($f.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) {
    [Console]::Out.Write($f.SelectedPath)
}
`
	out, err := p.execCmd(ctx, "powershell.exe", "-NoProfile", "-NonInteractive", "-Command", script)
	if err != nil || len(bytes.TrimSpace(out)) == 0 {
		return "", ErrFolderPickerCancelled
	}
	return string(bytes.TrimSpace(out)), nil
}

func (p *NativeFolderPicker) pickDarwin(ctx context.Context, initialDir string) (string, error) {
	script := `POSIX path of (choose folder with prompt "Selecione o diretório para compartilhar:")`
	out, err := p.execCmd(ctx, "osascript", "-e", script)
	if err != nil || len(bytes.TrimSpace(out)) == 0 {
		return "", ErrFolderPickerCancelled
	}
	return string(bytes.TrimSpace(out)), nil
}
