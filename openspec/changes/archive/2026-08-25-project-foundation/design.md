## Context

O repositório necessita de uma fundação profissional em Go para desenvolvimento de aplicações web com renderização HTML no servidor. Veja [proposal.md](proposal.md) para motivação e objetivos de negócio. As especificações de comportamento e critérios de aceitação PO/QA estão detalhados em [specs/project-foundation/spec.md](specs/project-foundation/spec.md).

## Goals / Non-Goals

**Goals:**
- Estabelecer a arquitetura padrão em Go (Clean Architecture / Ports & Adapters) garantindo isolamento total do domínio em `internal/`.
- Configurar a esteira de TDD/BDD com execução via `go test`, asserções idiomáticas, mocks de contratos e barreira de cobertura de código >= 80%.
- Implementar Harness Engineering via `Taskfile.yml` (ou `Makefile`) para padronizar comandos de build, teste, lint, segurança e setup.
- Configurar Loop Engineering com feedback rápido usando `air` para live-reload de código Go, templates HTML e assets.
- Aplicar Graph Engineering para manter o grafo de dependências entre módulos estritamente acíclico com injeção de dependência explícita.
- Estruturar diretrizes e regras de projeto em `.agent/rules/` e `openspec/config.yaml` para otimização de tokens e governança de futuras propostas de IA.
- Apresentar e padronizar o stack e plugins de frontend (HTMX + Alpine.js + Tailwind CSS / Templ).
- Implementar pipelines de CI no GitHub Actions e ganchos de pre-commit com Conventional Commits.

**Non-Goals:**
- Implementar regras de negócio específicas da aplicação final (ex: funcionalidades de SFTP ou regras de domínio comercial), focando estritamente na fundação e scaffold da engenharia.
- Adotar frameworks complexos de injeção de dependência baseados em reflection ou geração de código pesada em tempo de execução.

## Decisions

### 1. Estrutura de Pacotes Go (Clean Architecture / Ports & Adapters)
- **Decisão**: Utilizar o layout padrão de projeto Go:
  - `cmd/server/main.go`: Ponto de entrada, parse de flags/env e composição de dependências (composition root).
  - `internal/core/domain/`: Entidades de negócio puras sem dependências externas.
  - `internal/core/ports/`: Interfaces de entrada (use cases) e saída (repositórios, adaptadores).
  - `internal/core/services/`: Implementação da lógica dos casos de uso.
  - `internal/adapters/handlers/`: Controladores HTTP e renderizadores de templates.
  - `internal/adapters/repositories/`: Implementações de persistência/armazenamento.
  - `web/templates/`: Arquivos de templates HTML e layouts parciais.
  - `web/static/`: Folhas de estilo CSS, scripts JS e imagens.
  - `scripts/`: Scripts utilitários e harness de suporte.
- **Alternativas consideradas**:
  - *Monolito plano (package único)*: Rápido no início, mas inadequado para evolução com TDD e isolamento de contratos.

### 2. Framework e Práticas de Testes (TDD & BDD)
- **Decisão**:
  - Utilizar o pacote nativo `testing` do Go combinado com `github.com/stretchr/testify` (`assert` e `require`) para asserções limpas e expressivas.
  - Estruturar testes em subtestes BDD declarativos utilizando a convenção `t.Run("Given [context] When [action] Then [outcome]", func(t *testing.T) { ... })`.
  - Mocks manuais baseados em interfaces de `ports/` ou gerados via `mockery` para garantir determinismo e isolamento completo sem acoplamento a banco de dados real nos testes unitários.
  - Script de validação de cobertura `scripts/coverage.sh` que calcula a porcentagem total de cobertura via `go tool cover` e emite erro se for inferior a 80%.

### 3. Interface Centralizada de Comandos via Makefile
- **Decisão**: Adotar um **`Makefile` autodocumentado** como interface primária e universal de comandos da aplicação:
  - `make help` (ou apenas `make`): Exibe dinamicamente o menu com todos os comandos disponíveis e descrições formatadas.
  - `make setup`: Instala e checa ferramentas locais (`golangci-lint`, `air`, `govulncheck`, `tailwind`).
  - `make dev`: Inicia o loop de desenvolvimento com live-reload (Air) e rebuild automático de CSS/templates.
  - `make run`: Executa a aplicação diretamente a partir do código fonte.
  - `make test`: Executa suíte de testes com validação de barreira de cobertura (>= 80%).
  - `make test-unit`: Executa testes unitários rápidos.
  - `make lint`: Executa análise estática com `golangci-lint`.
  - `make fmt`: Formata código Go e templates HTML.
  - `make check`: Roda o pipeline completo de qualidade local (`fmt` + `lint` + `vulncheck` + `test`).
  - `make build`: Compila o binário de produção otimizado com stripping de símbolos.
  - `make clean`: Limpa binários compilados, relatórios de cobertura e arquivos temporários.

