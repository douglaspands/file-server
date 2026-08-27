## Context

O documento `docs/foundation-spec-prompt.md` atua como a referência mestre para inicialização e governança de projetos no ecossistema Antigravity CLI e OpenSpec. O **Pilar 6** deste prompt condensa os requisitos de Engenharia de Agentes e Operação do AGY. Anteriormente, o conceito de *Graph Engineering* no Pilar 6 e na regra `.agent/rules/agent_harness_engineering.md` foi incorretamente descrito com referências a camadas de arquitetura de software em Go (`ports -> domain -> services -> adapters -> cmd`).

No estado da arte de Inteligência Artificial e Sistemas Compostos por Agentes (Compound AI Systems), **Graph Engineering** é uma disciplina fundamental de IA referente à modelagem, orquestração e resolução de fluxos de execução como grafos de estado (State Graphs) e grafos acíclicos dirigidos (DAGs) de tarefas, raciocínio e subagentes.

## Goals / Non-Goals

**Goals:**
- Revisar integralmente o **Pilar 6** em `docs/foundation-spec-prompt.md` para cobrir de forma aprofundada e precisa a **Engenharia de Agentes**:
  1. *Prioridade Mandatória de Ferramentas Nativas & Tool Grounding*.
  2. *Uso Restrito do Terminal para Ciclo de Vida e Automação*.
  3. *Autonomia Operacional e Matriz de Permissões Pré-Autorizadas*.
  4. *Economia Ativa de Tokens, Curação de Contexto e Otimização de Janela*.
  5. *Disciplinas de Sistemas Compostos de IA: Harness Engineering, Loop Engineering e Graph Engineering*.
- Redefinir **Graph Engineering** no contexto de agentes de IA: modelagem de State Graphs, decomposição topológica de tarefas de planejamento e execução, controle de transições determinísticas de estado e orquestração de subagentes especializados.
- Alinhar o arquivo de regras operacionais `.agent/rules/agent_harness_engineering.md` com a nova definição de Graph Engineering.
- Garantir validação estrita no OpenSpec (`openspec validate --all`).

**Non-Goals:**
- Não alterar o código de aplicação Go nem as camadas arquiteturais da aplicação `file-server`.
- Não alterar os demais pilares de software de `docs/foundation-spec-prompt.md` (Pilares 1 a 5 e 7 a 10 mantêm sua integridade intacta).

## Decisions

### Decisão 1: Estruturação Robusta dos 5 Eixos do Pilar 6 em Engenharia de Agentes
- **Eixo 1 - Prioridade Mandatória de Ferramentas Nativas (Native Tool Grounding)**: O agente deve interagir com o workspace estritamente via ferramentas estruturadas (`write_to_file`, `replace_file_content`, `view_file`, `grep_search`, `find_by_name`, `list_dir`), eliminando I/O não estruturado via bash.
- **Eixo 2 - Uso Restrito de Terminal (`run_command`)**: Terminal restrito à execução de ferramentas de automação e ciclo de vida (`make`, `go`, `git`, `openspec`).
- **Eixo 3 - Autonomia Operacional & Matriz de Permissões**: Eliminação de atrito com pré-autorização de ferramentas de desenvolvimento em `.agent/settings.json`.
- **Eixo 4 - Economia Ativa de Tokens e Otimização de Janela de Contexto**: Leitura/edição cirúrgica com ranges de linhas, eliminação de duplicações de código no chat, resumos estruturados com links `[arquivo](file:///...)` e comandos enxutos.
- **Eixo 5 - Tríade de Engenharia de Agentes (Harness, Loop & Graph Engineering)**:
  - *Harness Engineering*: Scaffolding de regras determinísticas, guardrails, context curation e sandbox de segurança para governar a operação autônoma.
  - *Loop Engineering*: Ciclos ReAct / Reflection de auto-validação iterativa (inspecionar -> intervir pontualmente -> testar -> diagnosticar -> corrigir -> validar).
  - *Graph Engineering*: Modelagem do fluxo do agente como State Graphs / DAG de raciocínio, decomposição topológica de tarefas, transições determinísticas de estado e orquestração de subagentes.

*Alternativas consideradas*:
- Manter Graph Engineering focado na ordem de criação de pacotes Go: descartado, pois trata-se de arquitetura de software clássica e não de engenharia de agentes de IA.

### Decisão 2: Atualização Coerente de `.agent/rules/agent_harness_engineering.md`
- Atualizar a seção 3 do arquivo de regras para que reflita a orquestração de grafos de tarefas/raciocínio e transições de estado pelo Antigravity CLI.

## Risks / Trade-offs

- **[Risco] Divergência entre documentação mestre e regras locais**:
  - *Mitigação*: Atualização atômica em conjunto de `docs/foundation-spec-prompt.md` e `.agent/rules/agent_harness_engineering.md`.
