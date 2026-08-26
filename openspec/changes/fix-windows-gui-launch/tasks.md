## 1. Configuração da CLI e Desativação do Mousetrap

- [x] 1.1 Configurar `cobra.MousetrapHelpText = ""` na inicialização de `cmd/root.go` para permitir execução sem argumentos no Windows Explorer
- [x] 1.2 Validar a lógica de resolução de ambiente desktop no comando raiz para inicialização transparente da interface gráfica

## 2. Testes Unitários e Validação de Cobertura

- [x] 2.1 Implementar teste unitário em `cmd/cmd_test.go` validando que `cobra.MousetrapHelpText` está desabilitado
- [x] 2.2 Executar a suíte de testes e validar barreira de cobertura de testes $\ge 80\%$ via `./scripts/coverage.sh`

## 3. Validação de Qualidade e Compilação Multiplataforma

- [x] 3.1 Executar `make check` (formatação, lint estrito e testes)
- [x] 3.2 Compilar binários de produção via `make build-all` e verificar integridade dos executáveis Windows
