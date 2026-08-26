## Why

Atualmente, o `file-server` é inicializado exclusivamente via linha de comando (CLI) ou executa o servidor HTTP padrão caso nenhum argumento seja passado em terminais interativos. Para usuários de desktop que interagem via gerenciador de arquivos (clique duplo no executável) ou preferem uma experiência visual moderna e intuitiva sem necessidade de memorizar flags e comandos do terminal, é fundamental fornecer uma interface gráfica nativa (GUI) desktop como forma principal de uso da aplicação. Essa interface permite configurar visualmente os parâmetros de rede, diretório compartilhado, credenciais, certificados TLS e alternar entre os serviços HTTP/Web, FTP e SFTP com um único clique, além de acompanhar o status, visualizar de forma organizada e escalável múltiplos IPs e facilitar a captura e compartilhamento dessas URLs via cópia ou QR Code.

## What Changes

- **Inicialização Gráfica Automática como Forma Principal**:
  - Ao executar o binário sem argumentos ou ao clicar duas vezes no executável em ambientes desktop (detecção de TTY vs ambiente gráfico ou flag explícita `--gui`), a aplicação abre uma interface gráfica desktop nativa.
  - Suporte ao comportamento padrão de aplicações desktop: redimensionamento fluido, maximização de janela e encerramento completo e gracioso da aplicação ao fechar a janela.
- **Interface Gráfica Estilo GNOME / Moderna**:
  - Design clean, polido e moderno inspirado em GNOME Adwaita Dark / libadwaita, utilizando rigorosamente a paleta de cores do frontend web (`slate-900`, `slate-950`, `indigo-500`, `slate-100`, bordas `slate-800`).
  - Seletor de modo de operação em abas ou segmented button: **Web / HTTP**, **FTP** e **SFTP**.
  - Painel de configurações intuitivo com campos de formulário e seletores:
    - Seleção de diretório compartilhado com botão de navegação nativa de pastas (*Folder Picker*).
    - Configuração de Host (`0.0.0.0`, `127.0.0.1`, etc.) e Porta (com valores padrão contextuais `8080`, `2121`, `2222`).
    - Ativação de criptografia TLS/HTTPS/FTPS (autoassinado ou seleção de arquivos PEM de certificado e chave).
    - Configurações específicas de FTP/SFTP (usuário, gerador de senhas aleatórias seguras, chaves SSH e modo somente leitura).
- **Visualização Escalável, Captura e Compartilhamento de Múltiplos IPs**:
  - Layout dinâmico e responsivo para exibição de IPs com viewport rolável estilizada (`overflow-y-auto` com scrollbar dark personalizada) quando houver muitos adaptadores de rede (Wi-Fi, Ethernet, Docker, WSL, VPNs).
  - Destaque automático do IP principal da rede local e Loopback no topo com tags de identificação de interface (`Wi-Fi`, `Ethernet`, `VPN`, `Virtual/Bridge`).
  - Campo de busca rápida/filtro instantâneo de interfaces e botão para expandir/recolher interfaces secundárias.
  - Botão de cópia rápida em 1 clique para cada endereço com feedback visual imediato ("✓ Copiado!").
  - Ação de "Copiar Mensagem de Compartilhamento" pronta para envio em ferramentas de chat (Slack, Discord, WhatsApp, Teams).
  - Modal / Visualizador de **QR Code** integrado para acesso instantâneo via smartphones e tablets na mesma rede Wi-Fi.
- **Controle de Execução e Status em Tempo Real**:
  - Botão de ação primário destacado ("Iniciar Servidor" / "Parar Servidor") com feedback visual imediato de estado.
  - Botão direto "Abrir no Navegador" para acesso web instantâneo.
  - Painel de log de eventos integrado na parte inferior da interface, exibindo a saída do servidor em tempo real com auto-scroll.
- **Guia Principal no README com Prints da Interface**:
  - O `README.md` será reformulado para posicionar a interface gráfica desktop como a forma principal e padrão de utilização do File Server, incluindo prints e capturas visuais da interface, passo a passo de inicialização, configuração e compartilhamento.
  - Todas as seções existentes de uso via linha de comando (CLI), flags avançadas, APIs e guia de compilação/desenvolvimento serão mantidas organizadas nas seções seguintes para uso avançado e headless.

## Capabilities

### New Capabilities
- `desktop-gui-launcher`: Interface gráfica desktop (estilo GNOME / Adwaita Dark) para configuração visual, seleção de protocolos (HTTP, FTP, SFTP), inicialização do servidor com um clique, visualização escalável e captura simplificada de múltiplos IPs/URLs (com viewport rolável, tags de interface, busca, cópia em 1 clique e QR Code), visualização de status/logs, documentação visual no README como forma primária de uso e ciclo de vida padrão de desktop (maximizar, fechar e encerrar).

### Modified Capabilities
<!-- Nenhuma especificação existente teve seus requisitos de comportamento alterados; os modos CLI e servidores web/ftp/sftp permanecem 100% retrocompatíveis. -->

## Impact

- **Código e CLI**: Criação do pacote `cmd/gui.go` e integração na lógica do comando raiz (`cmd/root.go`) para detecção e inicialização da GUI quando nenhum argumento for passado.
- **Módulo de GUI Desktop**: Criação de pacote dedicado `internal/adapters/gui` responsável por gerenciar a janela desktop, binding de configurações para os serviços existentes (`ServerOptions`, `adapterftp.ServerOptions`, `adaptersftp.ServerOptions`), descoberta categorizada de adaptadores de rede, integração com clipboard e controle de logs em tempo real.
- **Frontend / Assets**: Criação de template e componentes do launcher incluindo container rolável de IPs, badges de tipo de rede, busca de adaptadores, botões de cópia e renderizador de QR Code via script leve integrado.
- **Documentação**: Atualização do `README.md` e inclusão de capturas de tela/prints ilustrativos da GUI na pasta de assets de documentação.
- **Dependências & Build**: Atualização do `Makefile` para suportar compilação com suporte à GUI desktop nativa e execução multiplataforma.
