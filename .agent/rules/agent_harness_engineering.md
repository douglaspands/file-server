# Diretrizes de Harness, Loop e Graph Engineering para Antigravity CLI

## 1. Harness Engineering (Execução Otimizada de Comandos)
- **Interface Primária**: Sempre utilize os alvos do `Makefile` para operações de desenvolvimento, testes e verificação (`make test`, `make lint`, `make check`, `make build`).
- **Eficiência de Tokens**: Evite comandos com saídas verbosas desnecessárias na janela de contexto.
- **Automação Idempotente**: Execute `make check` localmente antes de sinalizar a conclusão de qualquer tarefa.

## 2. Loop Engineering (Ciclos Rápidos de Auto-Validação)
- **Ciclo de Feedback Contínuo**:
  1. *Inspeção*: Compreender os requisitos e contratos existentes.
  2. *Implementação*: Fazer alterações mínimas e focadas.
  3. *Teste Automatizado*: Rodar testes unitários/cobertura imediatamente (`make test`).
  4. *Diagnóstico & Correção*: Caso haja falhas, analisar logs, corrigir e revalidar.
  5. *Validação Final*: Confirmar com `make check` e `openspec validate --all`.
- Não finalize tarefas sem evidências de testes passando e conformidade de cobertura >= 80%.

## 3. Graph Engineering (Navegação em DAG e Resolução Topológica)
- **Ordem de Implementação**:
  1. *Contratos / Interfaces (`internal/core/ports/`)* primeiro.
  2. *Entidades de Domínio (`internal/core/domain/`)*.
  3. *Serviços de Domínio (`internal/core/services/`)*.
  4. *Adaptadores e Handlers (`internal/adapters/`)*.
  5. *Composição / Entrypoints (`cmd/` e `main.go`)*.
- Nunca crie dependências circulares entre pacotes internos.
