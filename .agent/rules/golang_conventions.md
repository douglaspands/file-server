# Convenções de Arquitetura e Código em Golang

## 1. Estrutura e Isolamento
- Todo código de domínio e regras de negócio privadas devem residir em `internal/`.
- Mantenha injeção de dependência explícita via construtores (`New...`).
- Propague `context.Context` como primeiro parâmetro em chamadas de I/O e métodos de serviço.

## 2. Testes e Qualidade (TDD / BDD)
- Utilize `testing` + `testify` (`assert` e `require`).
- Estruture testes com convenção BDD: `t.Run("Given [context] When [action] Then [outcome]", func(t *testing.T) { ... })`.
- A cobertura de testes global deve ser mantida sempre em **>= 80%** (`make test-coverage`).

## 3. Frontend e Assets
- Todos os templates e assets estáticos devem ser embutidos via `go:embed` no pacote `web/`.
- Interfaces devem priorizar **HTMX + Alpine.js + Tailwind CSS** com renderização SSR.
