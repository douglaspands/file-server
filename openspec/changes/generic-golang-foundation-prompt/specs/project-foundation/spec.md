## MODIFIED Requirements

### Requirement: Estratégias do Antigravity CLI (Harness, Loop e Graph Engineering) e Otimização de Tokens
O repositório SHALL conter diretrizes, regras e configurações em `.agent/rules/`, `.agent/settings.json`, `AGENTS.md`, `GEMINI.md`, `docs/generic-golang-foundation-spec-prompt.md` e `openspec/config.yaml` que capacitem o Antigravity CLI a operar com autonomia máxima sem interrupções desnecessárias por confirmação de prompt, priorizando estritamente ferramentas nativas de arquivos em vez de comandos de terminal, otimizando o consumo de tokens na janela de contexto e aplicando integralmente as disciplinas de engenharia de IA e agentes (**Harness Engineering, Loop Engineering, Graph Engineering, Tooling Grounding e Token Economics**), tratando Graph Engineering puramente como modelagem de fluxo cognitivo/operacional de agentes e orquestração de grafos de execução (State Graphs / DAG de raciocínio e tarefas), além de disponibilizar um prompt mestre de fundação arquitetural genérico, modular, agnóstico a frameworks, com governança OpenSpec estruturada para PO e QA e parametrização de arquiteturas de compilação/release em `docs/generic-golang-foundation-spec-prompt.md`.
*(Visão PO: Garante autonomia operacional contínua, agilidade máxima sem interrupções manuais triviais e governança rigorosa de IA com economia drástica de tokens em todas as interações, além de prover um template mestre reutilizável para qualquer projeto Go com regras de negócio e critérios de aceitação perfeitamente legíveis. Visão QA: Garante validação determinística via ferramentas nativas de inspeção/edição, isolamento seguro do agente, orquestração estruturada de tarefas e cenários BDD/Gherkin preparados para frameworks de automação de testes).*

#### Scenario: Harness Engineering pelo Antigravity CLI
- **WHEN** o Antigravity CLI for inicializado e executar operações no repositório
- **THEN** o repositório deve fornecer um harness de IA composto por scaffolding de regras determinísticas (`.agent/rules/`, `AGENTS.md`, `GEMINI.md`), guardrails de segurança e contexto, configurações de permissões pré-autorizadas (`.agent/settings.json`) e restrições de ferramentas nativas para guiar o agente de forma segura, determinística e otimizada

#### Scenario: Loop Engineering pelo Antigravity CLI
- **WHEN** o Antigravity CLI implementar ou refatorar componentes da aplicação
- **THEN** o agente deve executar loops curtos de feedback contínuo (inspeção cirúrgica de requisitos -> intervenção pontual no código -> execução de testes automatizados -> diagnóstico de erros -> correção e validação final com `make check`) antes de submeter a solução

#### Scenario: Graph Engineering pelo Antigravity CLI
- **WHEN** o Antigravity CLI planejar, orquestrar ou executar tarefas e fluxos de raciocínio
- **THEN** o agente deve estruturar sua execução como um grafo direcionado acíclico (DAG / State Graph) de tarefas e decisões cognitivas, decompondo tarefas em nós ordenados topologicamente, controlando transições determinísticas de estado entre fases (planejamento, implementação, validação) e coordenando a execução de subagentes sem dependências circulares ou bloqueios

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

#### Scenario: Prompt mestre de fundação arquitetural genérico para Golang
- **WHEN** um desenvolvedor, PO ou QA consultar e utilizar `docs/generic-golang-foundation-spec-prompt.md`
- **THEN** o documento deve fornecer um prompt parametrizável e agnóstico de frameworks web/CLI específicos, sem imposições pré-moldadas de empacotamento ou live-reload (Pilar 5), com diretrizes estruturadas de governança OpenSpec para fácil entendimento de regras de negócio por PO e automação de testes por QA (Pilar 7), parametrização de plataformas alvo de compilação/release com fallback para arquitetura corrente ou questionamento proativo (Pilar 8), mantendo obrigatórios os pilares de Clean Architecture, TDD/BDD >= 80%, Makefile, Harness/Loop/Graph Engineering, CI/CD e governança Git com Squash
