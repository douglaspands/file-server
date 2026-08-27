## Why

No arquivo `docs/foundation-spec-prompt.md`, o **Pilar 6** do Prompt Mestre misturou conceitos de arquitetura de software clássica (como camadas Go e portas/adaptadores) na definição de *Graph Engineering*. No paradigma de Inteligência Artificial e Sistemas Compostos por Agentes (Compound AI / Agentic Systems), **Graph Engineering** é um conceito intrínseco de IA que diz respeito à orquestração de grafos de execução e fluxos de raciocínio de agentes (State Graphs / DAG de raciocínio, subagentes especializados, decomposição topológica de tarefas de planejamento, transições de estado determinísticas e roteamento).

Esta mudança se faz necessária para revisar integralmente o Pilar 6, garantindo que todos os seus tópicos abordem com precisão e profundidade a **Engenharia de Agentes** (Harness Engineering, Loop Engineering, Graph Engineering, Autonomia e Grounding de Ferramentas Nativas e Economia Ativa de Tokens e Contexto), alinhando a especificação `project-foundation`, as regras em `.agent/rules/` e o documento mestre de fundação.

## What Changes

- **Revisão do Pilar 6 em `docs/foundation-spec-prompt.md`**:
  - Redefinição completa de **Graph Engineering** no escopo de IA de Agentes (orquestração de State Graphs, decomposição de tarefas em DAG, roteamento e transições determinísticas de estado, coordenação e paralelismo de subagentes).
  - Aprimoramento e consolidação de **Harness Engineering** (scaffolding de regras determinísticas, guardrails, matriz de permissões, sandbox e restrições operacionais para o agente).
  - Aprimoramento e consolidação de **Loop Engineering** (ciclos ReAct, reflexão e auto-correção iterativa: inspeção -> intervenção mínima -> teste automatizado -> análise de erro -> correção e validação final).
  - Reforço nas diretrizes de **Prioridade Mandatória de Ferramentas Nativas** e **Uso Restrito de Terminal** (grounding estrito em tool calling estruturado vs comandos shell).
  - Reforço nas diretrizes de **Economia Ativa de Tokens e Otimização de Janela de Contexto** (leitura/edição cirúrgica de linhas, links markdown clicáveis, ausência de código redundante e saídas de comando enxutas).
- **Atualização da Especificação Principal `project-foundation`**:
  - Atualização do requisito e cenários de *Estratégias do Antigravity CLI (Harness, Loop e Graph Engineering) e Otimização de Tokens* para refletir Graph Engineering puramente como engenharia de agentes/IA.
- **Atualização da Regra Operacional `.agent/rules/agent_harness_engineering.md`**:
  - Ajuste da seção de Graph Engineering para focar em orquestração de grafos de tarefas/raciocínio, fluxo de estados de agentes e coordenação de subagentes.

## Capabilities

### Modified Capabilities

- `project-foundation`: Atualização do requisito de governança de IA e Antigravity CLI para alinhar a definição de Graph Engineering como conceito de IA de agentes (DAG de raciocínio/tarefas, orquestração de subagentes e transições de estado) e consolidar as boas práticas de engenharia de agentes.

## Impact

- **Documentação de Fundação**: Atualização do arquivo `docs/foundation-spec-prompt.md`.
- **Especificações OpenSpec**: Atualização do delta `specs/project-foundation/spec.md`.
- **Regras de Agente**: Alinhamento do arquivo `.agent/rules/agent_harness_engineering.md`.
- **Compatibilidade**: Nenhuma quebra de compatibilidade no código-fonte Go existente; as melhorias são conceituais, operacionais e de governança para agentes de IA.
