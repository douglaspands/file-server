# Prompt Mestre: Fundação de Projetos Go com OpenSpec e Antigravity CLI

Este documento contém o **Prompt Mestre de Fundação de Projetos** (`foundation-spec-prompt`), projetado para ser utilizado com o **Antigravity CLI** e o framework **OpenSpec** (por exemplo, via comando `/openspec-propose` ou como instrução inicial). Ele condensa todos os requisitos arquiteturais, de qualidade, automação, esteira de testes, CI/CD, governança, autonomia operacional e engenharia de agentes estabelecidos na especificação de fundação (`project-foundation`), permitindo replicar com precisão esse mesmo padrão de excelência para qualquer novo projeto Go.

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

A especificação deve ser gerada utilizando o padrão OpenSpec em Português do Brasil (PT-BR) e cobrir integralmente os 10 pilares fundamentais descritos a seguir, produzindo os artefatos: 'proposal.md', 'design.md', 'specs/project-foundation/spec.md', 'tasks.md', 'openspec/config.yaml', 'AGENTS.md', 'GEMINI.md', '.agent/settings.json' e as regras em '.agent/rules/'.

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
  * scripts/: Scripts de suporte, automação de cobertura e verificação.
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
PILAR 6: ENGENHARIA DE AGENTES, PRIORIDADE DE FERRAMENTAS E ECONOMIA DE TOKENS
================================================================================
- Documentar e configurar em 'AGENTS.md', 'GEMINI.md', '.agent/rules/' e '.agent/settings.json' as diretrizes de operação e autonomia do Antigravity (AGY):
  1. Prioridade Mandatória de Ferramentas Nativas (Native Tool Grounding):
     - Criar arquivos: Utilizar estritamente 'write_to_file'. Proibido executar 'cat << EOF', 'echo >' ou 'touch' via terminal.
     - Editar arquivos: Utilizar 'replace_file_content' para alterações pontuais e cirúrgicas. Proibido scripts 'sed', 'awk' ou 'cat >' no terminal.
     - Inspecionar / Ler arquivos: Utilizar 'view_file' especificando 'StartLine' e 'EndLine' para focar apenas no trecho relevante. Proibido 'cat', 'head', 'tail'.
     - Buscar código: Utilizar 'grep_search'. Proibido 'grep', 'rg' via terminal.
     - Localizar arquivos: Utilizar 'find_by_name'. Proibido 'find', 'ls -R' via terminal.
     - Listar diretórios: Utilizar 'list_dir'. Proibido 'ls' via terminal.
  2. Uso Restrito do Terminal ('run_command'):
     - O terminal deve ser utilizado exclusivamente para ferramentas do ciclo de vida: 'make' ('make test', 'make lint', 'make check', 'make build'), 'go' ('go test', 'go mod tidy'), 'git', 'openspec' e binários executáveis.
  3. Autonomia Operacional e Permissões (.agent/settings.json):
     - Configurar permissões pré-autorizadas ('allow') para todas as ferramentas cotidianas ('go', 'make', 'git', 'openspec', linters), eliminando prompts de confirmação repetitivos e assegurando execução autônoma contínua.
  4. Economia Ativa de Tokens, Curação de Contexto e Otimização de Janela:
     - Respostas do agente devem ser concisas e estruturadas em Markdown com links clicáveis '[arquivo](file:///caminho)'.
     - NUNCA duplicar blocos massivos de código já existentes no disco na resposta do chat.
     - Leitura cirúrgica de arquivos limitando linhas ('StartLine'/'EndLine') e edições pontuais com 'replace_file_content'.
     - Execução de comandos com saídas concisas e estruturadas (sem flags excessivamente verbosas).
  5. Disciplinas de Engenharia de Agentes em Sistemas Compostos (Harness, Loop & Graph Engineering):
     - Harness Engineering: Arcabouço operacional e scaffolding de IA composto por diretrizes determinísticas (.agent/rules/, AGENTS.md, GEMINI.md), matriz de permissões pré-autorizadas (.agent/settings.json), guardrails de segurança, contexto otimizado e limites de sandbox para governar a operação autônoma, previsível e segura do Antigravity CLI.
     - Loop Engineering: Ciclos cognitivos e iterativos de execução e auto-validação contínua do agente (ReAct / Reflection loops: inspeção cirúrgica de requisitos -> intervenção pontual no código -> execução de testes automatizados -> diagnóstico de erros -> correção e validação final com 'make check'), prevenindo loops infinitos e garantindo convergência rápida.
     - Graph Engineering (State Graphs & DAG de Raciocínio): Modelagem e orquestração do fluxo cognitivo/operacional do agente como um grafo direcionado acíclico (DAG) de tarefas, estados e decisões. Decomposição topológica de dependências entre tarefas de planejamento e execução, controle de transições determinísticas de estado entre fases (especificação -> implementação -> verificação -> reflexão), roteamento de fluxos e coordenação/paralelismo entre subagentes especializados sem dependências circulares ou bloqueios cognitivos.
     - Gestão de Contexto e Subagentes: Delegação de pesquisas e tarefas isoladas para subagentes dedicados para manter a janela de contexto do agente principal limpa, enxuta e focada.

