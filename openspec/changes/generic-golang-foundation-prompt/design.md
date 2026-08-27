## Context

O projeto definiu uma suíte de prompts mestres genéricos de fundação arquitetural para múltiplos ecossistemas (**Go**, **Python**, **Node.js** e **TypeScript**). Para atender desenvolvedores e equipes que utilizam diferentes ferramentas de linha de comando de IA para pair programming e automação (especialmente **Antigravity CLI** e **Claude Code**), é necessário prover arquivos dedicados e nomeados com clareza com relação à linguagem e ao harness utilizado (`-agy.md` e `-claude.md`). Veja mais em [proposal.md](file:///home/douglas/Workspace/gemini/sftp-server/openspec/changes/generic-golang-foundation-prompt/proposal.md).

## Goals / Non-Goals

**Goals:**
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
- Adaptar o **Pilar 6** (Engenharia de Agentes, Harness e Ferramentas Nativas) para o Claude Code:
  * Referência e geração de `CLAUDE.md` na raiz (com instruções de build, test, lint, governança e comandos) e `.claude/settings.json`.
  * Mapeamento de ferramentas nativas do Claude Code (`Write`, `Edit`, `View`, `GrepTool`, `GlobTool`, `LS`) e uso estritamente restrito da ferramenta `Bash` para ferramentas de ciclo de vida (`make`, `git`, `openspec`).
- Garantir consistência absoluta nos 10 pilares fundamentais, recomendação formal do framework **OpenSpec** (Pilar 7) e parametrização de plataformas alvo (`[PLATAFORMAS_ALVO]`) no Pilar 8 em todos os 8 documentos.

**Non-Goals:**
- Não alterar o código fonte ou binários em execução da aplicação `file-server`.
- Não diminuir a barreira de qualidade (cobertura >= 80% e `make check` continuam inegociáveis).

## Decisions

### Decisão 1: Matriz de 8 Arquivos Especializados (Linguagem x Harness)
- **Nomenclatura**: `generic-<linguagem>-foundation-spec-prompt-<harness>.md`.
- **Linguagens**: `golang`, `python`, `nodejs`, `typescript`.
- **Harnesses**:
  - `agy`: Otimizado para o Antigravity CLI (`AGENTS.md`, `GEMINI.md`, `.agent/rules/`, `.agent/settings.json`, ferramentas nativas `write_to_file`, `replace_file_content`, etc.).
  - `claude`: Otimizado para o Claude Code (`CLAUDE.md`, `.claude/settings.json`, ferramentas nativas `Write`, `Edit`, `View`, `GrepTool`, `GlobTool`, `LS`).

### Decisão 2: Adaptação Precisa do Pilar 6 para Claude Code
- No prompt do Claude Code, o Pilar 6 instrui o agente a produzir e respeitar:
  1. `CLAUDE.md` com guidelines de projeto, convenções de código, comandos rápidos (`make check`, `make test`) e regras inegociáveis.
  2. Uso prioritário das ferramentas nativas do Claude (`View`, `Edit`, `Write`, `GrepTool`, `GlobTool`, `LS`), evitando o uso do `Bash` para I/O básico em arquivos.
  3. Scaffolding de IA e Loop Engineering adaptado para o fluxo do Claude Code.

## Risks / Trade-offs

- **[Risco]** Manutenção de 8 arquivos de documentação.
  - *Mitigação*: Manter estrutura, seções e conteúdo dos pilares 1 a 5 e 7 a 10 perfeitamente alinhados, variando apenas o Pilar 6 e os artefatos de harness de IA (`AGENTS.md` vs `CLAUDE.md`).

## Migration Plan

1. Renomear os 4 arquivos existentes de prompt para `-agy.md`.
2. Criar os 4 novos arquivos de prompt para `-claude.md`.
3. Validar conformidade OpenSpec (`openspec validate --all`) e executar o quality gate (`make check`).
