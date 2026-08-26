SHELL := /bin/bash
.DEFAULT_GOAL := help

# Metadados de Compilação
VERSION ?= $(shell git describe --tags --exact-match 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE ?= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
MODULE := github.com/douglas/file-server

LDFLAGS := -X '$(MODULE)/internal/version.Version=$(VERSION)' \
           -X '$(MODULE)/internal/version.Commit=$(COMMIT)' \
           -X '$(MODULE)/internal/version.Date=$(DATE)' \
           -s -w

BIN_DIR := bin
DIST_DIR := dist
BINARY_NAME := file-server

export PATH := $(PATH):$(shell go env GOPATH)/bin:$(PWD)/$(BIN_DIR)

.PHONY: help
help: ## Exibe este menu interativo com os comandos disponíveis
	@echo "========================================================================"
	@echo "🛠️  File Server - Comandos Disponíveis (Makefile)"
	@echo "========================================================================"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo "========================================================================"

.PHONY: setup
setup: ## Instala e verifica ferramentas locais de desenvolvimento (linters, air, vulncheck)
	@echo "📦 Instalando e verificando ferramentas locais..."
	@./scripts/setup.sh

.PHONY: fmt
fmt: ## Formata todo o código fonte Go e templates
	@echo "🎨 Formatando código Go..."
	@go fmt ./...

.PHONY: lint
lint: ## Executa o linter estrito (golangci-lint)
	@echo "🔍 Executando golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "⚠️  golangci-lint não encontrado. Executando go vet..."; \
		go vet ./...; \
	fi

.PHONY: vulncheck
vulncheck: ## Executa auditoria de vulnerabilidades de segurança
	@echo "🛡️  Verificando vulnerabilidades de segurança..."
	@if command -v govulncheck >/dev/null 2>&1; then \
		govulncheck ./...; \
	else \
		echo "⚠️  govulncheck não encontrado. Execute 'make setup' para instalar."; \
	fi

.PHONY: test
test: test-coverage ## Executa a suíte de testes com validação de cobertura (>= 80%)

.PHONY: test-unit
test-unit: ## Executa testes unitários rápidos
	@echo "🧪 Executando testes unitários..."
	@go test -v -short ./...

.PHONY: test-coverage
test-coverage: ## Executa testes e valida barreira mínima de cobertura de 80%
	@./scripts/coverage.sh

.PHONY: check
check: fmt lint test ## Roda o pipeline de qualidade completo localmente (fmt + lint + test)
	@echo "✅ Todas as checagens de qualidade foram concluídas com sucesso!"

.PHONY: run
run: ## Executa a aplicação diretamente a partir do código fonte
	@go run main.go serve

.PHONY: dev
dev: ## Inicia o servidor em modo de desenvolvimento com live-reload (Air)
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "⚠️  Air não instalado. Iniciando modo padrão com go run..."; \
		go run main.go serve; \
	fi

.PHONY: build
build: ## Compila o binário de produção otimizado com assets embutidos
	@echo "🔨 Compilando binário para a plataforma atual (versão: $(VERSION))..."
	@mkdir -p $(BIN_DIR)
	@CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -trimpath -o $(BIN_DIR)/$(BINARY_NAME) main.go
	@echo "✅ Binário gerado em $(BIN_DIR)/$(BINARY_NAME)"

.PHONY: build-all
build-all: build-linux build-windows ## Compila binários para Linux e Windows (amd64 e arm64)
	@echo "🌍 Compilação cruzada para todas as plataformas concluída em $(DIST_DIR)/"

.PHONY: build-linux
build-linux: ## Compila binários para Linux (amd64 e arm64)
	@echo "🐧 Compilando binários para Linux..."
	@mkdir -p $(DIST_DIR)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -trimpath -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 main.go
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -trimpath -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 main.go

.PHONY: build-windows
build-windows: ## Compila binários para Windows (amd64 e arm64)
	@echo "🪟 Compilando binários para Windows..."
	@mkdir -p $(DIST_DIR)
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -trimpath -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe main.go
	@CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -trimpath -o $(DIST_DIR)/$(BINARY_NAME)-windows-arm64.exe main.go

.PHONY: clean
clean: ## Limpa binários compilados, relatórios de cobertura e temporários
	@echo "🧹 Limpando arquivos compilados e relatórios temporários..."
	@rm -rf $(BIN_DIR) $(DIST_DIR) coverage.out coverage.html tmp .air
	@echo "✅ Limpeza concluída."
