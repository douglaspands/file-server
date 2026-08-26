# Diretrizes do Projeto para o Antigravity (AGY)

## 📌 Identidade do Projeto
- **Nome**: File Server (`file-server`)
- **Propósito Central**: Servidor web de arquivos de alta performance para rede local (LAN), com navegação fluida, streaming direto com suporte a HTTP Range (206), downloads de pastas compactadas em ZIP via streaming sob demanda (zero resíduos em disco), uploads multipart/drag-and-drop, isolamento estrito de sandbox contra path traversal e criptografia em trânsito (TLS/HTTPS).
- **Padrões Técnicos**: Clean Architecture em Go, templates SSR embutidos (`embed.FS`), TDD/BDD com meta inegociável de cobertura &ge; 80% e Makefile universal.

---

## ⚡ Autonomia Operacional, Ferramentas Nativas e Economia de Tokens

1. **Prioridade Mandatória de Ferramentas Nativas**:
   - **Criar Arquivo**: Utilize `write_to_file`. É **ESTRITAMENTE PROIBIDO** criar arquivos via `cat << 'EOF' > ...`, `echo "..." > ...` ou `touch` no terminal.
   - **Editar Arquivo**: Utilize `replace_file_content` para edições cirúrgicas e pontuais. Nunca execute scripts `sed`, `awk` ou `cat >` via terminal.
   - **Inspecionar / Ler**: Utilize `view_file` especificando `StartLine` e `EndLine` para focar apenas no trecho relevante. Nunca execute `cat`, `head`, `tail` via terminal.
   - **Buscar Código e Arquivos**: Utilize `grep_search`, `find_by_name` e `list_dir`. Nunca execute `grep`, `find` ou `ls` via terminal.

2. **Uso Restrito do Terminal (`run_command`)**:
   - O terminal deve ser utilizado exclusivamente para ferramentas do ciclo de vida: `make` (`make test`, `make lint`, `make check`, `make build`), `go` (`go test`, `go mod tidy`), `git`, `openspec` e binários executáveis.

3. **Execução Autônoma e Eficiência de Tokens**:
   - Execute tarefas do plano sequencialmente sem interrupções desnecessárias por confirmação de prompt.
   - **Economia de Tokens**: Forneça respostas concisas, utilizando links Markdown no padrão `[arquivo](file:///caminho)`. **NUNCA** replique blocos inteiros de código já salvos em disco na resposta do chat.

---

## 🌿 Boas Práticas de Git e Controle de Branches

1. **Feature Branch por Especificação**:
   - A cada nova especificação criada no OpenSpec (`/opsx-propose`, `/openspec-propose` ou `openspec new change "<name>"`), deve ser criada uma branch de feature dedicada a partir da `main`:
     ```bash
     git checkout main
     git checkout -b feature/<change-name>
     ```
   - Todo o planejamento, desenvolvimento e testes ocorrem na branch `feature/<change-name>`.

2. **Desenvolvimento e Commits**:
   - Utilize mensagens no padrão Conventional Commits (`feat:`, `fix:`, `refactor:`, `test:`, `docs:`, `chore:`).
   - Mantenha `make test` e `make check` validados continuamente.

3. **Arquivamento, Merge e Pull Requests (Sempre Squash)**:
   - Ao concluir a implementação de uma mudança e executar o arquivamento (`openspec archive <change-name> --yes`):
     1. Realize o commit final na branch `feature/<change-name>`.
     2. **SOLICITAÇÃO OBRIGATÓRIA DE PERMISSÃO**: O AGY deve **SEMPRE solicitar a permissão explícita do usuário** antes de fazer o merge com a `main` ou abrir Pull Request (PR).
     3. **Tipo de Merge Obrigatório (Squash)**:
        - **Merge Local**: Quando autorizado a realizar o merge local com a `main`, deve ser executado **obrigatoriamente via Squash**:
          ```bash
          git checkout main
          git merge --squash feature/<change-name>
          git commit -m "feat(<change-name>): <resumo consolidado das mudanças>"
          ```
        - **Pull Request no GitHub (com Remote)**: Se houver repositório remoto configurado (GitHub) e for solicitada a abertura de PR (`gh pr create` ou link web), a estratégia de integração configurada deve ser **obrigatoriamente Squash and Merge** (`--squash`).
