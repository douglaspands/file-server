## Why

O arquivo `docs/foundation-spec-prompt.md` original possuía acoplamentos rígidos e específicos a certas tecnologias e frameworks em Go, limitando sua utilidade como gerador de fundações de software para novos projetos em múltiplos ecossistemas e diferentes ferramentas de IA de linha de comando.

Este projeto consolidou um estudo aprofundado de **Spec-Driven Development (SDD)** utilizando o framework **OpenSpec**. É necessário disponibilizar uma suíte completa de prompts mestres genéricos de fundação arquitetural e de engenharia de software para as principais linguagens e ecossistemas da indústria (**Golang**, **Python**, **Node.js (JavaScript puro / ESM)** e **TypeScript**), com versões dedicadas e otimizadas tanto para o **Antigravity CLI (`agy`)** quanto para o **Claude Code (`claude`)**, acompanhadas de um `docs/README.md` que catalogue e explique a finalidade de cada prompt e as decisões arquiteturais do estudo. Cada prompt preserva integralmente os 10 pilares inegociáveis de excelência (Clean Architecture/layout canônico, TDD/BDD com cobertura >= 80%, Makefile universal, CI/CD, documentação viva no README, governança Git aprofundada com `.gitignore` e squash, recomendação formal do framework OpenSpec com foco em PO/QA e Engenharia de Agentes com Harness, Loop e Graph).

## What Changes

- **Catálogo Central da Pasta Docs (`docs/README.md`)**:
  - Criar `docs/README.md` documentando o estudo de Spec-Driven Development (SDD), a metodologia de Engenharia de Agentes e a tabela comparativa e descritiva dos 8 prompts mestres.
- **Padronização de Nomenclatura com Sufixos de Harness (`-agy.md` e `-claude.md`)**:
  - Renomear os prompts existentes do Antigravity para `docs/generic-<linguagem>-foundation-spec-prompt-agy.md`.
  - Criar os prompts equivalentes para Claude Code sob `docs/generic-<linguagem>-foundation-spec-prompt-claude.md`.
- **Suíte Completa de 8 Prompts Mestres**:
  1. `docs/generic-golang-foundation-spec-prompt-agy.md` (Go + Antigravity)
  2. `docs/generic-golang-foundation-spec-prompt-claude.md` (Go + Claude Code)
  3. `docs/generic-python-foundation-spec-prompt-agy.md` (Python com `uv`, `ty`, `pyproject.toml` + Antigravity)
  4. `docs/generic-python-foundation-spec-prompt-claude.md` (Python com `uv`, `ty`, `pyproject.toml` + Claude Code)
  5. `docs/generic-nodejs-foundation-spec-prompt-agy.md` (Node.js ESM + Antigravity)
  6. `docs/generic-nodejs-foundation-spec-prompt-claude.md` (Node.js ESM + Claude Code)
  7. `docs/generic-typescript-foundation-spec-prompt-agy.md` (TypeScript strict + Antigravity)
  8. `docs/generic-typescript-foundation-spec-prompt-claude.md` (TypeScript strict + Claude Code)
- **Aprofundamento do Pilar 10 (Governança Git, Higiene de Repositório e `.gitignore`)**:
  - Mandato explícito para criação de arquivo `.gitignore` idiomático e completo (segredos `.env`, chaves criptográficas, binários, diretórios de build `dist/`, caches de linters/testes, relatórios de cobertura, ambientes virtuais e `node_modules`).
  - Justificativa da importância da governança Git na engenharia com agentes de IA (segurança, higiene de contexto para ferramentas de busca, histórico linear/bisectável na `main`, branches de feature isoladas e integração estritamente via Squash Merge com permissão do usuário).
- **Adaptação do Pilar 6 (Engenharia de Agentes e Harness)**:
  - **Versões AGY**: Governança via `AGENTS.md`, `GEMINI.md`, `.agent/rules/`, `.agent/settings.json` e ferramentas nativas (`write_to_file`, `replace_file_content`, `view_file`, `grep_search`, `find_by_name`, `list_dir`).
  - **Versões Claude Code**: Governança via `CLAUDE.md`, `.claude/settings.json` e ferramentas nativas do Claude (`Write`, `Edit`, `View`, `GrepTool`, `GlobTool`, `LS`), restringindo `Bash` a comandos do ciclo de vida.
- **Preservação e Governança Universal dos 10 Pilares em Todas as Linguagens e Harnesses**:
  - Clean Architecture / Ports & Adapters com isolamento de regras de domínio no core.
  - Modularidade de entrada e versionamento dinâmico.
  - Testes automatizados TDD/BDD com barreira >= 80% de cobertura.
  - Makefile universal e determinístico (`make check`, `make test`, `make build`, etc.).
  - Stack customizável deliberada com o solicitante do prompt (sem imposições pré-fabricadas).
  - Recomendação formal do framework **OpenSpec** com visão conjunta de negócio para o **PO** e cenários BDD/Gherkin para automação por **QA**.
  - Pipelines de CI/CD e Matriz de Release customizáveis (`[PLATAFORMAS_ALVO]`) com fallback seguro.
  - README vivo e exaustivo.

## Capabilities

### New Capabilities
<!-- Nenhuma nova capacidade introduzida. -->

### Modified Capabilities
- `project-foundation`: Atualização dos requisitos do prompt mestre de fundação arquitetural para disponibilizar a suíte completa de 8 prompts universais para Go, Python, Node.js e TypeScript com versões para Antigravity CLI e Claude Code, indexados e documentados em `docs/README.md`, com governança Git aprofundada e `.gitignore` mandatório.

## Impact

- **Documentação**: São disponibilizados 8 prompts mestres em `docs/` acompanhados do catálogo e guia de Spec-Driven Development em `docs/README.md`, com o Pilar 10 aprofundado.
- **Especificações OpenSpec**: `specs/project-foundation/spec.md` passa a referenciar os 8 prompts, o catálogo `docs/README.md` e os requisitos de `.gitignore`/Governança Git.
- **Código / Aplicação**: Nenhum impacto direto no código fonte em Go ou nos binários em execução da aplicação `file-server`.
