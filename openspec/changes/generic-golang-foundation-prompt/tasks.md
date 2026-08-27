## 1. Padronização e Nomenclatura dos Prompts do Antigravity (AGY)

- [x] 1.1 Renomear `docs/generic-golang-foundation-spec-prompt.md` para `docs/generic-golang-foundation-spec-prompt-agy.md`
- [x] 1.2 Renomear `docs/generic-python-foundation-spec-prompt.md` para `docs/generic-python-foundation-spec-prompt-agy.md`
- [x] 1.3 Renomear `docs/generic-nodejs-foundation-spec-prompt.md` para `docs/generic-nodejs-foundation-spec-prompt-agy.md`
- [x] 1.4 Renomear `docs/generic-typescript-foundation-spec-prompt.md` para `docs/generic-typescript-foundation-spec-prompt-agy.md`

## 2. Criação dos Prompts Mestres para Claude Code

- [x] 2.1 Criar `docs/generic-golang-foundation-spec-prompt-claude.md` adaptado ao Claude Code (`CLAUDE.md`, ferramentas nativas Write/Edit/View/Grep/Glob/LS)
- [x] 2.2 Criar `docs/generic-python-foundation-spec-prompt-claude.md` adaptado ao Claude Code com `uv`, `ty` e `pyproject.toml`
- [x] 2.3 Criar `docs/generic-nodejs-foundation-spec-prompt-claude.md` adaptado ao Claude Code com ESM nativo
- [x] 2.4 Criar `docs/generic-typescript-foundation-spec-prompt-claude.md` adaptado ao Claude Code com tipagem estrita

## 3. Catálogo de Documentação e Atualização de Referências

- [x] 3.1 Atualizar referências internas em `.agent/rules/agent_harness_engineering.md` e especificações para os novos nomes dos arquivos
- [x] 3.2 Criar `docs/README.md` com a visão geral do estudo de Spec-Driven Development, catálogo dos 8 prompts mestres em tabela comparativa e guia de utilização

## 4. Aprofundamento do Pilar 10 (Governança Git, Higiene e .gitignore)

- [x] 4.1 Atualizar o Pilar 10 nos 8 prompts mestres (`docs/generic-*-foundation-spec-prompt-*.md`) com `.gitignore` mandatório e justificativa da governança Git
- [x] 4.2 Atualizar `docs/README.md` refletindo as boas práticas do Pilar 10

## 5. Validação de Qualidade e Governança

- [x] 5.1 Validar a integridade das especificações OpenSpec através de `openspec validate --all`
- [x] 5.2 Executar o quality gate completo do projeto (`make check`) validando formatação, linters e suíte de testes com cobertura >= 80%
