## Context

Consulte `proposal.md` para a motivação. Atualmente, o arquivo `docs/foundation-spec-prompt.md` consolida o padrão de excelência de fundação de projetos em 10 pilares. O Pilar 4 já é dedicado integralmente à automação universal de comandos via `Makefile` (`make help`, `make test`, `make lint`, `make check`, `make build`, etc.). No entanto, o Pilar 6 ("Engenharia de Agentes, Prioridade de Ferramentas e Economia de Tokens") continha uma associação inadequada de Harness Engineering com a interface do `Makefile`.

Em Engenharia de Inteligência Artificial e Agentes Autônomos (como o Antigravity CLI / AGY), o **Harness** é o ambiente de sustentação, regras (`rules`), guardrails de segurança, sandboxing, gerenciamento de contexto, prioridade de ferramentas nativas e matriz de permissões que governa a execução confiável e autônoma do modelo de IA. A automação do projeto é apenas uma das interfaces consumidas pelo agente dentro de seu harness operacional.

## Goals / Non-Goals

**Goals:**
- Dissociar conceitualmente e textualmente o termo *Harness Engineering* da interface do `Makefile` em todos os documentos e especificações do projeto.
- Consolidar o Pilar 6 do prompt mestre (`docs/foundation-spec-prompt.md`) como a referência definitiva de boas práticas de IA, governança do Antigravity CLI, uso otimizado de contexto, guardrails e operação segura.
- Atualizar a diretriz `.agent/rules/agent_harness_engineering.md` para refletir as definições corretas de Harness de IA (scaffolding, guardrails, contexto e permissões).
- Atualizar as referências na especificação `openspec/specs/project-foundation/spec.md` e na documentação técnica (`README.md`, `AGENTS.md`, `GEMINI.md`).

**Non-Goals:**
- Não serão alterados os alvos ou comportamentos do `Makefile` (que permanecem sob a responsabilidade de automação do Pilar 4).
- Não serão removidas ferramentas nativas ou permissões já existentes em `.agent/settings.json`.
- Não serão alteradas lógicas Go ou funcionalidades de aplicação do servidor de arquivos.

## Decisions

### 1. Separação Conceitual: Automação da Aplicação vs Harness de Agentes de IA
- **Decisão**: Manter o `Makefile` estritamente como a interface universal de automação do ciclo de vida da aplicação (Pilar 4) e redefinir o *Harness Engineering* (Pilar 6) como a arquitetura de controle e scaffolding de IA que garante execução segura, determinística e de baixo custo de tokens para o Antigravity CLI.
- **Alternativa Considerada**: Manter o termo com dupla interpretação ("harness de build" e "harness de IA"). *Rejeitada* pois gera confusão semântica e dilui o propósito do Pilar 6, que é exclusivamente voltado para engenharia de agentes.

### 2. Definição Estruturada dos Componentes de Engenharia de IA no Pilar 6
- **Decisão**: Estruturar o Pilar 6 em torno dos quatro eixos canônicos de governança de IA:
  1. *AI Harness Engineering*: Scaffolding de regras (`.agent/rules/`, `AGENTS.md`, `GEMINI.md`), parametrização de permissões (`.agent/settings.json`), guardrails e restrições de sandbox.
  2. *Tooling Autonomy & Native Priority*: Prioridade mandatória de ferramentas nativas (`write_to_file`, `replace_file_content`, `view_file`, `grep_search`, `find_by_name`, `list_dir`) e restrição de terminal a comandos de ciclo de vida.
  3. *Token & Context Economy*: Otimizações cirúrgicas de leitura/edição e saídas concisas em Markdown com links clicáveis.
  4. *Loop & Graph Engineering*: Loops de auto-validação contínua do agente e execução topológica em grafo acíclico direcionado (DAG).

### 3. Alinhamento Transversal dos Artefatos
- **Decisão**: Atualizar sincronicamente `docs/foundation-spec-prompt.md`, `.agent/rules/agent_harness_engineering.md`, `openspec/specs/project-foundation/spec.md` e cabeçalhos em `README.md` durante a fase de aplicação para garantir consistência terminológica em 100% do repositório.

## Risks / Trade-offs

- **[Risco] Divergência entre especificações legadas arquivadas e a nova definição**:
  - *Mitigação*: Os arquivos sob `openspec/changes/archive/` são registros históricos imutáveis; a atualização focará no prompt mestre vivo (`docs/foundation-spec-prompt.md`), nas regras ativas (`.agent/rules/`), nas especificações consolidadas (`openspec/specs/`) e na documentação atual (`README.md`).
