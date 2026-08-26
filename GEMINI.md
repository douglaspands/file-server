# Diretrizes do Projeto para o Antigravity (AGY)

## 📌 Identidade do Projeto
- **Nome**: File Server (`file-server`)
- **Propósito Central**: Servidor web de arquivos de alta performance para rede local (LAN), com navegação fluida, streaming direto com suporte a HTTP Range (206), downloads de pastas compactadas em ZIP via streaming sob demanda (zero resíduos em disco), uploads multipart/drag-and-drop, isolamento estrito de sandbox contra path traversal e criptografia em trânsito (TLS/HTTPS).
- **Padrões Técnicos**: Clean Architecture em Go, templates SSR embutidos (`embed.FS`), TDD/BDD com meta inegociável de cobertura &ge; 80% e Makefile universal.

---

## 🌿 Boas Práticas de Git e Controle de Branches

1. **Feature Branch por Especificação**:
   - A cada nova especificação criada no OpenSpec (`/openspec-propose` ou `openspec new change "<name>"`), deve ser criada uma branch de feature dedicada a partir da `main`:
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
