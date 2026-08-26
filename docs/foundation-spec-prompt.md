# Prompt Mestre: Fundação de Projetos Go com OpenSpec e Antigravity CLI

Este documento contém o **Prompt Mestre de Fundação de Projetos** (`foundation-spec-prompt`), projetado para ser utilizado com o **Antigravity CLI** e o framework **OpenSpec** (por exemplo, via comando `/openspec-propose` ou como instrução inicial). Ele condensa todos os requisitos arquiteturais, de qualidade, automação, esteira de testes, CI/CD, governança e engenharia de agentes estabelecidos na especificação de fundação (`project-foundation`), permitindo replicar com precisão esse mesmo padrão de excelência para qualquer novo projeto Go.

---

## Como Utilizar Este Prompt

1. **Copie o texto** da seção [Prompt de Fundação para Novos Projetos](#prompt-de-fundação-para-novos-projetos) abaixo.
2. **Substitua as variáveis entre colchetes** pelos dados do seu novo projeto:
   - `[NOME_DO_PROJETO]`: Nome do projeto (ex: `file-server`, `auth-service`, `gateway-api`).
   - `[MODULO_GO]`: Caminho canônico do módulo Go (ex: `github.com/usuario/meu-projeto`).
   - `[DESCRICAO_DO_PROJETO]`: Resumo do propósito e valor de negócio da aplicação.
   - `[COMANDO_CLI]`: Nome do executável de linha de comando (ex: `file-server`, `app`).
3. **Execute no Antigravity CLI** ou passe para a ferramenta de IA como instrução para propor uma nova mudança:
   ```text
   /openspec-propose Crie a especificação de fundação arquitetural e de engenharia 'project-foundation' para o projeto [NOME_DO_PROJETO] seguindo as diretrizes abaixo:
   <cole o prompt>
   ```

---

## Prompt de Fundação para Novos Projetos

```text
Você é um Arquiteto de Software Principal, Engenheiro Líder em Go (Golang) e Especialista em Engenharia de Agentes de IA (Antigravity CLI).

Sua tarefa é criar a especificação completa de fundação arquitetural e de engenharia de software intitulada 'project-foundation' para o novo projeto '[NOME_DO_PROJETO]' (Módulo Go: '[MODULO_GO]', Executável CLI: '[COMANDO_CLI]'), cuja descrição é:
"[DESCRICAO_DO_PROJETO]"

A especificação deve ser gerada utilizando o padrão OpenSpec em Português do Brasil (PT-BR) e cobrir integralmente os 10 pilares fundamentais descritos a seguir, produzindo os artefatos: 'proposal.md', 'design.md', 'specs/project-foundation/spec.md', 'tasks.md', 'openspec/config.yaml' e as diretrizes em '.agent/rules/'.

================================================================================
PILAR 1: PADRÃO ARQUITETURAL E ESTRUTURA CANÔNICA EM GO (CLEAN ARCHITECTURE)
================================================================================
- Adotar o layout canônico de pastas da comunidade Go:
  * cmd/[COMANDO_CLI]/main.go: Ponto de entrada da aplicação, parse de flags e composição de dependências (composition root).
  * internal/core/domain/: Entidades de negócio puras, livres de dependências externas ou de infraestrutura.
  * internal/core/ports/: Interfaces e contratos de entrada (casos de uso) e saída (repositórios, adaptadores de I/O).
  * internal/core/services/: Implementação da lógica de negócios e casos de uso, dependendo apenas de domain e ports.
  * internal/adapters/handlers/: Controladores HTTP, adaptadores REST e renderizadores de templates web.
  * internal/adapters/repositories/: Adaptadores de persistência e armazenamento de dados.
  * internal/version/: Pacote centralizado para controle de versão, commit e build date.
  * internal/testutils/: Helpers, fixtures e utilitários para testes.
  * web/templates/ e web/static/: Templates HTML, estilos CSS, scripts JS e assets da interface.
  * scripts/: Scripts de suporte, harness de cobertura e verificação.
  * .github/workflows/: Pipelines de automação CI/CD e release multiplataforma.
- Regra de Isolamento: O código de domínio e regras de negócio privadas devem residir estritamente em 'internal/'.
- Injeção de Dependências: Realizada de forma explícita e manual nos construtores (NewService, NewHandler) no composition root, sem uso de frameworks pesados baseados em reflexão em tempo de execução.
- Propagação Idiomática: Propagação obrigatória de 'context.Context' como primeiro parâmetro em operações com I/O, concorrência ou ciclo de vida, com suporte a encerramento gracioso (graceful shutdown).
- Fornecer um serviço de domínio de referência inicial funcional e testado (ex: HealthCheckService com status e versão da aplicação).

================================================================================
PILAR 2: CLI MODULAR E EXTENSÍVEL COM VERSIONAMENTO DINÂMICO (COBRA + LDFLAGS)
================================================================================
- Arquitetura de CLI baseada em Cobra (github.com/spf13/cobra) estruturada em cmd/:
  * cmd/root.go: Comando raiz ('[COMANDO_CLI]'), flags globais (--config, --verbose).
  * cmd/version.go: Comando e argumento 'version' com suporte a flags (--json, -v, --version).
  * cmd/serve.go: Comando para iniciar o servidor web/HTTP com flags de porta e host configuráveis (--port, --host).
- Versionamento Dinâmico (Release vs Dev):
  * Pacote 'internal/version' com variáveis públicas: Version = "dev", Commit = "none", Date = "unknown".
  * Injeção de metadados em tempo de compilação via '-ldflags' no Makefile e GitHub Actions:
    - Ao compilar a partir de uma tag Git de release (ex: v1.0.0), exibir a versão semântica oficial: "[NOME_DO_PROJETO] version: v1.0.0 (commit: abc1234, built at: 2026-08-25T12:00:00Z)".
    - Ao executar em ambiente de desenvolvimento local ou sem tag, exibir identificação explícita de desenvolvimento: "[NOME_DO_PROJETO] version: dev (commit: abc1234, built at: 2026-08-25T12:00:00Z)".

================================================================================
PILAR 3: INFRAESTRUTURA DE TESTES AUTOMATIZADOS, TDD/BDD E BARREIRA >= 80%
================================================================================
- Ferramental: Utilizar o pacote nativo 'testing' do Go combinado com 'github.com/stretchr/testify' ('assert' e 'require') para asserções expressivas e limpas.
- Metodologia BDD: Estruturar cenários de teste através de subtestes BDD declarativos utilizando o padrão:
  t.Run("Given [context] When [action] Then [outcome]", func(t *testing.T) { ... })
- Desacoplamento e Mocks: Mocks determinísticos e isolados baseados exclusivamente nos contratos de interfaces em 'internal/core/ports/'.
- Barreira de Cobertura Inegociável:
  * Script automatizado 'scripts/coverage.sh' que calcula a cobertura global via 'go tool cover'.
  * O script deve falhar com código de erro se a cobertura de código for inferior a 80%.
  * Geração automática de relatórios 'coverage.out' e 'coverage.html'.

================================================================================
PILAR 4: INTERFACE UNIVERSAL DE COMANDOS VIA MAKEFILE AUTODOCUMENTADO
================================================================================
- Criar um 'Makefile' universal, autodocumentado e interativo como interface central de automação:
  * make help (alvo padrão): Exibe dinamicamente o menu com todos os comandos disponíveis e descrições formatadas.
  * make setup: Instala e valida ferramentas de desenvolvimento locais (golangci-lint, air, govulncheck, tailwind).
  * make dev: Inicia o loop de desenvolvimento com live-reload (Air) e recompilação automática de assets/templates.
  * make run: Executa a aplicação a partir do código fonte.
  * make test: Executa a suíte completa de testes com verificação de cobertura (barreira >= 80%).
  * make test-unit: Executa testes unitários rápidos.
  * make test-coverage: Gera o relatório HTML de cobertura e valida o limiar de 80%.
  * make lint: Executa o linter estrito 'golangci-lint' com regras de complexidade, bugs e boas práticas (.golangci.yml).
  * make fmt: Formata o código Go (gofmt, goimports) e templates.
  * make check: Quality Gate local unificado (fmt + lint + govulncheck + test com cobertura).
  * make build: Compila o binário de produção otimizado com stripping de símbolos (-s -w) e injeção de ldflags.
  * make build-all: Compila binários cruzados para Linux (amd64, arm64) e Windows (amd64, arm64).
  * make clean: Limpa binários em dist/, relatórios de cobertura e artefatos temporários.

================================================================================
PILAR 5: STACK DE FRONTEND WEB, LIVE-RELOAD E BINÁRIO 100% AUTOCONTIDO (GO:EMBED)
================================================================================
- Camada de Frontend SSR Hipermidiática:
  * Renderização no servidor com 'html/template' do Go com suporte a layouts parciais e blocos modulares.
  * HTMX para requisições AJAX, Server-Sent Events e atualizações parciais do DOM sem necessidade de SPAs pesadas.
  * Alpine.js para reatividade leve e comportamentos declarativos na interface (modais, menus, dropdowns).
  * Tailwind CSS (Standalone CLI) para estilização utilitária moderna sem dependência de ecossistema Node.js pesado.
- Empacotamento de Assets via 'go:embed' (Binário Único Autocontido):
  * O pacote 'web/' deve embutir todos os templates HTML e assets estáticos (CSS, JS, imagens) diretamente no executável via diretiva '//go:embed templates/* static/*'.
  * A aplicação compilada com 'make build' deve ser um arquivo único, totalmente executável em qualquer ambiente sem necessidade de copiar pastas externas 'templates/' ou 'static/'.
- Live-Reload no Desenvolvimento:
  * Configurar o Air ('.air.toml') monitorando alterações em arquivos Go, templates '.html' e folhas '.css'.
  * Script injetado de Live Reload no modo desenvolvimento ('make dev') recarregando a página no navegador com ciclo de feedback inferior a 2 segundos.

================================================================================
PILAR 6: ESTRATÉGIAS DO ANTIGRAVITY CLI (HARNESS, LOOP E GRAPH ENGINEERING)
================================================================================
- Documentar e configurar em '.agent/rules/' e 'openspec/config.yaml' as estratégias para desenvolvimento assistido por IA:
  1. Harness Engineering:
     - O Antigravity CLI deve utilizar o 'Makefile' como seu harness de automação ('make test', 'make lint', 'make check').
     - Utilizar comandos silenciosos/concisos para reduzir o gasto de tokens e manter o foco na janela de contexto.
  2. Loop Engineering:
     - Operar em ciclos curtos de auto-validação: Inspeção de Contratos -> Implementação Pontual -> Execução de Testes Automatizados -> Diagnóstico de Falhas -> Correção -> Validação Final com 'make check'.
     - Nenhuma tarefa deve ser marcada como concluída sem evidência de testes passando e cobertura >= 80%.
  3. Graph Engineering:
     - Planejar e implementar componentes navegando pelo grafo acíclico de dependências (DAG):
       Contratos/Interfaces ('ports/') -> Entidades de Domínio ('domain/') -> Serviços de Negócio ('services/') -> Adaptadores/Handlers ('adapters/') -> Composição ('cmd/').
     - Prevenir ativamente qualquer dependência circular entre módulos.

================================================================================
PILAR 7: GOVERNANÇA OPENSPEC (PO/QA, PT-BR E REGRAS PERMANENTES)
================================================================================
- Configurar 'openspec/config.yaml' com:
  * context: Declaração da stack (Go, Clean Architecture, HTMX, Tailwind, embed, TDD/BDD >= 80%, Makefile, PT-BR).
  * rules.proposal: Foco no 'porquê' e 'o que muda', declaração explícita de impacto e idioma PT-BR.
  * rules.specs: Escrita em PT-BR sob ótica CONJUNTA de PO (valor de negócio e critérios de aceitação) e QA (cenários de teste com 4 hashtags '#### Scenario:', bullets '- **WHEN**' e '- **THEN**', validação de bordas e erros). Mínimo de 50 caracteres na seção '## Purpose'.
  * rules.tasks: Seções numeradas '## N. Nome do Grupo', checkboxes '- [ ] N.M Descrição', tarefas explícitas de testes unitários/BDD com validação de cobertura >= 80% e 'make check'.
  * operations.apply.guidance: Executar 'make test' e 'make check' para validar alterações antes de concluir tarefas.
  * operations.archive.guidance: Ao arquivar uma especificação concluída, o agente DEVE criar o commit Git com Conventional Commits: 'feat(spec): archive <change-name> and apply changes'.

================================================================================
PILAR 8: PIPELINES DE CI/CD E RELEASE MULTIPLATAFORMA NO GITHUB ACTIONS
================================================================================
- Workflow de CI ('.github/workflows/ci.yml'):
  * Executado em Pull Requests e pushes na branch principal ('main').
  * Etapas: Checkout, Setup Go, golangci-lint, govulncheck, execução de testes com barreira de cobertura >= 80% e validação de integridade OpenSpec ('openspec validate --all').
- Workflow de Release Multiplataforma ('.github/workflows/release.yml'):
  * Disparado na criação de tags de release Git ('v*').
  * Cross-compilação com injeção de versão via '-ldflags' para:
    - Linux: amd64 e arm64 (binários empacotados em arquivos '.tar.gz').
    - Windows: amd64 e arm64 (binários '.exe' empacotados em arquivos '.zip').
  * Geração de checksums SHA256 e publicação automática de todos os artefatos na Release do GitHub.

================================================================================
PILAR 9: DOCUMENTAÇÃO VIVA E EXAUSTIVA NO README.MD
================================================================================
- Estruturar o 'README.md' raiz com as seguintes seções obrigatórias:
  1. Header & Badges: Título com badges SVG funcionais (Go Version, CI Quality Gate, Test Coverage >= 80%, Latest Release, License).
  2. Visão Geral da Aplicação: Propósito, valor entregue e diferenciais técnicos.
  3. Guia de Instalação e Uso:
     - Download dos binários pré-compilados e compilação a partir do código fonte.
     - Documentação completa e tabela com TODOS os comandos e flags da CLI ('[COMANDO_CLI]', 'version', 'serve', '--port', '--host', '--json', '--verbose', '--config', '--help') com exemplos práticos.
  4. Guia do Desenvolvedor:
     - Pré-requisitos e setup do ambiente ('make setup').
     - Arquitetura de software e layout de diretórios ('internal/core/', 'internal/adapters/').
     - Harness de comandos via Makefile ('make check', 'make dev', 'make test', 'make build-all').
     - Metodologia de testes TDD/BDD com barreira de 80%.
     - Live-reload no desenvolvimento web.
     - Diretrizes de IA (Antigravity) e padrão Conventional Commits.

================================================================================
PILAR 10: GOVERNANÇA GIT, CONVENTIONAL COMMITS E HOOKS
================================================================================
- Padrão estrito de Conventional Commits: 'feat:', 'fix:', 'refactor:', 'test:', 'docs:', 'chore:'.
- Configuração de ganchos de pre-commit para validação de formatação de código e mensagens de commit.
- Procedimento mandatório e automatizado de commit Git durante o arquivamento de especificações OpenSpec.
```

---

## Estrutura Esperada dos Artefatos OpenSpec Gerados

Quando o prompt acima for executado pelo Antigravity CLI, ele deve produzir a seguinte árvore de artefatos na pasta `openspec/changes/project-foundation/`:

```text
openspec/
├── config.yaml
└── changes/
    └── project-foundation/
        ├── proposal.md
        ├── design.md
        ├── tasks.md
        └── specs/
            └── project-foundation/
                └── spec.md
.agent/
└── rules/
    ├── agent_harness_engineering.md
    ├── golang_conventions.md
    └── archive_workflow.md
```

### Checklist dos Artefatos Gerados:
- [x] **`proposal.md`**: Define o *Why*, *What Changes*, *Capabilities* (`project-foundation`) e *Impact* (código, IA, governança).
- [x] **`design.md`**: Detalha *Context*, *Goals/Non-Goals*, *Decisions* (arquitetura Go, testes BDD, Makefile, Air, tokens, HTMX/Tailwind, Cobra/ldflags, release matrix, README, CI/CD) e *Risks/Trade-offs*.
- [x] **`specs/project-foundation/spec.md`**: Define os 10 requisitos com perspectivas explícitas de **PO** e **QA**, cada um acompanhado de cenários BDD `#### Scenario:` com `- **WHEN**` e `- **THEN**`.
- [x] **`tasks.md`**: Lista de 10 grupos de tarefas numeradas (`## 1.` a `## 10.`), com subtarefas granulares `- [ ] N.M`, incluindo verificações de compilação, testes unitários, linter, cobertura >= 80% e `openspec validate --all`.
- [x] **`openspec/config.yaml`**: Configuração central com `context`, `rules` (proposal, specs, design, tasks) e `operations` (apply e archive).
- [x] **`.agent/rules/`**: Diretrizes operacionais para Harness, Loop, Graph Engineering, convenções Go e fluxo de arquivamento com commit.
