## Why

O arquivo `docs/foundation-spec-prompt.md` atual possui acoplamentos rígidos e específicos a certas tecnologias e frameworks (como HTMX, Alpine.js, Tailwind CSS e Cobra CLI), além de assumir sempre uma estrutura de servidor web com interface SSR. Ao iniciar novos projetos Go de naturezas distintas (como microserviços headless, APIs gRPC/REST, ferramentas CLI sem interface web, daemons, bibliotecas ou TUIs), esse acoplamento impõe especificações desnecessárias que precisam ser manualmente expurgadas.

É necessário generalizar o prompt mestre de fundação para que ele sirva como alicerce universal para qualquer projeto em Go, delegando ao usuário a definição explícita do tipo de aplicação, módulos, binários, bibliotecas e stack tecnológica desejada, mantendo intactos todos os pilares essenciais de engenharia (Clean Architecture/layout canônico, TDD/BDD com cobertura >= 80%, Makefile universal, CI/CD, documentação viva no README, governança Git com squash e Engenharia de Agentes com Harness, Loop e Graph). Além disso, o arquivo deve ser renomeado para refletir sua natureza agnóstica e genérica: `docs/generic-golang-foundation-spec-prompt.md`.

## What Changes

- **Renomeação do Arquivo de Documentação**: Renomear `docs/foundation-spec-prompt.md` para `docs/generic-golang-foundation-spec-prompt.md`.
- **Desacoplamento de Frameworks Específicos**:
  - Remover imposições obrigatórias de frameworks web (HTMX, Alpine.js, Tailwind CSS) e CLI (Cobra) do corpo dos pilares fundamentais.
  - Transformar seções de stack e interface em blocos modulares/configuráveis definidos via variáveis e placeholders pelo usuário no momento da escrita do prompt.
- **Parametrização e Flexibilidade de Tecnologias**:
  - Adicionar placeholders e instruções claras para que o usuário informe: tipo de projeto (Web, API REST/gRPC, CLI, Daemon, Library, TUI), binários a compilar, dependências/frameworks escolhidos e módulos adicionais.
- **Preservação Integral dos Padrões de Excelência de Engenharia**:
  - Manter layout canônico e Clean Architecture em Go (`cmd/`, `internal/`).
  - Manter esteira automatizada de testes com TDD/BDD e barreira inegociável de cobertura >= 80%.
  - Manter Makefile universal autodocumentado (`make check`, `make test`, `make build`, etc.).
  - Manter governança de IA e Engenharia de Agentes do Antigravity (Harness Engineering, Loop Engineering, Graph Engineering, Native Tool Grounding e Economia Ativa de Tokens).
  - Manter governança OpenSpec (PT-BR, personas conjuntas de PO e QA com cenários BDD `#### Scenario:` e `WHEN/THEN`).
  - Manter pipelines de CI/CD e release multiplataforma no GitHub Actions.
  - Manter governança Git (Conventional Commits, branches de feature e merge exclusivamente Squash).
- **Atualização de Referências**: Atualizar referências ao arquivo renomeado nas especificações e documentações do projeto.

## Capabilities

### New Capabilities
<!-- Nenhuma nova capacidade introduzida. -->

### Modified Capabilities
- `project-foundation`: Atualização dos requisitos do prompt mestre de fundação arquitetural para padronizar o documento `docs/generic-golang-foundation-spec-prompt.md` como modelo genérico e agnóstico de tecnologias para qualquer projeto Go.

## Impact

- **Documentação**: `docs/foundation-spec-prompt.md` é renomeado e reformulado para `docs/generic-golang-foundation-spec-prompt.md`.
- **Especificações OpenSpec**: `specs/project-foundation/spec.md` passa a referenciar o novo nome e o comportamento agnóstico do prompt mestre.
- **Código / Aplicação**: Nenhum impacto direto no código fonte em Go ou nos binários em execução da aplicação `file-server`.
