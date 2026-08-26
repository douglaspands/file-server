# project-foundation Specification

## Purpose

Define a fundação arquitetural, de engenharia de software, esteira de testes automatizados (TDD/BDD), automações de harness/loop/graph, governança de IA e repositório para aplicações Go com renderização HTML.

## Requirements

### Requirement: Estrutura Modular e Padrão Idiomático de Projeto Go
O sistema SHALL seguir o layout padrão canônico da comunidade Go (`cmd/`, `internal/`, `pkg/`, `web/`, `scripts/`, `.github/`), garantindo encapsulamento estrito onde o código de domínio e regras de negócio residem exclusivamente em `internal/`.
*(Visão PO: Garante manutenibilidade e proteção das regras de negócio contra acoplamento externo. Visão QA: Garante isolamento para testes unitários e de integração sem vazamento de estado global).*

#### Scenario: Organização de pacotes e isolamento de domínio
- **WHEN** o desenvolvedor ou pipeline inspeciona a estrutura do repositório
- **THEN** os módulos privados e lógica de negócio devem residir sob o diretório `internal/`, a inicialização de binários sob `cmd/`, e recursos web estáticos/templates sob `web/`

#### Scenario: Compilação e checagem de integridade do módulo
- **WHEN** o comando de verificação e build do módulo Go for executado
- **THEN** o arquivo `go.mod` deve estar sincronizado sem dependências não utilizadas (`go mod tidy`) e o binário deve compilar com sucesso

### Requirement: Esteira de Testes Automatizados TDD, BDD e Cobertura Mínima
O sistema SHALL fornecer uma infraestrutura de testes automatizados suportando metodologias TDD e BDD com asserções claras, mocks desacoplados e medição rigorosa de cobertura de código com barreira mínima de 80%.
*(Visão PO: Garante entrega contínua com alta confiabilidade e redução drástica de regressões em produção. Visão QA: Garante validação de cenários felizes, bordas, falhas e contratos com rastreabilidade Given-When-Then).*

#### Scenario: Execução da suíte completa de testes com relatório de cobertura
- **WHEN** a suíte de testes unitários e de integração for disparada via ferramenta de automação
- **THEN** todos os testes devem ser executados, gerando relatório de cobertura (`coverage.out` / HTML) e falhando caso a cobertura global seja inferior a 80%

#### Scenario: Validação de cenários de comportamento BDD
- **WHEN** testes de especificação de comportamento (BDD) forem executados contra casos de uso do domínio
- **THEN** cada cenário deve validar as pré-condições (Given), ações do usuário/sistema (When) e os estados resultantes esperados (Then)

### Requirement: Interface de Comandos Centralizada via Makefile
O sistema SHALL disponibilizar um arquivo `Makefile` autodocumentado (e opcionalmente `Taskfile.yml`) contendo todos os comandos principais de uso, desenvolvimento e ciclo de vida da aplicação de fácil acesso (`make help`, `make setup`, `make dev`, `make run`, `make test`, `make lint`, `make check`, `make build`), garantindo idempotência e determinismo para humanos e agentes.
*(Visão PO: Padroniza o onboarding de desenvolvedores e agentes de IA com comandos intuitivos e de fácil descoberta, eliminando atrito operacional. Visão QA: Garante que os mesmos comandos executados localmente sejam executados no CI de forma determinística e reproduzível).*

#### Scenario: Descoberta e ajuda interativa de comandos via Makefile
- **WHEN** o desenvolvedor ou agente IA executar `make` ou `make help`
- **THEN** o sistema deve listar todos os comandos disponíveis com descrições claras de sua finalidade

#### Scenario: Execução do pipeline de verificação local
- **WHEN** o desenvolvedor ou agente IA executar o comando de verificação unificado (`make check` ou `task check`)
- **THEN** o sistema deve executar sequencialmente formatação, linter estrito (`golangci-lint`), checagem de vulnerabilidades (`govulncheck`) e suíte de testes com validação de cobertura

#### Scenario: Setup de ambiente reprodutível
- **WHEN** um novo ambiente for inicializado através do comando `make setup`
- **THEN** todas as ferramentas de desenvolvimento, linters e dependências necessárias devem ser instaladas e verificadas automaticamente

