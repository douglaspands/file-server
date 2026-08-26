## Why

Atualmente, o `file-server` é inicializado exclusivamente via linha de comando (CLI) ou executa o servidor HTTP padrão caso nenhum argumento seja passado em terminais interativos. Para usuários de desktop que interagem via gerenciador de arquivos (clique duplo no executável) ou preferem uma experiência visual moderna e intuitiva sem necessidade de memorizar flags e comandos do terminal, é fundamental fornecer uma interface gráfica nativa (GUI) desktop como forma principal de uso da aplicação. Essa interface permite configurar visualmente os parâmetros de rede, diretório compartilhado, credenciais, certificados TLS e alternar entre os serviços HTTP/Web, FTP e SFTP com um único clique, além de acompanhar o status, visualizar de forma organizada e escalável múltiplos IPs e facilitar a captura e compartilhamento dessas URLs via cópia ou QR Code.

Além disso, a aplicação possui identidade visual oficial com ícone de alta resolução embutido no executável binário e renderizado na barra de tarefas e na janela desktop, formatação limpa da versão sem duplicação de prefixos `"v"`, e tematização escura integrada da moldura da janela do navegador/Chrome nas cores da aplicação (`slate-950` / `#020617`).

## What Changes

- **Tematização da Janela do Navegador / Modo App**:
  - Suporte a tema escuro na janela do navegador via meta tags (`theme-color: #020617`, `color-scheme: dark`), Web App Manifest (`manifest.json`) e flags de modo escuro forçado do Chrome/Edge (`--force-dark-mode`, `--enable-features=WebUIDarkMode`), integrando a moldura e a barra de título da janela perfeitamente à paleta `slate-950`.
- **Normalização da Exibição de Versão**:
  - Ajuste nos formatos de log e templates para remover a concatenação redundante do prefixo `"v"`, evitando duplicidade (`vv1.1.0` → `v1.1.0`) e considerando a tag oficial do GitHub diretamente.
- **Identidade Visual e Ícone da Aplicação em Alta Resolução**:
  - Criação de ícone oficial clean, moderno e vetorial (SVG) com tema de servidor/energia de alta performance nas cores da paleta (`slate-900`/`slate-950`, `indigo-500`, `slate-100`).
  - Pacote de ícones multi-resolução (`icon.svg`, `icon-16.png`, `icon-32.png`, `icon-48.png`, `icon-128.png`, `icon-256.png`, `icon-512.png` e `.ico` multi-camadas).
  - Embutimento do ícone no arquivo executável binário (`.syso` para Windows) e vinculação nos cabeçalhos HTML para exibição na barra de tarefas, dock, janelas desktop e abas do navegador.
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
- `desktop-gui-launcher`: Interface gráfica desktop (estilo GNOME / Adwaita Dark) para configuração visual, seleção de protocolos (HTTP, FTP, SFTP), inicialização do servidor com um clique, visualização escalável e captura simplificada de múltiplos IPs/URLs (com viewport rolável, tags de interface, busca, cópia em 1 clique e QR Code), tematização da janela do navegador nas cores da aplicação, identidade visual com ícone em alta resolução no executável e na barra de tarefas, normalização de versão de release, visualização de status/logs, documentação visual no README como forma primária de uso e ciclo de vida padrão de desktop (maximizar, fechar e encerrar).

### Modified Capabilities
<!-- Nenhuma especificação existente teve seus requisitos de comportamento alterados; os modos CLI e servidores web/ftp/sftp permanecem 100% retrocompatíveis. -->

## Impact

- **Código e CLI**: Criação do pacote `cmd/gui.go` e integração na lógica do comando raiz (`cmd/root.go`) para detecção e inicialização da GUI quando nenhum argumento for passado.
- **Módulo de GUI Desktop**: Criação de pacote dedicado `internal/adapters/gui` responsável por gerenciar a janela desktop com flags de tema escuro, binding de configurações para os serviços existentes (`ServerOptions`, `adapterftp.ServerOptions`, `adaptersftp.ServerOptions`), descoberta categorizada de adaptadores de rede, integração com clipboard e controle de logs em tempo real.
- **Frontend / Assets**: Criação de template e componentes do launcher incluindo container rolável de IPs, badges de tipo de rede, busca de adaptadores, botões de cópia, renderizador de QR Code, `manifest.json` com `theme_color` e pacote de assets com ícone em alta resolução (`web/static/assets/icon.svg`, `icon.ico`, `icon-*.png`).
- **Documentação**: Atualização do `README.md` e inclusão de capturas de tela/prints ilustrativos da GUI na pasta de assets de documentação.
- **Dependências & Build**: Atualização do `Makefile` para suportar compilação com suporte à GUI desktop nativa, geração de recursos de ícone para Windows (`.syso`) e execução multiplataforma.
