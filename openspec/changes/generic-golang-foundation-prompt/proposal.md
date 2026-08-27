## Why

O arquivo `docs/foundation-spec-prompt.md` atual possuía acoplamentos rígidos e específicos a certas tecnologias e frameworks (como HTMX, Alpine.js, Tailwind CSS e Cobra CLI), além de assumir sempre uma estrutura de servidor web com interface SSR e decisões pré-fabricadas de empacotamento (`go:embed`) e live-reload (`Air`). Ao iniciar novos projetos Go de naturezas distintas (como microserviços headless, APIs gRPC/REST, ferramentas CLI sem interface web, daemons, bibliotecas ou TUIs), esse acoplamento impõe especificações desnecessárias que precisam ser manualmente expurgadas.

É necessário generalizar o prompt mestre de fundação para que ele sirva como alicerce universal para qualquer projeto em Go, delegando ao usuário a definição explícita do tipo de aplicação, módulos, binários, bibliotecas, empacotamento e stack tecnológica desejada, mantendo intactos todos os pilares essenciais de engenharia (Clean Architecture/layout canônico, TDD/BDD com cobertura >= 80%, Makefile universal, CI/CD, documentação viva no README, governança Git com squash e Engenharia de Agentes com Harness, Loop e Graph). Além disso, o arquivo foi renomeado para `docs/generic-golang-foundation-spec-prompt.md` e o Pilar 7 deve incluir boas práticas para que o Product Owner e o QA compreendam claramente as regras técnicas e de negócio, facilitando a automação de testes.

## What Changes

- **Renomeação do Arquivo de Documentação**: Renomear `docs/foundation-spec-prompt.md` para `docs/generic-golang-foundation-spec-prompt.md`.
- **Desacoplamento de Frameworks e Decisões Pré-Fabricadas (Pilar 5)**:
  - Remover imposições obrigatórias de frameworks web (HTMX, Alpine.js, Tailwind CSS) e CLI (Cobra) do corpo dos pilares fundamentais.
  - Remover decisões pré-definidas de empacotamento autocontido (`go:embed`) e loops de live-reload (`Air`), delegando essas escolhas à deliberação conjunta com o solicitante do prompt conforme a necessidade real do projeto.
  - Transformar seções de stack e interface em blocos modulares/configuráveis definidos via variáveis e placeholders pelo usuário.
- **Boas Práticas de Governança OpenSpec para PO e QA (Pilar 7)**:
  - Estruturação de especificações com linguagem ubíqua e clara de negócio para o Product Owner (PO), sem jargões de baixo nível que ofusquem o valor entregue.
  - Estruturação de cenários determinísticos BDD/Gherkin (`Given-When-Then`, com entradas, saídas, bordas e erros) prontos para serem consumidos por ferramentas de automação de testes pelo QA.
- **Parametrização e Flexibilidade de Tecnologias**:
  - Adicionar placeholders e instruções claras para que o usuário informe: tipo de projeto (Web, API REST/gRPC, CLI, Daemon, Library, TUI), binários a compilar, dependências/frameworks escolhidos e módulos adicionais.
- **Preservação Integral dos Padrões de Excelência de Engenharia**:
  - Manter layout canônico e Clean Architecture em Go (`cmd/`, `internal/`).
  - Manter esteira automatizada de testes com TDD/BDD e barreira inegociável de cobertura >= 80%.
  - Manter Makefile universal autodocumentado (`make check`, `make test`, `make build`, etc.).
  - Manter governança de IA e Engenharia de Agentes do Antigravity (Harness Engineering, Loop Engineering, Graph Engineering, Native Tool Grounding e Economia Ativa de Tokens).
  - Manter pipelines de CI/CD e release multiplataforma no GitHub Actions.
  - Manter governança Git (Conventional Commits, branches de feature e merge exclusivamente Squash).
- **Atualização de Referências**: Atualizar referências ao arquivo renomeado nas especificações e documentações do projeto.

## Capabilities

### New Capabilities
<!-- Nenhuma nova capacidade introduzida. -->

### Modified Capabilities
- `project-foundation`: Atualização dos requisitos do prompt mestre de fundação arquitetural para padronizar o documento `docs/generic-golang-foundation-spec-prompt.md` como modelo genérico, agnóstico de tecnologias e com governança PO/QA aprimorada para qualquer projeto Go.

## Impact

- **Documentação**: `docs/foundation-spec-prompt.md` é renomeado e reformulado para `docs/generic-golang-foundation-spec-prompt.md`.
- **Especificações OpenSpec**: `specs/project-foundation/spec.md` passa a referenciar o novo nome, o comportamento agnóstico do prompt mestre e as boas práticas de governança PO/QA.
- **Código / Aplicação**: Nenhum impacto direto no código fonte em Go ou nos binários em execução da aplicação `file-server`.
