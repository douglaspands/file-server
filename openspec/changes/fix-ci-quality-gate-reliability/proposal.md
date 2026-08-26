## Why

O pipeline de Integração Contínua (CI Quality Gate) no GitHub Actions tem falhado frequentemente durante a etapa `Run GolangCI-Lint` (ex: execução falha https://github.com/douglaspands/file-server/actions/runs/33021755985/job/98353546939) devido a incompatibilidades de regras padrão do linter (como a regra `goconst` acusando ocorrências repetidas de extensões no arquivo de domínio `file.go`) e à ausência de um arquivo `.golangci.yml` configurado na raiz do repositório. Além disso, o ambiente local falha em reproduzir os mesmos alertas antes do envio de Pull Requests quando o shim do `golangci-lint` entra em fallback para `go vet`.

Esta mudança estabelece uma configuração definitiva e padronizada do `golangci-lint` (`.golangci.yml`), corrige as ocorrências de strings repetidas no código de domínio e alinha o Makefile e o workflow de CI para determinismo e confiabilidade absoluta no Quality Gate.

## What Changes

- **Arquivo de Configuração do Linter (`.golangci.yml`)**: Criação da configuração oficial de linters na raiz do projeto com linters ativados (`errcheck`, `gosimple`, `govet`, `ineffassign`, `staticcheck`, `unused`, `gosec`, `goconst`, `revive`), timeouts e limites calibrados para prevenir falsos positivos em mapeamentos de extensões e testes.
- **Refatoração no Código de Domínio (`internal/core/domain/file.go`)**: Correção dos literais de string repetidos apontados pelo `goconst`, estruturando o mapeamento de extensões e categorias de forma idiomática e limpa.
- **Resolução Confiável do Linter no Makefile**: Atualização do target `lint` no `Makefile` para resolver o binário do `golangci-lint` diretamente do GOPATH ou PATH, garantindo paridade absoluta entre a verificação local (`make check`) e o GitHub Actions.
- **Alinhamento do Workflow de CI (`.github/workflows/ci.yml`)**: Configuração explícita da action do `golangci-lint` para usar a configuração do repositório sem discrepâncias de ambiente.

## Capabilities

### Modified Capabilities
- `project-foundation`: Atualiza o requisito de Governança de Git, GitHub Actions, Integração Contínua e Arquivamento com a garantia de configuração determinística e confiável do linter `.golangci.yml` e paridade entre `make check` e o CI Quality Gate.

## Impact

- **Configuração**: Criação do arquivo `.golangci.yml`.
- **Código Go**:
  - `internal/core/domain/file.go`: Refatoração do tratamento de extensões para conformidade estrita com os linters.
- **Makefile**: Ajuste do target `lint` para resolução robusta do executável `golangci-lint`.
- **CI / GitHub Actions**: `.github/workflows/ci.yml` alinhado para execução do `golangci-lint` de forma estável e determinística.
