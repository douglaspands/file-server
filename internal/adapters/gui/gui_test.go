package gui

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	adapterftp "github.com/douglas/file-server/internal/adapters/ftp"
	adaptersftp "github.com/douglas/file-server/internal/adapters/sftp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockRunner para testar o Controller
type MockRunner struct {
	runWebCalls  int
	runFTPCalls  int
	runSFTPCalls int
	webErr       error
	ftpErr       error
	sftpErr      error
}

func (m *MockRunner) RunWeb(ctx context.Context, host string, port int, targetDir string, useTLS bool, certFile, keyFile string) error {
	m.runWebCalls++
	if m.webErr != nil {
		return m.webErr
	}
	<-ctx.Done()
	return ctx.Err()
}

func (m *MockRunner) RunFTP(ctx context.Context, opts adapterftp.ServerOptions) error {
	m.runFTPCalls++
	if m.ftpErr != nil {
		return m.ftpErr
	}
	<-ctx.Done()
	return ctx.Err()
}

func (m *MockRunner) RunSFTP(ctx context.Context, opts adaptersftp.ServerOptions) error {
	m.runSFTPCalls++
	if m.sftpErr != nil {
		return m.sftpErr
	}
	<-ctx.Done()
	return ctx.Err()
}

func TestLogBroadcaster(t *testing.T) {
	b := NewLogBroadcaster(3)
	assert.NotNil(t, b)

	bDefault := NewLogBroadcaster(0)
	assert.Equal(t, 500, bDefault.capacity)

	// Test broadcast and capacity limitation
	b.Broadcast("msg 1")
	b.Broadcast("msg 2")
	b.Broadcast("msg 3")
	b.Broadcast("msg 4")

	history := b.GetHistory()
	assert.Len(t, history, 3)
	assert.Contains(t, history[0], "msg 2")
	assert.Contains(t, history[2], "msg 4")

	// Test Subscribe & Unsubscribe
	subHistory, ch := b.Subscribe()
	assert.Len(t, subHistory, 3)

	b.Broadcast("msg 5")
	select {
	case msg := <-ch:
		assert.Contains(t, msg, "msg 5")
	case <-time.After(1 * time.Second):
		t.Fatal("timeout esperando mensagem no canal")
	}

	b.Unsubscribe(ch)

	// Test io.Writer Write
	n, err := b.Write([]byte("linha de log 1\nlinha de log 2\n"))
	assert.NoError(t, err)
	assert.Equal(t, 30, n)

	// Test Clear
	b.Clear()
	assert.Empty(t, b.GetHistory())
}

