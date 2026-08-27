## Context

O arquivo `docs/foundation-spec-prompt.md` foi concebido originalmente com base nos requisitos específicos do projeto `file-server`, contendo acoplamentos a frameworks web (HTMX, Alpine.js, Tailwind CSS), CLI (Cobra), decisões pré-fixadas de empacotamento (`go:embed`), live-reload (`Air`) e plataformas fixas de compilação. Veja mais em [proposal.md](file:///home/douglas/Workspace/gemini/sftp-server/openspec/changes/generic-golang-foundation-prompt/proposal.md).

Para transformar esse prompt mestre em um modelo universal para qualquer novo projeto Go, é necessário desacoplar essas tecnologias específicas e parametrizar o template, delegando decisões arquiteturais e de compilação à negociação com o solicitante, incorporando boas práticas de governança OpenSpec para fácil leitura de regras de negócio pelo PO, automação de testes pelo QA e compilação/release orientada a plataformas alvo (`[PLATAFORMAS_ALVO]`).

## Goals / Non-Goals

**Goals:**
- Renomear o arquivo de documentação para `docs/generic-golang-foundation-spec-prompt.md`.
- Generalizar os 10 pilares fundamentais do prompt mestre para Golang, tornando-os agnósticos a frameworks web ou bibliotecas de CLI pré-definidas.
- Remover do Pilar 5 quaisquer decisões pré-definidas de empacotamento autocontido (`go:embed`) ou live-reload (`Air`), delegando essas escolhas ao solicitante do prompt conforme o contexto da aplicação.
- Aprimorar o Pilar 7 com boas práticas de governança OpenSpec para que o Product Owner (PO) compreenda claramente as regras técnicas e de negócio e o QA disponha de cenários testáveis preparados para automação de testes.
- Aprimorar o Pilar 8 com parametrização de plataformas de compilação e matriz de release (`[PLATAFORMAS_ALVO]`), aplicando fallback seguro para a arquitetura corrente ou instruindo o agente a questionar o usuário quando houver ambiguidade.
- Definir placeholders estruturados (`[NOME_DO_PROJETO]`, `[MODULO_GO]`, `[TIPO_DE_PROJETO]`, `[BINARIOS_OU_SERVICOS]`, `[STACK_E_FRAMEWORKS]`, `[PLATAFORMAS_ALVO]`, `[DESCRICAO_DO_PROJETO]`) permitindo ao desenvolvedor especificar todas as preferências ao invocar o prompt.
- Fornecer orientações e exemplos práticos de preenchimento para diferentes arquétipos de aplicação Go (Web/SSR, API REST/gRPC, CLI/TUI, Worker/Daemon, Biblioteca).
- Preservar rigorosamente os requisitos inegociáveis de engenharia: Clean Architecture/layout canônico, TDD/BDD com cobertura >= 80%, Makefile universal, Harness/Loop/Graph Engineering, CI/CD GitHub Actions e Git com Conventional Commits e merge via Squash.
- Atualizar todas as referências ao documento no repositório.

**Non-Goals:**
- Não alterar o código de produção ou binários do `file-server`.
- Não flexibilizar ou diminuir a barreira de qualidade (cobertura >= 80% e `make check` continuam mandatórios).
- Não remover as disciplinas de Engenharia de Agentes (Harness, Loop, Graph) ou governança OpenSpec.

## Decisions

### Decisão 1: Renomeação para `docs/generic-golang-foundation-spec-prompt.md`
- **Escolha**: Mover e renomear `docs/foundation-spec-prompt.md` para `docs/generic-golang-foundation-spec-prompt.md`.
- **Justificativa**: O novo nome explicita de forma inequívoca o propósito do arquivo como um gerador de especificações de fundação universal para Go.

### Decisão 2: Generalização dos Pilares de Stack, Entrada e Governança PO/QA
- **Escolha**:
  - **Pilar 2 (Modularidade de Entrada e Versionamento Dinâmico)**: Focado no composition root sob `cmd/` e injeção de versão via `-ldflags` (`internal/version`), deixando a biblioteca de parsing (Cobra, stdlib `flag`, urfave/cli, etc.) a critério do usuário.
  - **Pilar 5 (Stack Tecnológica e Dependências Customizáveis)**: Focado estritamente nas escolhas informadas em `[STACK_E_FRAMEWORKS]` e `[TIPO_DE_PROJETO]`, removendo imposições pré-definidas de `go:embed` ou live-reload (`Air`), que devem ser decididas conforme a necessidade real do projeto.
  - **Pilar 7 (Governança OpenSpec com Foco em PO e QA)**: Estruturar as regras para que:
    * O PO compreenda a motivação, o valor de negócio e as regras funcionais através de linguagem clara e critérios de aceitação objetivos, sem jargões de infraestrutura de baixo nível.
    * O QA obtenha especificações formais com cenários BDD/Gherkin (`Given-When-Then`), validação explícita de casos de borda, contratos de entrada/saída e tratamento de erros, viabilizando automação direta em ferramentas de teste.

### Decisão 3: Parametrização de Plataformas de Compilação e Release (Pilar 8)
- **Escolha**:
  - Adicionar placeholder `[PLATAFORMAS_ALVO]` no cabeçalho do prompt.
  - No Pilar 8, orientar a configuração da matriz de cross-compilação do Makefile (`make build-all`) e do GitHub Actions (`release.yml`) baseando-se estritamente nas plataformas informadas pelo usuário.
  - **Regra de Fallback e Boas Práticas**: Se `[PLATAFORMAS_ALVO]` não for preenchido, o agente deve assumir a arquitetura/sistema operacional corrente do host para `make build` ou questionar proativamente o usuário caso uma release multiplataforma seja explicitamente requerida.

### Decisão 4: Atualização Sistemática de Referências
- **Escolha**: Atualizar todas as ocorrências do nome antigo nos arquivos de regras (`.agent/rules/agent_harness_engineering.md`), configurações e documentações.
- **Justificativa**: Evita links quebrados e mantém coerência absoluta em todo o repositório.

## Risks / Trade-offs

- **[Risco]** Links quebrados apontando para o arquivo antigo.
  - *Mitigação*: Utilizar `grep_search` para rastrear todas as menções a `docs/foundation-spec-prompt.md` e atualizá-las para `docs/generic-golang-foundation-spec-prompt.md`.
- **[Risco]** Ambiguidade nas respostas da IA ao receber um prompt genérico.
  - *Mitigação*: Incluir no cabeçalho do documento uma tabela de placeholders e exemplos práticos de preenchimento para múltiplos arquétipos distintos com regras explícitas de fallback.

## Migration Plan

1. Atualizar o arquivo `docs/generic-golang-foundation-spec-prompt.md` com as refinarias do Pilar 5, Pilar 7 e Pilar 8.
2. Remover o arquivo legado `docs/foundation-spec-prompt.md`.
3. Atualizar referências no arquivo `.agent/rules/agent_harness_engineering.md` e em `specs/project-foundation/spec.md`.
4. Validar integridade do OpenSpec com `openspec validate --all` e do repositório com `make check`.
