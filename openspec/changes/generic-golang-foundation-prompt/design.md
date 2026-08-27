## Context

O arquivo `docs/foundation-spec-prompt.md` foi concebido originalmente com base nos requisitos específicos do projeto `file-server`. A generalização do conceito de prompt mestre de fundação arquitetural revelou a necessidade de prover modelos equivalentes de excelência para os principais ecossistemas de desenvolvimento utilizados em projetos corporativos e open-source: **Go**, **Python**, **Node.js (JavaScript puro)** e **TypeScript**. Veja mais em [proposal.md](file:///home/douglas/Workspace/gemini/sftp-server/openspec/changes/generic-golang-foundation-prompt/proposal.md).

## Goals / Non-Goals

**Goals:**
- Manter o arquivo `docs/generic-golang-foundation-spec-prompt.md` para projetos Go.
- Manter e aprimorar `docs/generic-python-foundation-spec-prompt.md` adaptado ao ecossistema Python moderno (`src` layout, `pyproject.toml` mandatório, `uv` para gestão de Python/venv/dependências, `ty` para checagem de tipos ultraveloz, `pytest` com cobertura >= 80%, `ruff`, empacotamento wheel/sdist).
- Manter `docs/generic-nodejs-foundation-spec-prompt.md` adaptado ao ecossistema Node.js / JavaScript moderno (ESM nativo, `package.json`, `node:test` ou `vitest` com cobertura >= 80%, `eslint`, `prettier`, `npm audit`).
- Manter `docs/generic-typescript-foundation-spec-prompt.md` adaptado ao ecossistema TypeScript com tipagem estrita (`strict: true`, `tsconfig` NodeNext, `vitest`, `tsup`/`tsc`, `eslint-typescript`, `prettier`, `make typecheck`).
- Garantir que todos os 4 prompts contenham rigorosamente os **10 pilares fundamentais de engenharia**, a recomendação formal do framework **OpenSpec** (Pilar 7) com foco em PO e QA, parametrização de plataformas alvo (`[PLATAFORMAS_ALVO]`) no Pilar 8, Makefile universal autodocumentado (Pilar 4) e Engenharia de Agentes com Harness/Loop/Graph (Pilar 6).

**Non-Goals:**
- Não alterar o código de produção ou binários do `file-server`.
- Não flexibilizar a barreira de qualidade (cobertura >= 80% e `make check` continuam mandatórios em todos os prompts).

## Decisions

### Decisão 1: Suíte de 4 Arquivos de Prompt Especializados por Linguagem
- **Escolha**:
  1. `docs/generic-golang-foundation-spec-prompt.md` (Go)
  2. `docs/generic-python-foundation-spec-prompt.md` (Python)
  3. `docs/generic-nodejs-foundation-spec-prompt.md` (Node.js / JS)
  4. `docs/generic-typescript-foundation-spec-prompt.md` (TypeScript)
- **Justificativa**: Cada ecossistema possui ferramentas canônicas distintas. Ter um documento dedicado por linguagem oferece instruções 100% idiomáticas e prontas para uso por desenvolvedores e agentes de IA sem ambiguidades de sintaxe ou configuração.

### Decisão 2: Adaptação Idiomática e Modernização do Python com `uv`, `ty` e `pyproject.toml`
- **Escolha**:
  - `uv` (Astral) como ferramenta padrão mandatória para gerenciar versões do Python (`uv python`), ambientes virtuais isolados (`uv venv`), adição e sincronização determinística de dependências e lockfiles (`uv add`, `uv lock`, `uv sync`) e execução ágil (`uv run`).
  - `ty` (Astral) como ferramenta de checagem estática de tipos de alta performance em substituição ao tradicional `mypy`, acelerando drasticamente o ciclo de feedback no `make check` e no loop do agente.
  - `pyproject.toml` como manifesto central e mandatório para metadados (PEP 621), build system, scripts e configurações de ferramentas (`[tool.ruff]`, `[tool.pytest.ini_options]`, `[tool.ty]`).

### Decisão 3: Recomendação Universal do Framework OpenSpec e PO/QA (Pilar 7)
- Todos os 4 prompts recomendam formalmente o OpenSpec, fornecendo diretrizes de linguagem ubíqua e acessível para o PO e cenários determinísticos BDD/Gherkin com `#### Scenario:` e `WHEN/THEN` para automação por QA.

## Risks / Trade-offs

- **[Risco]** Ferramenta `ty` ser recente no ecossistema.
  - *Mitigação*: O prompt orienta a configuração compatível com os padrões de type hints do Python (PEP 484/526) e execução via `uv run ty` ou fallback direto caso o projeto necessite de flags específicas.

## Migration Plan

1. Atualizar o prompt mestre de Python em `docs/generic-python-foundation-spec-prompt.md` com `uv`, `ty` e `pyproject.toml`.
2. Validar conformidade OpenSpec (`openspec validate --all`) e executar o quality gate (`make check`).
