## 1. Alinhamento de Versão do Go e Dependências

- [x] 1.1 Atualizar a diretiva `go` no `go.mod` para `1.24.0` e sincronizar dependências com `go mod tidy`
- [x] 1.2 Verificar a integridade do módulo com `go mod verify` e assegurar que não haja resíduos em `go.sum`

## 2. Atualização dos Workflows do GitHub Actions

- [x] 2.1 Atualizar `.github/workflows/ci.yml` para configurar `actions/setup-go@v5` com `go-version-file: 'go.mod'` e `actions/setup-node@v4` com `node-version: 22`
- [x] 2.2 Atualizar `.github/workflows/release.yml` para configurar `actions/setup-go@v5` com `go-version-file: 'go.mod'`

## 3. Validação do Quality Gate e Pipeline de Qualidade

- [x] 3.1 Executar a suíte de testes automatizados com validação de barreira de cobertura (>= 80%) via `make test`
- [x] 3.2 Executar checagem estática de linter via `make lint` e auditoria de segurança via `make vulncheck`
- [x] 3.3 Executar verificação unificada de qualidade via `make check` e compilação de binários via `make build-all`
- [x] 3.4 Validar conformidade e integridade da especificação com `openspec validate fix-ci-quality-gate --strict`
