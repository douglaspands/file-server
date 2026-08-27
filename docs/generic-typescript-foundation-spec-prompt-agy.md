# Prompt Mestre: Fundação Genérica de Projetos TypeScript com OpenSpec e Antigravity CLI

Este documento contém o **Prompt Mestre de Fundação de Projetos TypeScript** (`generic-typescript-foundation-spec-prompt`), projetado para ser utilizado com o **Antigravity CLI** e o framework **OpenSpec** (por exemplo, via comando `/openspec-propose` ou como instrução inicial).

É altamente **recomendado adotar o framework OpenSpec** como padrão de excelência para a governança contínua de especificações, documentação viva, rastreabilidade de decisões e alinhamento transparente entre visão de produto (PO) e qualidade técnica/automação (QA).

Ele condensa todos os requisitos de arquitetura limpa, tipagem estrita, qualidade de código, automação de comandos, esteira de testes com TDD/BDD, CI/CD, governança OpenSpec, autonomia operacional e Engenharia de Agentes (Harness, Loop e Graph Engineering), permitindo gerar uma fundação arquitetural de excelência para **qualquer tipo de projeto em TypeScript** (APIs REST/GraphQL/gRPC com Fastify/NestJS/Hono/Express, ferramentas de linha de comando CLI com Clipanion/Commander/Oclif, aplicações Web SSR/Fullstack, workers/microserviços assíncronos em background, ou bibliotecas NPM tipadas com ESM/CJS e `.d.ts`).

---

## Como Utilizar Este Prompt

