## 1. Criação e Estruturação do Prompt Mestre Genérico em Go

- [x] 1.1 Criar o arquivo `docs/generic-golang-foundation-spec-prompt.md` com estrutura universal e agnóstica para qualquer projeto Go
- [x] 1.2 Definir cabeçalho com guia de uso, tabela de variáveis/placeholders (`[NOME_DO_PROJETO]`, `[MODULO_GO]`, `[TIPO_DE_PROJETO]`, `[BINARIOS_OU_SERVICOS]`, `[STACK_E_FRAMEWORKS]`, `[PLATAFORMAS_ALVO]`, `[DESCRICAO_DO_PROJETO]`) e exemplos para múltiplos arquétipos
- [x] 1.3 Estruturar os 10 pilares fundamentais de engenharia (Clean Architecture, Entrypoints/ldflags, TDD/BDD >= 80%, Makefile universal, Stack customizável, Harness/Loop/Graph Engineering, OpenSpec PO/QA em PT-BR, CI/CD GitHub Actions, Living README e Governança Git com Squash)
- [x] 1.4 Refinar o Pilar 5 no prompt mestre removendo decisões pré-definidas de empacotamento autocontido (`go:embed`) e live-reload (`Air`), delegando essas escolhas ao solicitante do prompt
- [x] 1.5 Aprimorar o Pilar 7 no prompt mestre com boas práticas de governança OpenSpec para PO (regras de negócio e aceitação) e QA (preparado para ferramentas de automação de testes)
- [x] 1.6 Refinar o Pilar 8 no prompt mestre parametrizando arquiteturas de compilação/release (`[PLATAFORMAS_ALVO]`), definindo fallback para arquitetura corrente ou questionamento ao usuário, e incluindo boas práticas de prompt
- [x] 1.7 Acrescentar no documento e no Pilar 7 a recomendação formal do framework OpenSpec para boas práticas na governança contínua de especificações

## 2. Remoção do Arquivo Legado

- [x] 2.1 Remover o arquivo legado `docs/foundation-spec-prompt.md`

## 3. Atualização de Referências e Regras

- [x] 3.1 Atualizar referências a `docs/foundation-spec-prompt.md` em `.agent/rules/agent_harness_engineering.md` e eventuais outras documentações

## 4. Criação e Refinamento dos Prompts Mestres para Python, Node.js e TypeScript

- [x] 4.1 Criar o arquivo `docs/generic-python-foundation-spec-prompt.md` com os 10 pilares adaptados para o ecossistema Python
- [x] 4.2 Criar o arquivo `docs/generic-nodejs-foundation-spec-prompt.md` com os 10 pilares adaptados para o ecossistema Node.js / JavaScript moderno (ESM nativo, Clean Architecture, node:test/vitest >= 80%, Makefile com eslint/prettier, OpenSpec recomendado e CI/CD)
- [x] 4.3 Criar o arquivo `docs/generic-typescript-foundation-spec-prompt.md` com os 10 pilares adaptados para o ecossistema TypeScript (tipagem estrita, Clean Architecture, vitest >= 80%, Makefile com tsc/tsup/eslint, OpenSpec recomendado e CI/CD)
- [x] 4.4 Atualizar `docs/generic-python-foundation-spec-prompt.md` recomendando `uv` como gerenciador de versões/libs/venv, substituindo `mypy` por `ty` e exigindo o manifesto `pyproject.toml`

## 5. Validação de Qualidade e Governança

- [x] 5.1 Validar a integridade das especificações OpenSpec através de `openspec validate --all`
- [x] 5.2 Executar o quality gate completo do projeto (`make check`) validando formatação, linters e suíte de testes com cobertura >= 80%
