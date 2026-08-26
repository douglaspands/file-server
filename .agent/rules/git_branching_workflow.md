# Diretrizes de Git e Fluxo de Branches por Especificação (Git Feature Branching)

Este documento define as boas práticas de Git obrigatórias para o Antigravity (AGY) durante todo o ciclo de vida de especificações e mudanças no projeto.

---

## 🌿 1. Criação de Branch de Feature por Spec

A cada nova especificação ou mudança planejada no OpenSpec (`/openspec-propose`, `/opsx-propose` ou `openspec new change`):

1. **Origem na `main`**:
   Certifique-se de que a branch base `main` está atualizada e limpa.

2. **Criação da Branch**:
   Crie e alterne para uma branch dedicada com prefixo `feature/` correspondente ao nome da mudança:
   ```bash
   git checkout main
   git checkout -b feature/<change-name>
   ```

3. **Isolamento**:
   Todas as alterações de planejamento, código, testes e documentação relacionadas à especificação devem residir estritamente dentro da branch `feature/<change-name>`.

---

## 🛠️ 2. Desenvolvimento e Commits

Durante a fase de implementação (`/openspec-apply-change` ou `/opsx-apply`):

- **Conventional Commits**: Utilize o padrão Conventional Commits (`feat:`, `fix:`, `refactor:`, `test:`, `docs:`, `chore:`).
- **Validação Contínua**: Valide periodicamente com `make check` e `make test` mantendo cobertura &ge; 80%.

---

## 📦 3. Conclusão, Arquivamento, Merge e Pull Requests (Sempre Squash)

Ao concluir todas as tarefas da especificação:

1. **Validação Final e Arquivamento**:
   ```bash
   make check
   openspec archive <change-name> --yes
   ```

2. **Commit de Arquivamento na Feature Branch**:
   ```bash
   git add .
   git commit -m "feat(spec): archive <change-name> and apply changes"
   ```

3. **Solicitação Obrigatória de Permissão**:
   > ⚠️ **REGRA CRÍTICA PARA O AGY**: O merge com a branch `main` ou criação de Pull Request **NUNCA** deve ser executado de forma automática ou silenciosa.
   > O AGY deve **SEMPRE solicitar a permissão explícita do usuário** antes de realizar o merge ou abrir uma PR.

4. **Estratégia Obrigatória de Integração: SQUASH**:
   
   - **Cenário A: Merge Local Direto na `main` (Após Autorização)**:
     O merge local deve ser realizado **obrigatoriamente via Squash**, gerando um único commit limpo na `main`:
     ```bash
     git checkout main
     git merge --squash feature/<change-name>
     git commit -m "feat(<change-name>): <resumo consolidado das mudanças da spec>"
     ```

   - **Cenário B: Pull Request no GitHub (com Repositório Remoto Configurado)**:
     Quando houver remote origin configurado e for autorizada a abertura de PR:
     ```bash
     git push -u origin feature/<change-name>
     gh pr create --base main --head feature/<change-name> --title "feat(<change-name>): ..." --body "..."
     ```
     A estratégia de merge da PR no GitHub deve ser **estritamente Squash and Merge** (`gh pr merge --squash` ou selecionando a opção "Squash and merge" na interface web).
