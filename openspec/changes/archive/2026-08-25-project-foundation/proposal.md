## Why

O objetivo desta mudança é estabelecer a fundação arquitetural e operacional completa para o projeto **File Server** (`file-server`) utilizando a linguagem Go (Golang) com renderização HTML no servidor e tecnologias modernas de interface. 

Para assegurar longevidade, alta qualidade, facilidade de manutenção, excelente experiência de uso e máxima eficiência no desenvolvimento assistido por IA, é essencial padronizar desde o dia zero as práticas de Engenharia de Software em Golang (Clean Architecture, TDD/BDD, cobertura de testes >= 80%, Makefile, Git/GitHub, documentação viva com README.md e badges de status) e estruturar as diretrizes para que o **Antigravity CLI** utilize estratégias avançadas de desenvolvimento de agentes (**Harness Engineering, Loop Engineering, Graph Engineering** e otimização de tokens) orientadas por Spec-Driven Development (OpenSpec) em Português do Brasil sob a ótica conjunta de PO (Product Owner) e QA (Quality Assurance).

## What Changes

- **Renomeação do Projeto para `file-server`**: Identificação do projeto como `file-server` (módulo Go `github.com/douglas/file-server`, executável `file-server`).
- **Documentação de Referência e Boas Práticas (`README.md`)**:
  - Inclusão de badges de status relevantes (Go version, CI Quality Gate, Test Coverage, Release, License).
  - Descrição introdutória clara do que se trata a aplicação.
  - Passo a passo completo de instalação e uso de **todos os comandos e flags disponíveis** da CLI (`file-server --help`, `file-server version`, `file-server serve` com portas/hosts configuráveis).
  - Guia aprofundado para **desenvolvedores** (arquitetura Clean/Ports & Adapters, TDD/BDD, esteira de cobertura >= 80%, Makefile harness, live-reload com Air, convenções de commits e releases multiplataforma).
- **Estrutura Base do Projeto Go**: Adoção do padrão canônico de layout de repositórios Go (`cmd/`, `internal/`, `pkg/`, `web/`, `scripts/`, `.github/`).
- **Padrões de Engenharia & Arquitetura Go**: Definição de arquitetura em camadas desacopladas (Clean/Ports and Adapters), injeção de dependência explícita, tratamento idiomático de erros e propagação de `context.Context`.
- **Pipeline de Testes, TDD & BDD em Go**: Configuração de suíte de testes com suporte a TDD/BDD, mocks, fixtures, verificação de cobertura mínima (meta >= 80%) e relatórios automatizados.
- **Interface Universal de Comandos (Makefile)**: Criação de um **Makefile** universal e autodocumentado para disponibilizar os principais comandos de uso da aplicação com fácil acesso (`make help`, `make dev`, `make run`, `make test`, `make lint`, `make check`, `make build`, `make build-all`, `make setup`).
- **Estratégias de Desenvolvimento com Antigravity CLI (Harness, Loop e Graph Engineering)**:
  - **Harness Engineering**: O Antigravity CLI utiliza os alvos do `Makefile` e scripts de teste como seu *harness* de automação e execução determinística de tarefas com saídas enxutas e baixo consumo de tokens.
  - **Loop Engineering**: O Antigravity CLI opera em loops de feedback rápido (edição -> teste -> validação -> refinamento contínuo) para garantir implementação correta antes da entrega.
  - **Graph Engineering**: O Antigravity CLI utiliza a árvore de dependências do projeto e o grafo DAG de tarefas do OpenSpec para navegar, planejar e implementar componentes em ordem topológica correta, sem dependências circulares.
- **Fundação de CLI Extensível & Versionamento Dinâmico**:
  - Arquitetura modular de CLI baseada em **Cobra** preparada para receber comandos, subcomandos, argumentos posicionais e flags/opções em futuras especificações.
  - Implementação do comando/argumento `version`: exibe a versão semântica oficial ao ser compilado via tag de release, ou exibe identificador explícito de desenvolvimento (`dev` / hash do commit / data) quando executado em ambiente de desenvolvimento local via injeção de `-ldflags`.
- **Pipeline de Release Multiplataforma no GitHub Actions (Linux & Windows)**:
  - Criação de workflow `.github/workflows/release.yml` acionado por tags Git (`v*`) para cross-compilação automática de binários autocontidos para **Linux** (`amd64`, `arm64`) e **Windows** (`amd64`, `arm64`), gerando arquivos compactados (`.tar.gz` e `.zip`) e anexando-os diretamente à Release do GitHub.
- **Stack & Plugins de Frontend e Binário Autocontido**:
  - Sugestão primária: **HTMX + Alpine.js + Tailwind CSS** integrado com Go `html/template`.
  - **Empacotamento de Assets (`go:embed`)**: Todos os templates HTML, layouts, arquivos estáticos (CSS, JS, imagens) são embutidos diretamente no binário compilado, permitindo portabilidade total (binário único independente).
  - Configuração de live-reload no servidor e hot-rebuild de templates e estilos em modo de desenvolvimento.
- **Governança Git, GitHub, OpenSpec & Ciclo de Vida**:
  - Padrão de Conventional Commits e hooks de pre-commit.
  - Workflows do GitHub Actions para CI/CD (linting com `golangci-lint`, testes, auditoria de vulnerabilidades e verificação de specs).
  - Configuração permanente no `openspec/config.yaml` (`context`, `rules`, `operations`) garantindo que toda nova proposta/plano respeite estritamente: Go idiomático, TDD/BDD, cobertura >= 80%, Makefile, PT-BR e perspectivas conjuntas de PO e QA.
  - Procedimento automatizado e guiado de **commit Git mandatório no arquivamento da especificação** (`operations.archive.guidance`).

## Capabilities

### New Capabilities
- `project-foundation`: Estabelece a fundação do projeto File Server (`file-server`) em Go + HTML, cobrindo arquitetura Go, CLI extensível com comando version, pipeline de release multiplataforma (Linux/Windows), TDD/BDD, esteira de testes, estratégias de engenharia para o Antigravity (Harness, Loop e Graph Engineering), governança de repositório, documentação completa via README.md com badges, regras permanentes do OpenSpec, stack de frontend e empacotamento de assets (`go:embed`) em binário autocontido.

### Modified Capabilities
<!-- Nenhuma capacidade existente a modificar (projeto inicial). -->

## Impact

- **Código e Repositório**: Inicialização do módulo Go (`go.mod`) com nome `github.com/douglas/file-server`, estrutura de diretórios e comandos CLI (`cmd/`), configurações de ferramentas de lint (`.golangci.yml`), `Makefile` autodocumentado com suporte a injeção de versão via `ldflags`, workflows do GitHub Actions para CI (`.github/workflows/ci.yml`) e Release (`.github/workflows/release.yml`), documentação `README.md` estruturada com badges e instruções completas, e módulo `web/` com `embed.FS` para distribuição em arquivo único executável para Linux e Windows.
- **Estratégias de IA & DX (Antigravity/Gemini)**: Criação de diretrizes em `.agent/rules/` estruturando o harness, loops de feedback e grafo de resolução do agente, além da configuração normativa em `openspec/config.yaml` com suporte a commit no arquivamento.
- **Dependências**: Adição de dependências base do ecossistema Go (framework de CLI como Cobra, ferramentas de teste e utilitários) e ferramental de frontend (ex: Tailwind CLI / HTMX).
