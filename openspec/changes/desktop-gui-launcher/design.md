## Context

O `file-server` foi arquitetado seguindo Clean Architecture com portas e adaptadores em Go. Os serviços de servidor HTTP/Web (`cmd/serve.go`), FTP (`cmd/ftp.go`) e SFTP (`cmd/sftp.go`) são orquestrados via contextos (`context.Context`) com cancelamento gracioso (`signal.NotifyContext`).
Os templates e assets web já utilizam Tailwind CSS, Alpine.js e HTMX embutidos via `embed.FS`.

Para viabilizar a experiência de aplicativo desktop moderna, clean e com a estética GNOME / Adwaita Dark solicitada, a solução integrará uma janela desktop nativa/webview conectada a um controlador de launcher interno (`internal/adapters/gui`), reutilizando o design system do frontend existente, orquestrando o ciclo de vida dos servidores HTTP, FTP e SFTP, e fornecendo um container rolável inteligente e escalável para visualização de qualquer quantidade de interfaces de rede com cópia rápida e compartilhamento via QR Code.

Além disso, a aplicação terá identidade visual oficial com ícone de alta resolução embutido no executável binário e renderizado na barra de tarefas e janela desktop, normalização de versão limpa (sem duplicar a letra `"v"` de tags do Git), e o `README.md` apresentará a interface gráfica como o método principal de utilização.

## Goals / Non-Goals

**Goals:**
- Prover inicialização automática da interface gráfica desktop quando o binário for executado sem argumentos de linha de comando ou por clique duplo no gerenciador de arquivos.
- Oferecer uma interface gráfica moderna, clean e com a estética GNOME Adwaita Dark (fundo `slate-900`/`slate-950`, acentos `indigo-500`, contrastes `slate-100` e bordas sutis `slate-800`).
- Criar e integrar ícone oficial clean e moderno de alta resolução (SVG, PNGs multi-resolução e ICO) embutido no executável (.exe) e visível na barra de tarefas e na janela da aplicação.
- Exibir a versão da aplicação sem duplicar prefixos `"v"` (considerando tags do GitHub como `v1.1.0` diretamente).
- Permitir a seleção e parametrização completa dos modos de operação (Web HTTP/HTTPS, FTP/FTPS, SFTP sobre SSH).
- Fornecer diálogo nativo de seleção de diretórios (*Folder Picker*) do sistema operacional.
- Fornecer uma experiência fluida para máquinas com múltiplos adaptadores de rede (Wi-Fi, Ethernet, Docker bridges, WSL, Tailscale/WireGuard):
  - Viewport com rolagem vertical suave e barra estilizada dark.
  - Categorização automática por tags de interface (`Wi-Fi`, `Ethernet`, `VPN`, `Docker/Bridge`).
  - Destaque do IP principal no topo e opção de busca/filtro rápido.
  - Ação de "Copiar Todos os IPs" e "Copiar Mensagem Formatada para Mensageiros".
  - Gerador de **QR Code** integrado para cada endereço LAN.
- Exibir status em tempo real, botão "Abrir no Navegador" e painel dinâmico de logs com auto-scroll.
- Comportamento idêntico a qualquer aplicativo desktop nativo: redimensionável, maximizável e encerramento gracioso completo ao fechar a janela.
- Reestruturar o `README.md` destacando a interface gráfica no topo com prints e diagramas como a forma principal de uso, mantendo a documentação completa de CLI e desenvolvedor logo abaixo.

**Non-Goals:**
- Substituir ou alterar o funcionamento dos subcomandos CLI existentes (`serve`, `ftp`, `sftp`, `version`).
- Criar dependências pesadas de terceiros ou frameworks que quebrem a portabilidade e a compilação do binário Go.

## Decisions

### 1. Arquitetura da Interface Gráfica: Webview Nativo com Controlador em Go
- **Decisão**: Adotar uma janela desktop nativa acoplada a um controlador de estado em Go (`internal/adapters/gui`), renderizando uma página com design system GNOME / Adwaita Dark baseado em Tailwind CSS e Alpine.js.
- **Racional**: Garante 100% de coerência visual e reaproveitamento de componentes da interface web do projeto (`web/`), mantendo alta responsividade, baixo consumo de memória e inicialização instantânea.

