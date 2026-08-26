#!/usr/bin/env bash
set -euo pipefail

echo "========================================================================"
echo "🛠️  SFTP Server - Setup de Ferramentas de Desenvolvimento"
echo "========================================================================"

GOPATH_BIN="$(go env GOPATH)/bin"
mkdir -p "${GOPATH_BIN}"

install_tool() {
    local name="$1"
    local package="$2"

    if command -v "${name}" >/dev/null 2>&1 || [ -f "${GOPATH_BIN}/${name}" ]; then
        echo "✅ ${name} já está instalado."
    else
        echo "⬇️  Instalando ${name} (${package})..."
        go install "${package}@latest"
        echo "✅ ${name} instalado com sucesso em ${GOPATH_BIN}/${name}."
    fi
}

install_tool "golangci-lint" "github.com/golangci/golangci-lint/cmd/golangci-lint"
install_tool "air" "github.com/air-verse/air"
install_tool "govulncheck" "golang.org/x/vuln/cmd/govulncheck"

echo "========================================================================"
echo "✨ Todas as ferramentas de desenvolvimento estão prontas para uso!"
echo "💡 Certifique-se de que ${GOPATH_BIN} está no seu PATH."
echo "========================================================================"