================================================================================
PILAR 7: GOVERNANÇA OPENSPEC (PO/QA, PT-BR E REGRAS PERMANENTES)
================================================================================
- Configurar 'openspec/config.yaml' com:
  * context: Declaração da stack (Go, Clean Architecture, HTMX, Tailwind, embed, TDD/BDD >= 80%, Makefile, PT-BR, Autonomia e Ferramentas Nativas do AGY).
  * rules.proposal: Foco no 'porquê' e 'o que muda', declaração explícita de impacto e idioma PT-BR.
  * rules.specs: Escrita em PT-BR sob ótica CONJUNTA de PO (valor de negócio e critérios de aceitação) e QA (cenários de teste com 4 hashtags '#### Scenario:', bullets '- **WHEN**' e '- **THEN**', validação de bordas e erros). Mínimo de 50 caracteres na seção '## Purpose'.
  * rules.tasks: Seções numeradas '## N. Nome do Grupo', checkboxes '- [ ] N.M Descrição', tarefas explícitas de testes unitários/BDD com validação de cobertura >= 80% e 'make check'.
  * operations.apply.guidance: Executar 'make test' e 'make check' para validar alterações antes de concluir tarefas, priorizando ferramentas nativas.

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
     - Interface universal de comandos e automação via Makefile ('make check', 'make dev', 'make test', 'make build-all').
     - Metodologia de testes TDD/BDD com barreira de 80%.
     - Live-reload no desenvolvimento web.
     - Diretrizes de IA (Antigravity), autonomia de ferramentas e padrão Conventional Commits.

================================================================================
PILAR 10: GOVERNANÇA GIT, CONVENTIONAL COMMITS, BRANCHING E SQUASH MERGE
================================================================================
- Padrão estrito de Conventional Commits: 'feat:', 'fix:', 'refactor:', 'test:', 'docs:', 'chore:'.
- Feature Branch por Especificação:
  * Toda nova mudança do OpenSpec inicia em branch dedicada: 'git checkout main && git checkout -b feature/<change-name>'.
- Permissão Obrigatória do Usuário:
  * Ao concluir e arquivar uma especificação, o agente DEVE solicitar permissão explícita ao usuário antes de integrar na 'main' ou abrir Pull Request.
- Estratégia de Integração Exclusivamente SQUASH:
  * Merge Local: 'git checkout main && git merge --squash feature/<change-name> && git commit -m "feat(<change-name>): ..."'.
  * Pull Request: Integração configurada estritamente via Squash and Merge ('gh pr merge --squash').
```

---

## Estrutura Esperada dos Artefatos OpenSpec Gerados

Quando o prompt acima for executado pelo Antigravity CLI, ele deve produzir a seguinte árvore de artefatos no repositório:

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
AGENTS.md
GEMINI.md
.agent/
├── settings.json
└── rules/
    ├── agent_harness_engineering.md
    ├── agent_tooling_autonomy.md
    ├── git_branching_workflow.md
    ├── archive_workflow.md
    └── golang_conventions.md
```

### Checklist dos Artefatos Gerados:
- [x] **`proposal.md`**: Define o *Why*, *What Changes*, *Capabilities* (`project-foundation`) e *Impact* (código, IA, governança, autonomia).
- [x] **`design.md`**: Detalha *Context*, *Goals/Non-Goals*, *Decisions* (arquitetura Go, testes BDD, Makefile, Air, tokens, HTMX/Tailwind, Cobra/ldflags, release matrix, README, CI/CD, prioridade de ferramentas nativas) e *Risks/Trade-offs*.
- [x] **`specs/project-foundation/spec.md`**: Define os requisitos com perspectivas explícitas de **PO** e **QA**, cada um acompanhado de cenários BDD `#### Scenario:` com `- **WHEN**` e `- **THEN**`.
- [x] **`tasks.md`**: Lista de tarefas organizadas por seções numeradas (`## 1.` a `## 10.`), com subtarefas granulares `- [ ] N.M`, incluindo verificações de compilação, testes unitários, linter, cobertura >= 80% e `openspec validate --all`.
- [x] **`openspec/config.yaml`**: Configuração central com `context`, `rules` (proposal, specs, design, tasks) e `operations` (apply).
- [x] **`AGENTS.md` e `GEMINI.md`**: Diretrizes operacionais de alta prioridade para autonomia do AGY, ferramentas nativas, economia de tokens e fluxo Git squash.
- [x] **`.agent/settings.json`**: Permissões de desenvolvimento pré-autorizadas para máxima autonomia sem interrupções.
- [x] **`.agent/rules/`**: Diretrizes operacionais para Harness de IA (scaffolding e guardrails), Loop e Graph Engineering (State Graphs e DAG de tarefas), Autonomia de Ferramentas, Git Branching Workflow, Arquivamento com Squash e Convenções Go.
