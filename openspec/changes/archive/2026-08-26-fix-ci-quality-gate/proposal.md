## Why

No workflow de Integração Contínua do GitHub Actions (`CI Quality Gate`), a etapa `Run GolangCI-Lint` falhou (execução `32927115286`) com o erro:
`can't load config: the Go language version (go1.24) used to build golangci-lint is lower than the targeted Go version (1.25.0)`.

A causa raiz decorre do descompasso de versões: o `go.mod` declara a versão `go 1.25.0`, enquanto o binário oficial do `golangci-lint` (v1.64.8) utilizado pela `golangci/golangci-lint-action@v6` foi compilado com Go 1.24. Como o analisador do linter requer uma versão de compilação igual ou superior à versão alvo do `go.mod`, o processo de linting foi interrompido. Além disso, o setup do Go nos workflows (`ci.yml` e `release.yml`) referencia `go-version: '1.26.x'`, gerando divergência entre o ambiente local, o `go.mod` e o CI.

Esta alteração é necessária agora para restabelecer o Quality Gate automatizado do repositório, garantindo conformidade com os padrões de engenharia e impedindo que Pull Requests sejam bloqueados indevidamente.

## What Changes

- **Alinhamento da versão do Go no `go.mod`**: Ajustar a diretiva `go` no `go.mod` para `1.24.0`, garantindo total compatibilidade com todas as bibliotecas, sintaxe utilizada e analisadores estáticos da comunidade (incluindo `golangci-lint`).
- **Sincronização de versão do Go nos workflows do GitHub Actions**:
  - Atualizar `.github/workflows/ci.yml` para utilizar `go-version-file: 'go.mod'`, garantindo que o runner utilize exatamente a versão estabelecida no projeto.
  - Atualizar `.github/workflows/release.yml` para sincronizar com `go-version-file: 'go.mod'`, assegurando compilações de release homogêneas.
- **Ajuste na execução do GolangCI-Lint**: Garantir que o step `Run GolangCI-Lint` execute sem erros de versão sobre o `go.mod` alinhado.
- **Atualização do Node no setup do OpenSpec**: Atualizar `node-version: 22` em `ci.yml` para eliminar os avisos de depreciação do Node 20 reportados pelo GitHub Actions.
- **Garantia de Qualidade**: Validar execução de `make check`, `make lint` e `make test` com cobertura >= 80%.

## Capabilities

### New Capabilities

*(Nenhuma nova capacidade introduzida)*

### Modified Capabilities

- `project-foundation`: Atualiza os requisitos e cenários de governança de CI/Quality Gate para garantir compatibilidade da toolchain Go entre `go.mod`, linters estritos (`golangci-lint`) e GitHub Actions.

## Impact

- **Código & Configuração**: Alteração pontual em `go.mod`, `go.sum`, `.github/workflows/ci.yml` e `.github/workflows/release.yml`.
- **APIs e Domínio**: Zero impacto sobre APIs, contratos de domínio ou regras de negócio.
- **Esteira de CI**: Restabelecimento da aprovação automática do Quality Gate em pushes e Pull Requests.
