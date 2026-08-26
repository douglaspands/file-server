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
	case osWindows, osDarwin:
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
	case osWindows:
		cmd = execCommand("cmd", "/c", "start", "", url)
	case osDarwin:
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
		// Linux & Multiplataforma (Chromium-based com forçamento de tema escuro na moldura)
		{"google-chrome", []string{"--app=" + url, "--new-window", "--force-dark-mode", "--enable-features=WebUIDarkMode"}},
		{"chromium-browser", []string{"--app=" + url, "--new-window", "--force-dark-mode", "--enable-features=WebUIDarkMode"}},
		{"chromium", []string{"--app=" + url, "--new-window", "--force-dark-mode", "--enable-features=WebUIDarkMode"}},
		{"brave-browser", []string{"--app=" + url, "--new-window", "--force-dark-mode", "--enable-features=WebUIDarkMode"}},
		{"microsoft-edge", []string{"--app=" + url, "--new-window", "--force-dark-mode", "--enable-features=WebUIDarkMode"}},
		{"firefox", []string{"--new-window", url}},

		// Windows
		{"msedge", []string{"--app=" + url, "--new-window", "--force-dark-mode", "--enable-features=WebUIDarkMode"}},
		{"chrome", []string{"--app=" + url, "--new-window", "--force-dark-mode", "--enable-features=WebUIDarkMode"}},

		// macOS
		{"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", []string{"--app=" + url, "--force-dark-mode", "--enable-features=WebUIDarkMode"}},
	}

	for _, cand := range candidates {
		if _, err := lookPath(cand.name); err == nil {
			cmd := exec.CommandContext(ctx, cand.name, cand.args...) // #nosec G204
			if err := cmd.Start(); err == nil {
				return cmd, nil
			}
		}
	}

	// Fallback padrão para abrir URL diretamente
	var fallbackCmd *exec.Cmd
	switch runtime.GOOS {
	case osWindows:
		fallbackCmd = exec.CommandContext(ctx, "cmd", "/c", "start", "", url) // #nosec G204
	case osDarwin:
		fallbackCmd = exec.CommandContext(ctx, "open", url) // #nosec G204
	default:
		fallbackCmd = exec.CommandContext(ctx, "xdg-open", url) // #nosec G204
	}

	if err := fallbackCmd.Start(); err != nil {
		return nil, err
	}
	return fallbackCmd, nil
}
