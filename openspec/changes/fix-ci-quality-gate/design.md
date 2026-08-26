## Context

O projeto File Server adota Clean Architecture em Go, barreira de cobertura de testes inegociável (&ge; 80%) e pipeline de CI no GitHub Actions com validação de linters estritos (`golangci-lint`), auditoria de segurança (`govulncheck`) e integridade do OpenSpec.

Atualmente, o pipeline no GitHub Actions (`.github/workflows/ci.yml`) falha no step de linting com:
`can't load config: the Go language version (go1.24) used to build golangci-lint is lower than the targeted Go version (1.25.0)`.

Para a motivação detalhada e escopo funcional, consulte [proposal.md](file:///home/douglas/Workspace/gemini/sftp-server/openspec/changes/fix-ci-quality-gate/proposal.md).

## Goals / Non-Goals

**Goals:**
- Alinhar a versão do Go no `go.mod` para `1.24.0`, suportada nativamente pelos binários oficiais pré-compilados do `golangci-lint`.
- Configurar os workflows de CI e Release (`.github/workflows/ci.yml` e `.github/workflows/release.yml`) para utilizar `go-version-file: 'go.mod'`, garantindo sincronização automática e eliminando versionamento hardcoded conflitante.
- Atualizar a versão do Node.js de 20 para 22 LTS no workflow de CI para validação do OpenSpec, eliminando avisos de depreciação do runner.
- Validar que a esteira completa (`make check`, `make lint`, `make test`, `make vulncheck`) execute com sucesso.

**Non-Goals:**
- Não alterar nenhuma arquitetura, pacote de domínio (`internal/core/`), adaptadores (`internal/adapters/`) ou frontend (`web/`).
- Não alterar dependências de produção ou contratos da CLI.

## Decisions

### Decisão 1: Alinhar `go.mod` para `go 1.24.0`
- **Racional**: O binário oficial do `golangci-lint` v1.64.x é distribuído compilado com Go 1.24. Todas as funcionalidades atuais do projeto são 100% idiomáticas e suportadas em Go 1.24.
- **Alternativas consideradas**:
  - *Manter `go 1.25.0` e forçar `install-mode: goinstall` na action*: Descartado porque compilar o linter a partir dos fontes a cada execução do CI adiciona tempo desnecessário de build (~1 a 2 minutos por job) e viola a recomendação oficial dos mantenedores do linter.

### Decisão 2: Utilizar `go-version-file: 'go.mod'` no `actions/setup-go@v5`
- **Racional**: Elimina descompassos entre a versão declarada no módulo Go e a versão baixada pelo runner do GitHub Actions.
- **Alternativas consideradas**:
  - *Hardcode de `go-version: '1.24.x'`*: Descartado por exigir atualização manual em múltiplos locais caso a versão do Go seja elevada no futuro.

### Decisão 3: Atualizar `node-version: 22` no step de OpenSpec do CI
- **Racional**: O runner do GitHub Actions já emite avisos de depreciação para o Node 20. O Node 22 (LTS ativo) previne falhas futuras na validação via `npx openspec`.

## Risks / Trade-offs

- **[Risco]** Diferença na resolução de dependências no `go.mod` ao alterar a versão da toolchain.
  - **Mitigação**: Executar `go mod tidy` e `go mod verify` localmente e validar via `git diff --exit-code go.mod go.sum`.
- **[Risco]** Comportamento divergente no Makefile ou scripts locais.
  - **Mitigação**: Executar `make check` e `make build-all` garantindo paridade total entre ambiente local e CI.
