# Prompt Mestre: Fundação Genérica de Projetos Node.js / JavaScript Moderno com OpenSpec e Claude Code

Este documento contém o **Prompt Mestre de Fundação de Projetos Node.js / JavaScript Moderno** (`generic-nodejs-foundation-spec-prompt-claude`), projetado para ser utilizado com o **Claude Code** e o framework **OpenSpec** (por exemplo, via comando `/openspec-propose` ou como instrução inicial).

É altamente **recomendado adotar o framework OpenSpec** como padrão de excelência para a governança contínua de especificações, documentação viva, rastreabilidade de decisões e alinhamento transparente entre visão de produto (PO) e qualidade técnica/automação (QA).

Ele condensa todos os requisitos de arquitetura limpa em ESM nativo (`"type": "module"`), qualidade de código, automação de comandos via Makefile, esteira de testes com TDD/BDD, barreira inegociável de cobertura &ge; 80%, CI/CD, governança OpenSpec, autonomia operacional e Engenharia de Agentes (Harness, Loop e Graph Engineering com State Graphs e DAG de tarefas), permitindo gerar uma fundação arquitetural de excelência para **qualquer tipo de projeto em Node.js** (APIs REST/GraphQL/gRPC, ferramentas de linha de comando CLI, microsserviços, daemons/workers assíncronos de filas, aplicações Web SSR ou bibliotecas/módulos npm).

---

## Como Utilizar Este Prompt

