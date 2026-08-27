## Context

O arquivo `docs/foundation-spec-prompt.md` foi concebido originalmente com base nos requisitos específicos do projeto `file-server`. A generalização do conceito de prompt mestre de fundação arquitetural revelou a necessidade de prover modelos equivalentes de excelência para os principais ecossistemas de desenvolvimento utilizados em projetos corporativos e open-source: **Go**, **Python**, **Node.js (JavaScript puro)** e **TypeScript**. Veja mais em [proposal.md](file:///home/douglas/Workspace/gemini/sftp-server/openspec/changes/generic-golang-foundation-prompt/proposal.md).

## Goals / Non-Goals

**Goals:**
- Manter o arquivo `docs/generic-golang-foundation-spec-prompt.md` para projetos Go.
- Criar `docs/generic-python-foundation-spec-prompt.md` adaptado ao ecossistema Python (`src` layout, `pyproject.toml`, `pytest` com cobertura >= 80%, `ruff`, `mypy`, Type Hints, empacotamento wheel/sdist).
- Criar `docs/generic-nodejs-foundation-spec-prompt.md` adaptado ao ecossistema Node.js / JavaScript moderno (ESM nativo, `package.json`, `node:test` ou `vitest` com cobertura >= 80%, `eslint`, `prettier`, `npm audit`).
- Criar `docs/generic-typescript-foundation-spec-prompt.md` adaptado ao ecossistema TypeScript com tipagem estrita (`strict: true`, `tsconfig` NodeNext, `vitest`, `tsup`/`tsc`, `eslint-typescript`, `prettier`, `make typecheck`).
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
- **Justificativa**: Cada ecossistema possui ferramentas canônicas distintas (Go toolchain, Python/ruff/pytest, Node/npm, TypeScript/tsc/tsup). Ter um documento dedicado por linguagem oferece instruções 100% idiomáticas e prontas para uso por desenvolvedores e agentes de IA sem ambiguidades de sintaxe ou configuração.

### Decisão 2: Adaptação Idiomática dos 10 Pilares
- **Python**:
  - Clean Architecture em `src/<pacote>/` (`core/domain/`, `ports/`, `services/`, `adapters/`), tipagem estrita com Type Hints e Protocols (`typing.Protocol`), `pyproject.toml`, `pytest` + `pytest-cov` >= 80%, `ruff`, `mypy`, Makefile universal (`make check` = fmt + lint + typecheck + security + test).
- **Node.js (JavaScript puro)**:
  - Clean Architecture em `src/` (`core/domain/`, `ports/`, `services/`, `adapters/`), ESM nativo (`"type": "module"`), `node:test` ou `vitest` com cobertura >= 80%, `eslint`, `prettier`, `package.json` (`"bin"`, `"exports"`), Makefile universal.
- **TypeScript**:
  - Clean Architecture em `src/` (`core/domain/`, `ports/`, `services/`, `adapters/`), tipagem estrita sem `any` (`strict: true`, `NodeNext`), `vitest` com cobertura >= 80%, `eslint` (`@typescript-eslint`), `prettier`, `tsup`/`tsc`, Makefile universal (`make typecheck`, `make check`).

### Decisão 3: Recomendação Universal do Framework OpenSpec e PO/QA (Pilar 7)
- Todos os 4 prompts recomendam formalmente o OpenSpec, fornecendo diretrizes de linguagem ubíqua e acessível para o PO e cenários determinísticos BDD/Gherkin com `#### Scenario:` e `WHEN/THEN` para automação por QA.

## Risks / Trade-offs

- **[Risco]** Divergência entre os 10 pilares nos 4 arquivos de documentação.
  - *Mitigação*: Manter estrutura, seções e numeração de pilares rigorosamente idênticas em todos os documentos, adaptando apenas termos específicos da linguagem (ex: pacotes vs módulos, pytest vs testing vs vitest, golangci-lint vs ruff vs eslint).

## Migration Plan

1. Criar os prompts mestres para Python, Node.js e TypeScript.
2. Manter e validar o prompt de Golang.
3. Validar conformidade OpenSpec (`openspec validate --all`) e executar o quality gate (`make check`).