### Requirement: Live-Reloading e Feedback Visual no Desenvolvimento Web
O sistema SHALL fornecer capacidade de recarregamento rápido (live-reload) no desenvolvimento web através de monitoramento de alterações em código Go, templates HTML e assets estáticos.
*(Visão PO: Acelera a experimentação visual de interfaces e feedback imediato no desenvolvimento. Visão QA: Facilita testes exploratórios e validação visual de layout em tempo real).*

#### Scenario: Disparo de live-reload ao alterar arquivos
- **WHEN** qualquer arquivo `.go`, `.html`, `.css` ou template for alterado e salvo durante a execução em modo de desenvolvimento (`make dev`)
- **THEN** o servidor de desenvolvimento deve recompilar e reiniciar automaticamente em menos de 2 segundos, atualizando a interface web

### Requirement: Estratégias do Antigravity CLI (Harness, Loop e Graph Engineering) e Otimização de Tokens
O repositório SHALL conter diretrizes, regras e configurações em `.agent/rules/`, `.agent/settings.json`, `AGENTS.md`, `GEMINI.md` e `openspec/config.yaml` que capacitem o Antigravity CLI a operar com autonomia máxima sem interrupções desnecessárias por confirmação de prompt, priorizando estritamente ferramentas nativas de arquivos em vez de comandos de terminal, otimizando o consumo de tokens na janela de contexto e aplicando estratégias de engenharia de agente (**Harness Engineering, Loop Engineering, Graph Engineering e Tooling Autonomy**).
*(Visão PO: Garante autonomia operacional contínua, agilidade máxima sem interrupções manuais triviais e economia drástica de tokens em todas as interações. Visão QA: Garante validação determinística via ferramentas nativas de inspeção/edição e execução segura de comandos de ciclo de vida com rastreabilidade total).*

#### Scenario: Harness Engineering pelo Antigravity CLI
- **WHEN** o Antigravity CLI executar comandos de desenvolvimento, teste ou diagnóstico
- **THEN** o agente deve utilizar o Makefile como harness padronizado (`make test`, `make lint`, `make check`), consumindo saídas concisas e estruturadas para minimizar gasto de tokens

#### Scenario: Loop Engineering pelo Antigravity CLI
- **WHEN** o Antigravity CLI implementar ou refatorar componentes da aplicação
- **THEN** o agente deve executar loops curtos de feedback contínuo (código -> teste automatizado -> correção de falhas -> validação final) antes de submeter a solução

#### Scenario: Graph Engineering pelo Antigravity CLI
- **WHEN** o Antigravity CLI planejar ou implementar novos módulos e tarefas
- **THEN** o agente deve mapear o grafo de dependências (DAG) das tarefas e a árvore de chamadas do projeto, implementando componentes na ordem topológica correta (interfaces/ports antes de adaptadores) sem dependências circulares

#### Scenario: Governança contínua de regras para propostas futuras
- **WHEN** qualquer nova proposta ou plano for iniciado no repositório
- **THEN** o arquivo `openspec/config.yaml` deve injetar contexto e regras operacionais íntegras exigindo padrão Go, TDD/BDD, cobertura >= 80%, Makefile, PT-BR e perspectivas PO/QA sem avisos de parsing

#### Scenario: Prioridade mandatória de ferramentas nativas do AGY
- **WHEN** o Antigravity CLI precisar criar, inspecionar, editar ou pesquisar arquivos no repositório
- **THEN** o agente deve obrigatoriamente invocar as ferramentas nativas (`write_to_file`, `replace_file_content`, `view_file`, `grep_search`, `find_by_name`, `list_dir`) e abster-se estritamente de executar comandos bash equivalentes (`cat`, `echo >`, `sed`, `awk`, `touch`, `find`, `grep`, `ls`) via terminal

#### Scenario: Configuração de permissões para máxima autonomia
- **WHEN** o Antigravity CLI executar ferramentas essenciais e comandos de ciclo de vida do projeto (`go`, `make`, `git`, `openspec`, `golangci-lint`, `govulncheck`, scripts de teste)
- **THEN** o arquivo `.agent/settings.json` deve conter permissões pré-autorizadas (`allow`) para essas operações legítimas, eliminando solicitações de autorização repetitivas e garantindo execução fluida

#### Scenario: Economia de tokens e respostas concisas
- **WHEN** o Antigravity CLI interagir com o usuário ou processar arquivos
- **THEN** o agente deve responder de forma concisa utilizando links markdown `[arquivo](file:///...)`, aplicar leitura cirúrgica com intervalos de linhas (`StartLine`/`EndLine`) e edições pontuais com `replace_file_content`, evitando replicar blocos massivos de código na resposta

