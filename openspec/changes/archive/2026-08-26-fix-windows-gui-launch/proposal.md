## Why

No ambiente Windows, ao clicar duas vezes diretamente sobre o arquivo executável (`.exe`) a partir do Windows Explorer, a biblioteca Cobra intercepta a execução com sua rotina padrão de *mousetrap*, exibindo a mensagem `"This is a command line tool. You need to open cmd.exe and run it from there."` e encerrando o processo sem iniciar a interface gráfica desktop.

Esta mudança é necessária para permitir que usuários do Windows iniciem o File Server diretamente com duplo clique no executável de forma fluida e sem dependência de terminal, em conformidade com o propósito de ferramenta desktop amigável e acessível.

## What Changes

- Desativação global do comportamento de *mousetrap* da biblioteca Cobra (`cobra.MousetrapHelpText = ""`), permitindo que a execução sem argumentos no Windows Explorer prossiga normalmente para a detecção de ambiente desktop.
- Garantia de que a inicialização no Windows via duplo clique acione `gui.IsDesktopEnvironment()` e abra a janela da interface gráfica desktop (modo app ou navegador padrão).
- Inclusão de testes unitários e de integração validando que o *mousetrap* está desabilitado e que o fluxo de inicialização raiz responde corretamente.

## Capabilities

### New Capabilities
<!-- Nenhuma nova capacidade introduzida -->

### Modified Capabilities
- `desktop-gui-launcher`: Ajuste e reforço do cenário de inicialização automática no Windows Explorer via duplo clique sem interceptação de alerta de terminal.

## Impact

- **Código Afetado**: `cmd/root.go` (configuração de `cobra.MousetrapHelpText`), `cmd/cmd_test.go` ou `cmd/gui_test.go` (testes de inicialização).
- **Compatibilidade**: Total retrocompatibilidade com modo CLI headless e comandos via terminal (`file-server serve`, `file-server ftp`, `file-server sftp`).
- **Dependências**: Nenhuma dependência externa adicionada ou alterada.
