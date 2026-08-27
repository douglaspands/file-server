# Prompt Mestre: Fundação Genérica de Projetos Python com OpenSpec e Antigravity CLI

Este documento contém o **Prompt Mestre de Fundação de Projetos Python** (`generic-python-foundation-spec-prompt`), projetado para ser utilizado com o **Antigravity CLI** e o framework **OpenSpec** (por exemplo, via comando `/openspec-propose` ou como instrução inicial).

É altamente **recomendado adotar o framework OpenSpec** como padrão de excelência para a governança contínua de especificações, documentação viva, rastreabilidade de decisões e alinhamento transparente entre visão de produto (PO) e qualidade técnica/automação (QA).

Para o ecossistema Python moderno, é fortemente recomendado e padronizado:
- **`uv` (Astral)** como ferramenta universal e ultrarrápida para gerenciar versões do Python (`uv python`), ambientes virtuais isolados (`uv venv`), resolução de dependências com lockfile determinístico (`uv lock`, `uv sync`) e execução de comandos/ferramentas (`uv run`).
- **`ty` (Astral)** como checador estático de tipos de alta performance, substituindo o legado `mypy` com ganhos drásticos de velocidade de execução.
- **`pyproject.toml`** como manifesto central e mandatório (PEP 621 / PEP 517) para metadados, scripts executáveis e configuração unificada de todas as ferramentas (`[tool.ruff]`, `[tool.pytest.ini_options]`, `[tool.ty]`).

Ele condensa todos os requisitos de arquitetura limpa, tipagem estrita, qualidade de código, automação de comandos via Makefile, esteira de testes com TDD/BDD (`pytest`), CI/CD, governança OpenSpec, autonomia operacional e Engenharia de Agentes (Harness, Loop e Graph Engineering), permitindo gerar uma fundação arquitetural de excelência para **qualquer tipo de projeto em Python moderno** (APIs REST/GraphQL/gRPC, ferramentas de linha de comando CLI, aplicações Web/SSR, daemons/workers assíncronos, pipelines de dados, bibliotecas ou interfaces de terminal TUI).

---

## Como Utilizar Este Prompt