func TestNetworkCategorizationAndUrls(t *testing.T) {
	tests := []struct {
		name      string
		ip        net.IP
		wantType  InterfaceType
		wantLabel string
	}{
		{"lo", net.ParseIP("127.0.0.1"), TypeLoopback, "Loopback"},
		{"wlan0", net.ParseIP("192.168.1.10"), TypeWiFi, "Wi-Fi"},
		{"eth0", net.ParseIP("10.0.0.5"), TypeEthernet, "Ethernet"},
		{"tailscale0", net.ParseIP("100.64.0.1"), TypeVPN, "VPN"},
		{"docker0", net.ParseIP("172.17.0.1"), TypeDocker, "Docker / Bridge"},
		{"custom_net", net.ParseIP("192.168.100.1"), TypeOther, "Rede Local"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iType, label := CategorizeInterface(tt.name, tt.ip)
			assert.Equal(t, tt.wantType, iType)
			assert.Equal(t, tt.wantLabel, label)
		})
	}

	// Test BuildAccessURL
	assert.Equal(t, "http://127.0.0.1:8080", BuildAccessURL("127.0.0.1", 8080, "web", false))
	assert.Equal(t, "https://127.0.0.1:8443", BuildAccessURL("127.0.0.1", 8443, "web", true))
	assert.Equal(t, "ftp://127.0.0.1:2121", BuildAccessURL("127.0.0.1", 2121, "ftp", false))
	assert.Equal(t, "ftps://127.0.0.1:2121", BuildAccessURL("127.0.0.1", 2121, "ftp", true))
	assert.Equal(t, "sftp://127.0.0.1:2222", BuildAccessURL("127.0.0.1", 2222, "sftp", false))
	assert.Equal(t, "http://[::1]:8080", BuildAccessURL("::1", 8080, "web", false))

	// Test BuildAccessURLs
	localURL, lanURLs := BuildAccessURLs("0.0.0.0", 8080, "web", false)
	assert.Equal(t, "http://127.0.0.1:8080", localURL)
	assert.NotNil(t, lanURLs)

	localCustom, _ := BuildAccessURLs("192.168.1.100", 8080, "web", false)
	assert.Equal(t, "http://192.168.1.100:8080", localCustom)

	// Test DetectNetworkInterfaces
	ifaces := DetectNetworkInterfaces(8080, "web", false)
	require.NotEmpty(t, ifaces)
	assert.Equal(t, "lo (Localhost)", ifaces[0].Name)
	assert.True(t, ifaces[0].IsLoopback)
	assert.True(t, ifaces[0].IsRecommended)

	// Test FormatShareMessage
	cfg := ServerConfig{
		Protocol:  ProtocolWeb,
		TargetDir: "/tmp/share",
		Port:      8080,
		UseTLS:    true,
	}
	msg := FormatShareMessage(cfg, ifaces)
	assert.Contains(t, msg, "Web Segura (HTTPS)")
	assert.Contains(t, msg, "/tmp/share")
	assert.Contains(t, msg, "Links de Acesso:")

	// Test FTP share message with credentials
	ftpCfg := ServerConfig{
		Protocol:  ProtocolFTP,
		TargetDir: "/tmp/share",
		Port:      2121,
		User:      "usuario_ftp",
		Pass:      "senha123",
		ReadOnly:  true,
	}
	ftpMsg := FormatShareMessage(ftpCfg, ifaces)
	assert.Contains(t, ftpMsg, "usuario_ftp")
	assert.Contains(t, ftpMsg, "senha123")
	assert.Contains(t, ftpMsg, "Somente Leitura")

	sftpCfg := ServerConfig{
		Protocol:  ProtocolSFTP,
		TargetDir: "/tmp/share",
		Port:      2222,
		User:      "sftp_user",
	}
	sftpMsg := FormatShareMessage(sftpCfg, ifaces)
	assert.Contains(t, sftpMsg, "SFTP (SSH)")
}

