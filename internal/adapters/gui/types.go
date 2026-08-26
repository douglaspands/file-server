package gui

import (
	"context"
	"sync"
	"time"

	adapterftp "github.com/douglas/file-server/internal/adapters/ftp"
	adaptersftp "github.com/douglas/file-server/internal/adapters/sftp"
	"github.com/douglas/file-server/internal/core/domain"
)

// Protocol representa o protocolo de serviço selecionado no launcher.
type Protocol string

const (
	ProtocolWeb  Protocol = "web"
	ProtocolFTP  Protocol = "ftp"
	ProtocolSFTP Protocol = "sftp"
)

// InterfaceType categoriza o tipo de adaptador de rede detectado.
type InterfaceType string

const (
	TypeWiFi     InterfaceType = "wifi"
	TypeEthernet InterfaceType = "ethernet"
	TypeVPN      InterfaceType = "vpn"
	TypeDocker   InterfaceType = "docker"
	TypeLoopback InterfaceType = "loopback"
	TypeOther    InterfaceType = "other"
)

// NetworkInterface detalha uma interface de rede e seus acessos.
type NetworkInterface struct {
	Name          string        `json:"name"`
	IP            string        `json:"ip"`
	Type          InterfaceType `json:"type"`
	TypeLabel     string        `json:"typeLabel"`
	IsRecommended bool          `json:"isRecommended"`
	IsLoopback    bool          `json:"isLoopback"`
	URL           string        `json:"url"`
}

// ServerConfig encapsula todas as configurações unificadas da GUI para inicialização.
type ServerConfig struct {
	Protocol     Protocol `json:"protocol"`
	Host         string   `json:"host"`
	Port         int      `json:"port"`
	TargetDir    string   `json:"targetDir"`
	UseTLS       bool     `json:"useTLS"`
	TLSCert      string   `json:"tlsCert"`
	TLSKey       string   `json:"tlsKey"`
	User         string   `json:"user"`
	Pass         string   `json:"pass"`
	AuthKey      string   `json:"authKey"`
	HostKey      string   `json:"hostKey"`
	PassivePorts string   `json:"passivePorts"`
	ReadOnly     bool     `json:"readOnly"`
}

// ServerStatus detalha o estado atual de execução do servidor gerenciado.
type ServerStatus struct {
	IsRunning  bool                 `json:"isRunning"`
	Protocol   Protocol             `json:"protocol"`
	Host       string               `json:"host"`
	Port       int                  `json:"port"`
	TargetDir  string               `json:"targetDir"`
	UseTLS     bool                 `json:"useTLS"`
	LocalURL   string               `json:"localUrl"`
	LanURLs    []string             `json:"lanUrls"`
	Interfaces []NetworkInterface   `json:"interfaces"`
	User       string               `json:"user,omitempty"`
	ReadOnly   bool                 `json:"readOnly"`
	StartedAt  *time.Time           `json:"startedAt,omitempty"`
	Health     *domain.HealthStatus `json:"health,omitempty"`
}

// GUIOptions define as opções de configuração para o próprio servidor da GUI.
type GUIOptions struct {
	Host       string
	Port       int
	InitialDir string
	AutoOpen   bool
}

// ServerRunner define a interface para instanciar e executar os servidores em background.
type ServerRunner interface {
	RunWeb(ctx context.Context, host string, port int, targetDir string, useTLS bool, certFile, keyFile string) error
	RunFTP(ctx context.Context, opts adapterftp.ServerOptions) error
	RunSFTP(ctx context.Context, opts adaptersftp.ServerOptions) error
}

// Controller gerencia o estado dos servidores ativos, logs e diálogos do sistema.
type Controller struct {
	mu           sync.RWMutex
	status       ServerStatus
	config       ServerConfig
	activeCancel context.CancelFunc
	broadcaster  *LogBroadcaster
	runner       ServerRunner
	guiPort      int
	guiHost      string
}
