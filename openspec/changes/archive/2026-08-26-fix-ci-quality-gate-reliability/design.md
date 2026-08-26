## Context

O workflow de CI no GitHub Actions (`.github/workflows/ci.yml`) executa a etapa `Run GolangCI-Lint` via `golangci/golangci-lint-action@v6`. Na ausência de um arquivo `.golangci.yml` explícito no repositório, o linter adota configurações padrão estritas onde regras como `goconst` sinalizam literais repetidos (como a extensão `".txt"` ocorrendo 3 vezes em funções distintas no domínio).

Além disso, discrepâncias entre o ambiente local (onde shims do asdf ou ferramentas ausentes ativam o fallback para `go vet`) e a esteira de CI geram quebras inesperadas em Pull Requests.

## Goals / Non-Goals

**Goals:**
- Criar e padronizar o arquivo `.golangci.yml` na raiz do repositório com linters ativos e calibrados para evitar falsos positivos (`goconst`, `revive`, `gosec`, `errcheck`, `govet`, etc.).
- Refatorar o mapeamento de extensões em `internal/core/domain/file.go` eliminando repetições desnecessárias de literais de string através de constantes ou agrupamentos semânticos.
- Aprimorar o `Makefile` e o script `scripts/setup.sh` para garantir que `make lint` e `make check` encontrem e executem o `golangci-lint` localmente com paridade total em relação ao GitHub Actions.
- Garantir que todos os jobs de CI passem com sucesso (código de saída 0).

**Non-Goals:**
- Desativar linters essenciais ou relaxar barreiras de segurança e qualidade do projeto.

## Decisions

### 1. Criação do `.golangci.yml` com Configurações Estáveis
- **Decisão**: Adicionar `.golangci.yml` com linters consolidados (`errcheck`, `gosimple`, `govet`, `ineffassign`, `staticcheck`, `unused`, `gosec`, `goconst`, `revive`) e configurações específicas:
  - `goconst`: `min-len: 3`, `min-occurrences: 4`, `ignore-tests: true`.
  - `run.timeout: 5m`.
  - `issues.exclude-use-default: false`.
- **Justificativa**: Evita que atualizações de versão do linter apliquem regras arbitrárias e descalibradas sobre literais curtos ou mapas de extensões.

### 2. Refatoração de Literais de Extensão em `internal/core/domain/file.go`
- **Decisão**: Declarar constantes tipadas ou conjuntos para extensões comuns de texto/documentos no domínio para evitar repetição de strings mágicas no código de produção.
- **Justificativa**: Promove código mais idiomático, limpo e à prova de falhas em qualquer linter estático.

### 3. Melhoria na Resolução do Linter no `Makefile` e `scripts/setup.sh`
- **Decisão**: Atualizar o target `lint` do `Makefile` para verificar executáveis tanto no `PATH` quanto em `$(go env GOPATH)/bin` ou `~/.asdf/shims`, e atualizar `scripts/setup.sh` para rodar `asdf reshim golang` caso o ambiente utilize gerenciador de versão.
- **Justificativa**: Garante que o desenvolvedor execute a checagem real do linter antes de subir commits.

## Risks / Trade-offs

- **[Risco] Novos linters ativados gerarem alertas em código legado** → **Mitigação**: Testar a configuração contra toda a árvore do projeto (`./...`) garantindo zero avisos/erros antes de concluir a mudança.