func TestControllerStateAndLifecycle(t *testing.T) {
	tempDir := t.TempDir()
	mock := &MockRunner{}
	broadcaster := NewLogBroadcaster(100)

	ctrl := NewController(tempDir, mock, broadcaster)
	require.NotNil(t, ctrl)
	assert.Equal(t, broadcaster, ctrl.GetBroadcaster())

	// Default controller with nil runner and broadcaster
	ctrlNil := NewController("", nil, nil)
	assert.NotNil(t, ctrlNil)
	assert.NotNil(t, ctrlNil.GetBroadcaster())

	// Initial status
	status := ctrl.GetStatus()
	assert.False(t, status.IsRunning)
	assert.Equal(t, tempDir, status.TargetDir)
	assert.Equal(t, ProtocolWeb, status.Protocol)
	assert.Equal(t, tempDir, ctrl.GetConfig().TargetDir)

	// Set and Get GUI addr
	ctrl.SetGUIAddr("127.0.0.1", 9999)
	host, port := ctrl.GetGUIAddr()
	assert.Equal(t, "127.0.0.1", host)
	assert.Equal(t, 9999, port)

	// Start Web Server
	err := ctrl.StartServer(ServerConfig{
		Protocol:  ProtocolWeb,
		TargetDir: tempDir,
		Port:      8080,
	})
	assert.NoError(t, err)
	assert.True(t, ctrl.IsRunning())

	// Cannot start again while running
	err = ctrl.StartServer(ServerConfig{Protocol: ProtocolWeb, TargetDir: tempDir})
	assert.ErrorIs(t, err, ErrServerAlreadyRunning)

	// Stop Web Server
	err = ctrl.StopServer()
	assert.NoError(t, err)
	assert.False(t, ctrl.IsRunning())

	// Cannot stop when not running
	err = ctrl.StopServer()
	assert.ErrorIs(t, err, ErrServerNotRunning)

	// Start FTP Server with default user and generated pass
	err = ctrl.StartServer(ServerConfig{
		Protocol:  ProtocolFTP,
		TargetDir: tempDir,
		Port:      0,
	})
	assert.NoError(t, err)
	assert.True(t, ctrl.IsRunning())
	_ = ctrl.StopServer()

	// Start SFTP Server with default user and generated pass
	err = ctrl.StartServer(ServerConfig{
		Protocol:  ProtocolSFTP,
		TargetDir: tempDir,
		Port:      0,
	})
	assert.NoError(t, err)
	assert.True(t, ctrl.IsRunning())
	_ = ctrl.StopServer()

	// Start with invalid directory
	err = ctrl.StartServer(ServerConfig{
		Protocol:  ProtocolWeb,
		TargetDir: filepath.Join(tempDir, "nao_existe_12345"),
	})
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidTargetDir))
}

// MockFolderPicker para testes
type MockFolderPicker struct {
	chosenDir string
	err       error
}

func (m *MockFolderPicker) PickFolder(ctx context.Context, initialDir string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.chosenDir, nil
}

