## Context

O projeto consolidou um estudo de ponta em **Spec-Driven Development (SDD)** utilizando o framework **OpenSpec**, disponibilizando uma suíte completa de prompts mestres genéricos de fundação arquitetural para múltiplos ecossistemas (**Go**, **Python**, **Node.js** e **TypeScript**) adaptados para **Antigravity CLI (`agy`)** e **Claude Code (`claude`)**. Para orientar desenvolvedores e documentar adequadamente a pasta `docs/`, é necessário introduzir um `docs/README.md` como guia e catálogo explicativo. Veja mais em [proposal.md](file:///home/douglas/Workspace/gemini/sftp-server/openspec/changes/generic-golang-foundation-prompt/proposal.md).

## Goals / Non-Goals

**Goals:**
- Criar `docs/README.md` contendo:
  * Explicação detalhada sobre o estudo de Spec-Driven Development (SDD) com OpenSpec.
  * Pilares de Engenharia de Agentes (Harness, Loop e Graph Engineering com Native Tool Grounding e Economia de Tokens).
  * Tabela completa e comparativa dos 8 arquivos de prompt mestres da pasta com links e resumos técnicos.
  * Guia rápido de utilização dos prompts no Antigravity CLI e no Claude Code.
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
- Adaptar o **Pilar 6** (Engenharia de Agentes, Harness e Ferramentas Nativas) para o Claude Code (`CLAUDE.md`, `.claude/settings.json`, ferramentas nativas `Write`, `Edit`, `View`, `GrepTool`, `GlobTool`, `LS`).
- Garantir consistência absoluta nos 10 pilares fundamentais, recomendação formal do framework **OpenSpec** (Pilar 7) e parametrização de plataformas alvo (`[PLATAFORMAS_ALVO]`) no Pilar 8 em todos os 8 documentos.

**Non-Goals:**
- Não alterar o código fonte ou binários em execução da aplicação `file-server`.
- Não diminuir a barreira de qualidade (cobertura >= 80% e `make check` continuam inegociáveis).

## Decisions

### Decisão 1: Catálogo Central de Documentação em `docs/README.md`
- O diretório `docs/` conterá um README próprio, atuando como o ponto de entrada para o estudo de Spec-Driven Development e índice com tabela comparativa de todos os prompts mestres disponíveis.

### Decisão 2: Matriz de 8 Arquivos Especializados (Linguagem x Harness)
- **Nomenclatura**: `generic-<linguagem>-foundation-spec-prompt-<harness>.md`.
- **Linguagens**: `golang`, `python`, `nodejs`, `typescript`.
- **Harnesses**:
  - `agy`: Otimizado para o Antigravity CLI (`AGENTS.md`, `GEMINI.md`, `.agent/rules/`, `.agent/settings.json`, ferramentas nativas `write_to_file`, `replace_file_content`, etc.).
  - `claude`: Otimizado para o Claude Code (`CLAUDE.md`, `.claude/settings.json`, ferramentas nativas `Write`, `Edit`, `View`, `GrepTool`, `GlobTool`, `LS`).

### Decisão 3: Adaptação Precisa do Pilar 6 para Claude Code
- No prompt do Claude Code, o Pilar 6 instrui o agente a produzir e respeitar:
  1. `CLAUDE.md` com guidelines de projeto, convenções de código, comandos rápidos (`make check`, `make test`) e regras inegociáveis.
  2. Uso prioritário das ferramentas nativas do Claude (`View`, `Edit`, `Write`, `GrepTool`, `GlobTool`, `LS`), evitando o uso do `Bash` para I/O básico em arquivos.
  3. Scaffolding de IA e Loop Engineering adaptado para o fluxo do Claude Code.

## Risks / Trade-offs

- **[Risco]** Manutenção de 8 arquivos de documentação e do README de catálogo.
  - *Mitigação*: Manter estrutura, seções e conteúdo dos pilares 1 a 5 e 7 a 10 perfeitamente alinhados, variando apenas o Pilar 6 e os artefatos de harness de IA (`AGENTS.md` vs `CLAUDE.md`).

## Migration Plan

1. Renomear os 4 arquivos existentes de prompt para `-agy.md`.
2. Criar os 4 novos arquivos de prompt para `-claude.md`.
3. Criar `docs/README.md` como catálogo e guia de SDD.
4. Validar conformidade OpenSpec (`openspec validate --all`) e executar o quality gate (`make check`).
