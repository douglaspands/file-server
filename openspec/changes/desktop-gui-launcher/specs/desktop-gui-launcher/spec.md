## Purpose

Fornece uma interface gráfica nativa e moderna para desktop (estilo GNOME / Adwaita Dark) para inicialização interativa do File Server sem necessidade de argumentos de linha de comando ou ao clicar duas vezes no executável, permitindo configuração visual de serviços (Web/HTTP, FTP, SFTP), identidade visual oficial com ícone clean e moderno de alta resolução embutido no executável e visível na barra de tarefas/janela, formatação limpa da versão sem prefixos duplicados, visualização escalável, amigável e captura facilitada de múltiplos endereços IP/URLs para compartilhamento, documentação visual no README como forma principal de uso, e controle completo do ciclo de vida da aplicação.

## ADDED Requirements

### Requirement: Inicialização Automática da Interface Gráfica Desktop
O sistema DEVE abrir automaticamente a interface gráfica desktop (GUI) quando o executável for iniciado sem argumentos em ambiente gráfico ou quando o comando/flag explícito for informado.

#### Scenario: Execução sem argumentos ou clique duplo no executável
- **GIVEN** que o usuário está em um ambiente desktop (Linux/GNOME, macOS ou Windows)
- **WHEN** o usuário clica duas vezes no executável ou executa `file-server` no terminal sem parâmetros
- **THEN** o sistema DEVE inicializar e exibir a janela principal da interface gráfica desktop com as configurações padrão pré-carregadas.

#### Scenario: Execução explícita via comando gui
- **WHEN** o usuário executa `file-server gui` ou `file-server --gui`
- **THEN** o sistema DEVE iniciar a interface gráfica desktop mesmo a partir do terminal.

#### Scenario: Execução com subcomandos de linha de comando
- **WHEN** o usuário executa comandos específicos como `file-server serve`, `file-server ftp`, `file-server sftp` ou `file-server version`
- **THEN** o sistema DEVE executar no modo CLI headless sem exibir a interface gráfica, mantendo total retrocompatibilidade.

### Requirement: Identidade Visual e Ícone da Aplicação em Alta Resolução no Executável e Barra de Tarefas
O sistema DEVE fornecer um ícone oficial de alta resolução clean e moderno, visível tanto no arquivo executável compilado quanto na barra de tarefas e janela durante a execução.

#### Scenario: Exibição do ícone no executável do sistema operacional
- **GIVEN** que o binário executável foi compilado para a plataforma de destino
- **WHEN** o usuário visualiza o arquivo executável no gerenciador de arquivos (ex: Windows Explorer ou desktop Linux)
- **THEN** o sistema DEVE exibir o ícone oficial personalizado do File Server associado ao arquivo binário.

#### Scenario: Exibição do ícone na barra de tarefas e janela ativa
- **GIVEN** que a aplicação foi inicializada em modo GUI desktop ou explorador web
- **WHEN** a janela do aplicativo é aberta
- **THEN** o sistema operacional e a barra de tarefas/dock DEVEM exibir o ícone personalizado de alta resolução na barra de título e no alternador de tarefas.

#### Scenario: Disponibilização do pacote de assets em múltiplas resoluções
- **WHEN** os recursos visuais são carregados
- **THEN** o sistema DEVE fornecer o ícone em formato vetorial SVG e bitmaps multi-resolução (16x16, 32x32, 48x48, 128x128, 256x256, 512x512 e .ico).

### Requirement: Exibição Normalizada de Versão sem Duplicidade de Prefixo
O sistema DEVE formatar e exibir as informações de versão de forma limpa, respeitando a tag do GitHub sem duplicar a letra "v".

#### Scenario: Formatação limpa da versão sem prefixo duplicado
- **GIVEN** que o binário foi compilado com tag de release do Git (ex: v1.1.0)
- **WHEN** a versão é exibida nos banners de inicialização, templates HTML ou terminal
- **THEN** o sistema DEVE exibir a versão sem duplicar a letra "v" (exibindo "v1.1.0" ou "dev").

