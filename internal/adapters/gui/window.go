package gui

import (
	"context"
	"os"
	"os/exec"
	"runtime"
)

var (
	execCommand = exec.Command
	lookPath    = exec.LookPath
)

// IsDesktopEnvironment verifica se o sistema operacional possui interface gráfica disponível.
func IsDesktopEnvironment() bool {
	switch runtime.GOOS {
	case "windows", "darwin":
		return true
	default:
		// Linux / Unix: verifica variáveis de display
		if os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != "" || os.Getenv("XDG_CURRENT_DESKTOP") != "" {
			return true
		}
		return false
	}
}

// OpenURLInBrowser abre uma URL no navegador padrão do sistema operacional.
func OpenURLInBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = execCommand("cmd", "/c", "start", "", url)
	case "darwin":
		cmd = execCommand("open", url)
	default:
		cmd = execCommand("xdg-open", url)
	}

	return cmd.Start()
}

// LaunchDesktopWindow tenta abrir a URL em uma janela dedicada de aplicativo ou navegador desktop.
func LaunchDesktopWindow(ctx context.Context, url string) (*exec.Cmd, error) {
	// Lista de navegadores com suporte a modo janela de aplicativo (--app=)
	candidates := []struct {
		name string
		args []string
	}{
		// Linux & Multiplataforma
		{"google-chrome", []string{"--app=" + url, "--new-window"}},
		{"chromium-browser", []string{"--app=" + url, "--new-window"}},
		{"chromium", []string{"--app=" + url, "--new-window"}},
		{"brave-browser", []string{"--app=" + url, "--new-window"}},
		{"microsoft-edge", []string{"--app=" + url, "--new-window"}},
		{"firefox", []string{"--new-window", url}},

		// Windows
		{"msedge", []string{"--app=" + url, "--new-window"}},
		{"chrome", []string{"--app=" + url, "--new-window"}},

		// macOS
		{"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", []string{"--app=" + url}},
	}

	for _, cand := range candidates {
		if _, err := lookPath(cand.name); err == nil {
			cmd := exec.CommandContext(ctx, cand.name, cand.args...)
			if err := cmd.Start(); err == nil {
				return cmd, nil
			}
		}
	}

	// Fallback padrão para abrir URL diretamente
	var fallbackCmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		fallbackCmd = exec.CommandContext(ctx, "cmd", "/c", "start", "", url)
	case "darwin":
		fallbackCmd = exec.CommandContext(ctx, "open", url)
	default:
		fallbackCmd = exec.CommandContext(ctx, "xdg-open", url)
	}

	if err := fallbackCmd.Start(); err != nil {
		return nil, err
	}
	return fallbackCmd, nil
}