### 2. Identidade Visual e Embutimento Multiplataforma do Ícone da Aplicação
- **Decisão**:
  - Criação de master vetorial em SVG (`web/static/assets/icon.svg` e `docs/assets/icon.svg`) com design clean e moderno estilizando o símbolo do File Server (raio/energia em anil/índigo sob base geométrica escura).
  - Geração de pacote completo de ícones em alta resolução (`icon.ico`, `icon-16.png`, `icon-32.png`, `icon-48.png`, `icon-128.png`, `icon-256.png`, `icon-512.png`).
  - No Windows: embutimento via recurso de compilação `.syso` (compatível com `CGO_ENABLED=0`) garantindo o ícone visível no `.exe` pelo Windows Explorer.
  - No Desktop/Web: inclusão das tags `<link rel="icon">` e `<link rel="apple-touch-icon">` nos templates HTML (`gui_launcher.html` e `base.html`), fornecendo o ícone imediato na barra de tarefas, janelas e navegadores.

### 3. Normalização da Exibição de Versão
- **Decisão**:
  - Remover a adição manual de `"v"` nos formatos de log (`v%s` → `%s`), permitindo que a versão injetada via tag do Git/GitHub (ex: `v1.1.0`) ou default `dev` seja renderizada perfeitamente sem duplicidade (`vv1.1.0` → `v1.1.0`).

### 4. Tratamento Escalável de Múltiplos IPs & Categorização de Interfaces
- **Decisão**: 
  - O backend Go analisa as interfaces via `net.Interfaces()` e `services.GetLANIPAddresses()`, associando nome de interface (ex: `wlan0`, `eth0`, `docker0`, `tailscale0`) e classificando seu tipo (`wifi`, `ethernet`, `vpn`, `virtual`).
  - A interface web/desktop encapsula a lista de IPs em um container com altura controlada (`max-h-48` / `max-h-56` com `overflow-y-auto`), exibindo badges inteligentes, barra de pesquisa/filtro por nome ou IP, e botão para copiar todos os endereços de uma única vez.
  - Geração de QR Code leve no frontend para cada IP individual, facilitando a conexão imediata de celulares na rede local.

### 5. Estratégia de Documentação no README com Prints Visuais
- **Decisão**:
  - Salvar capturas de tela/prints ilustrativos de alta qualidade da interface gráfica em `docs/assets/gui-launcher.svg`.
  - Reestruturar o `README.md` iniciando com "🚀 Como Usar (Interface Gráfica Desktop)", exibindo o screenshot da aplicação, explicando o fluxo visual (escolher pasta, selecionar protocolo, iniciar e copiar links/QR Code).
  - Manter as seções de "💻 Modo Linha de Comando (CLI)", tabelas de flags, APIs REST e guia de desenvolvimento organizadas nas seções seguintes.

### 6. Diálogo Nativo de Seleção de Diretórios (Folder Picker)
- **Decisão**: Utilizar integração nativa do sistema operacional (Zenity/KDialog no Linux, PowerShell/FolderBrowserDialog no Windows, AppleScript/NSOpenPanel no macOS) acionada via API REST/SSE interna do launcher.

### 7. Orquestração e Ciclo de Vida dos Servidores
- **Decisão**: O controlador da GUI gerencia um `context.WithCancel` para cada servidor ativo. Ao clicar em "Iniciar Servidor", a rotina do serviço escolhido é iniciada em uma goroutine; ao clicar em "Parar" ou ao fechar a janela desktop, o contexto é cancelado, executando o `Shutdown` gracioso de 5 segundos.

### 8. Transmissão de Logs em Tempo Real
- **Decisão**: Implementar um `LogBroadcaster` em memória no pacote `internal/adapters/gui` que captura as mensagens de log dos servidores e as transmite para o painel de logs da interface via Server-Sent Events (SSE) ou HTMX/Alpine.

## Risks / Trade-offs

- **[Risco] Execução em ambientes Linux sem servidor gráfico / display (ex: SSH ou servidores puros)**:
  - *Mitigação*: Detectar se a variável `$DISPLAY` / `$WAYLAND_DISPLAY` está presente antes de abrir a janela gráfica; caso contrário, emitir mensagem orientando o uso das flags CLI ou iniciar no modo padrão headless.
- **[Risco] Encerramento abrupto da janela desktop**:
  - *Mitigação*: O listener de fechamento de janela captura o evento de encerramento (`WM_DELETE_WINDOW` / `close`), aciona o cancelamento do contexto de todos os servidores e finaliza a aplicação (`os.Exit(0)`).

## Migration Plan

- Não há impacto de quebra (breaking change).
- Ao executar `file-server` diretamente sem parâmetros em ambiente gráfico desktop, a GUI é iniciada.
- Todos os comandos anteriores como `file-server serve`, `file-server ftp` ou execução em scripts continuam funcionando sem qualquer alteração.