### Requirement: Design Visual Estilo GNOME e Paleta de Cores Unificada
A interface gráfica DEVE adotar um design clean, moderno e polido no padrão GNOME / Adwaita Dark, utilizando estritamente a paleta de cores do frontend web (`slate-900`, `slate-950`, `indigo-500`, `slate-100`, bordas `slate-800`).

#### Scenario: Aplicação do tema visual escuro integrado
- **WHEN** a janela da interface gráfica é renderizada
- **THEN** o fundo DEVE utilizar tons escuros (`slate-900`/`slate-950`), textos em alto contraste (`slate-100`/`slate-300`), destaques e botões de ação em anil (`indigo-500`/`indigo-600`) e bordas sutis (`slate-800`).

#### Scenario: Navegação por abas ou seletor de protocolo
- **WHEN** o usuário visualiza a barra superior ou controle segmentado da janela
- **THEN** o sistema DEVE apresentar seletores claros para os 3 modos de serviço: "Web (HTTP/HTTPS)", "FTP / FTPS" e "SFTP (SSH)".

### Requirement: Seleção e Configuração Visual de Parâmetros
A interface gráfica DEVE disponibilizar controles visuais intuitivos para todos os parâmetros de configuração suportados pelo servidor.

#### Scenario: Seleção de diretório compartilhado com diálogo nativo
- **WHEN** o usuário clica no botão de busca/navegação de diretório (ícone de pasta)
- **THEN** o sistema DEVE abrir a caixa de diálogo nativa do sistema operacional para seleção de pasta e preencher o campo de caminho com o diretório escolhido.

#### Scenario: Configuração de portas e host
- **WHEN** o usuário seleciona um modo de serviço (Web, FTP ou SFTP)
- **THEN** o sistema DEVE preencher os campos com os valores padrão recomendados (`8080` para Web, `2121` para FTP, `2222` para SFTP, host `0.0.0.0`) e permitir alteração livre pelo usuário.

#### Scenario: Configuração de segurança e TLS
- **WHEN** o usuário ativa o seletor de TLS/Criptografia
- **THEN** o sistema DEVE permitir escolher entre certificado autoassinado instantâneo ou selecionar arquivos PEM customizados de certificado e chave privada.

#### Scenario: Configuração de autenticação e geração de senhas seguras
- **WHEN** o usuário acessa as abas de FTP ou SFTP
- **THEN** o sistema DEVE fornecer campos de usuário, senha, botão "Gerar Senha Segura" (gerador aleatório de 12 caracteres) e opção de modo somente leitura.

### Requirement: Visualização Estruturada, Escalável e Compartilhamento de Múltiplos IPs
A interface gráfica DEVE acomodar de forma elegante e amigável qualquer quantidade de interfaces de rede (desde 1 até dezenas de adaptadores como Wi-Fi, Ethernet, Docker, WSL, VPNs) sem quebrar o layout da janela.

#### Scenario: Rolagem suave e viewport adaptável para múltiplos IPs
- **GIVEN** que o sistema detectou múltiplos IPs de rede (3 ou mais interfaces)
- **WHEN** o usuário visualiza a lista de endereços ativos
- **THEN** o sistema DEVE acomodar os cards dentro de um container com rolagem vertical suave estilizada (`max-h-48` / `max-h-56`) e contador de adaptadores (ex: "4 interfaces ativas"), impedindo que a janela principal seja deformada ou esticada.

#### Scenario: Destaque inteligente e categorização por tipo de interface
- **WHEN** a lista de interfaces é renderizada
- **THEN** o sistema DEVE posicionar o Loopback Local e a interface principal física (LAN/Wi-Fi) no topo com badge "Recomendado" / "Principal", identificando visualmente cada placa com tags sutis (`Wi-Fi`, `Ethernet`, `VPN`, `Docker/Bridge`).

#### Scenario: Filtro ou busca rápida de interfaces
- **GIVEN** muitas interfaces de rede detectadas
- **WHEN** o usuário digita no campo de busca rápida ou clica nas tags de filtro
- **THEN** a lista DEVE filtrar instantaneamente apenas os endereços correspondentes.

