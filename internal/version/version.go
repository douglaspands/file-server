package version

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	// Version é a versão semântica da aplicação (injetada via -ldflags em releases, ou 'dev' por padrão).
	Version = "dev"
	// Commit é o hash do commit git da compilação.
	Commit = "none"
	// Date é a data/hora UTC da compilação.
	Date = "unknown"
)

// Info encapsula as informações estruturadas de versão.
type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"build_date"`
	GoVersion string `json:"go_version"`
	Platform  string `json:"platform"`
	IsDev     bool   `json:"is_dev"`
}

// Get retorna os metadados de versão da aplicação.
func Get() Info {
	isDev := Version == "dev" || strings.HasPrefix(Version, "v0.0.0-dev") || Version == ""
	v := Version
	if v == "" {
		v = "dev"
	}
	return Info{
		Version:   v,
		Commit:    Commit,
		BuildDate: Date,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		IsDev:     isDev,
	}
}

// String formata os detalhes de versão para exibição no terminal.
func (i Info) String() string {
	if i.IsDev {
		return fmt.Sprintf("File Server version: %s (development build, commit: %s, date: %s, %s, %s)",
			i.Version, i.Commit, i.BuildDate, i.GoVersion, i.Platform)
	}
	return fmt.Sprintf("File Server version: %s (commit: %s, date: %s, %s, %s)",
		i.Version, i.Commit, i.BuildDate, i.GoVersion, i.Platform)
}
