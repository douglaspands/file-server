## 1. Revisão do Pilar 6 em docs/foundation-spec-prompt.md

- [x] 1.1 Redefinir Graph Engineering no Pilar 6 de `docs/foundation-spec-prompt.md` como modelagem de State Graphs, decomposição de tarefas em DAG, controle de transições determinísticas de estado e orquestração de subagentes
- [x] 1.2 Revisar e consolidar as seções de Harness Engineering, Loop Engineering, Prioridade de Ferramentas Nativas, Autonomia e Economia de Tokens no Pilar 6 de `docs/foundation-spec-prompt.md`
- [x] 1.3 Atualizar o checklist final de artefatos em `docs/foundation-spec-prompt.md` para refletir as diretrizes aprimoradas de engenharia de agentes

## 2. Alinhamento das Diretrizes em .agent/rules/

- [x] 2.1 Atualizar a seção de Graph Engineering em `.agent/rules/agent_harness_engineering.md` para focar em orquestração de DAG de raciocínio/tarefas, State Graphs e subagentes
- [x] 2.2 Revisar a coerência geral de `.agent/rules/agent_harness_engineering.md` com as práticas de economia e autonomia

## 3. Validação e Quality Gate

- [x] 3.1 Executar `openspec validate --all` para certificar a conformidade e integridade das especificações
- [x] 3.2 Executar `make check` para assegurar integridade completa do projeto (formatação, linters e testes com cobertura >= 80%)