#### Scenario: Cópia rápida para área de transferência em 1 clique
- **WHEN** o usuário clica no botão de cópia ao lado de qualquer endereço IP ou URL
- **THEN** o sistema DEVE copiar o endereço correspondente diretamente para a área de transferência do sistema operacional e fornecer feedback visual instantâneo (ex: "✓ Copiado!").

#### Scenario: Copiar mensagem pronta para ferramentas de comunicação
- **WHEN** o usuário clica no botão "Compartilhar Link" / "Copiar Texto de Acesso"
- **THEN** o sistema DEVE copiar uma mensagem formatada contendo o nome do compartilhamento, tipo de protocolo, URLs de acesso categorizadas e credenciais temporárias (caso ativas), pronta para colar em mensageiros como Slack, WhatsApp, Teams ou Discord.

#### Scenario: Exibição de QR Code para conexões móveis
- **WHEN** o usuário clica no botão "QR Code" associado a um endereço LAN
- **THEN** o sistema DEVE abrir um modal ou popover exibindo o QR Code gerado para a URL, permitindo que smartphones e tablets na mesma rede Wi-Fi acessem instantaneamente ao escanear a tela.

### Requirement: Controle de Execução e Monitoramento em Tempo Real
A interface DEVE fornecer botões claros de inicialização/interrupção e exibir o status de execução, links de acesso e logs em tempo real.

#### Scenario: Iniciar servidor com um clique
- **GIVEN** configurações válidas preenchidas na interface
- **WHEN** o usuário clica no botão principal "Iniciar Servidor"
- **THEN** o sistema DEVE iniciar a rotina do servidor em background, alterar o estado do botão para "Parar Servidor", exibir indicador de status "🟢 Em Execução" e renderizar a lista de interfaces de rede ativas.

#### Scenario: Acesso rápido ao serviço iniciado
- **GIVEN** que o servidor Web foi iniciado na interface gráfica
- **WHEN** o usuário clica no botão "Abrir no Navegador" ou clica na URL de acesso exibida
- **THEN** o sistema DEVE abrir a URL correspondente no navegador padrão do sistema operacional.

#### Scenario: Visualização de logs do servidor
- **WHEN** o servidor emite mensagens de log, acessos ou erros
- **THEN** o painel de logs integrado da janela DEVE atualizar instantaneamente com as novas entradas formatadas.

#### Scenario: Parar servidor
- **GIVEN** que o servidor está em execução
- **WHEN** o usuário clica no botão "Parar Servidor"
- **THEN** o sistema DEVE realizar o encerramento gracioso do serviço, limpar o status ativo e permitir nova inicialização ou alteração de configurações.

### Requirement: Ciclo de Vida e Comportamento Padrão de Janela Desktop
A janela da aplicação DEVE responder a todos os controles padrão do gerenciador de janelas do desktop (minimizar, maximizar e fechar).

#### Scenario: Maximização e restauração da janela
- **WHEN** o usuário clica no botão de maximizar da janela ou arrasta para o topo da tela
- **THEN** a janela DEVE expandir ocupando a tela de forma responsiva, mantendo a proporção e alinhamento dos elementos visuais.

#### Scenario: Fechamento da janela e encerramento da aplicação
- **WHEN** o usuário clica no botão de fechar (X) da janela ou pressiona o atalho de fechar (Alt+F4 / Ctrl+Q)
- **THEN** o sistema DEVE interromper graciosamente quaisquer servidores ativos e finalizar completamente o processo da aplicação sem deixar resíduos ou processos órfãos em background.

### Requirement: Documentação Visual e Guia de Uso no README
A documentação principal do projeto no `README.md` DEVE apresentar a interface gráfica desktop como a forma primária e recomendada de utilização, incluindo capturas de tela e passo a passo ilustrado.

#### Scenario: Estrutura do README destacando a interface gráfica
- **WHEN** o usuário lê o `README.md` do projeto
- **THEN** a seção inicial de "Como Usar / Guia Rápido" DEVE exibir prints/capturas da interface gráfica desktop, explicando o uso sem terminal, e posicionar o guia de comandos CLI, APIs e referências técnicas nas seções subsequentes.
