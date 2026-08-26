## MODIFIED Requirements

### Requirement: Inicialização Automática da Interface Gráfica Desktop
O sistema DEVE abrir automaticamente a interface gráfica desktop (GUI) quando o executável for iniciado sem argumentos em ambiente gráfico ou quando o comando/flag explícito for informado, sem ser bloqueado ou interceptado por proteções de terminal CLI como o mousetrap no Windows Explorer.

#### Scenario: Execução sem argumentos ou clique duplo no executável
- **GIVEN** que o usuário está em um ambiente desktop (Linux/GNOME, macOS ou Windows)
- **WHEN** o usuário clica duas vezes no executável diretamente no Windows Explorer ou executa `file-server` no terminal sem parâmetros
- **THEN** o sistema DEVE inicializar e exibir a janela principal da interface gráfica desktop com as configurações padrão pré-carregadas sem exibir mensagens de erro de linha de comando ou exigir abertura prévia de `cmd.exe`.

#### Scenario: Execução explícita via comando gui
- **WHEN** o usuário executa `file-server gui` ou `file-server --gui`
- **THEN** o sistema DEVE iniciar a interface gráfica desktop mesmo a partir do terminal.

#### Scenario: Execução com subcomandos de linha de comando
- **WHEN** o usuário executa comandos específicos como `file-server serve`, `file-server ftp`, `file-server sftp` ou `file-server version`
- **THEN** o sistema DEVE executar no modo CLI headless sem exibir a interface gráfica, mantendo total retrocompatibilidade.
