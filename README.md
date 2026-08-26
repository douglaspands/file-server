# ⚡ File Server

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go&logoColor=white)](https://golang.org)
[![CI Quality Gate](https://img.shields.io/badge/CI-Passing-success?style=flat&logo=githubactions&logoColor=white)](.github/workflows/ci.yml)
[![Test Coverage](https://img.shields.io/badge/Coverage-83.6%25-brightgreen?style=flat)](scripts/coverage.sh)
[![Architecture](https://img.shields.io/badge/Architecture-Clean%20%2F%20Ports%20%26%20Adapters-informational?style=flat)](#-arquitetura-e-estrutura-de-pacotes)
[![Platforms](https://img.shields.io/badge/Platforms-Linux%20%7C%20Windows-blueviolet?style=flat)](#-compilação-multiplataforma)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat)](LICENSE)

> Servidor de arquivos e aplicação web moderna de alto desempenho desenvolvida em Go (Golang), com renderização no servidor (SSR), reatividade hipermidiática via HTMX, CLI modular extensível e empacotamento integral de assets (`go:embed`) em um **único binário executável 100% autocontido e portátil**.

---

## 📖 Visão Geral

O **File Server** é uma solução completa para gerenciamento e disponibilização de arquivos projetada sob os mais altos padrões de Engenharia de Software.

### ✨ Principais Destaques

- 🚀 **Binário Único Autocontido**: Todos os templates HTML, layouts, estilos CSS e scripts JS são compilados dentro do executável via `go:embed`. Transporte apenas um arquivo para qualquer servidor Linux ou Windows.
- 💻 **CLI Modular e Extensível**: Inicialização estruturada via Cobra, com suporte a múltiplos comandos, argumentos, flags e versionamento dinâmico via `-ldflags`.
- 🌐 **Interface Web Moderna e Reativa**: Renderização no servidor (SSR) veloz e leve combinada com **HTMX**, **Alpine.js** e **Tailwind CSS** para interatividade parcial sem recarregamento de página.
- 🧪 **Excelência em Qualidade (TDD / BDD)**: Suíte rigorosa de testes unitários e de integração com cobertura de código validada continuamente (**meta inegociável &ge; 80%**).
- ⚙️ **Harness Engineering Universal**: `Makefile` autodocumentado como ponto de contato unificado para desenvolvimento, testes, linters e builds.
- 🤖 **IA & Spec-Driven Development**: Governança perpétua para Antigravity CLI via OpenSpec, integrando Harness, Loop e Graph Engineering.

---

## 🚀 Guia de Uso (Passo a Passo)

### 1. Obtenção do Executável

Você pode compilar localmente ou utilizar o binário da release:

```bash
# Compilar o binário para seu sistema operacional
make build

# O binário será gerado em bin/file-server
./bin/file-server --help
```

---

### 2. Catálogo Completo de Comandos da CLI

A aplicação oferece uma interface de terminal rica com subcomandos e flags configuráveis:

| Comando / Flag | Descrição | Exemplo de Uso |
| :--- | :--- | :--- |
| `file-server` | Comando raiz; exibe a visão geral da aplicação | `./bin/file-server` |
| `file-server --help` (`-h`) | Exibe a ajuda interativa com todos os comandos e flags | `./bin/file-server --help` |
| `file-server version` | Exibe a versão semântica, commit, data de compilação e plataforma | `./bin/file-server version` |
| `file-server version --json` (`-j`) | Exibe os metadados completos de versão em formato JSON estruturado | `./bin/file-server version -j` |
| `file-server serve` | Inicia o servidor HTTP web e API na porta padrão (`8080`) | `./bin/file-server serve` |
| `--port` (`-p`) | Define a porta TCP na qual o servidor irá escutar (default: `8080`) | `./bin/file-server serve -p 3000` |
| `--host` | Define o endereço IP/host de escuta (default: `0.0.0.0`) | `./bin/file-server serve --host 127.0.0.1` |
| `--config` | Caminho para arquivo customizado de configuração YAML | `./bin/file-server --config config.yaml` |
| `--verbose` (`-v`) | Habilita logs detalhados de diagnóstico | `./bin/file-server serve -v` |

---

### 3. Exemplos Práticos de Uso

#### Iniciar o Servidor Web e Acessar no Navegador
```bash
# Inicia o servidor na porta 8080
./bin/file-server serve --port=8080

# Acesse no seu navegador:
# http://localhost:8080
```

#### Consultar a Versão da Aplicação
```bash
# Versão textual formatada
./bin/file-server version

# Versão em formato JSON (ideal para scripts e esteiras)
./bin/file-server version --json
```

#### Endpoints de API Disponíveis
- `GET /`: Interface web responsiva com status do servidor.
- `GET /partials/health`: Fragmento HTML para atualização dinâmica via HTMX.
- `GET /api/health`: Status de saúde e métricas de uptime em JSON.
- `GET /api/version`: Metadados estruturados de versão e build em JSON.

---

## 🛠️ Guia para Desenvolvedores

Esta seção detalha os padrões arquiteturais, ferramentas e fluxos de trabalho para contribuir com o projeto.

### 📋 Pré-requisitos

- **Go**: Versão `1.26` ou superior instalada ([golang.org](https://golang.org)).
- **Make**: Utilitário Make (`/usr/bin/make`) para automação de tarefas.
- **Git**: Controle de versão.

---

### 📦 Setup do Ambiente de Desenvolvimento

Execute o comando de setup para instalar e verificar todas as ferramentas recomendadas (`golangci-lint`, `air` para live-reload e `govulncheck`):

```bash
make setup
```

---

### 🏗️ Arquitetura e Estrutura de Pacotes

O projeto adota os princípios de **Clean Architecture** e **Ports & Adapters**, assegurando o encapsulamento estrito das regras de negócio em `internal/`:

```
file-server/
├── cmd/                          # Camada de entrada da CLI (Cobra)
│   ├── root.go                   # Comando base e flags globais
│   ├── version.go                # Subcomando 'version'
│   └── serve.go                  # Subcomando 'serve' (Composition Root)
├── internal/                     # Código privado do domínio (não importável externamente)
│   ├── core/
│   │   ├── domain/               # Entidades de negócio e modelos puros
│   │   ├── ports/                # Interfaces de entrada/saída (contratos)
│   │   └── services/             # Implementação da lógica de casos de uso
│   ├── adapters/
│   │   └── handlers/             # Controladores HTTP, HTML e endpoints REST
│   └── version/                  # Metadados de compilação (ldflags)
├── web/                          # Recursos de frontend
│   ├── templates/                # Layouts base, páginas e partials HTML
│   ├── static/                   # Folhas de estilo CSS e scripts JS
│   └── web.go                    # Empacotamento embutido (go:embed embed.FS)
├── scripts/                      # Scripts do harness (cobertura, setup)
├── .agent/rules/                 # Diretrizes de IA (Harness, Loop e Graph Engineering)
├── .github/workflows/            # CI/CD (Quality Gate e Release Multiplataforma)
├── .githooks/                    # Ganchos de pre-commit e Conventional Commits
├── openspec/                     # Especificações normativas (Spec-Driven Development)
└── Makefile                      # Interface universal de comandos
```

---

### 🕹️ Harness de Comandos (Makefile)

O `Makefile` autodocumentado centraliza todas as operações do ciclo de vida:

```bash
# Exibir o menu interativo com todos os comandos disponíveis
make help
```

| Alvo | Finalidade |
| :--- | :--- |
| `make help` | Exibe o menu com todos os comandos disponíveis e descrições |
| `make setup` | Instala e checa ferramentas locais (`golangci-lint`, `air`, `govulncheck`) |
| `make dev` | Inicia o servidor em modo de desenvolvimento com live-reload (Air) |
| `make run` | Executa a aplicação diretamente do código fonte |
| `make fmt` | Formata o código Go e templates |
| `make lint` | Executa análise estática com `golangci-lint` (ou `go vet`) |
| `make vulncheck` | Executa auditoria de vulnerabilidades de segurança com `govulncheck` |
| `make test` | Executa toda a suíte de testes com validação de barreira de cobertura (&ge; 80%) |
| `make test-unit` | Executa testes unitários rápidos |
| `make test-coverage`| Roda `./scripts/coverage.sh`, gera relatório HTML e valida a meta de 80% |
| `make check` | Executa o pipeline completo de qualidade local (`fmt` + `lint` + `test`) |
| `make build` | Compila o binário de produção otimizado com assets embutidos para a plataforma atual |
| `make build-all` | Realiza compilação cruzada para Linux e Windows (`amd64` e `arm64`) |
| `make clean` | Limpa diretórios de build (`bin/`, `dist/`), relatórios e temporários |

---

### 🧪 Testes Automatizados (TDD & BDD)

Os testes seguem a convenção BDD declarativa estruturada em subtestes:

```go
func TestHealthService(t *testing.T) {
    t.Run("Given initialized health service When checking health Then returns healthy status", func(t *testing.T) {
        svc := services.NewHealthService()
        status, err := svc.Check(context.Background())

        require.NoError(t, err)
        assert.Equal(t, "healthy", status.Status)
    })
}
```

Para rodar os testes e inspecionar o relatório visual de cobertura:
```bash
make test-coverage
# Abra o arquivo coverage.html no navegador para verificar linhas cobertas
```

---

### ⚡ Live-Reloading no Frontend

Para desenvolver a interface com recarregamento instantâneo a cada alteração em arquivos `.go`, `.html` ou `.css`:

```bash
make dev
```

---

### 🌍 Compilação Multiplataforma (Linux & Windows)

Para gerar executáveis para todas as arquiteturas suportadas com injeção automática de versão via `ldflags`:

```bash
make build-all
```

Arquivos gerados em `dist/`:
- `file-server-linux-amd64` (Linux x86_64)
- `file-server-linux-arm64` (Linux ARM64)
- `file-server-windows-amd64.exe` (Windows x86_64)
- `file-server-windows-arm64.exe` (Windows ARM64)

---

### 🏷️ Convenções de Commit e Githooks

O repositório adota o padrão **Conventional Commits**:

- `feat:` Nova funcionalidade
- `fix:` Correção de bug
- `refactor:` Refatoração de código sem alteração de comportamento
- `test:` Inclusão ou ajuste de testes
- `docs:` Alterações em documentações
- `chore:` Tarefas de manutenção ou dependências

Os ganchos em `.githooks/` validam a formatação do código e o padrão da mensagem antes de cada commit.

---

## 📄 Licença

Este projeto está sob a licença [MIT](LICENSE).
