## 1. Configuração Estável do Linter (.golangci.yml)

- [x] 1.1 Criar `.golangci.yml` na raiz do projeto com linters consolidados e parâmetros calibrados para `goconst`, `revive` e `gosec`
- [x] 1.2 Atualizar `.github/workflows/ci.yml` para assegurar execução padronizada com `.golangci.yml`

## 2. Refatoração de Literais de Extensão no Domínio

- [x] 2.1 Refatorar `DetectCategory` e `IsViewableFormat` em `internal/core/domain/file.go` eliminando repetições de strings apontadas pelo `goconst`
- [x] 2.2 Validar que todos os testes de domínio em `internal/core/domain/file_test.go` continuam passando

## 3. Harmonização do Makefile e Ferramentas Locais

- [x] 3.1 Atualizar o target `lint` do `Makefile` para resolver o binário do `golangci-lint` em múltiplos caminhos (PATH, GOPATH/bin) com diagnósticos claros
- [x] 3.2 Atualizar `scripts/setup.sh` para garantir compatibilidade com asdf e shims

## 4. Validação de Qualidade e Cobertura

- [x] 4.1 Executar `golangci-lint run ./...` e certificar zero erros/avisos
- [x] 4.2 Executar `make check` e `make test` validando cobertura $\ge 80\%$
