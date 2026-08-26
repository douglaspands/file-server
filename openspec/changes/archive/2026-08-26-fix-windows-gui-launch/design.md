## Context

A biblioteca `spf13/cobra` possui um mecanismo nativo para Windows (`inconshreveable/mousetrap`) que detecta quando um executável foi iniciado a partir de um clique duplo no Windows Explorer. Por padrão, a variável global `cobra.MousetrapHelpText` contém o aviso `"This is a command line tool. You need to open cmd.exe and run it from there."`.

Quando essa variável está preenchida, o Cobra intercepta a execução no Windows antes de invocar a função `RunE` do `RootCmd`, abortando o processo e impedindo que a lógica de detecção de ambiente desktop (`gui.IsDesktopEnvironment()`) e inicialização da interface gráfica seja executada.

Ver `proposal.md` para motivação detalhada.

## Goals / Non-Goals

**Goals:**
- Desativar a interceptação do *mousetrap* do Cobra definindo `cobra.MousetrapHelpText = ""` no pacote `cmd`.
- Garantir que o duplo clique no Windows Explorer invoque normalmente o comando raiz (`RootCmd.RunE`), disparando a interface gráfica desktop via `RunGUIWithOptions`.
- Manter o suporte completo a todos os subcomandos e flags CLI (`serve`, `ftp`, `sftp`, `version`, `--help`) sem qualquer regressão.
- Validar a configuração através de testes unitários com cobertura $\ge 80\%$.

**Non-Goals:**
- Não compilar binários separados para CLI e GUI; o `file-server` permanece um executável unificado e portátil.
- Não forçar `-ldflags="-H=windowsgui"` na compilação, o que suprimiria a saída padrão (stdout/stderr) quando executado em terminais (`cmd.exe`, PowerShell).

## Decisions

### Decisão 1: Configurar `cobra.MousetrapHelpText = ""` na inicialização da CLI
- **Justificativa**: Ao zerar `cobra.MousetrapHelpText` na função `init()` de `cmd/root.go`, o Cobra ignora a verificação do Explorer e delega o controle imediatamente para a execução do comando raiz.
- **Alternativas consideradas**:
  - *Separar em dois binários (`file-server.exe` e `file-server-gui.exe`)*: Aumentaria a complexidade de compilação, distribuição e manutenção.
  - *Flag `-H=windowsgui`*: Impediria o uso conveniente como ferramenta de linha de comando no terminal.

### Decisão 2: Estratégia de Testes Unitários e Validação
- Adicionar teste unitário em `cmd/cmd_test.go` verificando explicitamente que `cobra.MousetrapHelpText` está vazio, garantindo que alterações futuras não reativem o bloqueio acidentalmente.
- Validar suíte de testes com `make check` e `make test-coverage`.

## Risks / Trade-offs

- **[Janela de terminal em background no Windows]** → Ao dar duplo clique em um binário híbrido de console no Windows, o sistema abre uma janela de terminal que permanece ativa durante a execução do servidor local enquanto a janela gráfica/navegador é exibida. Essa janela fornece visibilidade dos logs do processo e encerra graciosamente com Ctrl+C ou fechamento da janela.
