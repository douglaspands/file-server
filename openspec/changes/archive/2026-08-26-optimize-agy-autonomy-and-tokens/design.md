## Context

Ver `proposal.md` para motivação e `specs/project-foundation/spec.md` para os requisitos funcionais. O repositório utiliza o ecossistema Antigravity (AGY) com regras em `AGENTS.md`, `GEMINI.md`, `.agent/rules/`, `.agent/settings.json` e `openspec/config.yaml`. Entretanto, a ausência de regras expressas de prioridade de ferramentas nativas gerava invocações de comandos bash (como `cat << 'EOF'`, `echo >`, `sed`) que acionavam prompts de autorização no terminal e inflavam o consumo de tokens.

## Goals / Non-Goals

**Goals:**
- Estabelecer a prioridade inegociável de ferramentas nativas do AGY (`write_to_file`, `replace_file_content`, `view_file`, `grep_search`, `find_by_name`, `list_dir`) para todas as manipulações de arquivos e pesquisas no código.
- Configurar `.agent/settings.json` com permissões automáticas abrangentes para comandos seguros de desenvolvimento (`go`, `make`, `git`, `openspec`, `golangci-lint`, `govulncheck`, `scripts`), eliminando interrupções.
- Criar regras dedicadas de governança (`.agent/rules/agent_tooling_autonomy.md`) e enriquecer `agent_harness_engineering.md`, `AGENTS.md` e `GEMINI.md` com diretrizes de economia ativa de tokens e respostas concisas.
- Corrigir a estrutura de orientação de `archive` em `openspec/config.yaml` para remover warnings no parser do OpenSpec.

**Non-Goals:**
- Não alterar a arquitetura, regras de negócio ou contratos do servidor de arquivos Go (`internal/`, `cmd/`, `web/`).
- Não autorizar comandos destrutivos perigosos (ex: `rm -rf /`, `mkfs`, `sudo`), mantendo a integridade do sistema.

## Decisions

### Decisão 1: Hierarquia de Ferramentas Nativas vs Comandos de Terminal
- **Abordagem**: Definir explicitamente que ferramentas nativas (`write_to_file`, `replace_file_content`, `view_file`, `grep_search`, `find_by_name`, `list_dir`) têm precedência absoluta. Comandos de terminal (`run_command`) devem ser restritos exclusivamente a execução de binários do ciclo de vida (`make`, `go test`, `git`, `openspec`).
- **Alternativa Considerada**: Permitir comandos bash com operadores de redirecionamento (`>`, `>>`, `|`). Rejeitada porque subshells com pipes acionam checagens de segurança no AGY, produzem saídas prolixas e consomem muito mais tokens.

### Decisão 2: Estratégia de Permissões em `.agent/settings.json`
- **Abordagem**: Expandir a lista `"allow"` com padrões globais e subcomandos para todas as ferramentas legítimas de engenharia do repositório, mantendo `"ask"` para operações remotas/irreversíveis (`git push`, `rm`, `sudo`) e `"deny"` para comandos destrutivos de sistema.
- **Alternativa Considerada**: Conceder permissão irrestrita total (`*`). Rejeitada para preservar safety gates contra deleção acidental de arquivos de sistema ou push não autorizado.

### Decisão 3: Práticas Ativas de Economia de Tokens
- **Abordagem**:
  1. **Leitura Cirúrgica**: Uso de `StartLine` e `EndLine` no `view_file` ao inspecionar blocos específicos de código.
  2. **Edição Pontual**: Uso preferencial de `replace_file_content` para diffs contíguos em vez de reescrever arquivos inteiros.
  3. **Outputs Enxutos**: Respostas do agente estruturadas em Markdown direto, referenciando arquivos via links clicáveis (`[arquivo](file:///caminho)`) sem duplicar listagens ou blocos de código já salvos.
  4. **Comandos Concisos**: Execução de comandos com flags silenciosas/padrão, evitando logs excessivamente verbosos (ex: evitar `-v` salvo em diagnósticos de falha).

### Decisão 4: Centralização e Consistência de Regras
- **Abordagem**: Sincronizar as diretrizes em `AGENTS.md`, `GEMINI.md`, `.agent/rules/agent_harness_engineering.md` e novo `.agent/rules/agent_tooling_autonomy.md`, assegurando que agentes com diferentes carregamentos de contexto sigam exatamente o mesmo protocolo.

## Risks / Trade-offs

- **[Risco] O agente tentar usar comandos bash por hábito em prompts complexos** → **Mitigação**: Instruções explícitas com cláusulas "NUNCA" e "SEMPRE" nas regras de alta prioridade (`AGENTS.md`, `GEMINI.md` e `.agent/rules/`).
- **[Risco] Avisos de sintaxe no `openspec/config.yaml`** → **Mitigação**: Ajustar o formato do campo `guidance` para array de strings ou estrutura suportada nativamente pelo OpenSpec CLI.
