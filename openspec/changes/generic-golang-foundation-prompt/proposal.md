## Why

O arquivo `docs/foundation-spec-prompt.md` original possuía acoplamentos rígidos e específicos a certas tecnologias e frameworks em Go, limitando sua utilidade como gerador de fundações de software para novos projetos em múltiplos ecossistemas.

É necessário disponibilizar uma suíte completa de prompts mestres genéricos de fundação arquitetural e de engenharia de software para as principais linguagens e ecossistemas da indústria: **Golang**, **Python**, **Node.js (JavaScript puro / ESM)** e **TypeScript**. Cada prompt preserva integralmente os 10 pilares inegociáveis de excelência (Clean Architecture/layout canônico, TDD/BDD com cobertura >= 80%, Makefile universal, CI/CD, documentação viva no README, governança Git com squash, recomendação formal do framework OpenSpec com foco em PO/QA e Engenharia de Agentes com Harness, Loop e Graph), adaptando as convenções, ferramentas de build, linters e runners de teste de forma idiomática para cada linguagem.

No ecossistema Python, é essencial recomendar ferramentas modernas de alta performance da nova geração: **`uv`** (Astral) como gerenciador unificado de versões do Python, ambientes virtuais, dependências e lockfiles; **`ty`** como checador estático de tipos ultraveloz (substituindo o legadomypy); e a padronização mandatória do **`pyproject.toml`** como manifesto central de configuração e dependências do projeto.

## What Changes

- **Renomeação do Arquivo Go**: Renomear `docs/foundation-spec-prompt.md` para `docs/generic-golang-foundation-spec-prompt.md` e generalizar seus 10 pilares.
- **Criação e Refinamento do Prompt Mestre para Python**:
  - Criar `docs/generic-python-foundation-spec-prompt.md` adaptado ao ecossistema Python moderno (`src` layout, `pyproject.toml` mandatório, `uv` para gestão de Python/venv/dependências, `ty` para checagem estática de tipos de alta performance, `pytest` com cobertura >= 80%, `ruff`, empacotamento wheels/sdist e Makefile universal).
- **Criação do Prompt Mestre para Node.js (JavaScript ESM)**:
  - Criar `docs/generic-nodejs-foundation-spec-prompt.md` adaptado a JavaScript moderno (`"type": "module"`, `package.json`, `node:test` ou `vitest` com cobertura >= 80%, `eslint`, `prettier`, `npm audit` e Makefile universal).
- **Criação do Prompt Mestre para TypeScript**:
  - Criar `docs/generic-typescript-foundation-spec-prompt.md` adaptado a TypeScript com tipagem estrita (`strict: true`, `tsconfig` NodeNext, `vitest`, `tsup`/`tsc`, `eslint-typescript`, `prettier`, verificação estática de tipos `make typecheck` e Makefile universal).
- **Preservação e Governança Universal dos 10 Pilares em Todas as Linguagens**:
  - Clean Architecture / Ports & Adapters com isolamento de regras de domínio no core.
  - Modularidade de entrada e versionamento dinâmico.
  - Testes automatizados TDD/BDD com barreira >= 80% de cobertura.
  - Makefile universal e determinístico (`make check`, `make test`, `make build`, etc.).
  - Stack customizável deliberada com o solicitante do prompt (sem imposições pré-fabricadas).
  - Engenharia de Agentes (Harness, Loop e Graph Engineering, Native Tool Grounding e Economia de Tokens).
  - Recomendação formal do framework **OpenSpec** com visão conjunta de negócio para o **PO** e cenários BDD/Gherkin para automação por **QA**.
  - Pipelines de CI/CD e Matriz de Release customizáveis (`[PLATAFORMAS_ALVO]`) com fallback seguro.
  - README vivo e exaustivo.
  - Governança Git com Conventional Commits, Feature Branches e Merge Squash.

## Capabilities

### New Capabilities
<!-- Nenhuma nova capacidade introduzida. -->

### Modified Capabilities
- `project-foundation`: Atualização dos requisitos do prompt mestre de fundação arquitetural para disponibilizar a suíte de prompts universais e agnósticos para Go, Python (com `uv`, `ty` e `pyproject.toml`), Node.js e TypeScript, todos recomendando o framework OpenSpec para governança de especificações.

## Impact

- **Documentação**: São disponibilizados 4 prompts mestres universais em `docs/`: `generic-golang-foundation-spec-prompt.md`, `generic-python-foundation-spec-prompt.md`, `generic-nodejs-foundation-spec-prompt.md` e `generic-typescript-foundation-spec-prompt.md`.
- **Especificações OpenSpec**: `specs/project-foundation/spec.md` passa a referenciar os prompts multilíngues genéricos.
- **Código / Aplicação**: Nenhum impacto direto no código fonte em Go ou nos binários em execução da aplicação `file-server`.