### Requirement: Arquitetura de Frontend, Assets Embutidos e Experiência do Usuário
O sistema SHALL integrar uma camada moderna de frontend para renderização do lado do servidor (SSR) combinada com interatividade hipermidiática declarativa (recomendações: HTMX, Alpine.js e Tailwind CSS), fornecendo interfaces responsivas e empacotamento integral de assets (`go:embed`) para geração de binários únicos 100% autocontidos e portáteis.
*(Visão PO: Proporciona uma experiência de usuário (UX) fluida com interfaces ricas e facilidade total de implantação/transporte em um único arquivo executável sem arquivos externos. Visão QA: Permite testes de renderização de componentes HTML isolados, validação de integridade dos assets embutidos em tempo de build e testes de requisições parciais HTMX).*

#### Scenario: Renderização de componentes HTML dinâmicos via servidor
- **WHEN** o usuário interagir com um elemento dinâmico na interface web
- **THEN** o servidor Go deve renderizar fragmentos HTML semânticos e atualizá-los no DOM sem recarregar a página inteira

#### Scenario: Processamento e compilação de assets de interface
- **WHEN** o pipeline de build de frontend for executado
- **THEN** as classes de estilo (Tailwind) e scripts estáticos devem ser minificados e gerados no diretório de distribuição web

#### Scenario: Empacotamento de templates e assets em binário único autocontido
- **WHEN** o projeto for compilado para distribuição (`make build`)
- **THEN** todos os templates HTML, arquivos CSS, scripts JS e assets estáticos devem ser embutidos diretamente no executável compilado via `embed.FS`, permitindo a execução da aplicação completa em qualquer diretório sem requerer arquivos web externos

### Requirement: Estrutura de CLI Extensível e Versionamento Dinâmico
O sistema SHALL ser inicializado por meio de uma interface de linha de comando (CLI) extensível e modular, suportando comandos, subcomandos, argumentos posicionais e flags (`options`), além de fornecer um comando e argumento `version` que reporte o número da versão semântica oficial quando compilado para release ou sinalize explicitamente o estado de desenvolvimento quando em ambiente local.
*(Visão PO: Permite que o usuário e outros sistemas descubram facilmente a versão e usem a aplicação com flexibilidade via terminal. Visão QA: Garante validação determinística de compatibilidade de versões e testes de argumentos de CLI).*

#### Scenario: Exibição de versão em compilação oficial de release
- **WHEN** o usuário executar a aplicação com o comando ou flag de versão (ex: `app version` ou `app --version`) sobre um binário compilado a partir de uma release tag
- **THEN** o sistema deve imprimir a versão semântica oficial da release (ex: `v1.0.0`), hash do commit e data de compilação

#### Scenario: Exibição de versão em ambiente de desenvolvimento local
- **WHEN** o usuário executar a aplicação em ambiente de desenvolvimento local sem tag de release
- **THEN** o sistema deve exibir uma indicação explícita de desenvolvimento (ex: `dev`, `v0.0.0-dev` ou hash local) informando que o binário não é uma release oficial

#### Scenario: Extensibilidade para novos comandos e opções de CLI
- **WHEN** novas especificações adicionarem argumentos e opções de linha de comando
- **THEN** a estrutura modular da CLI (baseada no padrão Cobra) deve permitir registrar novos comandos e flags sem acoplamento com a lógica de inicialização base

### Requirement: Pipeline de Release Multiplataforma no GitHub Actions (Linux & Windows)
O sistema SHALL disponibilizar um workflow automatizado no GitHub Actions disparado pela criação de tags de release Git (`v*`), responsável por compilar o projeto de forma cruzada para as plataformas **Linux** (`amd64`, `arm64`) e **Windows** (`amd64`, `arm64`), gerando arquivos compactados (`.tar.gz` e `.zip`) e publicando-os automaticamente nos assets da Release do GitHub.
*(Visão PO: Garante distribuição imediata e sem atrito para os principais sistemas operacionais do mercado. Visão QA: Garante que os binários para Windows e Linux sejam compilados e testados de forma reproduzível na esteira oficial de CI).*

#### Scenario: Publicação automática de binários em tag de release
- **WHEN** uma nova tag de versão no formato `v*` for criada e enviada ao repositório GitHub
- **THEN** o workflow de release deve compilar os binários com assets embutidos e injeção de versão via `ldflags` para Linux e Windows (x86_64 e ARM64), gerando os arquivos de distribuição e publicando-os na Release do GitHub

