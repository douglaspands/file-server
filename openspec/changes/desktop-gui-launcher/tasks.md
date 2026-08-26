## 1. Estrutura do Controlador GUI e Gerenciamento de Estado

- [x] 1.1 Criar pacote `internal/adapters/gui` com tipos de modelo, opções de configuração e gerenciador de estado dos servidores.
- [x] 1.2 Implementar sistema de broadcast de logs em memória (`LogBroadcaster`) com buffer circular para streaming de eventos em tempo real.
- [x] 1.3 Implementar serviço de seleção de diretórios via diálogo nativo do sistema operacional (*folder picker*).
- [x] 1.4 Implementar detecção categorizada de adaptadores de rede (Wi-Fi, Ethernet, VPN, Docker, Bridge) e endpoints de listagem/formatação de mensagens de compartilhamento.

## 2. Interface Desktop GNOME Adwaita Dark (Frontend)

- [x] 2.1 Criar template `web/templates/pages/gui_launcher.html` com layout inspirado em GNOME Adwaita Dark e paleta de cores unificada (`slate-900`, `slate-950`, `indigo-500`, `slate-100`).
- [x] 2.2 Implementar componentes reativos em Alpine.js/HTMX para seleção de protocolos (Web/HTTP, FTP, SFTP), preenchimento automático de portas e validações de formulário.
- [x] 2.3 Implementar viewport rolável estilizada (`overflow-y-auto`) para múltiplos IPs com busca rápida, tags de tipo de adaptador (`Wi-Fi`, `Ethernet`, `VPN`, `Docker`), botão "Copiar Todos" e botões de cópia individual em 1 clique ("✓ Copiado!").
- [x] 2.4 Implementar modal de visualização de QR Code e terminal de logs com auto-scroll integrado.

## 3. Janela Desktop Nativa e Ciclo de Vida do Aplicativo

- [x] 3.1 Implementar gerenciador de janela desktop nativa/webview configurado com redimensionamento, suporte a maximização e dimensões mínimas.
- [x] 3.2 Implementar hook de encerramento seguro capturando o fechamento da janela para realizar o shutdown gracioso de servidores e finalizar o processo.

## 4. Integração no CLI e Inicialização Padrão

- [x] 4.1 Criar comando `file-server gui` e flag `--gui` no Cobra (`cmd/gui.go`).
- [x] 4.2 Atualizar `cmd/root.go` para inicializar a GUI automaticamente ao executar o binário sem argumentos ou por clique duplo em ambiente gráfico.

## 5. Testes, Qualidade e Documentação

- [x] 5.1 Criar testes unitários para o controlador GUI, detecção/categorização de interfaces de rede e broadcasters em `internal/adapters/gui/gui_test.go`.
- [x] 5.2 Criar testes unitários para os comandos CLI em `cmd/gui_test.go` e integração em `cmd/cmd_test.go`.
- [x] 5.3 Executar validação de cobertura inegociável $\ge$ 80% e suite completa com `make check` e `make test`.
- [x] 5.4 Salvar prints visuais da interface gráfica em `docs/assets/` e reformular o `README.md` destacando a interface desktop como o guia principal de uso, mantendo a documentação completa de CLI e desenvolvimento organizada logo abaixo.

## 6. Identidade Visual e Ícone da Aplicação no Executável e Barra de Tarefas

- [x] 6.1 Criar master vetorial (SVG) e pacote de ícones em alta resolução (.ico, .png 16/32/48/128/256/512) em `web/static/assets/` e `docs/assets/`.
- [x] 6.2 Vincular os ícones nos cabeçalhos dos templates HTML (`gui_launcher.html` e `base.html`) para renderização na barra de tarefas e janela.
- [x] 6.3 Configurar recurso de ícone embutido (.syso) para o binário executável do Windows no `Makefile` e `cmd/`.
- [x] 6.4 Validar compilação cruzada (`make check` e `make build-all`) e atualizar documentação visual no `README.md`.
