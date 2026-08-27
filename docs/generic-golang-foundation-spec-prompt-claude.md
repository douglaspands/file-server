# Prompt Mestre: Fundação Genérica de Projetos Go (Golang) com OpenSpec e Claude Code

Este documento contém o **Prompt Mestre de Fundação de Projetos Go** (`generic-golang-foundation-spec-prompt-claude`), projetado para ser utilizado com o **Claude Code** (Anthropic) e o framework **OpenSpec** (por exemplo, via comando `/openspec-propose` ou como instrução inicial).

É altamente **recomendado adotar o framework OpenSpec** como padrão de excelência para a governança contínua de especificações, documentação viva, rastreabilidade de decisões e alinhamento transparente entre visão de produto (PO) e qualidade técnica/automação (QA).

Ele condensa todos os requisitos de arquitetura limpa, qualidade de código, automação de comandos, esteira de testes com TDD/BDD, CI/CD, governança OpenSpec, autonomia operacional e Engenharia de Agentes (Harness, Loop e Graph Engineering com State Graphs e DAG de tarefas), permitindo gerar uma fundação arquitetural de excelência para **qualquer tipo de projeto em Go** (APIs REST/gRPC, ferramentas de linha de comando CLI, aplicações Web, daemons/workers em background, bibliotecas ou interfaces de terminal TUI).

---

## Como Utilizar Este Prompt