### Requirement: Documentação Abrangente e Padrão de Engenharia no README.md
O sistema SHALL fornecer um arquivo `README.md` completo, estruturado e profissional na raiz do repositório, contendo badges relevantes de status/tecnologias, descrição clara e objetiva da aplicação, passo a passo detalhado de instalação e uso com todos os comandos e flags da CLI explicados, seguido por uma seção dedicada a desenvolvedores detalhando arquitetura, testes TDD/BDD, Makefile harness, live-reload e fluxo de contribuição.
*(Visão PO: Garante clareza imediata sobre o valor da aplicação, facilidade de adoção por usuários finais e documentação acessível de todos os comandos. Visão QA: Garante que todas as instruções de execução, flags de terminal e comandos de desenvolvimento estejam descritos com exatidão e sem discrepâncias).*

#### Scenario: Apresentação da aplicação e badges de status
- **WHEN** qualquer usuário ou desenvolvedor acessar o repositório
- **THEN** o `README.md` deve iniciar com título, badges de status (Go, CI, Coverage, Release, License) e uma descrição clara e objetiva do que é e para que serve o File Server

#### Scenario: Guia de uso com detalhamento de todos os comandos da CLI
- **WHEN** o usuário consultar a seção de uso da aplicação
- **THEN** o documento deve fornecer um passo a passo prático de inicialização e uma tabela/lista com **todos os comandos e flags disponíveis** (`file-server`, `version`, `serve`, `--port`, `--host`, `--json`, `--help`, etc.) com explicações de uso e exemplos práticos

#### Scenario: Guia técnico e boas práticas para desenvolvedores
- **WHEN** um desenvolvedor consultar a documentação para contribuir ou estender o projeto
- **THEN** o documento deve detalhar a arquitetura em camadas (`internal/`), esteira de testes TDD/BDD com barreira de cobertura (>= 80%), comandos do `Makefile`, fluxo de live-reload, governança do Antigravity CLI e padrão Conventional Commits

### Requirement: Governança de Git, GitHub Actions, Integração Contínua e Arquivamento
O sistema SHALL aplicar convenções de Conventional Commits, ganchos de pre-commit, workflows do GitHub Actions com toolchain Go e linters perfeitamente alinhados através de um arquivo de configuração `.golangci.yml` determinístico na raiz do repositório para validação contínua estrita sem falsos positivos, e procedimento mandatório de commit Git na conclusão/arquivamento de qualquer especificação OpenSpec.
*(Visão PO: Garante histórico de mudanças limpo, rastreabilidade de valor de negócio, esteira de CI sem falsos negativos e persistência garantida do código após a conclusão de uma spec. Visão QA: Funciona como Quality Gate automatizado determinístico, garantindo que o compilador, o linter estrito e a suíte de testes executem de forma idêntica localmente e no CI).*

#### Scenario: Validação automática em Pull Request via GitHub Actions
- **WHEN** um Pull Request ou push for enviado ao repositório GitHub
- **THEN** o workflow de CI deve sincronizar a versão do Go a partir do `go.mod`, executar o `golangci-lint` utilizando a configuração oficial `.golangci.yml` do repositório, auditar vulnerabilidades via `govulncheck`, rodar os testes automatizados com validação de cobertura >= 80% e verificar a integridade do OpenSpec

#### Scenario: Compatibilidade e integridade do linter no Quality Gate
- **WHEN** a etapa de linting (`Run GolangCI-Lint`) for executada no pipeline de CI ou localmente via `make lint` e `make check`
- **THEN** a execução do `golangci-lint` deve respeitar o arquivo `.golangci.yml` da raiz do projeto, executando linters estáticos sem falsos positivos com sucesso e código de saída zero

#### Scenario: Validação local antes do commit (Pre-commit hook)
- **WHEN** um desenvolvedor tentar efetuar um commit local
- **THEN** os hooks configurados devem verificar a formatação do código e a mensagem de commit de acordo com a especificação Conventional Commits

#### Scenario: Commit Git obrigatório no arquivamento da especificação
- **WHEN** o comando de arquivamento de especificação (`openspec archive` ou `/openspec-archive-change`) for executado após a conclusão da implementação
- **THEN** as orientações operacionais (`operations.archive.guidance`) devem instruir o agente a criar o commit com Conventional Commits contendo as alterações e a sincronização da spec principal
