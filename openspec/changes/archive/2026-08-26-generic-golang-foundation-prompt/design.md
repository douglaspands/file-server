## Context

O projeto consolidou um estudo de ponta em **Spec-Driven Development (SDD)** utilizando o framework **OpenSpec**, disponibilizando uma suíte completa de prompts mestres genéricos de fundação arquitetural para múltiplos ecossistemas (**Go**, **Python**, **Node.js** e **TypeScript**) adaptados para **Antigravity CLI (`agy`)** e **Claude Code (`claude`)**. Para orientar desenvolvedores e documentar adequadamente a pasta `docs/`, foi introduzido um `docs/README.md` como guia e catálogo explicativo. Para elevar a segurança, integridade e higiene dos repositórios gerados pelos prompts, o **Pilar 10 (Governança Git)** deve ser aprofundado com ênfase na criação mandatória do `.gitignore` e nas práticas de versionamento com agentes de IA. Veja mais em [proposal.md](file:///home/douglas/Workspace/gemini/sftp-server/openspec/changes/generic-golang-foundation-prompt/proposal.md).

## Goals / Non-Goals

**Goals:**
- Aprofundar o **Pilar 10** em todos os 8 prompts mestres e no `docs/README.md`, incluindo:
  * Criação mandatória e estruturada do arquivo `.gitignore` idiomático por ecossistema (segredos, binários, caches, relatórios de cobertura, dependências e IDEs).
  * Explicitar a importância e justificativa da governança Git na engenharia com agentes de IA (segurança, higiene de contexto, histórico linear/bisectável na `main`, branches isoladas e merge squash com autorização).
- Manter o catálogo e índice em `docs/README.md` alinhado e atualizado com as boas práticas do Pilar 10.
- Padronizar os nomes dos 4 prompts do Antigravity CLI com o sufixo `-agy.md`:
  * `docs/generic-golang-foundation-spec-prompt-agy.md`
  * `docs/generic-python-foundation-spec-prompt-agy.md`
  * `docs/generic-nodejs-foundation-spec-prompt-agy.md`
  * `docs/generic-typescript-foundation-spec-prompt-agy.md`
- Criar os 4 prompts equivalentes otimizados para o **Claude Code** com o sufixo `-claude.md`:
  * `docs/generic-golang-foundation-spec-prompt-claude.md`
  * `docs/generic-python-foundation-spec-prompt-claude.md`
  * `docs/generic-nodejs-foundation-spec-prompt-claude.md`
  * `docs/generic-typescript-foundation-spec-prompt-claude.md`

**Non-Goals:**
- Não alterar o código fonte ou binários em execução da aplicação `file-server`.
- Não diminuir a barreira de qualidade (cobertura >= 80% e `make check` continuam inegociáveis).

## Decisions

### Decisão 1: Aprofundamento do Pilar 10 (Governança Git e `.gitignore`)
- Todos os 8 prompts exigirão a criação inicial de um `.gitignore` completo cobrindo:
  1. Segredos e arquivos sensíveis (`.env`, certificados, chaves privadas).
  2. Caches e artefatos de compilação/build/distribuição.
  3. Relatórios de cobertura e profiling de testes.
  4. Diretórios de dependências (`.venv`, `node_modules`, `vendor`).
  5. Arquivos de configuração de IDEs e do sistema operacional.
- O prompt detalhará as justificativas técnicas e operacionais para a governança Git estrita em projetos auxiliados por agentes de IA.

### Decisão 2: Matriz de 8 Arquivos Especializados (Linguagem x Harness)
- **Nomenclatura**: `generic-<linguagem>-foundation-spec-prompt-<harness>.md`.
- **Linguagens**: `golang`, `python`, `nodejs`, `typescript`.
- **Harnesses**: `agy` e `claude`.

## Risks / Trade-offs

- **[Risco]** Manutenção da consistência entre 8 prompts e o catálogo README.
  - *Mitigação*: Padronizar o texto do Pilar 10 com variações estritamente idiomáticas por ecossistema.

## Migration Plan

1. Atualizar o Pilar 10 em todos os 8 arquivos de prompt mestres em `docs/`.
2. Atualizar a seção do Pilar 10 no `docs/README.md`.
3. Validar conformidade OpenSpec (`openspec validate --all`) e executar o quality gate (`make check`).
