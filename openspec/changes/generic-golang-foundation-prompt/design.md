## Context

O arquivo `docs/foundation-spec-prompt.md` foi concebido originalmente com base nos requisitos específicos do projeto `file-server`, contendo acoplamentos a frameworks web (HTMX, Alpine.js, Tailwind CSS) e CLI (Cobra). Veja mais em [proposal.md](file:///home/douglas/Workspace/gemini/sftp-server/openspec/changes/generic-golang-foundation-prompt/proposal.md).

Para transformar esse prompt mestre em um modelo universal para qualquer novo projeto Go, é necessário desacoplar essas tecnologias específicas e parametrizar o template, preservando a totalidade dos pilares de governança, arquitetura limpa, esteira de testes, automação e engenharia de agentes.

## Goals / Non-Goals

**Goals:**
- Renomear o arquivo de documentação para `docs/generic-golang-foundation-spec-prompt.md`.
- Generalizar os 10 pilares fundamentais do prompt mestre para Golang, tornando-os agnósticos a frameworks web ou bibliotecas de CLI pré-definidas.
- Definir placeholders estruturados (`[NOME_DO_PROJETO]`, `[MODULO_GO]`, `[TIPO_DE_PROJETO]`, `[BINARIOS_OU_SERVICOS]`, `[STACK_E_FRAMEWORKS]`, `[DESCRICAO_DO_PROJETO]`) permitindo ao desenvolvedor especificar as tecnologias desejadas ao invocar o prompt.
- Fornecer orientações e exemplos práticos de preenchimento para diferentes arquétipos de aplicação Go (Web/SSR, API REST/gRPC, CLI/TUI, Worker/Daemon, Biblioteca).
- Preservar rigorosamente os requisitos inegociáveis de engenharia: Clean Architecture/layout canônico, TDD/BDD com cobertura >= 80%, Makefile universal, Harness/Loop/Graph Engineering, governança OpenSpec (PO/QA em PT-BR), CI/CD GitHub Actions e Git com Conventional Commits e merge via Squash.
- Atualizar todas as referências ao documento no repositório.

**Non-Goals:**
- Não alterar o código de produção ou binários do `file-server`.
- Não flexibilizar ou diminuir a barreira de qualidade (cobertura >= 80% e `make check` continuam mandatórios).
- Não remover as disciplinas de Engenharia de Agentes (Harness, Loop, Graph) ou governança OpenSpec.

## Decisions

### Decisão 1: Renomeação para `docs/generic-golang-foundation-spec-prompt.md`
- **Escolha**: Mover e renomear `docs/foundation-spec-prompt.md` para `docs/generic-golang-foundation-spec-prompt.md`.
- **Justificativa**: O novo nome explicita de forma inequívoca o propósito do arquivo como um gerador de especificações de fundação universal para Go.
- **Alternativas consideradas**: Manter o nome antigo (causaria confusão quanto ao escopo) ou abreviar para `docs/generic-go-prompt.md` (menos descritivo).

### Decisão 2: Generalização dos Pilares de Stack e Entrada
- **Escolha**:
  - **Pilar 2 (Modularidade de Entrada e Versionamento Dinâmico)**: Focado no composition root sob `cmd/` e injeção de versão via `-ldflags` (`internal/version`), deixando a biblioteca de parsing (Cobra, stdlib `flag`, urfave/cli, etc.) a critério do usuário.
  - **Pilar 5 (Stack Tecnológica e Dependências Customizáveis)**: Focado em diretrizes de modularidade, baixa dependência externa, separação de camadas e empacotamento autocontido (usando `go:embed` quando aplicável), permitindo ao usuário definir no prompt se a aplicação terá Web, REST, gRPC, TUI, workers, etc.
- **Justificativa**: Permite que o mesmo prompt sirva tanto para um microserviço headless quanto para uma aplicação web monolítica ou ferramenta de terminal pura.
- **Alternativas consideradas**: Criar múltiplos arquivos de prompt para cada tipo de projeto (Web, CLI, API) - descartado por redundância e custo de manutenção; um único prompt parametrizável é muito mais sustentável.

### Decisão 3: Atualização Sistemática de Referências
- **Escolha**: Atualizar todas as ocorrências do nome antigo nos arquivos de regras (`.agent/rules/agent_harness_engineering.md`), configurações e documentações.
- **Justificativa**: Evita links quebrados e mantém coerência absoluta em todo o repositório.

## Risks / Trade-offs

- **[Risco]** Links quebrados apontando para o arquivo antigo.
  - *Mitigação*: Utilizar `grep_search` para rastrear todas as menções a `docs/foundation-spec-prompt.md` e atualizá-las para `docs/generic-golang-foundation-spec-prompt.md`.
- **[Risco]** Ambiguidade nas respostas da IA ao receber um prompt genérico.
  - *Mitigação*: Incluir no cabeçalho do documento uma tabela de placeholders e exemplos práticos de preenchimento para 4 arquétipos distintos (API REST, Web SSR, Ferramenta CLI, Worker assíncrono).

## Migration Plan

1. Criar o novo arquivo `docs/generic-golang-foundation-spec-prompt.md` com a versão generalizada e modular.
2. Remover o arquivo legado `docs/foundation-spec-prompt.md`.
3. Atualizar referências no arquivo `.agent/rules/agent_harness_engineering.md` e em `specs/project-foundation/spec.md`.
4. Validar integridade do OpenSpec com `openspec validate --all` e do repositório com `make check`.