func TestServerEndpoints(t *testing.T) {
	tempDir := t.TempDir()
	mock := &MockRunner{}
	broadcaster := NewLogBroadcaster(100)
	ctrl := NewController(tempDir, mock, broadcaster)
	picker := &MockFolderPicker{chosenDir: tempDir}

	server, err := NewServer(ctrl, picker)
	require.NoError(t, err)

	mux := http.NewServeMux()
	err = server.RegisterRoutes(mux)
	require.NoError(t, err)

	// GET /
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "File Server")

	// GET /not-found
	req = httptest.NewRequest(http.MethodGet, "/invalid", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// GET /api/status
	req = httptest.NewRequest(http.MethodGet, "/api/status", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var st ServerStatus
	err = json.NewDecoder(w.Body).Decode(&st)
	assert.NoError(t, err)
	assert.False(t, st.IsRunning)

	// GET /api/interfaces
	req = httptest.NewRequest(http.MethodGet, "/api/interfaces?port=8080&protocol=web&tls=true", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var ifaces []NetworkInterface
	err = json.NewDecoder(w.Body).Decode(&ifaces)
	assert.NoError(t, err)
	assert.NotEmpty(t, ifaces)

	// GET /api/share-message
	req = httptest.NewRequest(http.MethodGet, "/api/share-message", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var shareResp map[string]interface{}
	err = json.NewDecoder(w.Body).Decode(&shareResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, shareResp["message"])

	// POST /api/picker/folder (success)
	reqBody := `{"currentDir":"` + tempDir + `"}`
	req = httptest.NewRequest(http.MethodPost, "/api/picker/folder", strings.NewReader(reqBody))
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var pickerResp map[string]interface{}
	err = json.NewDecoder(w.Body).Decode(&pickerResp)
	assert.NoError(t, err)
	assert.True(t, pickerResp["success"].(bool))
	assert.Equal(t, tempDir, pickerResp["path"])

	// POST /api/picker/folder (picker error)
	errPicker := &MockFolderPicker{err: errors.New("falha")}
	serverErrPicker, _ := NewServer(ctrl, errPicker)
	muxErr := http.NewServeMux()
	_ = serverErrPicker.RegisterRoutes(muxErr)
	req = httptest.NewRequest(http.MethodPost, "/api/picker/folder", strings.NewReader(reqBody))
	w = httptest.NewRecorder()
	muxErr.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// POST /api/picker/folder (wrong method)
	req = httptest.NewRequest(http.MethodGet, "/api/picker/folder", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)

	// POST /api/server/start (valid)
	startBody := fmt.Sprintf(`{"protocol":"web","port":8080,"targetDir":%q}`, tempDir)
	req = httptest.NewRequest(http.MethodPost, "/api/server/start", strings.NewReader(startBody))
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// POST /api/server/start (invalid json)
	req = httptest.NewRequest(http.MethodPost, "/api/server/start", strings.NewReader("invalid-json"))
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// POST /api/server/start (already running error)
	req = httptest.NewRequest(http.MethodPost, "/api/server/start", strings.NewReader(startBody))
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// POST /api/server/start (wrong method)
	req = httptest.NewRequest(http.MethodGet, "/api/server/start", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)

	// POST /api/server/stop (valid)
	req = httptest.NewRequest(http.MethodPost, "/api/server/stop", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// POST /api/server/stop (already stopped error)
	req = httptest.NewRequest(http.MethodPost, "/api/server/stop", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// POST /api/server/stop (wrong method)
	req = httptest.NewRequest(http.MethodGet, "/api/server/stop", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)

	// POST /api/app/open-browser
	openBody := `{"url":"http://127.0.0.1:8080"}`
	req = httptest.NewRequest(http.MethodPost, "/api/app/open-browser", strings.NewReader(openBody))
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// POST /api/app/open-browser (wrong method)
	req = httptest.NewRequest(http.MethodGet, "/api/app/open-browser", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)

	// Test Server Start & Stop
	guiURL, err := server.Start("127.0.0.1", 0)
	assert.NoError(t, err)
	assert.Contains(t, guiURL, "http://127.0.0.1:")
	err = server.Stop()
	assert.NoError(t, err)
}

func TestNativeFolderPicker(t *testing.T) {
	tempDir := t.TempDir()

	_ = NewNativeFolderPicker()

	// Mock picker com zenity
	picker := &NativeFolderPicker{
		lookPath: func(file string) (string, error) {
			if file == "zenity" {
				return "/usr/bin/zenity", nil
			}
			return "", exec.ErrNotFound
		},
		execCmd: func(ctx context.Context, name string, args ...string) ([]byte, error) {
			if name == "zenity" {
				return []byte(tempDir + "\n"), nil
			}
			return nil, errors.New("command failed")
		},
	}

	res, err := picker.PickFolder(context.Background(), tempDir)
	assert.NoError(t, err)
	assert.Equal(t, tempDir, res)

	// Teste com kdialog
	kdialogPicker := &NativeFolderPicker{
		lookPath: func(file string) (string, error) {
			if file == "kdialog" {
				return "/usr/bin/kdialog", nil
			}
			return "", exec.ErrNotFound
		},
		execCmd: func(ctx context.Context, name string, args ...string) ([]byte, error) {
			if name == "kdialog" {
				return []byte(tempDir + "\n"), nil
			}
			return nil, errors.New("command failed")
		},
	}
	resK, err := kdialogPicker.pickLinux(context.Background(), tempDir)
	assert.NoError(t, err)
	assert.Equal(t, tempDir, resK)

	// Teste com qarma
	qarmaPicker := &NativeFolderPicker{
		lookPath: func(file string) (string, error) {
			if file == "qarma" {
				return "/usr/bin/qarma", nil
			}
			return "", exec.ErrNotFound
		},
		execCmd: func(ctx context.Context, name string, args ...string) ([]byte, error) {
			if name == "qarma" {
				return []byte(tempDir + "\n"), nil
			}
			return nil, errors.New("command failed")
		},
	}
	resQ, err := qarmaPicker.pickLinux(context.Background(), tempDir)
	assert.NoError(t, err)
	assert.Equal(t, tempDir, resQ)

	// Teste pickWindows
	winPicker := &NativeFolderPicker{
		execCmd: func(ctx context.Context, name string, args ...string) ([]byte, error) {
			return []byte(tempDir + "\n"), nil
		},
	}
	resW, err := winPicker.pickWindows(context.Background(), tempDir)
	assert.NoError(t, err)
	assert.Equal(t, tempDir, resW)

	// Teste pickDarwin
	macPicker := &NativeFolderPicker{
		execCmd: func(ctx context.Context, name string, args ...string) ([]byte, error) {
			return []byte(tempDir + "\n"), nil
		},
	}
	resM, err := macPicker.pickDarwin(context.Background(), tempDir)
	assert.NoError(t, err)
	assert.Equal(t, tempDir, resM)

	// Teste de cancelamento / retorno vazio
	cancelPicker := &NativeFolderPicker{
		lookPath: func(file string) (string, error) {
			return "/usr/bin/zenity", nil
		},
		execCmd: func(ctx context.Context, name string, args ...string) ([]byte, error) {
			return []byte(""), nil
		},
	}
	_, err = cancelPicker.PickFolder(context.Background(), tempDir)
	assert.ErrorIs(t, err, ErrFolderPickerCancelled)

	// Teste de ferramenta não encontrada
	noToolPicker := &NativeFolderPicker{
		lookPath: func(file string) (string, error) {
			return "", exec.ErrNotFound
		},
		execCmd: func(ctx context.Context, name string, args ...string) ([]byte, error) {
			return nil, exec.ErrNotFound
		},
	}
	_, err = noToolPicker.PickFolder(context.Background(), tempDir)
	assert.Error(t, err)
}

func TestWindowAndDesktopEnvironment(t *testing.T) {
	_ = IsDesktopEnvironment()

	os.Setenv("DISPLAY", ":0")
	assert.True(t, IsDesktopEnvironment())

	// Mock execCommand e lookPath para testar LaunchDesktopWindow
	oldLookPath := lookPath
	oldExecCommand := execCommand
	defer func() {
		lookPath = oldLookPath
		execCommand = oldExecCommand
	}()

	lookPath = func(file string) (string, error) {
		if file == "google-chrome" {
			return "/usr/bin/google-chrome", nil
		}
		return "", exec.ErrNotFound
	}

	execCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("echo", "opened")
	}

	cmd, err := LaunchDesktopWindow(context.Background(), "http://127.0.0.1:8080")
	if err == nil && cmd != nil {
		_ = cmd.Wait()
	}

	_ = OpenURLInBrowser("http://127.0.0.1:8080")
}

func TestDefaultRunner(t *testing.T) {
	runner := NewDefaultRunner()
	assert.NotNil(t, runner)

	tempDir := t.TempDir()
	ctx, cancel := context.WithCancel(context.Background())

	// Inicia Web Server com context cancelado imediatamente
	cancel()
	err := runner.RunWeb(ctx, "127.0.0.1", 0, tempDir, false, "", "")
	assert.True(t, err == nil || errors.Is(err, context.Canceled) || errors.Is(err, http.ErrServerClosed))

	// Inicia FTP com context cancelado
	ctxFTP, cancelFTP := context.WithCancel(context.Background())
	cancelFTP()
	_ = runner.RunFTP(ctxFTP, adapterftp.ServerOptions{
		Host:      "127.0.0.1",
		Port:      0,
		TargetDir: tempDir,
	})

	// Inicia SFTP com context cancelado
	ctxSFTP, cancelSFTP := context.WithCancel(context.Background())
	cancelSFTP()
	_ = runner.RunSFTP(ctxSFTP, adaptersftp.ServerOptions{
		Host:      "127.0.0.1",
		Port:      0,
		TargetDir: tempDir,
	})
}
