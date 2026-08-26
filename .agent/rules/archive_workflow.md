# Fluxo de Arquivamento e Commit de Especificações (OpenSpec)

Ao concluir a implementação de uma mudança planejada no OpenSpec:

1. **Validação**:
   Execute `make check` e `openspec validate --all` para garantir que todos os testes passem (cobertura >= 80%) e todas as especificações estejam íntegras.

2. **Arquivamento**:
   Execute o arquivamento da mudança via comando:
   ```bash
   openspec archive <change-name> --yes
   ```
   Isso sincroniza as especificações delta para a pasta principal `openspec/specs/` e move a mudança para `openspec/changes/archive/`.

3. **Commit Git Obrigatório**:
   Crie o commit no Git com mensagem padronizada no formato Conventional Commits:
   ```bash
   git add .
   git commit -m "feat(spec): archive <change-name> and apply changes"
   ```
