## Context

O `file-server` foi arquitetado seguindo Clean Architecture com portas e adaptadores em Go. Os serviços de servidor HTTP/Web (`cmd/serve.go`), FTP (`cmd/ftp.go`) e SFTP (`cmd/sftp.go`) são orquestrados via contextos (`context.Context`) com cancelamento gracioso (`signal.NotifyContext`).
Os templates e assets web já utilizam Tailwind CSS, Alpine.js e HTMX embutidos via `embed.FS`.

Para viabilizar a experiência de aplicativo desktop moderna, clean e com a estética GNOME / Adwaita Dark solicitada, a solução integrará uma janela desktop nativa/webview conectada a um controlador de launcher interno (`internal/adapters/gui`), reutilizando o design system do frontend existente, orquestrando o ciclo de vida dos servidores HTTP, FTP e SFTP, e fornecendo um container rolável inteligente e escalável para visualização de qualquer quantidade de interfaces de rede com cópia rápida e compartilhamento via QR Code.

Além disso, a documentação central do projeto (`README.md`) será reestruturada para apresentar a interface gráfica como o método principal e mais amigável de utilização, suportada por capturas visuais da GUI.

## Goals / Non-Goals

**Goals:**
- Prover inicialização automática da interface gráfica desktop quando o binário for executado sem argumentos de linha de comando ou por clique duplo no gerenciador de arquivos.
- Oferecer uma interface gráfica moderna, clean e com a estética GNOME Adwaita Dark (fundo `slate-900`/`slate-950`, acentos `indigo-500`, contrastes `slate-100` e bordas sutis `slate-800`).
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
- **Alternativas consideradas**:
  - *Fyne / Gio*: Embora sejam bibliotecas Go puras, exigiriam duplicar toda a estilização CSS/HTML em código Go e dificultariam atingir exatamente a paleta e acabamento visual do frontend web existente.
  - *Electron*: Descartado por gerar binários pesados (> 100MB) e alto consumo de memória.

### 2. Tratamento Escalável de Múltiplos IPs & Categorização de Interfaces
- **Decisão**: 
  - O backend Go analisa as interfaces via `net.Interfaces()` e `services.GetLANIPAddresses()`, associando nome de interface (ex: `wlan0`, `eth0`, `docker0`, `tailscale0`) e classificando seu tipo (`wifi`, `ethernet`, `vpn`, `virtual`).
  - A interface web/desktop encapsula a lista de IPs em um container com altura controlada (`max-h-48` / `max-h-56` com `overflow-y-auto`), exibindo badges inteligentes, barra de pesquisa/filtro por nome ou IP, e botão para copiar todos os endereços de uma única vez.
  - Geração de QR Code leve no frontend para cada IP individual, facilitando a conexão imediata de celulares na rede local.
- **Racional**: Garante que a aplicação se comporte com perfeição tanto em laptops simples (com 1 placa Wi-Fi) quanto em estações de trabalho de desenvolvedores avançados com dezenas de interfaces virtuais, pontes Docker e VPNs, sem poluir ou quebrar a janela.

### 3. Estratégia de Documentação no README com Prints Visuais
- **Decisão**:
  - Salvar capturas de tela/prints ilustrativos de alta qualidade da interface gráfica em `docs/assets/gui-launcher.png`.
  - Reestruturar o `README.md` iniciando com "🚀 Como Usar (Interface Gráfica Desktop)", exibindo o screenshot da aplicação, explicando o fluxo visual (escolher pasta, selecionar protocolo, iniciar e copiar links/QR Code).
  - Manter as seções de "💻 Modo Linha de Comando (CLI)", tabelas de flags, APIs REST e guia de desenvolvimento organizadas nas seções seguintes.
- **Racional**: Proporciona excelente primeira impressão para usuários visuais e desktop, sem perder o detalhamento técnico exigido por administradores de sistema e desenvolvedores.

### 4. Diálogo Nativo de Seleção de Diretórios (Folder Picker)
- **Decisão**: Utilizar integração nativa do sistema operacional (Zenity/KDialog no Linux, PowerShell/FolderBrowserDialog no Windows, AppleScript/NSOpenPanel no macOS) acionada via API REST/SSE interna do launcher.
- **Racional**: Proporciona a experiência familiar de seleção de pastas nativa do GNOME/Desktop sem necessidade de CGo ou bibliotecas pesadas de terceiros.

### 5. Orquestração e Ciclo de Vida dos Servidores
- **Decisão**: O controlador da GUI gerencia um `context.WithCancel` para cada servidor ativo. Ao clicar em "Iniciar Servidor", a rotina do serviço escolhido é iniciada em uma goroutine; ao clicar em "Parar" ou ao fechar a janela desktop, o contexto é cancelado, executando o `Shutdown` gracioso de 5 segundos.
- **Racional**: Reutiliza 100% dos adaptadores e serviços existentes (`services.NewFileService`, `adapterftp.NewServer`, `adaptersftp.NewServer`), garantindo segurança, integridade e isolamento.

### 6. Transmissão de Logs em Tempo Real
- **Decisão**: Implementar um `LogBroadcaster` em memória no pacote `internal/adapters/gui` que captura as mensagens de log dos servidores e as transmite para o painel de logs da interface via Server-Sent Events (SSE) ou HTMX/Alpine.
- **Racional**: Permite ao usuário inspecionar requisições, downloads, conexões e eventuais erros diretamente na janela do aplicativo sem precisar de um terminal aberto.

## Risks / Trade-offs

- **[Risco] Execução em ambientes Linux sem servidor gráfico / display (ex: SSH ou servidores puros)**:
  - *Mitigação*: Detectar se a variável `$DISPLAY` / `$WAYLAND_DISPLAY` está presente antes de abrir a janela gráfica; caso contrário, emitir mensagem orientando o uso das flags CLI ou iniciar no modo padrão headless.
- **[Risco] Encerramento abrupto da janela desktop**:
  - *Mitigação*: O listener de fechamento de janela captura o evento de encerramento (`WM_DELETE_WINDOW` / `close`), aciona o cancelamento do contexto de todos os servidores e finaliza a aplicação (`os.Exit(0)`).

## Migration Plan

- Não há impacto de quebra (breaking change).
- Ao executar `file-server` diretamente sem parâmetros em ambiente gráfico desktop, a GUI é iniciada.
- Todos os comandos anteriores como `file-server serve`, `file-server ftp` ou execução em scripts continuam funcionando sem qualquer alteração.
