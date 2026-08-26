# Fluxo de Arquivamento, Commit e Merge de Especificações (OpenSpec)

Ao concluir a implementação de uma mudança planejada no OpenSpec:

1. **Validação**:
   Execute `make check` e `openspec validate --all` para garantir que todos os testes passem (cobertura >= 80%) e todas as especificações estejam íntegras.

2. **Arquivamento**:
   Execute o arquivamento da mudança via comando:
   ```bash
   openspec archive <change-name> --yes
   ```
   Isso sincroniza as especificações delta para a pasta principal `openspec/specs/` e move a mudança para `openspec/changes/archive/`.

3. **Commit Git na Feature Branch**:
   Crie o commit no Git na branch da feature com mensagem padronizada no formato Conventional Commits:
   ```bash
   git add .
   git commit -m "feat(spec): archive <change-name> and apply changes"
   ```

4. **Solicitação de Permissão para Merge / Pull Request**:
   Após o arquivamento e commit na branch `feature/<change-name>`, o AGY deve **obrigatoriamente perguntar ao usuário se deseja realizar o merge local com a branch `main` ou abrir uma Pull Request no GitHub**.

5. **Estratégia Obrigatória de Integração: SQUASH**:
   - Se o usuário autorizar o **merge local**:
     ```bash
     git checkout main
     git merge --squash feature/<change-name>
     git commit -m "feat(<change-name>): <resumo consolidado das mudanças da spec>"
     ```
   - Se houver **repositório remoto (GitHub)** e for solicitada PR:
     A Pull Request deve ser aberta e mesclada **exclusivamente com Squash and Merge** (`--squash`).
