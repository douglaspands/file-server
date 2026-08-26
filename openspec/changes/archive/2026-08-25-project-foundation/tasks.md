## 1. Inicialização do Módulo Go, Estrutura de Diretórios e CLI Modular

- [x] 1.1 Inicializar o módulo Go (`go.mod`) com a versão moderna do Go e adicionar dependência do Cobra CLI (`github.com/spf13/cobra`).
- [x] 1.2 Criar a árvore de diretórios padrão: `cmd/` (`root.go`, `version.go`, `serve.go`), `internal/core/domain/`, `internal/core/ports/`, `internal/core/services/`, `internal/adapters/handlers/`, `internal/version/`, `web/templates/`, `web/static/` e `scripts/`.
- [x] 1.3 Implementar o pacote `internal/version` e o comando CLI `version` com injeção de versão via `-ldflags` e detecção explícita de ambiente `dev`.
- [x] 1.4 Configurar arquivo `.gitignore` abrangente para Go, binários (`.exe`, ELF), coverage e assets compilados.
- [x] 1.5 Atualizar o nome do módulo Go e identificadores da aplicação para `file-server` (`github.com/douglas/file-server`).

## 2. Camada de Domínio, Portas e Injeção de Dependências em Go

- [x] 2.1 Definir interfaces de portas (`internal/core/ports/`) para serviços e adaptadores, garantindo desacoplamento estrito.
- [x] 2.2 Implementar entidade de domínio e serviço modelo (ex: `HealthCheckService` e status da aplicação) em `internal/core/`.
- [x] 2.3 Implementar adaptadores HTTP e handlers REST/HTML em `internal/adapters/handlers/`.
- [x] 2.4 Configurar o composition root no comando `cmd/serve.go` instanciando dependências de forma explícita e configurando encerramento gracioso (graceful shutdown) com `context.Context`.

## 3. Infraestrutura de Testes (TDD, BDD e Barreira de Cobertura)

- [x] 3.1 Adicionar dependências de testes (`testify` para assert/require) e configurar helpers de fixtures/mocks.
- [x] 3.2 Implementar testes unitários com abordagem TDD e subtestes BDD (`t.Run("Given... When... Then...")`) para os comandos CLI, serviço de domínio e handlers HTTP.
- [x] 3.3 Criar o script `scripts/coverage.sh` para executar testes, gerar `coverage.out`, exibir sumário de cobertura e validar barreira mínima de 80%.

## 4. Interface Centralizada de Comandos via Makefile

- [x] 4.1 Criar `Makefile` autodocumentado (com `make help` interativo) disponibilizando os comandos essenciais: `help`, `setup`, `dev`, `run`, `test`, `test-unit`, `test-coverage`, `lint`, `fmt`, `check`, `build`, `build-all` (Linux e Windows com `ldflags`) e `clean`.
- [x] 4.2 Configurar o linter estrito `.golangci.yml` cobrindo bugs, complexidade, concorrência e estilo idiomático.
- [x] 4.3 Implementar script de setup de ambiente para checar e instalar ferramentas de desenvolvimento (`golangci-lint`, `air`, `govulncheck`).

## 5. Live-Reloading no Desenvolvimento Web

- [x] 5.1 Configurar o arquivo `.air.toml` para monitorar alterações em arquivos Go, templates HTML (`.html`) e folhas de estilo (`.css`).
- [x] 5.2 Implementar mecanismo leve de recarregamento do navegador (Live Reload) injetável em ambiente de desenvolvimento local (`make dev`).
- [x] 5.3 Validar que o ciclo de feedback de recompilação execute em menos de 2 segundos.

## 6. Camada de Frontend Web (HTML Templates, HTMX, Alpine.js, Tailwind e Assets Embutidos)

- [x] 6.1 Estruturar os templates HTML base (`web/templates/layouts/base.html`) e páginas parciais em `web/templates/`.
- [x] 6.2 Integrar assets de frontend modernos (HTMX e Alpine.js) e configuração do Tailwind CSS (Standalone CLI).
- [x] 6.3 Criar página inicial modelo responsiva com componentes dinâmicos via HTMX para validar o fluxo de renderização SSR e interatividade parcial.
- [x] 6.4 Implementar empacotamento de templates e assets estáticos via `go:embed` (`embed.FS`) no pacote `web/`, permitindo que o binário compilado (`make build`) rode como executável único e 100% autocontido.

## 7. Estratégias do Antigravity CLI (Harness, Loop e Graph Engineering) e Governança OpenSpec

- [x] 7.1 Criar arquivos de regras em `.agent/rules/` e diretrizes de projeto capacitando o Antigravity CLI a utilizar Harness (Makefile/testes concisos), Loops de auto-validação contínua e Graph Engineering (resolução topológica de tarefas DAG).
- [x] 7.2 Configurar `openspec/config.yaml` com `context`, `rules` (proposal, specs em PT-BR PO/QA, tasks com cobertura) e `operations` (`apply` com `make check` e `archive` com commit obrigatório).

## 8. Governança Git, GitHub Actions (CI & Release Multiplataforma) e Arquivamento

- [x] 8.1 Configurar ganchos de pre-commit (`.githooks/` ou `pre-commit`) para formatação de código e validação de Conventional Commits.
- [x] 8.2 Criar o workflow do GitHub Actions `.github/workflows/ci.yml` para validar lint, testes com barreira de cobertura >= 80%, análise de vulnerabilidades (`govulncheck`) e integridade OpenSpec (`openspec validate --all`).
- [x] 8.3 Criar o workflow do GitHub Actions `.github/workflows/release.yml` para compilação cruzada automática (Linux e Windows `amd64`/`arm64`) e publicação de assets na Release ao criar tags `v*`.
- [x] 8.4 Implementar e documentar o fluxo de arquivamento garantindo commit Git automático e estruturado ao concluir cada especificação.

## 9. Documentação do Projeto e README.md

- [x] 9.1 Criar o arquivo `README.md` abrangente com badges relevantes de status/tecnologias, descrição da aplicação, passo a passo de uso com todos os comandos e flags da CLI explicados e guia dedicado para desenvolvedores.
- [x] 9.2 Validar a completude, clareza e precisão dos comandos documentados no `README.md`.

## 10. Verificação Final da Fundação

- [x] 10.1 Executar a suíte de verificação completa do harness (`make check` ou `task check`) validando compilação, linting sem advertências e cobertura de testes >= 80%.
- [x] 10.2 Testar a compilação cruzada (`make build-all`), validação do comando `version` e integridade das especificações com `openspec validate --all`.