1. **Copie o texto** da seção [Prompt de Fundação para Novos Projetos Go com Claude Code](#prompt-de-fundação-para-novos-projetos-go-com-claude-code) abaixo.
2. **Substitua as variáveis entre colchetes** pelos dados e tecnologias desejadas para o seu projeto:
   - `[NOME_DO_PROJETO]`: Nome do projeto (ex: `auth-service`, `data-pipeline`, `backup-cli`, `file-server`).
   - `[MODULO_GO]`: Caminho canônico do módulo Go (ex: `github.com/usuario/meu-projeto`).
   - `[TIPO_DE_PROJETO]`: Arquétipo da aplicação (ex: `API REST Headless`, `Ferramenta de Linha de Comando CLI`, `Serviço Web com SSR`, `Worker Assíncrono / Daemon`, `Biblioteca Go`, `Terminal UI (TUI)`).
   - `[BINARIOS_OU_SERVICOS]`: Nomes dos pontos de entrada/executáveis compilados sob `cmd/` (ex: `cmd/backup-cli`, `cmd/api` e `cmd/worker`, `cmd/server`).
   - `[STACK_E_FRAMEWORKS]`: Bibliotecas, frameworks e drivers específicos que você deseja utilizar no projeto (ex: `Chi router + pgx para PostgreSQL`, `Cobra CLI`, `Go net/http stdlib + embed.FS + Tailwind CSS`, `Bubble Tea para TUI`, ou `Apenas biblioteca padrão de Go`).
   - `[PLATAFORMAS_ALVO]`: Sistemas operacionais e arquiteturas alvo para compilação e release (ex: `Linux (amd64, arm64), macOS (arm64, amd64), Windows (amd64)` ou `Apenas ambiente corrente`). *(Caso não informado, o padrão assumirá a arquitetura corrente ou o agente questionará o usuário se houver necessidade de distribuição multiplataforma)*.
   - `[DESCRICAO_DO_PROJETO]`: Resumo objetivo do propósito, valor de negócio e responsabilidades da aplicação.
3. **Execute no Claude Code** ou passe para a ferramenta de IA como instrução para propor uma nova mudança:
   ```text
   /openspec-propose Crie a especificação de fundação arquitetural e de engenharia 'project-foundation' para o projeto [NOME_DO_PROJETO] seguindo as diretrizes abaixo:
   <cole o prompt preenchido>
   ```

---

### Exemplos Práticos de Preenchimento

#### Exemplo 1: Ferramenta de Linha de Comando (CLI)
- `[NOME_DO_PROJETO]`: `s3-sync-cli`
- `[MODULO_GO]`: `github.com/empresa/s3-sync-cli`
- `[TIPO_DE_PROJETO]`: `Ferramenta de Linha de Comando (CLI)`
- `[BINARIOS_OU_SERVICOS]`: `cmd/s3-sync/main.go` (binário `s3-sync`)
- `[STACK_E_FRAMEWORKS]`: `github.com/spf13/cobra para subcomandos e flags, AWS SDK for Go v2 para operações S3`
- `[PLATAFORMAS_ALVO]`: `Linux (amd64, arm64), macOS (arm64, amd64), Windows (amd64)`
- `[DESCRICAO_DO_PROJETO]`: `Utilitário de linha de comando para sincronização bidirecional de diretórios locais com buckets Amazon S3 com suporte a filtros glob e concorrência configurável.`

#### Exemplo 2: Microserviço / API REST Headless
- `[NOME_DO_PROJETO]`: `order-service`
- `[MODULO_GO]`: `github.com/empresa/order-service`
- `[TIPO_DE_PROJETO]`: `API REST e Event Consumer`
- `[BINARIOS_OU_SERVICOS]`: `cmd/api/main.go` (binário `order-api`) e `cmd/worker/main.go` (binário `order-worker`)
- `[STACK_E_FRAMEWORKS]`: `Chi router para endpoints HTTP, pgx/v5 para PostgreSQL com pool de conexões, RabbitMQ amqp091 para mensageria assíncrona`
- `[PLATAFORMAS_ALVO]`: `Linux (amd64, arm64)`
- `[DESCRICAO_DO_PROJETO]`: `Microserviço de gestão de pedidos de e-commerce, processamento de pagamentos e consumo de eventos de estoque com persistência relacional.`

#### Exemplo 3: Aplicação Web com Renderização no Servidor (SSR)
- `[NOME_DO_PROJETO]`: `metrics-dashboard`
- `[MODULO_GO]`: `github.com/empresa/metrics-dashboard`
- `[TIPO_DE_PROJETO]`: `Serviço Web com Interface SSR e Streaming`
- `[BINARIOS_OU_SERVICOS]`: `cmd/dashboard/main.go` (binário `metrics-dashboard`)
- `[STACK_E_FRAMEWORKS]`: `Go html/template embutido com go:embed, HTMX para atualizações parciais do DOM, Tailwind CSS standalone CLI para estilização`
- `[PLATAFORMAS_ALVO]`: `Linux (amd64, arm64), Windows (amd64, arm64)`
- `[DESCRICAO_DO_PROJETO]`: `Painel de monitoramento e visualização de métricas de infraestrutura em tempo real compilado em binário único 100% autocontido.`

---

## Prompt de Fundação para Novos Projetos Go com Claude Code

```text
Você é um Arquiteto de Software Principal, Engenheiro Líder em Go (Golang) e Especialista em Engenharia de Agentes de IA (Claude Code).

Sua tarefa é criar a especificação completa de fundação arquitetural e de engenharia de software intitulada 'project-foundation' para o novo projeto '[NOME_DO_PROJETO]' (Módulo Go: '[MODULO_GO]', Tipo: '[TIPO_DE_PROJETO]', Binários/Pontos de Entrada: '[BINARIOS_OU_SERVICOS]', Stack & Tecnologias: '[STACK_E_FRAMEWORKS]', Plataformas Alvo: '[PLATAFORMAS_ALVO]'), cuja descrição é:
"[DESCRICAO_DO_PROJETO]"

A especificação deve ser gerada utilizando o padrão OpenSpec em Português do Brasil (PT-BR) e cobrir integralmente os 10 pilares fundamentais descritos a seguir, produzindo os artefatos: 'proposal.md', 'design.md', 'specs/project-foundation/spec.md', 'tasks.md', 'openspec/config.yaml', 'CLAUDE.md', '.claude/settings.json' e as regras complementares sob '.claude/rules/'.

================================================================================
PILAR 1: PADRÃO ARQUITETURAL E ESTRUTURA CANÔNICA EM GO (CLEAN ARCHITECTURE)
================================================================================
- Adotar o layout canônico de pastas da comunidade Go adaptado para a finalidade do projeto:
  * cmd/[BINARIOS_OU_SERVICOS]/: Pontos de entrada da aplicação, parse de flags/configurações e composição explícita de dependências (composition root).
  * internal/core/domain/: Entidades e modelos de negócio puros, livres de acoplamentos externos ou de infraestrutura.
  * internal/core/ports/: Interfaces e contratos formais de entrada (casos de uso/serviços de aplicação) e saída (repositórios, adaptadores de I/O, clientes externos).
  * internal/core/services/: Implementação dos casos de uso e regras de negócio, dependendo unicamente de domain e ports.
  * internal/adapters/: Adaptadores de entrada (HTTP handlers, CLI commands, gRPC servers, workers) e saída (bancos de dados, filas, sistemas de arquivos, APIs externas).
  * internal/version/: Pacote centralizado para metadados de versão, commit e data de compilação.
  * internal/testutils/: Helpers, fixtures e utilitários auxiliares de teste.
  * pkg/: (Opcional) Código explicitamente público que pode ser consumido por módulos externos como biblioteca.
  * scripts/: Scripts de suporte, automação de cobertura e verificação de integridade.
  * .github/workflows/: Pipelines de automação CI/CD e release multiplataforma.
- Regra de Isolamento Estrito: Toda lógica interna, regras de negócio e adaptadores privados devem residir sob 'internal/'.
- Injeção de Dependências Explícita: Realizada manualmente via construtores tipados (NewService, NewAdapter) no composition root em 'cmd/', evitando reflexão em tempo de execução ou frameworks de DI invasivos.
- Propagação Idiomática de Contexto: Uso obrigatório de 'context.Context' como primeiro parâmetro em funções com operações de I/O, concorrência ou ciclo de vida, garantindo suporte nativo a encerramento gracioso (graceful shutdown).
- Fornecer um serviço de domínio de referência inicial funcional e testado (ex: HealthCheckService com status e versão da aplicação).

================================================================================
PILAR 2: PONTOS DE ENTRADA, MODULARIDADE E VERSIONAMENTO DINÂMICO (LDFLAGS)
================================================================================
- Estruturação modular dos executáveis e subcomandos em 'cmd/':
  * Cada binário declarado em '[BINARIOS_OU_SERVICOS]' deve possuir seu próprio subdiretório em 'cmd/' com um 'main.go' enxuto.
  * Suporte a comando ou flag 'version' para inspeção rápida da compilação.
  * Adoção de biblioteca de CLI/parsing conforme especificado em '[STACK_E_FRAMEWORKS]' (ex: Cobra, stdlib flag, urfave/cli).
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
- Criar um 'Makefile' universal, autodocumentado e determinístico como interface central de automação:
  * make help (alvo padrão): Exibe dinamicamente o menu com todos os comandos disponíveis e descrições formatadas.
  * make setup: Instala e valida ferramentas de desenvolvimento locais (golangci-lint, govulncheck, e ferramentas específicas do projeto).
  * make dev: Inicia o ambiente de desenvolvimento local (com live-reload ou watch quando aplicável).
  * make run: Executa a aplicação/serviço a partir do código fonte.
  * make test: Executa a suíte completa de testes com verificação de cobertura (barreira >= 80%).
  * make test-unit: Executa testes unitários rápidos.
  * make test-coverage: Gera o relatório HTML de cobertura e valida o limiar de 80%.
  * make lint: Executa o linter estrito 'golangci-lint' com regras de qualidade, complexidade e bugs (.golangci.yml).
  * make fmt: Formata o código Go (gofmt, goimports).
  * make check: Quality Gate local unificado (fmt + lint + govulncheck + test com cobertura).
  * make build: Compila o(s) binário(s) de produção otimizado(s) com stripping de símbolos (-s -w) e injeção de ldflags para o ambiente local/corrente.
  * make build-all: Compila binários para todas as plataformas e arquiteturas definidas em '[PLATAFORMAS_ALVO]' (ou compilação local se nenhuma matriz multiplataforma for especificada).
  * make clean: Limpa binários compilados em dist/, relatórios de cobertura e artefatos temporários.

================================================================================
PILAR 5: STACK DE APLICAÇÃO, DEPENDÊNCIAS E DECISÕES ARQUITETURAIS CUSTOMIZÁVEIS
================================================================================
- As decisões de frameworks, bibliotecas externas, drivers, estratégias de empacotamento, persistência e ferramentas de desenvolvimento são deliberadas em conjunto com o solicitante do prompt conforme o contexto e escopo do projeto em '[STACK_E_FRAMEWORKS]' e '[TIPO_DE_PROJETO]'.
- Princípio da Parcimônia:
  * Priorizar a biblioteca padrão de Go (stdlib) sempre que atender com elegância e eficiência aos requisitos.
  * Quando bibliotecas de terceiros forem necessárias, selecionar pacotes maduros, seguros, ativamente mantidos e com baixo acoplamento/árvore de dependências enxuta.
- Decisões Sob Demanda (sem imposições pré-fabricadas):
  * Estratégias como empacotamento embutido (//go:embed), servidores HTTP/gRPC, ferramentas de live-reload no desenvolvimento local (Air), ORMs ou migrations SQL NÃO devem ser impostas arbitrariamente: só devem ser adotadas se fizerem sentido para o projeto e forem expressamente acordadas no prompt.

================================================================================
PILAR 6: ENGENHARIA DE AGENTES, PRIORIDADE DE FERRAMENTAS E ECONOMIA DE TOKENS (CLAUDE CODE)
================================================================================
- Documentar e configurar em 'CLAUDE.md', '.claude/settings.json' e '.claude/rules/' as diretrizes de operação e autonomia do Claude Code:
  1. Prioridade Mandatória de Ferramentas Nativas (Native Tool Grounding):
     - Criar arquivos: Utilizar estritamente a ferramenta nativa de escrita de arquivos ('Write'). É proibido executar 'cat << EOF', 'echo >' ou 'touch' via terminal/Bash.
     - Editar arquivos: Utilizar a ferramenta nativa de edição ('Edit') para alterações pontuais e cirúrgicas. É proibido utilizar scripts 'sed', 'awk' ou redirecionamentos 'cat >' no terminal/Bash.
     - Inspecionar / Ler arquivos: Utilizar a ferramenta nativa de visualização ('View') especificando os limites de linhas quando aplicável para focar no trecho relevante. É proibido executar 'cat', 'head', 'tail' via terminal/Bash.
     - Buscar código: Utilizar a ferramenta nativa de busca em conteúdo ('GrepTool' / grep nativo). É proibido executar 'grep' ou 'ripgrep' arbitrariamente via Bash.
     - Localizar arquivos: Utilizar a ferramenta nativa de busca por padrão/nome ('GlobTool' / glob nativo). É proibido executar 'find' ou 'ls -R' via Bash.
     - Listar diretórios: Utilizar a ferramenta nativa de listagem ('LS'). É proibido executar 'ls' via terminal/Bash.
  2. Uso Restrito do Terminal ('Bash'):
     - O terminal/Bash deve ser utilizado exclusivamente para ferramentas do ciclo de vida: 'make' ('make test', 'make lint', 'make check', 'make build'), 'go' ('go test', 'go mod tidy'), 'git', 'openspec' e binários executáveis.
  3. Autonomia Operacional e Permissões (.claude/settings.json / CLAUDE.md):
     - Configurar permissões e diretrizes pré-autorizadas para ferramentas cotidianas ('go', 'make', 'git', 'openspec', linters), eliminando fricção e prompts repetitivos para assegurar execução autônoma contínua.
  4. Economia Ativa de Tokens, Curação de Contexto e Otimização de Janela:
     - Respostas do agente devem ser concisas e estruturadas em Markdown com referências diretas a arquivos.
     - NUNCA duplicar blocos massivos de código já existentes no disco na resposta do chat.
     - Leitura cirúrgica de arquivos e edições pontuais com a ferramenta nativa 'Edit'.
     - Execução de comandos com saídas concisas e estruturadas (sem flags excessivamente verbosas).
  5. Disciplinas de Engenharia de Agentes em Sistemas Compostos (Harness, Loop & Graph Engineering):
     - Harness Engineering: Arcabouço operacional e scaffolding de IA composto por diretrizes determinísticas (CLAUDE.md, .claude/settings.json, .claude/rules/), matriz de permissões, guardrails de segurança, contexto otimizado e limites de sandbox para governar a operação autônoma, previsível e segura do Claude Code.
     - Loop Engineering: Ciclos cognitivos e iterativos de execução e auto-validação contínua do agente (ReAct / Reflection loops: inspeção cirúrgica de requisitos -> intervenção pontual no código -> execução de testes automatizados -> diagnóstico de erros -> correção e validação final com 'make check'), prevenindo loops infinitos e garantindo convergência rápida.
     - Graph Engineering (State Graphs & DAG de Raciocínio): Modelagem e orquestração do fluxo cognitivo/operacional do agente como um grafo direcionado acíclico (DAG) de tarefas, estados e decisões. Decomposição topológica de dependências entre tarefas de planejamento e execução, controle de transições determinísticas de estado entre fases (especificação -> implementação -> verificação -> reflexão), roteamento de fluxos e coordenação/paralelismo entre subtarefas especializadas sem dependências circulares ou bloqueios cognitivos.
     - Gestão de Contexto e Subtarefas: Gestão enxuta do contexto de trabalho para manter a janela cognitiva limpa, focada e eficiente.

================================================================================
PILAR 7: GOVERNANÇA DE ESPECIFICAÇÕES COM OPENSPEC (RECOMENDADO PARA BOAS PRÁTICAS)
================================================================================
- Recomendação Mandatória de Boas Práticas: É altamente recomendado adotar o framework 'OpenSpec' como a ferramenta padrão de governança de especificações, rastreabilidade de mudanças e alinhamento contínuo entre Product Owner (PO) e Quality Assurance (QA).
- Configurar 'openspec/config.yaml' promovendo entendimento mútuo entre PO e QA:
  * context: Declaração da stack do projeto (Go, Clean Architecture, tecnologias escolhidas, TDD/BDD >= 80%, Makefile, PT-BR, Autonomia e Ferramentas Nativas do Claude Code).
  * rules.proposal: Foco no 'porquê' (motivação de negócio) e 'o que muda' (escopo funcional e técnico), declaração explícita de impacto e idioma PT-BR.
  * rules.specs: Escrita em PT-BR sob ótica CONJUNTA e colaborativa de PO e QA:
    - **Diretrizes para o Product Owner (PO)**:
      * Redação em linguagem clara, ubíqua e acessível ao negócio, descrevendo o valor entregue sem detalhes internos de implementação (como nomes de funções ou classes) que obscureçam a regra.
      * Seção '## Purpose' obrigatória (mínimo de 50 caracteres) explicando claramente o objetivo da capacidade para o produto e seus usuários.
      * Requisitos funcionais ('### Requirement: <Nome>') com critérios de aceitação objetivos e verificáveis.
    - **Diretrizes para o Quality Assurance (QA) e Automação de Testes**:
      * Estruturação formal de cenários BDD/Gherkin com 4 hashtags '#### Scenario: <Nome>' e bullets padronizados '- **WHEN**' e '- **THEN**' (e opcionalmente '- **GIVEN**').
      * Cobertura determinística de fluxos principais (caminho feliz), fluxos alternativos, validação de limites/bordas e tratamento de erros e exceções.
      * Cenários redigidos de forma precisa para que ferramentas de automação de testes (como Godog/Cucumber para Go, Playwright, testes de contrato ou testes E2E) consigam traduzir e validar os comportamentos de ponta a ponta sem ambiguidades.
  * rules.tasks: Seções numeradas '## N. Nome do Grupo', checkboxes '- [ ] N.M Descrição', tarefas explícitas de testes unitários/BDD com validação de cobertura >= 80% e 'make check'.
  * operations.apply.guidance: Executar 'make test' e 'make check' para validar alterações antes de concluir tarefas, priorizando ferramentas nativas.

================================================================================
PILAR 8: PIPELINES DE CI/CD E MATRIZ DE RELEASE (PLATAFORMAS CUSTOMIZÁVEIS)
================================================================================
- Workflow de CI ('.github/workflows/ci.yml'):
  * Executado em Pull Requests e pushes na branch principal ('main').
  * Etapas: Checkout, Setup Go, golangci-lint, govulncheck, execução de testes com barreira de cobertura >= 80% e validação de integridade OpenSpec ('openspec validate --all').
- Workflow de Release Multiplataforma ('.github/workflows/release.yml'):
  * Disparado na criação de tags de release Git ('v*').
  * Matriz de Compilação parametrizada estritamente conforme '[PLATAFORMAS_ALVO]':
    - Cross-compilação com injeção de versão via '-ldflags' para cada par SO/Arquitetura especificado (ex: Linux amd64/arm64, Darwin arm64/amd64, Windows amd64/arm64).
    - Empacotamento adequado para cada plataforma (arquivos '.tar.gz' para Unix/Linux/macOS e '.zip' para Windows com binários '.exe').
    - Geração de checksums SHA256 e publicação automática de todos os artefatos compilados na Release do GitHub.
- Boas Práticas de Prompt e Regra de Fallback:
  * Caso as plataformas alvo ('[PLATAFORMAS_ALVO]') NÃO sejam informadas pelo usuário no prompt:
    1. O agente deve assumir como padrão a compilação exclusiva para a plataforma/arquitetura corrente do host.
    2. Se o projeto demandar distribuição de binários para terceiros e houver ambiguidade, o agente DEVE questionar proativamente o usuário antes de assumir matrizes arbitrárias.

================================================================================
PILAR 9: DOCUMENTAÇÃO VIVA E EXAUSTIVA NO README.MD
================================================================================
- Estruturar o 'README.md' raiz com as seguintes seções obrigatórias:
  1. Header & Badges: Título com badges funcionais (Go Version, CI Quality Gate, Test Coverage >= 80%, Latest Release, License).
  2. Visão Geral da Aplicação: Propósito, valor entregue e diferenciais técnicos.
  3. Guia de Instalação e Uso:
     - Download dos binários pré-compilados e compilação a partir do código fonte.
     - Documentação completa e tabela com comandos, flags, parâmetros de configuração e exemplos práticos de execução.
  4. Guia do Desenvolvedor:
     - Pré-requisitos e setup do ambiente ('make setup').
     - Arquitetura de software e layout de diretórios ('cmd/', 'internal/core/', 'internal/adapters/').
     - Interface universal de comandos e automação via Makefile ('make check', 'make dev', 'make test', 'make build-all').
     - Metodologia de testes TDD/BDD com barreira de 80%.
     - Diretrizes de IA (Claude Code), autonomia de ferramentas e padrão Conventional Commits.

================================================================================
PILAR 10: GOVERNANÇA GIT, HIGIENE DE REPOSITÓRIO E SQUASH MERGE
================================================================================
- Criação e Configuração Mandatória de '.gitignore' Idiomático:
  * Criar imediatamente na raiz do repositório um arquivo '.gitignore' abrangente, cobrindo:
    - Binários compilados e executáveis: 'bin/', 'dist/', '*.exe', '*.exe~', '*.dll', '*.so', '*.dylib'.
    - Relatórios e saídas de teste/cobertura: 'coverage.out', 'coverage.html', 'coverage.txt', '*.prof'.
    - Segredos e variáveis de ambiente: '.env', '.env.*', '*.pem', '*.key', '*.crt' (exceto templates públicos como '.env.example').
    - Caches e ferramentas de build: '.cache/', '.air.tmp/', 'vendor/'.
    - Configurações de IDEs e do Sistema Operacional: '.vscode/', '.idea/', '*.swp', '.DS_Store', 'Thumbs.db'.
- Importância e Boas Práticas de Governança Git em Engenharia com Agentes de IA:
  * Higiene e Economia de Contexto: Manter o repositório rigorosamente livre de artefatos de compilação, relatórios temporários e dependências locais evita que ferramentas de busca e I/O do Claude Code ('GrepTool', 'GlobTool', 'LS') indexem ou leiam lixo binário, protegendo a janela de contexto de alucinações e desperdício de tokens.
  * Segurança Operacional Inegociável: Prevenção estrita contra commits acidentais de credenciais, certificados e segredos na árvore Git.
  * Padrão Estrito de Conventional Commits: Todos os commits devem seguir a convenção semântica ('feat:', 'fix:', 'refactor:', 'test:', 'docs:', 'chore:', 'perf:', 'ci:'), assegurando rastreabilidade clara e automação facilitada de changelogs e releases.
  * Feature Branch por Especificação: Toda nova mudança proposta pelo OpenSpec deve ser desenvolvida em branch dedicada a partir da 'main' ('git checkout -b feature/<change-name>'), garantindo que intervenções de IA ocorram em ambiente isolado (sandbox).
  * Permissão Obrigatória do Usuário: O agente de IA NUNCA deve realizar merge na 'main' ou abrir Pull Request sem autorização explícita e prévia do usuário.
  * Estratégia de Integração Exclusivamente SQUASH:
    - Merge Local: 'git checkout main && git merge --squash feature/<change-name> && git commit -m "feat(<change-name>): <resumo consolidado das mudanças>"'.
    - Pull Request (GitHub): Configurar e executar a integração exclusivamente via Squash and Merge ('gh pr merge --squash').
    - Racional: O Squash Merge consolida múltiplos micro-commits de desenvolvimento e experimentação em um único commit atômico e testado na 'main', garantindo histórico limpo, linear, legível e 100% bisectável ('git bisect') e reversível ('git revert').
```

---

## Estrutura Esperada dos Artefatos OpenSpec Gerados

Quando o prompt acima for executado pelo Claude Code, ele deve produzir a seguinte árvore de artefatos no repositório:

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
CLAUDE.md
.claude/
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
- [x] **`design.md`**: Detalha *Context*, *Goals/Non-Goals*, *Decisions* (arquitetura Go, testes BDD, Makefile, stack escolhida, ldflags, release matrix, README, CI/CD, prioridade de ferramentas nativas do Claude Code) e *Risks/Trade-offs*.
- [x] **`specs/project-foundation/spec.md`**: Define os requisitos com perspectivas explícitas de **PO** e **QA**, cada um acompanhado de cenários BDD `#### Scenario:` com `- **WHEN**` e `- **THEN**`.
- [x] **`tasks.md`**: Lista de tarefas organizadas por seções numeradas (`## N. Nome do Grupo`), com subtarefas granulares `- [ ] N.M`, incluindo verificações de compilação, testes unitários, linter, cobertura >= 80% e `openspec validate --all`.
- [x] **`openspec/config.yaml`**: Configuração central com `context`, `rules` (proposal, specs, design, tasks) e `operations` (apply).
- [x] **`CLAUDE.md`**: Diretrizes operacionais de alta prioridade para autonomia do Claude Code, ferramentas nativas (`Write`, `Edit`, `View`, `GrepTool`, `GlobTool`, `LS`), uso restrito de `Bash`, economia de tokens e fluxo Git squash.
- [x] **`.claude/settings.json`**: Permissões e configurações pré-autorizadas para máxima autonomia sem interrupções.
- [x] **`.claude/rules/`**: Diretrizes operacionais para Harness de IA (scaffolding e guardrails), Loop e Graph Engineering (State Graphs e DAG de tarefas), Autonomia de Ferramentas, Git Branching Workflow, Arquivamento com Squash e Convenções Go.
