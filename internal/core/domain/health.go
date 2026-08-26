package domain

import (
	"time"
)

// HealthStatus representa o estado de saúde da aplicação.
type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version"`
	Uptime    string            `json:"uptime"`
	Details   map[string]string `json:"details,omitempty"`
}

// SystemInfo representa os dados gerais e ambiente do servidor.
type SystemInfo struct {
	AppName   string    `json:"app_name"`
	Version   string    `json:"version"`
	GoVersion string    `json:"go_version"`
	OS        string    `json:"os"`
	Arch      string    `json:"arch"`
	NumCPU    int       `json:"num_cpu"`
	StartedAt time.Time `json:"started_at"`
}
