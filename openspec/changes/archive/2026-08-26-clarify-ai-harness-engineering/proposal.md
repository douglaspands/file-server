## Why

No documento `docs/foundation-spec-prompt.md` (Pilar 6) e em artefatos derivados de governança de IA (`.agent/rules/agent_harness_engineering.md`, especificação `project-foundation`), o conceito de **Harness Engineering** foi erroneamente associado ao `Makefile` e automação de comandos. Na engenharia de software e inteligência artificial, o `Makefile` é estritamente uma interface de automação do ciclo de vida da aplicação (coberta pelo Pilar 4), enquanto o **Harness** no contexto de agentes de IA refere-se ao ecossistema de scaffolding, regras operacionais, restrições de ferramentas, configurações de permissões, sandbox e guardrails de contexto que direcionam o comportamento do modelo (Antigravity CLI / AGY) de forma segura, determinística e otimizada.

Esta mudança faz-se necessária agora para alinhar conceitualmente o prompt mestre de fundação e todas as diretrizes do projeto aos princípios reais de Engenharia de IA, garantindo que o Pilar 6 trate exclusivamente de boas práticas de IA e governança do agente AGY.

## What Changes

- **Redefinição do Pilar 6 em `docs/foundation-spec-prompt.md`**: Dissociar Harness Engineering do Makefile (que já pertence ao Pilar 4) e redefinir Harness Engineering como o arcabouço de controle, scaffolding, guardrails, governança de contexto e parametrização segura para o Antigravity CLI.
- **Revisão da regra `.agent/rules/agent_harness_engineering.md`**: Atualizar a seção de Harness Engineering para focar em scaffolding de IA, guardrails, restrições nativas e isolamento seguro de execução do agente, mantendo a referência a comandos apenas como ferramentas consumidas dentro do harness.
- **Atualização da especificação `project-foundation`**: Ajustar os requisitos e cenários de BDD em `openspec/specs/project-foundation/spec.md` para refletir com precisão a distinção entre a interface de comandos/automação (Makefile) e a engenharia de harness de IA (governança, scaffolding e restrições do AGY).
- **Ajustes de coerência na documentação (`README.md`, `AGENTS.md`, `GEMINI.md`)**: Harmonizar referências textuais onde o termo "harness" possa estar incorretamente aplicado à automação tradicional, preservando a clareza conceitual em todo o repositório.

## Capabilities

### New Capabilities
<!-- Nenhuma nova capability funcional criada; trata-se de refinamento de governança e fundação -->

### Modified Capabilities
- `project-foundation`: Atualizar a definição e cenários de teste/especificação do requisito de Estratégias do Antigravity CLI (Harness, Loop e Graph Engineering), dissociando Harness Engineering de automação de Makefile e consolidando-o como arcabouço de governança e scaffolding de IA.

## Impact

- **Documentação de Fundação (`docs/foundation-spec-prompt.md`)**: Atualização do Pilar 6 do prompt mestre e do checklist de artefatos.
- **Diretrizes de Agente (`.agent/rules/agent_harness_engineering.md`)**: Ajuste conceitual na seção 1 (Harness Engineering de IA).
- **Especificações OpenSpec (`openspec/specs/project-foundation/spec.md`)**: Atualização dos cenários de aceitação BDD referentes a Harness Engineering.
- **Documentação do Projeto (`README.md`)**: Revisão de cabeçalhos e termos que misturavam "harness de automação" com ferramentas de build.
