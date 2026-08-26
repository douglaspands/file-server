## Why

O Antigravity (AGY) tem interrompido frequentemente o fluxo de desenvolvimento ao solicitar autorizações do usuário para operações triviais (como escrita de arquivos via `cat`, `echo >` ou manipulações via Bash) em vez de utilizar suas ferramentas nativas especializadas (`write_to_file`, `replace_file_content`, `view_file`, `grep_search`, `find_by_name`, `list_dir`). Além do atrito operacional com prompts de confirmação repetitivos, o uso inadequado de comandos de shell verbosos e respostas volumosas consome tokens desnecessários na janela de contexto. É necessário ajustar a governança, as regras de repositório e as permissões de execução para conferir autonomia máxima ao AGY e otimizar o consumo de tokens.

## What Changes

- **Prioridade Mandatória de Ferramentas Nativas**: Estabelecer regra rígida para que o AGY utilize ferramentas nativas (`write_to_file`, `replace_file_content`, `view_file`, `grep_search`, `find_by_name`, `list_dir`) para criação, edição, inspeção e busca de arquivos, proibindo o uso de comandos bash (`cat`, `echo >`, `sed`, `awk`, `touch`, `find`, `grep`, `ls`) para essas tarefas.
- **Otimização de Permissões para Autonomia**: Ajustar `.agent/settings.json` para autorizar de forma determinística e segura todas as ferramentas e comandos de desenvolvimento cotidianos (`go`, `make`, `git`, `openspec`, `golangci-lint`, `govulncheck`), eliminando solicitações de autorização desnecessárias.
- **Estratégias de Economia de Tokens**:
  - Instruir o agente a adotar respostas concisas, estruturadas em Markdown com links clicáveis (`file://`), sem duplicar blocos maciços de código já salvos em disco.
  - Leitura seletiva de arquivos utilizando `StartLine` e `EndLine` no `view_file` para evitar carregar centenas de linhas desnecessárias no contexto.
  - Edição cirúrgica com `replace_file_content` em vez de reescrever arquivos completos.
  - Execução de comandos com flags de saída concisa (evitando `-v` ou outputs excessivamente prolixos em execuções de rotina).
- **Atualização das Regras do Repositório**: Atualizar `AGENTS.md`, `GEMINI.md` e os arquivos sob `.agent/rules/` (`agent_harness_engineering.md` e novo `agent_tooling_autonomy.md`) consolidando essas práticas como diretrizes de alta prioridade.
- **Correção da Configuração do OpenSpec**: Ajustar a estrutura de orientações operacionais em `openspec/config.yaml` eliminando avisos de formato.

## Capabilities

### Modified Capabilities
- `project-foundation`: Atualizar os requisitos de governança do Antigravity CLI para incluir a prioridade mandatória de ferramentas nativas, matriz de permissões para autonomia operacional e regras ativas de economia de tokens e otimização de saídas.

## Impact

- **Arquivos Afetados**: `.agent/settings.json`, `.agent/rules/agent_harness_engineering.md`, `.agent/rules/agent_tooling_autonomy.md`, `AGENTS.md`, `GEMINI.md`, `openspec/config.yaml`.
- **APIs / Código Go**: Nenhum impacto negativo ou quebra de contratos na base de código Go existente.
- **Fluxo do Desenvolvedor / IA**: Eliminação de prompts de confirmação desnecessários, maior velocidade nas iterações e redução drástica do consumo de tokens na janela de contexto.