1. **Copie o texto** da seção [Prompt de Fundação para Novos Projetos Python](#prompt-de-fundação-para-novos-projetos-python) abaixo.
2. **Substitua as variáveis entre colchetes** pelos dados e tecnologias desejadas para o seu projeto:
   - `[NOME_DO_PROJETO]`: Nome do projeto no repositório/distribuição (ex: `data-pipeline`, `backup-cli`, `order-service`, `file-server`).
   - `[PACOTE_PYTHON]`: Nome canônico do pacote/módulo Python sob `src/` em snake_case (ex: `data_pipeline`, `backup_cli`, `order_service`, `file_server`).
   - `[TIPO_DE_PROJETO]`: Arquétipo da aplicação (ex: `API REST Headless`, `Ferramenta de Linha de Comando (CLI)`, `Serviço Web com SSR / Dashboard`, `Worker Assíncrono / Daemon / Pipeline de Dados`, `Biblioteca / Pacote Python`, `Terminal UI (TUI)`).
   - `[ENTRYPOINTS_OU_SERVICOS]`: Nomes e caminhos dos pontos de entrada/scripts configurados em `pyproject.toml` `[project.scripts]` ou módulos executáveis (ex: `backup-cli = "backup_cli.main:main"`, `order-api = "order_service.adapters.api.server:main"`, `order-worker = "order_service.adapters.workers.main:main"`).
   - `[STACK_E_FRAMEWORKS]`: Bibliotecas, frameworks, ORMs e drivers específicos que você deseja utilizar no projeto (ex: `FastAPI + SQLAlchemy 2.0 (async) + Pydantic v2 + Alembic`, `Typer + Rich + Boto3`, `Flask + Jinja2 + HTMX + Tailwind CSS`, `Celery + Redis + psycopg3`, ou `Apenas biblioteca padrão de Python`).
   - `[PLATAFORMAS_ALVO]`: Versões do Python suportadas e plataformas/arquiteturas alvo para empacotamento e distribuição (ex: `Python 3.11, 3.12, 3.13 em Linux (x86_64, aarch64), macOS (arm64, x86_64), Windows (x86_64)` ou `Python >=3.11 no ambiente corrente`). *(Caso não informado, o padrão assumirá a versão do Python e arquitetura correntes do host ou o agente questionará o usuário se houver necessidade de matriz multiplataforma)*.
   - `[DESCRICAO_DO_PROJETO]`: Resumo objetivo do propósito, valor de negócio e responsabilidades da aplicação.
3. **Execute no Antigravity CLI** ou passe para a ferramenta de IA como instrução para propor uma nova mudança:
   ```text
   /openspec-propose Crie a especificação de fundação arquitetural e de engenharia 'project-foundation' para o projeto [NOME_DO_PROJETO] seguindo as diretrizes abaixo:
   <cole o prompt preenchido>
   ```

---

### Exemplos Práticos de Preenchimento

#### Exemplo 1: Ferramenta de Linha de Comando (CLI)
- `[NOME_DO_PROJETO]`: `s3-sync-cli`
- `[PACOTE_PYTHON]`: `s3_sync_cli`
- `[TIPO_DE_PROJETO]`: `Ferramenta de Linha de Comando (CLI)`
- `[ENTRYPOINTS_OU_SERVICOS]`: `s3-sync = "s3_sync_cli.cli:app"`
- `[STACK_E_FRAMEWORKS]`: `Typer para comandos/flags, Rich para formatação de terminal e barras de progresso, Boto3 para operações AWS S3`
- `[PLATAFORMAS_ALVO]`: `Python 3.11, 3.12, 3.13 em Linux, macOS e Windows`
- `[DESCRICAO_DO_PROJETO]`: `Utilitário CLI de alta performance para sincronização bidirecional de diretórios locais com buckets Amazon S3 com suporte a filtros glob, concorrência assíncrona e relatórios visuais.`

#### Exemplo 2: Microserviço / API REST Headless
- `[NOME_DO_PROJETO]`: `order-service`
- `[PACOTE_PYTHON]`: `order_service`
- `[TIPO_DE_PROJETO]`: `API REST e Consumidor de Eventos`
- `[ENTRYPOINTS_OU_SERVICOS]`: `order-api = "order_service.adapters.api.server:main"`, `order-worker = "order_service.adapters.workers.main:main"`
- `[STACK_E_FRAMEWORKS]`: `FastAPI para endpoints HTTP, Uvicorn como ASGI server, SQLAlchemy 2.0 assíncrono com asyncpg para PostgreSQL, Alembic para migrações, Pydantic v2 para schemas de validação`
- `[PLATAFORMAS_ALVO]`: `Python 3.11, 3.12 em Linux (x86_64, aarch64)`
- `[DESCRICAO_DO_PROJETO]`: `Microserviço assíncrono de gestão de pedidos de e-commerce, processamento de pagamentos e consumo de eventos de estoque com persistência relacional e rastreamento distribuído.`

#### Exemplo 3: Worker Assíncrono / Processamento de Tarefas
- `[NOME_DO_PROJETO]`: `video-transcoder`
- `[PACOTE_PYTHON]`: `video_transcoder`
- `[TIPO_DE_PROJETO]`: `Worker Assíncrono de Fila`
- `[ENTRYPOINTS_OU_SERVICOS]`: `transcoder-worker = "video_transcoder.adapters.workers.celery_worker:main"`
- `[STACK_E_FRAMEWORKS]`: `Celery para orquestração de tarefas distribuídas, Redis como broker e backend de resultados, ffmpeg-python para transcodificação de mídia, MinIO/S3 para armazenamento de artefatos`
- `[PLATAFORMAS_ALVO]`: `Python 3.11, 3.12 em Linux (x86_64, aarch64)`
- `[DESCRICAO_DO_PROJETO]`: `Sistema escalável de processamento assíncrono de vídeos, conversão de codecs, extração de thumbnails e empacotamento HLS/DASH com monitoramento de métricas.`

#### Exemplo 4: Aplicação Web com Renderização no Servidor (SSR)
- `[NOME_DO_PROJETO]`: `metrics-dashboard`
- `[PACOTE_PYTHON]`: `metrics_dashboard`
- `[TIPO_DE_PROJETO]`: `Serviço Web com Interface SSR e Streaming`
- `[ENTRYPOINTS_OU_SERVICOS]`: `metrics-dashboard = "metrics_dashboard.main:main"`
- `[STACK_E_FRAMEWORKS]`: `FastAPI + Jinja2 para renderização server-side de templates, HTMX para interatividade sem JavaScript pesado, Tailwind CSS standalone para estilização, SQLite3/aiosqlite para persistência local`
- `[PLATAFORMAS_ALVO]`: `Python >=3.11 em Linux, macOS e Windows`
- `[DESCRICAO_DO_PROJETO]`: `Painel web de visualização e monitoramento de telemetria e métricas em tempo real com baixo consumo de memória e interface reativa.`

---

## Prompt de Fundação para Novos Projetos Python

```text
Você é um Arquiteto de Software Principal, Engenheiro Líder em Python e Especialista em Engenharia de Agentes de IA (Antigravity CLI).

Sua tarefa é criar a especificação completa de fundação arquitetural e de engenharia de software intitulada 'project-foundation' para o novo projeto '[NOME_DO_PROJETO]' (Pacote Python: '[PACOTE_PYTHON]', Tipo: '[TIPO_DE_PROJETO]', Entrypoints/Serviços: '[ENTRYPOINTS_OU_SERVICOS]', Stack & Tecnologias: '[STACK_E_FRAMEWORKS]', Plataformas Alvo: '[PLATAFORMAS_ALVO]'), cuja descrição é:
"[DESCRICAO_DO_PROJETO]"

A especificação deve ser gerada utilizando o padrão OpenSpec em Português do Brasil (PT-BR) e cobrir integralmente os 10 pilares fundamentais descritos a seguir, produzindo os artefatos: 'proposal.md', 'design.md', 'specs/project-foundation/spec.md', 'tasks.md', 'openspec/config.yaml', 'AGENTS.md', 'GEMINI.md', '.agent/settings.json' e as regras em '.agent/rules/'.

================================================================================
PILAR 1: PADRÃO ARQUITETURAL EM PYTHON E LAYOUT CANÔNICO (CLEAN ARCHITECTURE)
================================================================================
- Adotar o layout canônico moderno baseado no padrão 'src-layout' (PEP 621 / PEP 517/518), garantindo isolamento estrito contra importações acidentais de pacotes não instalados:
  * pyproject.toml: Arquivo de configuração centralizado, padronizado e MANDATÓRIO para metadados do projeto (PEP 621), build system (hatchling, flit-core ou setuptools), scripts executáveis, dependências de produção/desenvolvimento e configuração de ferramentas (ruff, ty, pytest, coverage).
  * src/[PACOTE_PYTHON]/main.py ou cli.py: Ponto de entrada da aplicação, parse de argumentos/configurações e composição explícita de dependências (composition root).
  * src/[PACOTE_PYTHON]/core/domain/: Entidades, modelos de negócio puros e objetos de valor (dataclasses imutáveis ou modelos Pydantic v2), enums, exceções de domínio e regras invariantes, totalmente livres de acoplamentos externos ou de frameworks de infraestrutura.
  * src/[PACOTE_PYTHON]/core/ports/: Interfaces e contratos formais de entrada (casos de uso/serviços de aplicação) e saída (repositórios, adaptadores de I/O, clientes externos, mensageria) definidos explicitamente via 'typing.Protocol' (PEP 544) ou classes abstratas ('abc.ABC').
  * src/[PACOTE_PYTHON]/core/services/: Implementação dos casos de uso e orquestração das regras de negócio, dependendo unicamente de domain e ports.
  * src/[PACOTE_PYTHON]/adapters/: Adaptadores de entrada (routers HTTP/FastAPI/Flask, comandos CLI Typer/Click, workers Celery, controllers) e adaptadores de saída (repositórios SQLAlchemy/SQLModel, clientes Redis, clientes HTTP httpx/requests, adaptadores de sistema de arquivos).
  * src/[PACOTE_PYTHON]/version.py: Pacote/módulo centralizado para metadados de versão semântica, commit Git e data de build.
  * tests/: Estrutura espelhada de testes dividida em 'unit/', 'integration/', 'e2e/' e 'conftest.py' centralizado para fixtures reutilizáveis.
  * scripts/: Scripts auxiliares de automação, cálculo estrito de cobertura de testes ('scripts/coverage.sh') e checagens de qualidade.
  * .github/workflows/: Pipelines de automação CI/CD e release para PyPI / GitHub Releases.
- Tipagem Estrita Obrigatória (Type Hints) e Checagem Ultrarrápida com 'ty':
  * Adoção de Type Hints estritos (PEP 484, PEP 585, PEP 604) em 100% das assinaturas de funções, métodos e atributos de classes.
  * Validação estrita por checagem estática de tipos de alta performance utilizando a ferramenta 'ty' da Astral ('ty src tests' ou 'uv run ty') em vez de ferramentas legadas mais lentas.
- Injeção de Dependências Explícita (Composition Root):
  * Realizada manualmente nos pontos de entrada (main/cli/factories) via construtores tipados, sem o uso de frameworks invasivos de DI ou reflexão oculta em tempo de execução.
- Encerramento Gracioso (Graceful Shutdown) e Gerenciamento de Contexto:
  * Tratamento nativo de sinais de término do sistema operacional (SIGTERM, SIGINT) através de 'signal' ou 'asyncio' / context managers assíncronos ('asynccontextmanager').
  * Liberação determinística de recursos (fechamento de pools de banco de dados, encerramento de conexões de rede e finalização de tarefas em background) utilizando blocos 'try...finally' e gerenciadores de contexto ('with' / 'async with').
- Fornecer um serviço de domínio de referência inicial funcional e testado (ex: HealthCheckService retornando status de saúde, integridade e metadados de versão da aplicação).

================================================================================
PILAR 2: PONTOS DE ENTRADA, MODULARIDADE E VERSIONAMENTO DINÂMICO
================================================================================
- Estruturação modular dos executáveis e subcomandos via 'pyproject.toml':
  * Cada comando/executável declarado em '[ENTRYPOINTS_OU_SERVICOS]' deve ser registrado formalmente na seção '[project.scripts]' do 'pyproject.toml' apontando para funções de entrada limpas e modulares (ex: `meu-comando = "[PACOTE_PYTHON].main:main"`).
  * Suporte a comando ou flag de versão ('--version' ou subcomando 'version') para inspeção rápida em linha de comando.
  * Adoção de biblioteca de CLI/parsing conforme especificado em '[STACK_E_FRAMEWORKS]' (ex: Typer, Click, argparse stdlib).
- Versionamento Dinâmico (Release vs Dev):
  * Módulo 'src/[PACOTE_PYTHON]/version.py' com resolução dinâmica em tempo de execução utilizando 'importlib.metadata.version("[NOME_DO_PROJETO]")' com fallback resiliente para ambiente de desenvolvimento local:
    - Ao executar a partir de um pacote instalado/tag de release (ex: v1.0.0), exibir a versão semântica oficial: "[NOME_DO_PROJETO] version: v1.0.0 (commit: abc1234, built at: 2026-08-26T12:00:00Z)".
    - Ao executar em ambiente de desenvolvimento local ou a partir do código fonte sem tag, exibir identificação explícita de desenvolvimento: "[NOME_DO_PROJETO] version: dev (commit: abc1234, built at: 2026-08-26T12:00:00Z)".

================================================================================
PILAR 3: INFRAESTRUTURA DE TESTES AUTOMATIZADOS, TDD/BDD E BARREIRA >= 80%
================================================================================
- Ferramental: Utilizar o 'pytest' como test runner universal, integrado com 'pytest-cov' para medição de cobertura, 'pytest-mock' ou 'unittest.mock' para isolamento e 'pytest-asyncio' quando houver código assíncrono.
- Metodologia BDD: Estruturar cenários de teste através de nomenclatura e funções BDD declarativas utilizando o padrão:
  def test_given_[contexto]_when_[acao]_then_[resultado_esperado]():
      # Given
      ...
      # When
      ...
      # Then
      ...
- Desacoplamento e Mocks: Mocks determinísticos e isolados baseados estritamente nos contratos de interface definidos em 'src/[PACOTE_PYTHON]/core/ports/', evitando dependências externas ou I/O real em testes unitários.
- Barreira de Cobertura Inegociável:
  * Script automatizado 'scripts/coverage.sh' que executa a suíte de testes com medição de cobertura sobre 'src/[PACOTE_PYTHON]/'.
  * O script deve falhar com código de saída diferente de zero se a cobertura global for inferior a 80% ('--cov-fail-under=80').
  * Geração automática de relatórios de cobertura em terminal e formato HTML em 'dist/coverage/'.

================================================================================
PILAR 4: INTERFACE UNIVERSAL DE COMANDOS VIA MAKEFILE AUTODOCUMENTADO
================================================================================
- Criar um 'Makefile' universal, autodocumentado e determinístico como interface central de automação apoiado pelo 'uv':
  * make help (alvo padrão): Exibe dinamicamente o menu com todos os comandos disponíveis e descrições formatadas.
  * make setup: Instala/sincroniza a versão do Python, cria o ambiente virtual isolado e instala todas as dependências do 'pyproject.toml' com resolução determinística via 'uv sync'.
  * make dev: Inicia o ambiente de desenvolvimento local (com live-reload via Uvicorn/Flask ou monitor de testes pytest-watch quando aplicável).
  * make run: Executa a aplicação/serviço a partir do ponto de entrada principal ('uv run python -m [PACOTE_PYTHON]').
  * make test: Executa a suíte completa de testes com verificação de cobertura (barreira inegociável >= 80% via 'uv run pytest').
  * make test-unit: Executa testes unitários rápidos ('uv run pytest tests/unit').
  * make test-coverage: Gera o relatório HTML de cobertura em dist/coverage/ e valida o limiar de 80%.
  * make typecheck: Executa a checagem estática de tipos ultrarrápida com 'ty' ('uv run ty' ou 'ty src tests').
  * make lint: Executa linters estritos de código com 'ruff' ('uv run ruff check .').
  * make fmt: Formata o código automaticamente com 'ruff' ('uv run ruff format .' e 'uv run ruff check --fix .').
  * make check: Quality Gate local unificado (fmt + lint + typecheck com 'ty' + pip-audit para auditoria de vulnerabilidades + test com cobertura >= 80%).
  * make build: Constrói os pacotes de distribuição de produção (Wheel '.whl' e Source Distribution '.tar.gz') sob 'dist/' utilizando 'uv build'.
  * make clean: Limpa ambientes virtuais (.venv), caches temporários (.pytest_cache, .ruff_cache, .ty_cache, __pycache__), relatórios de cobertura e artefatos de compilação em dist/ e build/.

================================================================================
PILAR 5: STACK DE APLICAÇÃO, DEPENDÊNCIAS E DECISÕES ARQUITETURAIS CUSTOMIZÁVEIS
================================================================================
- Gestão de Dependências e Configuração Padronizada com 'pyproject.toml' e 'uv':
  * É MANDATÓRIO centralizar todas as dependências de produção e desenvolvimento no arquivo 'pyproject.toml' (PEP 621), gerenciadas via 'uv add' / 'uv lock'.
  * O 'uv' é a ferramenta padrão e recomendada para gerenciamento de versões do Python ('uv python install / pin'), criação de venv e bloqueio de dependências ('uv.lock').
- As decisões de frameworks, bibliotecas externas, drivers, persistência e ferramentas de desenvolvimento são deliberadas em conjunto com o solicitante do prompt conforme o contexto e escopo do projeto em '[STACK_E_FRAMEWORKS]' e '[TIPO_DE_PROJETO]'.
- Princípio da Parcimônia:
  * Priorizar a biblioteca padrão de Python (stdlib: 'asyncio', 'pathlib', 'dataclasses', 'typing', 'logging', 'sqlite3', 'argparse') sempre que atender com elegância, manutenibilidade e eficiência aos requisitos.
  * Quando bibliotecas de terceiros forem necessárias, selecionar pacotes consolidados, seguros, modernos, ativamente mantidos e com baixo acoplamento.
- Decisões Sob Demanda (sem imposições pré-fabricadas):
  * Frameworks pesados, ORMs complexos, brokers de mensageria (Celery/RabbitMQ), bancos relacionais ou NoSQL, ferramentas de template SSR (Jinja2) NÃO devem ser impostos arbitrariamente: só devem ser adotados se fizerem sentido para o projeto e forem expressamente acordados no prompt.

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
     - O terminal deve ser utilizado exclusivamente para ferramentas do ciclo de vida: 'make' ('make test', 'make lint', 'make check', 'make build'), 'uv', 'python', 'pytest', 'ruff', 'ty', 'git', 'openspec' e binários executáveis.
  3. Autonomia Operacional e Permissões (.agent/settings.json):
     - Configurar permissões pré-autorizadas ('allow') para todas as ferramentas cotidianas ('uv', 'python', 'pytest', 'make', 'git', 'openspec', 'ruff', 'ty'), eliminando prompts de confirmação repetitivos e assegurando execução autônoma contínua.
  4. Economia Ativa de Tokens, Curação de Contexto e Otimização de Janela:
     - Respostas do agente devem ser concisas e estruturadas em Markdown com links clicáveis '[arquivo](file:///caminho)'.
     - NUNCA duplicar blocos massivos de código já existentes no disco na resposta do chat.
     - Leitura cirúrgica de arquivos limitando linhas ('StartLine'/'EndLine') e edições pontuais com 'replace_file_content'.
     - Execução de comandos com saídas concisas e estruturadas (sem flags excessivamente verbosas).
  5. Disciplinas de Engenharia de Agentes em Sistemas Compostos (Harness, Loop & Graph Engineering):
     - Harness Engineering: Arcabouço operacional e scaffolding de IA composto por diretrizes determinísticas (.agent/rules/, AGENTS.md, GEMINI.md), matriz de permissões pré-autorizadas (.agent/settings.json), guardrails de segurança, contexto otimizado e limites de sandbox para governar a operação autônoma, previsível e segura do Antigravity CLI.
     - Loop Engineering: Ciclos cognitivos e iterativos de execução e auto-validação contínua do agente (ReAct / Reflection loops: inspeção cirúrgica de requisitos -> intervenção pontual no código -> execução de testes automatizados com pytest -> diagnóstico de erros -> correção e validação final com 'make check'), prevenindo loops infinitos e garantindo convergência rápida.
     - Graph Engineering (State Graphs & DAG de Raciocínio): Modelagem e orquestração do fluxo cognitivo/operacional do agente como um grafo direcionado acíclico (DAG) de tarefas, estados e decisões. Decomposição topológica de dependências entre tarefas de planejamento e execução, controle de transições determinísticas de estado entre fases (especificação -> implementação -> verificação -> reflexão), roteamento de fluxos e coordenação/paralelismo entre subagentes especializados sem dependências circulares ou bloqueios cognitivos.
     - Gestão de Contexto e Subagentes: Delegação de pesquisas e tarefas isoladas para subagentes dedicados para manter a janela de contexto do agente principal limpa, enxuta e focada.

================================================================================
PILAR 7: GOVERNANÇA DE ESPECIFICAÇÕES COM OPENSPEC (RECOMENDADO PARA BOAS PRÁTICAS)
================================================================================
- Recomendação Mandatória de Boas Práticas: É altamente recomendado adotar o framework 'OpenSpec' como a ferramenta padrão de governança de especificações, rastreabilidade de mudanças e alinhamento contínuo entre Product Owner (PO) e Quality Assurance (QA).
- Configurar 'openspec/config.yaml' promovendo entendimento mútuo entre PO e QA:
  * context: Declaração da stack do projeto (Python, Clean Architecture, tecnologias escolhidas, TDD/BDD >= 80%, Makefile, PT-BR, Autonomia e Ferramentas Nativas do AGY).
  * rules.proposal: Foco no 'porquê' (motivação de negócio) e 'o que muda' (escopo funcional e técnico), declaração explícita de impacto e idioma PT-BR.
  * rules.specs: Escrita em PT-BR sob ótica CONJUNTA e colaborativa de PO e QA:
    - **Diretrizes para o Product Owner (PO)**:
      * Redação em linguagem clara, ubíqua e acessível ao negócio, descrevendo o valor entregue sem detalhes internos de implementação (como nomes de funções ou classes) que obscureçam a regra.
      * Seção '## Purpose' obrigatória (mínimo de 50 caracteres) explicando claramente o objetivo da capacidade para o produto e seus usuários.
      * Requisitos funcionais ('### Requirement: <Nome>') com critérios de aceitação objetivos e verificáveis.
    - **Diretrizes para o Quality Assurance (QA) e Automação de Testes**:
      * Estruturação formal de cenários BDD/Gherkin com 4 hashtags '#### Scenario: <Nome>' e bullets padronizados '- **WHEN**' e '- **THEN**' (e opcionalmente '- **GIVEN**').
      * Cobertura determinística de fluxos principais (caminho feliz), fluxos alternativos, validação de limites/bordas e tratamento de erros e exceções.
      * Cenários redigidos de forma precisa para que ferramentas de automação de testes (como pytest, pytest-bdd, behave, testes de contrato ou testes E2E) consigam traduzir e validar os comportamentos de ponta a ponta sem ambiguidades.
  * rules.tasks: Seções numeradas '## N. Nome do Grupo', checkboxes '- [ ] N.M Descrição', tarefas explícitas de testes unitários/BDD com validação de cobertura >= 80% e 'make check'.
  * operations.apply.guidance: Executar 'make test' e 'make check' para validar alterações antes de concluir tarefas, priorizando ferramentas nativas.

================================================================================
PILAR 8: PIPELINES DE CI/CD E MATRIZ DE RELEASE (PLATAFORMAS CUSTOMIZÁVEIS)
================================================================================
- Workflow de CI ('.github/workflows/ci.yml'):
  * Executado em Pull Requests e pushes na branch principal ('main').
  * Matriz de execução testando versões suportadas do Python (ex: 3.11, 3.12, 3.13) e sistemas operacionais declarados em '[PLATAFORMAS_ALVO]'.
  * Etapas: Checkout, Setup Python / uv, Cache de dependências, Instalação determinística via 'uv sync', Execução de linter e checagem de tipos (ruff e ty), Auditoria de vulnerabilidades de segurança ('pip-audit'), Execução de testes com barreira de cobertura >= 80% e validação de integridade OpenSpec ('openspec validate --all').
- Workflow de Release e Publicação ('.github/workflows/release.yml'):
  * Disparado na criação de tags de release Git ('v*').
  * Construção automatizada de pacotes de distribuição padrão (Wheel '.whl' e Source Distribution '.tar.gz') via 'uv build'.
  * Geração de checksums SHA256 de integridade.
  * Publicação automática dos artefatos na Release do GitHub e/ou no PyPI (utilizando autenticação segura via Trusted Publishing / OIDC do GitHub Actions).
- Boas Práticas de Prompt e Regra de Fallback:
  * Caso as plataformas e versões alvo ('[PLATAFORMAS_ALVO]') NÃO sejam informadas pelo usuário no prompt:
    1. O agente deve assumir como padrão a versão de Python e ambiente corrente do host (ex: Python >=3.11).
    2. Se o projeto demandar distribuição de bibliotecas ou CLIs para terceiros e houver ambiguidade, o agente DEVE questionar proativamente o usuário antes de assumir matrizes arbitrárias.

================================================================================
PILAR 9: DOCUMENTAÇÃO VIVA E EXAUSTIVA NO README.MD
================================================================================
- Estruturar o 'README.md' raiz com as seguintes seções obrigatórias:
  1. Header & Badges: Título com badges funcionais (Python Version, CI Quality Gate, Test Coverage >= 80%, PyPI/Latest Release, License).
  2. Visão Geral da Aplicação: Propósito, valor entregue e diferenciais técnicos.
  3. Guia de Instalação e Uso:
     - Instalação via gerenciador de pacotes ('uv tool install' ou 'pip install') ou execução via código fonte.
     - Documentação completa e tabela com comandos, flags de linha de comando, variáveis de ambiente (.env), parâmetros de configuração e exemplos práticos de execução.
  4. Guia do Desenvolvedor:
     - Pré-requisitos e setup do ambiente virtual de desenvolvimento ('make setup' via 'uv').
     - Arquitetura de software e layout de diretórios ('src/[PACOTE_PYTHON]/core/', 'src/[PACOTE_PYTHON]/adapters/').
     - Interface universal de comandos e automação via Makefile ('make check', 'make dev', 'make test', 'make build').
     - Metodologia de testes TDD/BDD com barreira de 80%.
     - Diretrizes de IA (Antigravity), autonomia de ferramentas e padrão Conventional Commits.

================================================================================
PILAR 10: GOVERNANÇA GIT, HIGIENE DE REPOSITÓRIO E SQUASH MERGE
================================================================================
- Criação e Configuração Mandatória de '.gitignore' Idiomático:
  * Criar imediatamente na raiz do repositório um arquivo '.gitignore' abrangente para Python/uv, cobrindo:
    - Ambientes virtuais e pacotes: '.venv/', 'venv/', 'env/', 'ENV/', '__pypackages__/'.
    - Caches e compilação Python: '__pycache__/', '*.py[cod]', '*$py.class', '.pytest_cache/', '.ruff_cache/', '.ty_cache/'.
    - Distribuição e builds: 'dist/', 'build/', '*.egg-info/', '*.egg', '*.whl', '*.tar.gz'.
    - Relatórios e saídas de teste/cobertura: 'coverage.html', '.coverage', '.coverage.*', 'coverage/', 'htmlcov/', '.tox/', '.nox/'.
    - Segredos e variáveis de ambiente: '.env', '.env.*', '*.pem', '*.key', '*.crt' (exceto templates públicos como '.env.example').
    - Configurações de IDEs e do Sistema Operacional: '.vscode/', '.idea/', '*.swp', '.DS_Store', 'Thumbs.db'.
- Importância e Boas Práticas de Governança Git em Engenharia com Agentes de IA:
  * Higiene e Economia de Contexto: Manter o repositório rigorosamente livre de diretórios virtuais ('.venv'), caches pesados de lint/tipagem e artefatos de build evita que ferramentas de busca e I/O do agente ('grep_search', 'find_by_name', 'list_dir') indexem milhares de arquivos desnecessários, protegendo a janela de contexto de alucinações e desperdício de tokens.
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
    └── python_conventions.md
```

### Checklist dos Artefatos Gerados:
- [x] **`proposal.md`**: Define o *Why*, *What Changes*, *Capabilities* (`project-foundation`) e *Impact* (código, IA, governança, autonomia).
- [x] **`design.md`**: Detalha *Context*, *Goals/Non-Goals*, *Decisions* (arquitetura Python src-layout, pyproject.toml, uv, ty, testes BDD, Makefile, stack escolhida, release matrix, README, CI/CD, prioridade de ferramentas nativas) e *Risks/Trade-offs*.
- [x] **`specs/project-foundation/spec.md`**: Define os requisitos com perspectivas explícitas de **PO** e **QA**, cada um acompanhado de cenários BDD `#### Scenario:` com `- **WHEN**` e `- **THEN**`.
- [x] **`tasks.md`**: Lista de tarefas organizadas por seções numeradas (`## N. Nome do Grupo`), com subtarefas granulares `- [ ] N.M`, incluindo verificações de setup uv, tipagem ty, testes unitários pytest, linter ruff, cobertura >= 80% e `openspec validate --all`.
- [x] **`openspec/config.yaml`**: Configuração central com `context`, `rules` (proposal, specs, design, tasks) e `operations` (apply).
- [x] **`AGENTS.md` e `GEMINI.md`**: Diretrizes operacionais de alta prioridade para autonomia do AGY, ferramentas nativas, economia de tokens e fluxo Git squash.
- [x] **`.agent/settings.json`**: Permissões de desenvolvimento pré-autorizadas para máxima autonomia sem interrupções.
- [x] **`.agent/rules/`**: Diretrizes operacionais para Harness de IA (scaffolding e guardrails), Loop e Graph Engineering (State Graphs e DAG de tarefas), Autonomia de Ferramentas, Git Branching Workflow, Arquivamento com Squash e Convenções Python.