1. **Copie o texto** da seção [Prompt de Fundação para Novos Projetos TypeScript](#prompt-de-fundação-para-novos-projetos-typescript) abaixo.
2. **Substitua as variáveis entre colchetes** pelos dados e tecnologias desejadas para o seu projeto:
   - `[NOME_DO_PROJETO]`: Nome do projeto (ex: `billing-service`, `data-sync-cli`, `notification-worker`, `ts-validator-kit`).
   - `[PACOTE_TS]`: Nome do pacote / namespace npm (ex: `@empresa/billing-service`, `@empresa/ts-validator-kit`, `data-sync-cli`).
   - `[TIPO_DE_PROJETO]`: Arquétipo da aplicação (ex: `API REST / GraphQL / gRPC`, `Ferramenta de Linha de Comando (CLI)`, `Worker Assíncrono / Background Job`, `Biblioteca NPM Compartilhada`, `Serviço Web com SSR / Fullstack`).
   - `[ENTRYPOINTS_OU_SERVICOS]`: Nomes dos pontos de entrada/executáveis sob `src/` ou `bin/` (ex: `src/cli.ts` gerando `bin/data-sync`, `src/index.ts` e `src/worker.ts`, `src/server.ts`).
   - `[STACK_E_FRAMEWORKS]`: Bibliotecas, frameworks e drivers específicos que você deseja utilizar no projeto (ex: `Fastify + Zod + Prisma ORM`, `NestJS + TypeORM`, `Hono + Drizzle ORM`, `Clipanion + Zod + Chalk + Ora`, `tsup + Vitest` para biblioteca NPM, ou `Apenas Node.js runtime + TypeScript puro`).
   - `[PLATAFORMAS_ALVO]`: Versões do Node.js/Bun/Deno e sistemas operacionais alvo para execução e release (ex: `Node.js 20.x, 22.x LTS em Linux (amd64, arm64), macOS (arm64, amd64), Windows (x64)` ou `Apenas ambiente corrente`). *(Caso não informado, o padrão assumirá a versão e arquitetura corrente ou o agente questionará o usuário se houver necessidade de distribuição multiplataforma/multi-runtime)*.
   - `[DESCRICAO_DO_PROJETO]`: Resumo objetivo do propósito, valor de negócio e responsabilidades da aplicação.
3. **Execute no Antigravity CLI** ou passe para a ferramenta de IA como instrução para propor uma nova mudança:
   ```text
   /openspec-propose Crie a especificação de fundação arquitetural e de engenharia 'project-foundation' para o projeto [NOME_DO_PROJETO] seguindo as diretrizes abaixo:
   <cole o prompt preenchido>
   ```

---

### Exemplos Práticos de Preenchimento

#### Exemplo 1: Ferramenta de Linha de Comando (CLI)
- `[NOME_DO_PROJETO]`: `cloud-deploy-cli`
- `[PACOTE_TS]`: `@empresa/cloud-deploy-cli`
- `[TIPO_DE_PROJETO]`: `Ferramenta de Linha de Comando (CLI)`
- `[ENTRYPOINTS_OU_SERVICOS]`: `src/cli.ts` (binário `cloud-deploy` via campo `"bin"` do `package.json`)
- `[STACK_E_FRAMEWORKS]`: `Clipanion para parsing tipado de comandos e flags, Zod para validação de esquemas de configuração, Chalk e Ora para UI interativa no terminal, tsup para build ESM/CJS`
- `[PLATAFORMAS_ALVO]`: `Node.js 20.x e 22.x LTS (Linux, macOS, Windows)`
- `[DESCRICAO_DO_PROJETO]`: `Utilitário de linha de comando para automação de provisionamento e deploy de infraestrutura em múltiplos provedores de nuvem com validação estrita de manifestos.`

#### Exemplo 2: Microserviço / API REST Headless
- `[NOME_DO_PROJETO]`: `billing-service`
- `[PACOTE_TS]`: `@empresa/billing-service`
- `[TIPO_DE_PROJETO]`: `API REST e Event Consumer`
- `[ENTRYPOINTS_OU_SERVICOS]`: `src/index.ts` (servidor HTTP Fastify) e `src/worker.ts` (processamento de webhooks/filas)
- `[STACK_E_FRAMEWORKS]`: `Fastify v4+ com fastify-type-provider-zod, Prisma ORM com PostgreSQL, BullMQ + IORedis para mensageria assíncrona, Vitest para testes com v8 coverage`
- `[PLATAFORMAS_ALVO]`: `Node.js 20.x e 22.x LTS (Linux amd64/arm64)`
- `[DESCRICAO_DO_PROJETO]`: `Microserviço de faturamento, recorrência e gestão de assinaturas com emissão de faturas, conciliação bancária e processamento de pagamentos via webhooks.`

#### Exemplo 3: Worker Assíncrono / Background Job
- `[NOME_DO_PROJETO]`: `notification-worker`
- `[PACOTE_TS]`: `@empresa/notification-worker`
- `[TIPO_DE_PROJETO]`: `Worker Assíncrono / Daemon de Mensageria`
- `[ENTRYPOINTS_OU_SERVICOS]`: `src/worker.ts` (entrypoint principal do consumer)
- `[STACK_E_FRAMEWORKS]`: `BullMQ com Redis, Zod para contratos de payloads de eventos, Nodemailer e Twilio SDK para canais de notificação, Pino para structured logging JSON`
- `[PLATAFORMAS_ALVO]`: `Node.js 20.x e 22.x LTS (Linux amd64/arm64)`
- `[DESCRICAO_DO_PROJETO]`: `Worker assíncrono de alta performance para processamento e entrega resiliente de notificações transacionais (E-mail, SMS, Push e WhatsApp) com retry backoff exponencial e DLQ.`

#### Exemplo 4: Biblioteca NPM Compartilhada Tipada
- `[NOME_DO_PROJETO]`: `ts-validator-kit`
- `[PACOTE_TS]`: `@empresa/ts-validator-kit`
- `[TIPO_DE_PROJETO]`: `Biblioteca NPM Compartilhada Tipada`
- `[ENTRYPOINTS_OU_SERVICOS]`: `src/index.ts` (export principal com suporte a subpaths via `"exports"` no `package.json`)
- `[STACK_E_FRAMEWORKS]`: `tsup para geração de bundles ESM e CJS com declarações .d.ts e source maps, TypeBox/Zod para schemas ultrarrápidos, Vitest para suíte de testes unitários`
- `[PLATAFORMAS_ALVO]`: `Node.js 18.x, 20.x, 22.x, Deno, Bun e Browsers modernos`
- `[DESCRICAO_DO_PROJETO]`: `Biblioteca utilitária isomórfica de validação, sanitização e transformação de dados com inferência estática de tipos de alta performance em tempo de compilação.`

---

## Prompt de Fundação para Novos Projetos TypeScript

```text
Você é um Arquiteto de Software Principal, Engenheiro Líder em TypeScript / Node.js e Especialista em Engenharia de Agentes de IA (Antigravity CLI).

Sua tarefa é criar a especificação completa de fundação arquitetural e de engenharia de software intitulada 'project-foundation' para o novo projeto '[NOME_DO_PROJETO]' (Pacote NPM: '[PACOTE_TS]', Tipo: '[TIPO_DE_PROJETO]', Entrypoints/Serviços: '[ENTRYPOINTS_OU_SERVICOS]', Stack & Tecnologias: '[STACK_E_FRAMEWORKS]', Plataformas Alvo: '[PLATAFORMAS_ALVO]'), cuja descrição é:
"[DESCRICAO_DO_PROJETO]"

A especificação deve ser gerada utilizando o padrão OpenSpec em Português do Brasil (PT-BR) e cobrir integralmente os 10 pilares fundamentais descritos a seguir, produzindo os artefatos: 'proposal.md', 'design.md', 'specs/project-foundation/spec.md', 'tasks.md', 'openspec/config.yaml', 'AGENTS.md', 'GEMINI.md', '.agent/settings.json' e as regras em '.agent/rules/'.

================================================================================
PILAR 1: PADRÃO ARQUITETURAL EM TYPESCRIPT E LAYOUT CANÔNICO (CLEAN ARCHITECTURE)
================================================================================
- Adotar a estrutura canônica de Clean Architecture / Ports & Adapters adaptada para TypeScript moderno:
  * src/core/domain/: Entidades, objetos de valor (Value Objects) e tipos de negócio puros, sem dependências de infraestrutura ou frameworks externos.
  * src/core/ports/: Interfaces e contratos formais de entrada (casos de uso / application services) e saída (repositórios, adaptadores de I/O, gateways, clientes externos).
  * src/core/services/: Implementação dos casos de uso e regras de negócio, dependendo unicamente de domain e ports.
  * src/adapters/: Adaptadores de entrada (HTTP controllers/handlers, CLI commands, resolvers, consumers) e saída (bancos de dados, ORMs, HTTP clients, caches, sistemas de arquivos).
  * src/index.ts: Ponto de entrada principal da biblioteca/serviço, exportando contratos ou compondo o composition root.
  * src/cli.ts: Ponto de entrada específico para ferramentas de linha de comando (quando aplicável).
  * src/version.ts: Módulo centralizado para metadados de versão, commit e data de build.
  * tests/ ou src/**/*.spec.ts: Suítes de testes unitários, de integração e BDD com fixtures e helpers dedicados.
  * scripts/: Scripts de automação de cobertura, verificação de integridade e build.
  * .github/workflows/: Pipelines de automação CI/CD e release.
- Configuração Estrita de TypeScript ('tsconfig.json'):
  * "strict": true (com noImplicitAny, strictNullChecks, strictFunctionTypes, etc. ativados).
  * "target": "ES2022" (ou superior) e "moduleResolution": "NodeNext" ou "Bundler" para suporte nativo e robusto a ECMAScript Modules (ESM).
  * "declaration": true e "declarationMap": true para geração automática de tipos (.d.ts).
  * "noUncheckedIndexedAccess": true para máxima segurança com acessos a arrays e objetos dinâmicos.
- Validação em Runtime e Tipagem Estática:
  * Adoção de bibliotecas de validação de esquemas em runtime com inferência estática de tipos (ex: Zod, TypeBox ou Valibot conforme '[STACK_E_FRAMEWORKS]'), garantindo que dados externos (inputs HTTP, variáveis de ambiente, argumentos de CLI) sejam validados na fronteira dos adaptadores antes de atingir o domínio.
- Injeção de Dependências Explícita e Tipada:
  * Realizada manualmente via construtores tipados no composition root (em 'src/index.ts' ou 'src/cli.ts'), eliminando magic strings, reflexão em runtime ou acoplamento a decorators/containers de DI pesados.
- Gerenciamento Assíncrono e Encerramento Gracioso (Graceful Shutdown):
  * Suporte a 'AbortSignal' para cancelamento cooperativo de operações assíncronas e handlers de 'SIGTERM' / 'SIGINT' para liberação determinística de conexões, pools e recursos de I/O.
- Fornecer um serviço de domínio de referência inicial funcional e testado (ex: HealthCheckService com status, versão, uptime e ambiente).

================================================================================
PILAR 2: PONTOS DE ENTRADA, MODULARIDADE E VERSIONAMENTO DINÂMICO
================================================================================
- Estruturação Modular do 'package.json':
  * Campos modernos de manifesto configurados com precisão:
    - "type": "module" para suporte nativo a ESM.
    - "main", "module", "types" e subpath exports via "exports": { ".": { "import": "./dist/index.js", "types": "./dist/index.d.ts" } }.
    - Campo "bin" para executáveis de CLI declarados em '[ENTRYPOINTS_OU_SERVICOS]' (ex: "bin": { "[NOME_DO_PROJETO]": "./dist/cli.js" }).
    - "files": ["dist"] para publicação limpa e sem resíduos no NPM Registry.
- CLI Modular Tipada:
  * Ponto de entrada 'src/cli.ts' com parsing robusto, validação de flags e subcomandos utilizando a biblioteca definida em '[STACK_E_FRAMEWORKS]' (ex: Clipanion, Commander, Oclif).
  * Suporte nativo a subcomando/flag '--version' e '-v' exibindo informações detalhadas de compilação.
- Versionamento Dinâmico (Release vs Dev):
  * Módulo 'src/version.ts' centralizado que expõe: version, commit, buildDate.
  * Injeção de metadados durante o build via bundler ('tsup' / 'esbuild' define/banner) ou leitura segura de metadados do 'package.json':
    - Em release oficial: "[NOME_DO_PROJETO] version: v1.0.0 (commit: abc1234, built at: 2026-08-25T12:00:00Z)".
    - Em desenvolvimento local: "[NOME_DO_PROJETO] version: dev (commit: abc1234, built at: 2026-08-25T12:00:00Z)".

================================================================================
PILAR 3: INFRAESTRUTURA DE TESTES AUTOMATIZADOS, TDD/BDD E BARREIRA >= 80%
================================================================================
- Ferramental Moderno de Testes:
  * Utilizar 'Vitest' (com provedor de cobertura v8) ou 'Jest' (com '@swc/jest' / 'ts-jest') conforme definido em '[STACK_E_FRAMEWORKS]'.
- Metodologia BDD Declarativa:
  * Estruturar cenários de teste através de blocos BDD declarativos utilizando o padrão:
    test('Given [context] When [action] Then [outcome]', async () => { ... })
    ou blocos aninhados:
    describe('Given [context]', () => { it('When [action] Then [outcome]', async () => { ... }) })
- Desacoplamento e Mocks Tipados:
  * Mocks determinísticos e fortemente tipados baseados exclusivamente nos contratos de interfaces em 'src/core/ports/' (utilizando 'vi.fn()' do Vitest ou 'jest.fn()').
- Barreira de Cobertura Inegociável (Quality Gate >= 80%):
  * Script automatizado 'scripts/coverage.sh' ou configuração estrita de thresholds no 'vitest.config.ts' / 'jest.config.ts'.
  * O teste DEVE falhar com código de erro se a cobertura global de código for inferior a 80% em linhas, funções, branches e statements.
  * Geração automática de relatórios em formato texto e HTML no diretório 'coverage/'.

================================================================================
PILAR 4: INTERFACE UNIVERSAL DE COMANDOS VIA MAKEFILE AUTODOCUMENTADO
================================================================================
- Criar um 'Makefile' universal, autodocumentado e determinístico como interface central de automação:
  * make help (alvo padrão): Exibe dinamicamente o menu com todos os comandos disponíveis e descrições formatadas.
  * make setup: Instala dependências via gerenciador de pacotes definido (pnpm / npm / yarn) e valida o ambiente.
  * make dev: Inicia o ambiente de desenvolvimento local com live-reload / watch (ex: 'tsx watch' ou 'tsup --watch').
  * make run: Executa a aplicação/serviço a partir do código compilado ou runtime TypeScript direto.
  * make test: Executa a suíte completa de testes com verificação de cobertura (barreira >= 80%).
  * make test-unit: Executa testes unitários rápidos.
  * make test-coverage: Gera o relatório HTML de cobertura e valida o limiar de 80%.
  * make lint: Executa o linter estrito 'ESLint' com regras de '@typescript-eslint'.
  * make typecheck: Executa a checagem estática de tipos do compilador TypeScript ('tsc --noEmit').
  * make fmt: Formata o código TypeScript e JSON via 'Prettier' / ESLint.
  * make check: Quality Gate local unificado (typecheck + fmt-check + lint + audit + test com barreira >= 80%).
  * make build: Compila o projeto gerando os artefatos otimizados de produção em 'dist/' (bundles ESM/CJS e declarações '.d.ts' via tsup ou tsc).
  * make clean: Limpa artefatos compilados em 'dist/', diretório 'coverage/' e caches temporários.

================================================================================
PILAR 5: STACK DE APLICAÇÃO, DEPENDÊNCIAS E DECISÕES ARQUITETURAIS CUSTOMIZÁVEIS
================================================================================
- As decisões de frameworks, bibliotecas externas, drivers de banco de dados, estratégias de empacotamento e ferramentas de desenvolvimento são deliberadas em conjunto com o solicitante do prompt conforme o contexto e escopo do projeto em '[STACK_E_FRAMEWORKS]' e '[TIPO_DE_PROJETO]'.
- Princípio da Parcimônia:
  * Priorizar módulos e APIs nativas modernas do ecossistema Node.js / JavaScript (Fetch API, Web Streams, Crypto, AbortController) sempre que atenderem com elegância e eficiência aos requisitos.
  * Selecionar pacotes de terceiros maduros, seguros, ativamente mantidos, com suporte nativo a ESM e com tipagem TypeScript de primeira classe (First-Class Types).
- Decisões Sob Demanda (sem imposições pré-fabricadas):
  * Ferramentas e frameworks como Fastify, Hono, NestJS, Express, Prisma, Drizzle, TypeORM, BullMQ, Clipanion, Commander, etc., NÃO devem ser impostos arbitrariamente: só devem ser adotados se fizerem sentido para o projeto e forem expressamente acordados no prompt.

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
     - O terminal deve ser utilizado exclusivamente para ferramentas do ciclo de vida: 'make' ('make test', 'make lint', 'make typecheck', 'make check', 'make build'), gerenciadores de pacotes ('pnpm', 'npm', 'yarn'), 'git', 'openspec' e binários executáveis.
  3. Autonomia Operacional e Permissões (.agent/settings.json):
     - Configurar permissões pré-autorizadas ('allow') para todas as ferramentas cotidianas ('pnpm', 'npm', 'make', 'git', 'openspec', 'tsc', linters), eliminando prompts de confirmação repetitivos e assegurando execução autônoma contínua.
  4. Economia Ativa de Tokens, Curação de Contexto e Otimização de Janela:
     - Respostas do agente devem ser concisas e estruturadas em Markdown com links clicáveis '[arquivo](file:///caminho)'.
     - NUNCA duplicar blocos massivos de código já existentes no disco na resposta do chat.
     - Leitura cirúrgica de arquivos limitando linhas ('StartLine'/'EndLine') e edições pontuais com 'replace_file_content'.
     - Execução de comandos com saídas concisas e estruturadas (sem flags excessivamente verbosas).
  5. Disciplinas de Engenharia de Agentes em Sistemas Compostos (Harness, Loop & Graph Engineering):
     - Harness Engineering: Arcabouço operacional e scaffolding de IA composto por diretrizes determinísticas (.agent/rules/, AGENTS.md, GEMINI.md), matriz de permissões pré-autorizadas (.agent/settings.json), guardrails de tipagem estrita, contexto otimizado e limites de sandbox para governar a operação autônoma, previsível e segura do Antigravity CLI.
     - Loop Engineering: Ciclos cognitivos e iterativos de execução e auto-validação contínua do agente (ReAct / Reflection loops: inspeção cirúrgica de requisitos -> intervenção pontual no código -> execução de typecheck e testes -> diagnóstico de erros -> correção e validação final com 'make check'), prevenindo loops infinitos e garantindo convergência rápida.
     - Graph Engineering (State Graphs & DAG de Raciocínio): Modelagem e orquestração do fluxo cognitivo/operacional do agente como um grafo direcionado acíclico (DAG) de tarefas, estados e decisões. Decomposição topológica de dependências entre tarefas de planejamento e execução, controle de transições determinísticas de estado entre fases (especificação -> implementação -> verificação -> reflexão), roteamento de fluxos e coordenação/paralelismo entre subagentes especializados sem dependências circulares ou bloqueios cognitivos.
     - Gestão de Contexto e Subagentes: Delegação de pesquisas e tarefas isoladas para subagentes dedicados para manter a janela de contexto do agente principal limpa, enxuta e focada.

================================================================================
PILAR 7: GOVERNANÇA DE ESPECIFICAÇÕES COM OPENSPEC (RECOMENDADO PARA BOAS PRÁTICAS)
================================================================================
- Recomendação Mandatória de Boas Práticas: É altamente recomendado adotar o framework 'OpenSpec' como a ferramenta padrão de governança de especificações, rastreabilidade de mudanças e alinhamento contínuo entre Product Owner (PO) e Quality Assurance (QA).
- Configurar 'openspec/config.yaml' promovendo entendimento mútuo entre PO e QA:
  * context: Declaração da stack do projeto (TypeScript, Clean Architecture, tecnologias escolhidas, TDD/BDD >= 80%, Makefile, PT-BR, Autonomia e Ferramentas Nativas do AGY).
  * rules.proposal: Foco no 'porquê' (motivação de negócio) e 'o que muda' (escopo funcional e técnico), declaração explícita de impacto e idioma PT-BR.
  * rules.specs: Escrita em PT-BR sob ótica CONJUNTA e colaborativa de PO e QA:
    - **Diretrizes para o Product Owner (PO)**:
      * Redação em linguagem clara, ubíqua e acessível ao negócio, descrevendo o valor entregue sem detalhes internos de implementação (como nomes de funções ou classes) que obscureçam a regra.
      * Seção '## Purpose' obrigatória (mínimo de 50 caracteres) explicando claramente o objetivo da capacidade para o produto e seus usuários.
      * Requisitos funcionais ('### Requirement: <Nome>') com critérios de aceitação objetivos e verificáveis.
    - **Diretrizes para o Quality Assurance (QA) e Automação de Testes**:
      * Estruturação formal de cenários BDD/Gherkin com 4 hashtags '#### Scenario: <Nome>' e bullets padronizados '- **WHEN**' e '- **THEN**' (e opcionalmente '- **GIVEN**').
      * Cobertura determinística de fluxos principais (caminho feliz), fluxos alternativos, validação de limites/bordas e tratamento de erros e exceções.
      * Cenários redigidos de forma precisa para que ferramentas de automação de testes (como Vitest, Jest, Playwright, Cucumber-js ou testes de contrato) consigam traduzir e validar os comportamentos de ponta a ponta sem ambiguidades.
  * rules.tasks: Seções numeradas '## N. Nome do Grupo', checkboxes '- [ ] N.M Descrição', tarefas explícitas de typecheck, testes unitários/BDD com validação de cobertura >= 80% e 'make check'.
  * operations.apply.guidance: Executar 'make test' e 'make check' para validar alterações antes de concluir tarefas, priorizando ferramentas nativas.

================================================================================
PILAR 8: PIPELINES DE CI/CD E MATRIZ DE RELEASE (PLATAFORMAS CUSTOMIZÁVEIS)
================================================================================
- Workflow de CI ('.github/workflows/ci.yml'):
  * Executado em Pull Requests e pushes na branch principal ('main').
  * Matriz de execução testando nas versões de Node.js especificadas em '[PLATAFORMAS_ALVO]' (ex: Node.js 20.x, 22.x).
  * Etapas: Checkout, Setup Node.js / Package Manager, instalação de dependências congeladas (pnpm install --frozen-lockfile / npm ci), make typecheck, make lint, execução de testes com barreira de cobertura >= 80% e validação de integridade OpenSpec ('openspec validate --all').
- Workflow de Release ('.github/workflows/release.yml'):
  * Disparado na criação de tags de release Git ('v*').
  * Compilação dos artefatos finais de produção ('make build').
  * Publicação do pacote tipado no NPM Registry / GitHub Packages com proveniência e tags de release semânticas, OU empacotamento de binários executáveis standalone (via pkg/sea/bun) conforme '[TIPO_DE_PROJETO]'.
  * Geração de checksums SHA256 e publicação automática de notas de versão na Release do GitHub.
- Boas Práticas de Prompt e Regra de Fallback:
  * Caso as plataformas alvo ('[PLATAFORMAS_ALVO]') NÃO sejam informadas pelo usuário no prompt:
    1. O agente deve assumir como padrão a versão LTS corrente do Node.js no ambiente do host.
    2. Se o projeto demandar distribuição de binários ou suporte multi-runtime (Node.js, Bun, Deno, Navegador) e houver ambiguidade, o agente DEVE questionar proativamente o usuário antes de assumir matrizes arbitrárias.

================================================================================
PILAR 9: DOCUMENTAÇÃO VIVA E EXAUSTIVA NO README.MD
================================================================================
- Estruturar o 'README.md' raiz com as seguintes seções obrigatórias:
  1. Header & Badges: Título com badges funcionais (Node.js Version, TypeScript Strict, CI Quality Gate, Test Coverage >= 80%, NPM Version / Latest Release, License).
  2. Visão Geral da Aplicação: Propósito, valor entregue e diferenciais técnicos.
  3. Guia de Instalação e Uso:
     - Instalação via gerenciador de pacotes ('pnpm add [PACOTE_TS]' ou 'npm install [PACOTE_TS]') ou execução global.
     - Documentação completa e tabela com comandos, flags, variáveis de ambiente, parâmetros de configuração e exemplos práticos de execução.
  4. Guia do Desenvolvedor:
     - Pré-requisitos e setup do ambiente ('make setup').
     - Arquitetura de software e layout de diretórios ('src/core/', 'src/adapters/').
     - Interface universal de comandos e automação via Makefile ('make check', 'make dev', 'make test', 'make build').
     - Metodologia de testes TDD/BDD com barreira de 80%.
     - Diretrizes de IA (Antigravity), autonomia de ferramentas e padrão Conventional Commits.

================================================================================
PILAR 10: GOVERNANÇA GIT, HIGIENE DE REPOSITÓRIO E SQUASH MERGE
================================================================================
- Criação e Configuração Mandatória de '.gitignore' Idiomático:
  * Criar imediatamente na raiz do repositório um arquivo '.gitignore' abrangente para TypeScript/Node.js, cobrindo:
    - Dependências de pacotes: 'node_modules/', '.pnpm-store/', '.yarn/'.
    - Distribuição e builds: 'dist/', 'build/', 'out/', '*.tsbuildinfo', '.turbo/'.
    - Relatórios e saídas de teste/cobertura: 'coverage/', '.nyc_output/', '*.lcov'.
    - Logs e arquivos temporários: 'npm-debug.log*', 'yarn-debug.log*', 'yarn-error.log*', 'pnpm-debug.log*', '*.log'.
    - Segredos e variáveis de ambiente: '.env', '.env.*', '*.pem', '*.key', '*.crt' (exceto templates públicos como '.env.example').
    - Configurações de IDEs e do Sistema Operacional: '.vscode/', '.idea/', '*.swp', '.DS_Store', 'Thumbs.db'.
- Importância e Boas Práticas de Governança Git em Engenharia com Agentes de IA:
  * Higiene e Economia de Contexto: Manter o repositório rigorosamente livre de 'node_modules' e diretórios compilados 'dist/' evita que ferramentas de busca e I/O do agente ('grep_search', 'find_by_name', 'list_dir') indexem dezenas de milhares de arquivos de tipos ou dependências, protegendo a janela de contexto de alucinações e desperdício de tokens.
  * Segurança Operacional Inegociável: Prevenção estrita contra commits acidentais de credenciais, chaves de API e segredos na árvore Git.
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
    └── typescript_conventions.md
```

### Checklist dos Artefatos Gerados:
- [x] **`proposal.md`**: Define o *Why*, *What Changes*, *Capabilities* (`project-foundation`) e *Impact* (código, tipagem, IA, governança, autonomia).
- [x] **`design.md`**: Detalha *Context*, *Goals/Non-Goals*, *Decisions* (Clean Architecture em TypeScript, `tsconfig.json` estrito, testes BDD com Vitest/Jest, Makefile, stack escolhida, empacotamento ESM/CJS com `.d.ts`, release matrix, README, CI/CD, prioridade de ferramentas nativas) e *Risks/Trade-offs*.
- [x] **`specs/project-foundation/spec.md`**: Define os requisitos com perspectivas explícitas de **PO** e **QA**, cada um acompanhado de cenários BDD `#### Scenario:` com `- **WHEN**` e `- **THEN**`.
- [x] **`tasks.md`**: Lista de tarefas organizadas por seções numeradas (`## N. Nome do Grupo`), com subtarefas granulares `- [ ] N.M`, incluindo verificações de compilação/typecheck, testes unitários, linter, cobertura >= 80% e `openspec validate --all`.
- [x] **`openspec/config.yaml`**: Configuração central com `context`, `rules` (proposal, specs, design, tasks) e `operations` (apply).
- [x] **`AGENTS.md` e `GEMINI.md`**: Diretrizes operacionais de alta prioridade para autonomia do AGY, ferramentas nativas, economia de tokens e fluxo Git squash.
- [x] **`.agent/settings.json`**: Permissões de desenvolvimento pré-autorizadas para máxima autonomia sem interrupções.
- [x] **`.agent/rules/`**: Diretrizes operacionais para Harness de IA (scaffolding e guardrails), Loop e Graph Engineering (State Graphs e DAG de tarefas), Autonomia de Ferramentas, Git Branching Workflow, Arquivamento com Squash e Convenções TypeScript.
