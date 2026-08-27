## MODIFIED Requirements

### Requirement: Estratégias do Antigravity CLI (Harness, Loop e Graph Engineering) e Otimização de Tokens
O repositório SHALL conter diretrizes, regras e configurações em `.agent/rules/`, `.agent/settings.json`, `AGENTS.md`, `GEMINI.md` e `openspec/config.yaml` que capacitem o Antigravity CLI a operar com autonomia máxima sem interrupções desnecessárias por confirmação de prompt, priorizando estritamente ferramentas nativas de arquivos em vez de comandos de terminal, otimizando o consumo de tokens na janela de contexto e aplicando estratégias de engenharia de IA e agente (**Harness Engineering, Loop Engineering, Graph Engineering e Tooling Autonomy**), tratando Harness Engineering estritamente como o arcabouço de regras, scaffolding, guardrails e configurações de segurança e contexto para a IA.
*(Visão PO: Garante autonomia operacional contínua, agilidade máxima sem interrupções manuais triviais e governança rigorosa de IA com economia drástica de tokens em todas as interações. Visão QA: Garante validação determinística via ferramentas nativas de inspeção/edição, isolamento seguro do agente e execução de ciclo de vida com rastreabilidade total).*

#### Scenario: Harness Engineering pelo Antigravity CLI
- **WHEN** o Antigravity CLI for inicializado e executar operações no repositório
- **THEN** o repositório deve fornecer um harness de IA composto por scaffolding de regras determinísticas (`.agent/rules/`, `AGENTS.md`, `GEMINI.md`), guardrails de segurança e contexto, configurações de permissões pré-autorizadas (`.agent/settings.json`) e restrições de ferramentas nativas para guiar o agente de forma segura, determinística e otimizada

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
