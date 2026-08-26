# Diretrizes de Autonomia de Ferramentas e Prioridade Nativa para Antigravity (AGY)

Este documento estabelece as regras obrigatórias de seleção de ferramentas para o Antigravity (AGY) no repositório, garantindo autonomia máxima, eliminação de prompts de autorização desnecessários e economia ativa de tokens.

---

## ⚡ 1. Prioridade Mandatória de Ferramentas Nativas

O AGY possui ferramentas nativas especializadas e determinísticas para manipulação do workspace. É **ESTRITAMENTE PROIBIDO** utilizar comandos bash equivalentes via `run_command` quando houver uma ferramenta nativa disponível.

| Operação Desejada | Ferramenta Nativa Obrigatória | Comandos Shell PROIBIDOS via `run_command` |
| :--- | :--- | :--- |
| **Criar novo arquivo** | `write_to_file` | `cat << 'EOF' > ...`, `echo "..." > ...`, `touch` |
| **Editar arquivo existente** | `replace_file_content` *(ou `write_to_file` se recriação total)* | `sed`, `awk`, `cat > ...`, scripts inline de stream |
| **Ler / inspecionar arquivo** | `view_file` | `cat`, `head`, `tail`, `less`, `more` |
| **Pesquisar texto no código** | `grep_search` | `grep`, `egrep`, `rg`, `ag` |
| **Localizar arquivos por padrão** | `find_by_name` | `find`, `locate`, `ls -R` |
| **Listar conteúdo de diretório** | `list_dir` | `ls`, `dir`, `tree` |

---

## 🛠️ 2. Escopo Exclusivo de Execução de Terminal (`run_command`)

A ferramenta `run_command` deve ser utilizada **exclusivamente** para:
1. **Comandos de Ciclo de Vida e Build**: `make` (`make test`, `make lint`, `make check`, `make build`, etc.) e `go` (`go test`, `go mod tidy`, `go build`).
2. **Controle de Versão Git**: `git status`, `git branch`, `git checkout`, `git add`, `git commit`, `git merge --squash`, etc.
3. **Gerenciamento do OpenSpec**: `openspec new change`, `openspec status`, `openspec instructions`, `openspec validate`, `openspec archive`.
4. **Execução de Binários e Scripts de Teste**: `./bin/file-server`, `./scripts/coverage.sh`.

---

## 🎯 3. Autonomia Operacional Sem Interrupções

- **Execução Direta**: Execute as tarefas do plano sequencialmente com total autonomia, sem solicitar confirmação ao usuário para etapas já aprovadas na especificação.
- **Segurança Determinística**: O arquivo `.agent/settings.json` já contém todas as permissões necessárias liberadas para o fluxo de desenvolvimento.

---

## 💡 4. Diretrizes Ativas de Economia de Tokens

1. **Leitura Cirúrgica**: Ao inspecionar partes de arquivos, sempre forneça `StartLine` e `EndLine` no `view_file` para evitar ler centenas de linhas irrelevantes na janela de contexto.
2. **Edição Pontual**: Prefira `replace_file_content` substituindo blocos contíguos específicos em vez de regravar arquivos inteiros com `write_to_file`.
3. **Respostas Enxutas**: Estruture as respostas de forma concisa. **NUNCA** replique blocos massivos de código ou conteúdo de arquivos na resposta do chat — utilize links clicáveis no formato `[nome_do_arquivo](file:///caminho/absoluto)`.
4. **Comandos com Saída Concisa**: Evite flags que gerem saídas prolixas (ex: evite `go test -v` a menos que seja estritamente necessário para isolar uma falha pontual).