### 4. Live-Reloading e Feedback Visual no Desenvolvimento Web
- **Decisão**: Configurar `air` (`.air.toml`) monitorando diretórios de código Go (`internal/`, `cmd/`) e diretórios web (`web/templates/`, `web/static/`), acionando recompilação incremental e recarregamento automático no navegador via script de live-reload injetado em modo de desenvolvimento (`make dev`).

### 5. Arquitetura Go e Injeção Explícita de Dependências
- **Decisão**:
  - Código Go com desacoplamento estrito e sem dependências circulares: `adapters` -> `ports` <- `services` -> `domain`.
  - Injeção de dependência explícita manual no composition root (`cmd/server/main.go`), tornando todo o fluxo de inicialização rastreável e determinístico sem uso de frameworks pesados de reflexão.

### 6. Estratégias do Antigravity CLI (Harness, Loop e Graph Engineering) & Governança OpenSpec
- **Decisão**:
  - Configurar em `.agent/rules/` e nas regras do projeto as práticas e estratégias empregadas pelo **Antigravity CLI** para desenvolver a aplicação:
    1. **Harness Engineering**: O Antigravity CLI utiliza o `Makefile` como seu *harness* de automação e teste, executando comandos determinísticos com flags silenciosas/concisas (`make test`, `make lint`, `make check`), otimizando drasticamente o consumo de tokens na janela de contexto.
    2. **Loop Engineering**: O Antigravity CLI opera em loops rápidos de feedback contínuo (Inspeção -> Implementação -> Execução de Testes Automatizados -> Diagnóstico de Erros -> Refinamento -> Validação Final), garantindo resolução de problemas com evidências antes de entregar tarefas.
    3. **Graph Engineering**: O Antigravity CLI utiliza representações em grafo (DAG) das tarefas e das dependências entre arquivos/módulos para planejar a sequência topológica de edições (criando contratos/portas antes das implementações), eliminando retrabalho e bloqueios de dependência.
  - Configurar `openspec/config.yaml` de forma robusta e definitiva:
    - **`context`**: Documenta a stack (Go, HTML, HTMX, Tailwind), arquitetura Clean/Ports & Adapters, TDD/BDD, barreira de cobertura >= 80%, Makefile e convenções.
    - **`rules.proposal`**: Exige alinhamento com a arquitetura e declaração explícita de impacto.
    - **`rules.specs`**: Exige escrita em PT-BR e perspectiva conjunta PO (regras de negócio e aceitação) e QA (cenários BDD com 4 hashtags `#### Scenario:` e validações de borda).
    - **`rules.tasks`**: Exige tarefas granulares e rastreáveis incluindo testes unitários/integração para cada caso de uso.
    - **`operations.apply.guidance`**: Instruções para que a aplicação sempre execute `make check` antes de finalizar a implementação.
    - **`operations.archive.guidance`**: Instrução mandatória para que o agente crie o commit Git com Conventional Commits ao arquivar qualquer especificação concluída.

### 7. Stack de Frontend, Assets Embutidos (`go:embed`) e Plugins
- **Decisão e Opções para Frontend**:
  - **Opção Recomendada (HTMX + Alpine.js + Tailwind CSS)**:
    - **HTMX**: Fornece AJAX, WebSockets e Server-Sent Events direto no HTML, permitindo interfaces dinâmicas com renderização no servidor Go.
    - **Alpine.js**: Fornece reatividade leve no cliente (dropdowns, modais, validação de inputs) sem necessidade de build React/Vue.
    - **Tailwind CSS (Standalone CLI)**: Utilitários CSS modernos compilados localmente sem depender de ecossistema Node.js pesado.
  - **Empacotamento de Assets via `go:embed` (Binário Único Autocontido)**:
    - O pacote `web/` expõe uma interface `io/fs.FS` utilizando a diretiva `//go:embed templates/* static/*` (`embed.FS`).
    - Todos os templates HTML, arquivos CSS, scripts JS e imagens são compilados diretamente dentro do binário final executável (`make build`), permitindo que o executável seja transportado e executado em qualquer ambiente sem pastas ou arquivos externos adicionais.
    - Em ambiente de desenvolvimento local (`make dev`), o carregador de templates pode alternar dinamicamente para leitura do sistema de arquivos (`os.DirFS`), permitindo edição instantânea com live-reload.
  - **Alternativa Tipada**: **Templ (`a-h/templ`)** para templates Go tipados com verificação em tempo de compilação.
  - **Plugins & Ferramental Recomendado**:
    - Extensão Tailwind CSS IntelliSense / HTMX syntax highlighting.
    - Script leve de Live Reload para templates.

### 8. Estrutura de CLI Extensível, Injeção de Versão e Release Multiplataforma
- **Decisão - CLI Modular**:
  - Utilizar o padrão de comandos baseado em **Cobra** (`github.com/spf13/cobra`), estruturando os comandos em `cmd/` com binário nomeado `file-server`:
    - `cmd/root.go`: Comando raiz (`file-server`), flags globais (ex: `--config`, `--verbose`).
    - `cmd/version.go`: Subcomando e flag `version` (`file-server version`, `file-server -v`, `file-server --version`, `file-server version --json`).
    - `cmd/serve.go`: Subcomando para iniciar o servidor web/HTTP (`file-server serve --port=8080 --host=0.0.0.0`).
