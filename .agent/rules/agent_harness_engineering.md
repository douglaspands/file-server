# Diretrizes de Harness, Loop, Graph Engineering e Economia de Tokens para Antigravity CLI

## 1. AI Harness Engineering (Scaffolding, Guardrails e Governança de Agente)
- **Scaffolding de Regras & Contexto**: O harness do AGY é composto por regras determinísticas (`.agent/rules/`, `AGENTS.md`, `GEMINI.md`) e restrições de sistema que estabelecem limites claros e instruções precisas para a IA.
- **Guardrails e Sandbox**: O agente opera com isolamento estrito no workspace, pré-autorizações de segurança (`.agent/settings.json`) e estrita observância a ferramentas nativas de arquivos.
- **Prioridade de Ferramentas Nativas**: Utilize as ferramentas nativas (`write_to_file`, `replace_file_content`, `view_file`, `grep_search`, `find_by_name`, `list_dir`) para manipulação de arquivos e código. Nunca execute comandos bash para ler, escrever ou buscar arquivos.
- **Uso Controlado de Ferramentas de Ciclo de Vida**: O agente consome alvos de automação (`make test`, `make lint`, `make check`, `make build`) estritamente como ferramentas de verificação dentro de seu harness operacional.

## 2. Loop Engineering (Ciclos Rápidos de Auto-Validação)
- **Ciclo de Feedback Contínuo**:
  1. *Inspeção*: Compreender os requisitos e contratos existentes via `view_file` (com limites de linhas).
  2. *Implementação*: Fazer alterações mínimas e focadas com `replace_file_content` ou `write_to_file`.
  3. *Teste Automatizado*: Rodar testes unitários/cobertura imediatamente (`make test`).
  4. *Diagnóstico & Correção*: Caso haja falhas, analisar logs de forma direcionada, corrigir e revalidar.
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

## 4. Token & Context Economy (Eficiência da Janela de Contexto)
- **Leitura Cirúrgica**: Use `StartLine` e `EndLine` no `view_file` para evitar ler arquivos extensos na íntegra.
- **Edição Cirúrgica**: Prefira `replace_file_content` para substituição de trechos específicos sem reescrever o arquivo completo.
- **Outputs e Respostas Concisas**: Responda com resumos objetivos e links no padrão `[arquivo](file:///caminho)`, sem duplicar blocos de código já salvos no disco.
- **Execução Enxuta**: Utilize comandos com saídas estruturadas e evite flags verbosas desnecessárias (ex: `-v` em testes).