1. **Copie o texto** da seção [Prompt de Fundação para Novos Projetos Node.js](#prompt-de-fundação-para-novos-projetos-nodejs) abaixo.
2. **Substitua as variáveis entre colchetes** pelos dados e tecnologias desejadas para o seu projeto:
   - `[NOME_DO_PROJETO]`: Nome do projeto (ex: `task-runner`, `payment-api`, `backup-cli`, `queue-worker`, `file-server`).
   - `[PACOTE_NODE]`: Nome do pacote no `package.json` (ex: `@empresa/task-runner`, `backup-cli`, `@org/payment-service`).
   - `[TIPO_DE_PROJETO]`: Arquétipo da aplicação (ex: `API REST Headless`, `Ferramenta de Linha de Comando CLI`, `Serviço Web com SSR`, `Worker Assíncrono / Queue Consumer`, `Biblioteca / Módulo npm`, `Microsserviço de Backend`).
   - `[ENTRYPOINTS_OU_SERVICOS]`: Nomes dos pontos de entrada/serviços executáveis sob `bin/` ou `src/` (ex: `bin/cli.js`, `src/main.js`, `src/api.js` e `src/worker.js`, `src/server.js`).
   - `[STACK_E_FRAMEWORKS]`: Bibliotecas, frameworks e drivers específicos que você deseja utilizar no projeto (ex: `Fastify + pg (node-postgres)`, `Express + Prisma/Drizzle`, `Hono`, `Commander.js + chalk`, `BullMQ + ioredis`, ou `Apenas módulos nativos do Node.js node:*`).
   - `[PLATAFORMAS_ALVO]`: Versões de runtime Node.js (ex: `Node.js 20 LTS, Node.js 22 LTS`) e/ou plataformas/arquiteturas de distribuição (ex: `Linux (amd64, arm64), macOS (arm64, amd64), Windows (amd64)` ou `Apenas ambiente corrente / Runtime Node.js LTS`). *(Caso não informado, o padrão assumirá o runtime Node.js corrente do host ou o agente questionará o usuário se houver necessidade de empacotamento standalone multiplataforma via SEA/pkg/docker)*.
   - `[DESCRICAO_DO_PROJETO]`: Resumo objetivo do propósito, valor de negócio e responsabilidades da aplicação.
3. **Execute no Claude Code** ou passe para a ferramenta de IA como instrução para propor uma nova mudança:
   ```text
   /openspec-propose Crie a especificação de fundação arquitetural e de engenharia 'project-foundation' para o projeto [NOME_DO_PROJETO] seguindo as diretrizes abaixo:
   <cole o prompt preenchido>
   ```

---

### Exemplos Práticos de Preenchimento

#### Exemplo 1: Ferramenta de Linha de Comando (CLI)
- `[NOME_DO_PROJETO]`: `backup-cli`
- `[PACOTE_NODE]`: `@empresa/backup-cli`
- `[TIPO_DE_PROJETO]`: `Ferramenta de Linha de Comando (CLI)`
- `[ENTRYPOINTS_OU_SERVICOS]`: `bin/backup-cli.js` (executável `backup-cli`)
- `[STACK_E_FRAMEWORKS]`: `Commander.js para subcomandos e flags, chalk para colorização de terminal, @aws-sdk/client-s3 para upload em nuvem`
- `[PLATAFORMAS_ALVO]`: `Node.js 20 LTS e 22 LTS em Linux (x64/arm64), macOS (x64/arm64), Windows (x64)`
- `[DESCRICAO_DO_PROJETO]`: `Utilitário de linha de comando para backup e sincronização compactada de diretórios locais com repositórios S3 com suporte a filtros glob e relatório de progresso.`

#### Exemplo 2: Microserviço / API REST Headless
- `[NOME_DO_PROJETO]`: `order-service`
- `[PACOTE_NODE]`: `@empresa/order-service`
- `[TIPO_DE_PROJETO]`: `API REST e Event Consumer`
- `[ENTRYPOINTS_OU_SERVICOS]`: `src/api.js` (serviço HTTP REST) e `src/worker.js` (processador de filas de eventos)
- `[STACK_E_FRAMEWORKS]`: `Fastify para endpoints HTTP de alta performance, pg (node-postgres) com pool de conexões, BullMQ + ioredis para mensageria assíncrona`
- `[PLATAFORMAS_ALVO]`: `Node.js 20 LTS e 22 LTS (Linux amd64/arm64)`
- `[DESCRICAO_DO_PROJETO]`: `Microsserviço de gestão de pedidos de e-commerce, processamento assíncrono de pagamentos e consumo de eventos de estoque com persistência relacional.`

#### Exemplo 3: Aplicação Web / Serviço com Renderização no Servidor (SSR)
- `[NOME_DO_PROJETO]`: `status-monitor`
- `[PACOTE_NODE]`: `@empresa/status-monitor`
- `[TIPO_DE_PROJETO]`: `Serviço Web com Interface SSR e Streaming`
- `[ENTRYPOINTS_OU_SERVICOS]`: `src/server.js` (serviço web `status-monitor`)
- `[STACK_E_FRAMEWORKS]`: `Express + EJS (ou Fastify + @fastify/view), Server-Sent Events (SSE) para atualizações em tempo real, Tailwind CSS standalone CLI para estilização`
- `[PLATAFORMAS_ALVO]`: `Node.js 20 LTS e 22 LTS (Linux, macOS, Windows)`
- `[DESCRICAO_DO_PROJETO]`: `Painel web de monitoramento de integridade e telemetria de serviços de infraestrutura em tempo real com streaming de eventos e interface responsiva.`

#### Exemplo 4: Worker de Filas / Processamento Assíncrono
- `[NOME_DO_PROJETO]`: `notification-worker`
- `[PACOTE_NODE]`: `@empresa/notification-worker`
- `[TIPO_DE_PROJETO]`: `Worker Assíncrono / Queue Consumer`
- `[ENTRYPOINTS_OU_SERVICOS]`: `src/worker.js` (serviço worker `notification-worker`)
- `[STACK_E_FRAMEWORKS]`: `BullMQ + ioredis para consumo de jobs, nodemailer para disparo transacional de e-mails, Handlebars para compilação de templates de notificação`
- `[PLATAFORMAS_ALVO]`: `Node.js 20 LTS e 22 LTS (Linux amd64/arm64)`
- `[DESCRICAO_DO_PROJETO]`: `Processador de tarefas em background para entrega escalável de notificações multicanal (e-mail, push, webhooks) com controle de retentativas e dead-letter queues.`

---

## Prompt de Fundação para Novos Projetos Node.js

```text
Você é um Arquiteto de Software Principal, Engenheiro Líder em Node.js / JavaScript Moderno (ESM) e Especialista em Engenharia de Agentes de IA (Claude Code).

Sua tarefa é criar a especificação completa de fundação arquitetural e de engenharia de software intitulada 'project-foundation' para o novo projeto '[NOME_DO_PROJETO]' (Pacote: '[PACOTE_NODE]', Tipo: '[TIPO_DE_PROJETO]', Entrypoints/Serviços: '[ENTRYPOINTS_OU_SERVICOS]', Stack & Tecnologias: '[STACK_E_FRAMEWORKS]', Plataformas Alvo: '[PLATAFORMAS_ALVO]'), cuja descrição é:
"[DESCRICAO_DO_PROJETO]"

A especificação deve ser gerada utilizando o padrão OpenSpec em Português do Brasil (PT-BR) e cobrir integralmente os 10 pilares fundamentais descritos a seguir, produzindo os artefatos: 'proposal.md', 'design.md', 'specs/project-foundation/spec.md', 'tasks.md', 'openspec/config.yaml', 'CLAUDE.md' e '.claude/settings.json'.

================================================================================
PILAR 1: PADRÃO ARQUITETURAL EM NODE.JS E LAYOUT CANÔNICO (CLEAN ARCHITECTURE ESM)
================================================================================
- Adotar o padrão Clean Architecture / Ports & Adapters com ECMAScript Modules nativo ("type": "module" no package.json):
  * bin/: Scripts executáveis CLI com shebang (#!/usr/bin/env node) e permissões de execução (quando aplicável).
  * src/core/domain/: Entidades, objetos de valor e modelos de negócio puros em JavaScript ESM, totalmente desacoplados de frameworks, bancos ou I/O.
  * src/core/ports/: Contratos de interfaces e abstrações formais (documentadas com JSDoc @interface / @typedef ou classes base) de entrada (casos de uso) e saída (repositórios, adaptadores de I/O, clientes HTTP, brokers de fila).
  * src/core/services/: Implementação dos casos de uso e regras de negócio da aplicação, dependendo exclusivamente de domain e ports.
  * src/adapters/: Adaptadores de entrada (HTTP controllers/routers Fastify/Express/Hono, comandos CLI, handlers de eventos/filas) e saída (bancos de dados, filesystem, APIs externas, mensageria).
  * src/config/: Carregamento, tipagem e validação de configurações e variáveis de ambiente (process.env).
  * src/utils/ ou src/shared/: Helpers, formatadores e utilitários compartilhados transversais.
  * src/main.js (ou src/server.js, src/cli.js): Composition Root enxuto, responsável por instanciar adaptadores e injetar dependências nos serviços.
  * tests/ (ou test/): Estrutura de testes dividida em 'unit/', 'integration/' e 'fixtures/'.
  * scripts/: Scripts de automação, validação de cobertura ('scripts/coverage.sh') e suporte ao ciclo de vida.
  * .github/workflows/: Pipelines de automação CI/CD e release.
- Regra de Isolamento e Encapsulamento:
  * Toda lógica de negócio reside sob 'src/core/'.
  * Os adaptadores de infraestrutura e frameworks residem sob 'src/adapters/'.
- Injeção de Dependências Explícita (Composition Root):
  * Realizada manualmente via funções construtoras / factories (ex: 'createUserService({ userRepository, emailAdapter })'), evitando frameworks mágicos de injeção em runtime ou estado global mutável.
- Cancelamento Assíncrono com AbortSignal:
  * Uso idiomático de 'AbortSignal' / 'AbortController' como padrão para cancelamento de operações assíncronas, requisições HTTP, streams e I/O.
- Graceful Shutdown Robusto:
  * Manipuladores determinísticos para 'process.on("SIGTERM")' e 'process.on("SIGINT")', garantindo o fechamento seguro de conexões de banco de dados, encerramento de servidores HTTP, esvaziamento de filas pendentes e encerramento com timeout de proteção.
- Serviço de Domínio de Referência Inicial:
  * Fornecer um serviço de referência funcional e testado (ex: HealthCheckService informando status, uptime, versão lida do package.json e ambiente).

================================================================================
PILAR 2: PONTOS DE ENTRADA, MODULARIDADE E VERSIONAMENTO DINÂMICO
================================================================================
- Estruturação de Metadados e Entrypoints no 'package.json':
  * Declaração explícita de `"type": "module"`.
  * Configuração de `"engines": { "node": ">=20.0.0" }`.
  * Campo `"bin"` mapeando os executáveis declarados em '[ENTRYPOINTS_OU_SERVICOS]' para seus scripts correspondentes em 'bin/' ou 'src/'.
  * Campo `"exports"` definindo os pontos de entrada públicos da aplicação/módulo e garantindo encapsulamento de módulos internos.
  * Declaração de scripts npm canônicos delegados ou sincronizados com o Makefile ('start', 'dev', 'test', 'lint', 'build').
- Parsing de Linha de Comando e Modularidade:
  * Adoção de biblioteca de CLI/parsing conforme especificado em '[STACK_E_FRAMEWORKS]' (ex: Commander.js, 'node:util parseArgs' nativo, Yargs).
  * Suporte a subcomandos modulares e flag/comando '--version' e '--help'.
- Versionamento Dinâmico:
  * Leitura dinâmica da versão semântica diretamente do 'package.json' (via 'import { createRequire } from "node:module"' ou leitura assíncrona de 'package.json' relativo a 'import.meta.url'), evitando hardcoding de versões no código fonte.
  * Exibição formatada: "[NOME_DO_PROJETO] vX.Y.Z (Node.js: [versao], env: [NODE_ENV])".

================================================================================
PILAR 3: INFRAESTRUTURA DE TESTES AUTOMATIZADOS, TDD/BDD E BARREIRA >= 80%
================================================================================
- Ferramental de Testes:
  * Utilizar o executor nativo 'node:test' + 'node:assert/strict' do Node.js ou suíte moderna em ESM puro conforme '[STACK_E_FRAMEWORKS]' (ex: Vitest ou Jest em ESM).
- Metodologia BDD:
  * Estruturar cenários de teste através de subtestes BDD declarativos utilizando o padrão:
    test('Given [context] When [action] Then [outcome]', async (t) => { ... })
    ou blocos 'describe / it' estruturados.
- Desacoplamento e Mocks Determinísticos:
  * Mocks determinísticos e isolados baseados estritamente nos contratos de interfaces em 'src/core/ports/' (utilizando 't.mock' nativo do 'node:test' ou funções de fábrica simuladas em 'tests/fixtures/').
- Barreira de Cobertura Inegociável:
  * Script automatizado 'scripts/coverage.sh' que executa a suíte de testes com cálculo de cobertura global via 'c8' ou cobertura nativa do 'node:test --experimental-test-coverage' / 'v8'.
  * O script DEVE falhar com código de saída diferente de zero caso a cobertura global de código (linhas, branches ou funções) for inferior a 80%.
  * Geração automática de relatórios em terminal e em formato LCOV/HTML no diretório 'coverage/'.

================================================================================
PILAR 4: INTERFACE UNIVERSAL DE COMANDOS VIA MAKEFILE AUTODOCUMENTADO
================================================================================
- Criar um 'Makefile' universal, autodocumentado e determinístico como interface central de automação:
  * make help (alvo padrão): Exibe dinamicamente o menu com todos os comandos disponíveis e suas respectivas descrições.
  * make setup: Instala dependências determinísticas via 'npm ci' (ou gerenciador configurado como pnpm/yarn) e valida o ambiente Node.js.
  * make dev: Inicia o ambiente de desenvolvimento com live-reload/watch ('node --watch' ou ferramenta correspondente da stack).
  * make run: Executa a aplicação/serviço a partir do código fonte ('node src/main.js' ou entrypoint equivalente).
  * make test: Executa a suíte completa de testes validando a barreira de cobertura (>= 80%).
  * make test-unit: Executa os testes unitários rápidos de domínio e serviços.
  * make test-coverage: Gera o relatório HTML de cobertura e valida o limiar de 80%.
  * make lint: Executa o linter estrito 'eslint' (ESLint 9 com Flat Config 'eslint.config.js') para detecção de bugs, problemas de escopo e qualidade de código.
  * make fmt: Formata o código do projeto utilizando 'prettier' ('prettier --write .').
  * make fmt-check: Valida a conformidade de formatação sem alterar arquivos ('prettier --check .').
  * make check: Quality Gate local unificado (fmt-check + lint + npm audit para vulnerabilidades + test com barreira de cobertura >= 80%).
  * make build: Prepara artefatos de distribuição, valida compilação/empacotamento ou gera Single Executable Application (SEA) / bundle quando aplicável.
  * make clean: Remove diretórios e artefatos temporários ('coverage/', 'dist/', '.nyc_output/', logs).

================================================================================
PILAR 5: STACK DE APLICAÇÃO, DEPENDÊNCIAS E DECISÕES ARQUITETURAIS CUSTOMIZÁVEIS
================================================================================
- As decisões de frameworks, bibliotecas externas, drivers de banco, mensageria e ferramentas são deliberadas conforme o escopo e contexto do projeto em '[STACK_E_FRAMEWORKS]' e '[TIPO_DE_PROJETO]'.
- Princípio da Parcimônia no Ecossistema Node.js:
  * Priorizar os módulos nativos do Node.js ('node:fs/promises', 'node:path', 'node:http', 'node:crypto', 'node:test', 'node:util', 'node:events', 'node:stream', 'node:worker_threads') sempre que atenderem com elegância, robustez e performance aos requisitos.
  * Quando bibliotecas de terceiros forem necessárias, selecionar pacotes maduros, seguros, ativamente mantidos, com poucas dependências transitivas e suporte nativo a ESM.
- Decisões Sob Demanda (sem imposições pré-fabricadas):
  * ORMs pesados, bundlers complexos (Webpack/Vite/Rollup) ou frameworks específicos NÃO devem ser impostos arbitrariamente: só devem ser incluídos se fizerem sentido para o projeto e forem acordados expressamente no prompt.

================================================================================
PILAR 6: ENGENHARIA DE AGENTES, PRIORIDADE DE FERRAMENTAS E ECONOMIA DE TOKENS (CLAUDE CODE)
================================================================================
- Documentar e configurar em 'CLAUDE.md' e '.claude/settings.json' as diretrizes de operação e autonomia do Claude Code:
  1. Prioridade Mandatória de Ferramentas Nativas (Native Tool Grounding):
     - Criar arquivos: Utilizar estritamente a ferramenta nativa de escrita ('Write'). É proibido executar comandos 'cat << EOF', 'echo >' ou 'touch' via 'Bash'.
     - Editar arquivos: Utilizar 'Edit' para alterações pontuais e cirúrgicas. É proibido executar scripts 'sed', 'awk' ou 'cat >' via 'Bash'.
     - Inspecionar / Ler arquivos: Utilizar 'View' especificando números de linha para focar apenas no trecho relevante. É proibido executar 'cat', 'head', 'tail' via 'Bash'.
     - Buscar código: Utilizar 'GrepTool' ('Grep'). É proibido executar 'grep', 'rg' via 'Bash'.
     - Localizar arquivos: Utilizar 'GlobTool' ('Glob'). É proibido executar 'find', 'ls -R' via 'Bash'.
     - Listar diretórios: Utilizar 'LS'. É proibido executar 'ls' via 'Bash'.
  2. Uso Restrito da Ferramenta de Terminal ('Bash'):
     - A ferramenta 'Bash' deve ser utilizada exclusivamente para ferramentas do ciclo de vida: 'make' ('make test', 'make lint', 'make check', 'make build'), 'npm' / 'pnpm' ('npm test', 'npm ci'), 'node', 'git', 'openspec' e binários executáveis.
  3. Autonomia Operacional e Permissões (.claude/settings.json):
     - Configurar permissões pré-autorizadas ('allowedTools' / auto-approve) para todas as ferramentas cotidianas ('npm', 'node', 'make', 'git', 'openspec', 'eslint', 'prettier'), eliminando prompts de confirmação repetitivos e assegurando execução autônoma contínua.
  4. Economia Ativa de Tokens, Curação de Contexto e Otimização de Janela:
     - Respostas do agente devem ser concisas e estruturadas em Markdown com links clicáveis '[arquivo](file:///caminho)'.
     - NUNCA duplicar blocos massivos de código já existentes no disco na resposta do chat.
     - Leitura cirúrgica de arquivos limitando visualizações parciais com 'View' e edições pontuais com 'Edit'.
     - Execução de comandos com saídas concisas e estruturadas (sem flags excessivamente verbosas).
  5. Disciplinas de Engenharia de Agentes em Sistemas Compostos (Harness, Loop & Graph Engineering):
     - Harness Engineering: Arcabouço operacional e scaffolding de IA composto por diretrizes determinísticas ('CLAUDE.md'), matriz de permissões pré-autorizadas ('.claude/settings.json'), guardrails de segurança, contexto otimizado e limites de sandbox para governar a operação autônoma, previsível e segura do Claude Code.
     - Loop Engineering: Ciclos cognitivos e iterativos de execução e auto-validação contínua do agente (ReAct / Reflection loops: inspeção cirúrgica de requisitos -> intervenção pontual no código -> execução de testes automatizados -> diagnóstico de erros -> correção e validação final com 'make check'), prevenindo loops infinitos e garantindo convergência rápida.
     - Graph Engineering (State Graphs & DAG de Raciocínio): Modelagem e orquestração do fluxo cognitivo/operacional do agente como um grafo direcionado acíclico (DAG) de tarefas, estados e decisões. Decomposição topológica de dependências entre tarefas de planejamento e execução, controle de transições determinísticas de estado entre fases (especificação -> implementação -> verificação -> reflexão), roteamento de fluxos e coordenação/paralelismo entre subtarefas especializadas sem dependências circulares ou bloqueios cognitivos.
     - Gestão de Contexto: Curação contínua e compactação de histórico para manter a janela de contexto do agente limpa, enxuta e focada.

================================================================================
PILAR 7: GOVERNANÇA DE ESPECIFICAÇÕES COM OPENSPEC (RECOMENDADO PARA BOAS PRÁTICAS)
================================================================================
- Recomendação Mandatória de Boas Práticas: É altamente recomendado adotar o framework 'OpenSpec' como a ferramenta padrão de governança de especificações, rastreabilidade de mudanças e alinhamento contínuo entre Product Owner (PO) e Quality Assurance (QA).
- Configurar 'openspec/config.yaml' promovendo entendimento mútuo entre PO e QA:
  * context: Declaração da stack do projeto (Node.js ESM, Clean Architecture, tecnologias escolhidas, TDD/BDD >= 80%, Makefile, PT-BR, Autonomia e Ferramentas Nativas do Claude Code).
  * rules.proposal: Foco no 'porquê' (motivação de negócio) e 'o que muda' (escopo funcional e técnico), declaração explícita de impacto e idioma PT-BR.
  * rules.specs: Escrita em PT-BR sob ótica CONJUNTA e colaborativa de PO e QA:
    - **Diretrizes para o Product Owner (PO)**:
      * Redação em linguagem clara, ubíqua e acessível ao negócio, descrevendo o valor entregue sem detalhes internos de implementação (como nomes de funções ou classes) que obscureçam a regra.
      * Seção '## Purpose' obrigatória (mínimo de 50 caracteres) explicando claramente o objetivo da capacidade para o produto e seus usuários.
      * Requisitos funcionais ('### Requirement: <Nome>') com critérios de aceitação objetivos e verificáveis.
    - **Diretrizes para o Quality Assurance (QA) e Automação de Testes**:
      * Estruturação formal de cenários BDD/Gherkin com 4 hashtags '#### Scenario: <Nome>' e bullets padronizados '- **WHEN**' e '- **THEN**' (e opcionalmente '- **GIVEN**').
      * Cobertura determinística de fluxos principais (caminho feliz), fluxos alternativos, validação de limites/bordas e tratamento de erros e exceções.
      * Cenários redigidos de forma precisa para que ferramentas de automação de testes (como Playwright, Cucumber.js, testes de contrato ou testes de integração) consigam traduzir e validar os comportamentos de ponta a ponta sem ambiguidades.
  * rules.tasks: Seções numeradas '## N. Nome do Grupo', checkboxes '- [ ] N.M Descrição', tarefas explícitas de testes unitários/BDD com validação de cobertura >= 80% e 'make check'.
  * operations.apply.guidance: Executar 'make test' e 'make check' para validar alterações antes de concluir tarefas, priorizando ferramentas nativas.

================================================================================
PILAR 8: PIPELINES DE CI/CD E MATRIZ DE RELEASE (PLATAFORMAS CUSTOMIZÁVEIS)
================================================================================
- Workflow de CI ('.github/workflows/ci.yml'):
  * Executado em Pull Requests e pushes na branch principal ('main').
  * Matriz de Versões Node.js: Executar testes nas versões Node.js LTS suportadas (ex: Node 20.x e Node 22.x).
  * Etapas: Checkout, Setup Node.js com cache de dependências ('npm' / 'pnpm'), 'npm ci', validação de formatação ('make fmt-check'), linter ('make lint'), auditoria de segurança ('npm audit'), execução de testes com barreira de cobertura >= 80% ('scripts/coverage.sh') e validação de integridade OpenSpec ('openspec validate --all').
- Workflow de Release ('.github/workflows/release.yml'):
  * Disparado na criação de tags de release Git ('v*').
  * Parametrizado estritamente conforme '[PLATAFORMAS_ALVO]':
    - Publicação no registro de pacotes npm (quando aplicável).
    - Empacotamento de artefatos (.tar.gz / .zip) ou compilação de binários executáveis standalone (Node.js Single Executable Application - SEA ou contêineres Docker) conforme as plataformas alvo definidas.
    - Geração de checksums SHA256 e publicação automática de todos os artefatos compilados na Release do GitHub.
- Boas Práticas de Prompt e Regra de Fallback:
  * Caso as plataformas alvo ('[PLATAFORMAS_ALVO]') NÃO sejam informadas pelo usuário no prompt:
    1. O agente deve assumir como padrão o runtime Node.js LTS no ambiente corrente do host.
    2. Se o projeto demandar distribuição de binários/pacotes para terceiros e houver ambiguidade, o agente DEVE questionar proativamente o usuário antes de assumir matrizes arbitrárias.

================================================================================
PILAR 9: DOCUMENTAÇÃO VIVA E EXAUSTIVA NO README.MD
================================================================================
- Estruturar o 'README.md' raiz com as seguintes seções obrigatórias:
  1. Header & Badges: Título com badges funcionais (Node.js Version >= 20, CI Quality Gate, Test Coverage >= 80%, npm/release version, License).
  2. Visão Geral da Aplicação: Propósito, valor entregue e diferenciais técnicos.
  3. Guia de Instalação e Uso:
     - Instalação via gerenciador de pacotes (npm/pnpm/yarn) ou clone do repositório.
     - Documentação completa com comandos, flags de CLI, endpoints de API ou variáveis de ambiente de configuração.
     - Exemplos práticos de uso e execução.
  4. Guia do Desenvolvedor:
     - Pré-requisitos (Node.js LTS >= 20.0.0, npm/pnpm, Make) e setup do ambiente ('make setup').
     - Arquitetura de software e layout de diretórios ('src/core/', 'src/adapters/').
     - Interface universal de comandos e automação via Makefile ('make check', 'make dev', 'make test', 'make lint').
     - Metodologia de testes TDD/BDD com barreira inegociável de 80%.
     - Diretrizes de IA (Claude Code), autonomia de ferramentas e padrão Conventional Commits.

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
└── settings.json
```

### Checklist dos Artefatos Gerados:
- [x] **`proposal.md`**: Define o *Why*, *What Changes*, *Capabilities* (`project-foundation`) e *Impact* (código, IA, governança, autonomia).
- [x] **`design.md`**: Detalha *Context*, *Goals/Non-Goals*, *Decisions* (arquitetura Node.js ESM, Clean Architecture, testes BDD com barreira &ge; 80%, Makefile, stack escolhida, release matrix, README, CI/CD, prioridade de ferramentas nativas) e *Risks/Trade-offs*.
- [x] **`specs/project-foundation/spec.md`**: Define os requisitos com perspectivas explícitas de **PO** e **QA**, cada um acompanhado de cenários BDD `#### Scenario:` com `- **WHEN**` e `- **THEN**`.
- [x] **`tasks.md`**: Lista de tarefas organizadas por seções numeradas (`## N. Nome do Grupo`), com subtarefas granulares `- [ ] N.M`, incluindo verificações de setup, testes unitários, linter, cobertura &ge; 80% e `openspec validate --all`.
- [x] **`openspec/config.yaml`**: Configuração central com `context`, `rules` (proposal, specs, design, tasks) e `operations` (apply).
- [x] **`CLAUDE.md`**: Diretrizes operacionais de alta prioridade para autonomia do Claude Code, ferramentas nativas (`Write`, `Edit`, `View`, `GrepTool`, `GlobTool`, `LS`), uso restrito de `Bash`, economia de tokens, Harness/Loop/Graph Engineering e fluxo Git squash.
- [x] **`.claude/settings.json`**: Permissões de desenvolvimento pré-autorizadas para máxima autonomia sem interrupções.