- **Decisão - Versionamento Dinâmico (Release vs Dev)**:
  - Criar o pacote `internal/version` contendo as variáveis `Version = "dev"`, `Commit = "none"`, `Date = "unknown"`.
  - No `Makefile`, injetar as variáveis via `ldflags`:
    ```makefile
    VERSION ?= $(shell git describe --tags --exact-match 2>/dev/null || echo "dev")
    COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
    DATE ?= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
    LDFLAGS = -X 'github.com/douglas/file-server/internal/version.Version=$(VERSION)' \
              -X 'github.com/douglas/file-server/internal/version.Commit=$(COMMIT)' \
              -X 'github.com/douglas/file-server/internal/version.Date=$(DATE)' -s -w
    ```
  - Quando o usuário rodar `file-server version` sem release tag, a saída reporta claramente o estado `dev` (ex: `File Server version: dev (commit: abc1234, built at: 2026-08-25)`). Em release tag oficial, exibe a tag exata (ex: `File Server version: v1.0.0`).
- **Decisão - Pipeline de Release Multiplataforma (.github/workflows/release.yml)**:
  - Disparado em `push: tags: ['v*']`.
  - Matrix de compilação cruzada para:
    - **Linux**: `GOOS=linux GOARCH=amd64` e `GOOS=linux GOARCH=arm64` (binário `file-server` empacotado em `.tar.gz`).
    - **Windows**: `GOOS=windows GOARCH=amd64` e `GOOS=windows GOARCH=arm64` (binário `file-server.exe` empacotado em `.zip`).
  - Geração de checksums SHA256 e publicação automática dos artefatos na Release do GitHub.

### 9. Documentação e Padrão de Engenharia no README.md
- **Decisão**:
  - O `README.md` será a referência central do repositório, organizado rigorosamente na seguinte estrutura:
    1. **Header & Badges**: Título do projeto com badges SVG (Go Version, Build/CI Status, Test Coverage >= 80%, Latest Release, License).
    2. **Visão Geral**: Descrição clara do propósito da aplicação (File Server moderno e autocontido com interface web e CLI).
    3. **Passo a Passo de Instalação e Uso**:
       - Download do binário pré-compilado ou compilação local.
       - Guia completo e tabela de **todos os comandos e flags disponíveis** da CLI (`file-server`, `version`, `serve` e opções `-p`, `-h`, `-j`, `-v`, `--config`).
       - Exemplos práticos de uso e navegação na interface web embutida.
    4. **Guia do Desenvolvedor**:
       - Pré-requisitos e setup do ambiente local (`make setup`).
       - Arquitetura do projeto e padrão de pastas (`internal/core/`, `internal/adapters/`, `cmd/`, `web/`).
       - Harness de comandos via `Makefile` (`make check`, `make dev`, `make test`, `make build-all`).
       - Metodologia de testes (TDD/BDD com meta de cobertura >= 80%).
       - Live-reloading (`make dev` com Air e SSE hook).
       - Diretrizes do Antigravity CLI e Conventional Commits.

### 10. Governança Git, GitHub Actions e Ciclo de Vida de Especificações
- **Decisão**:
  - Configurar `.github/workflows/ci.yml` executando em cada Pull Request e push na `main`:
    1. Checkout e Setup de Go.
    2. `golangci-lint`.
    3. `govulncheck`.
    4. Execução de testes com validação de barreira de cobertura de 80%.
    5. `openspec validate --all` para integridade das especificações.
  - Configuração de Conventional Commits (`feat:`, `fix:`, `refactor:`, `test:`, `docs:`, `chore:`).
  - **Fluxo de Commit no Arquivamento**: Ao finalizar a implementação de qualquer change e disparar o arquivamento (`/openspec-archive-change` ou `openspec archive`), o assistente deve automaticamente preparar o commit no Git contendo todos os arquivos implementados, specs sincronizadas e mensagem padronizada no formato Conventional Commits.

## Risks / Trade-offs

- **[Complexidade inicial da estrutura de camadas e CLI]** → *Mitigação*: Manter comandos concisos e um exemplo canônico de caso de uso (ex: `HealthCheck` / status) pronto na fundação para servir de modelo de referência.
- **[Barreira de 80% de cobertura bloqueando PRs]** → *Mitigação*: Facilitar a criação de testes com helpers, fixtures e testes de tabela (table-driven tests) padrão no Go.
- **[Dependência de ferramentas locais (Task, Air, Linters)]** → *Mitigação*: O comando `make setup` instala automaticamente todas as ferramentas no `$GOPATH/bin` ou diretório local `bin/`.
