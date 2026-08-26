## Why

O compartilhamento e transferência rápida de arquivos em redes locais (LAN) frequentemente esbarra em ferramentas complexas de configurar, interfaces antiquadas ou servidores pesados que comprometem a velocidade e a segurança. Há uma necessidade premente de uma solução leve, portátil, de alta performance e com interface visual moderna e intuitiva, que possa ser inicializada instantaneamente via CLI em qualquer diretório (assumindo a pasta atual por padrão ou recebendo um diretório via argumento) e permita aos usuários na rede navegar, enviar e baixar arquivos e diretórios inteiros (compactados em ZIP sob demanda) com segurança estrita (sandboxing absoluto dentro do diretório raiz).

## What Changes

- **CLI Flexível com Diretório Raiz Dinâmico e Descoberta de IPs**: Inicialização padrão apontando para o diretório atual de execução (`cwd`), com suporte completo a caminhos relativos (`./`, `../`), atalhos de diretório home (`~`, `~/...`) e caminhos absolutos, passados diretamente como argumento posicional ou através da flag `--dir` / `-d`, exibindo automaticamente no terminal os endereços IP da máquina local (LAN) para facilitar o acesso de outros dispositivos na rede.
- **Servidor Web de Alta Performance**: Implementação de pipeline de transferência otimizado para rede local com suporte a streaming direto, headers `Range` para download resumível, concorrência nativa em Go e baixo consumo de memória.
- **Interface Web Moderna, Intuitiva e Responsiva Multi-Dispositivo**: Interface visual elegante, moderna e fluida utilizando Tailwind CSS, Alpine.js e HTMX, totalmente adaptada para smartphones, tablets, notebooks e PCs/desktops (com áreas de toque confortáveis, breakpoints inteligentes, visualização adaptativa de colunas e área de upload acessível tanto via drag-and-drop quanto via seletor nativo).
- **Upload de Arquivos Simples e Robusto**: Capacidade de fazer upload de arquivos individuais ou múltiplos diretamente para o diretório atualmente navegado com área de *drag & drop* e feedback visual em tempo real.
- **Download de Arquivos e Pastas Compactadas (ZIP)**: Download direto de arquivos individuais e download sob demanda de diretórios inteiros compactados em formato `.zip` via streaming direto (zero resíduos em disco ou com limpeza estrita e determinística de temporários após a resposta).
- **Proteção Rigorosa contra Path Traversal (Sandboxing)**: Validação estrita de limites que proíbe terminantemente o acesso a qualquer caminho, link simbólico ou referência relativa (`..`) fora do diretório raiz configurado na inicialização da aplicação.
- **Camada de Criptografia em Trânsito de Alta Performance (TLS/HTTPS & Anti-Sniffing)**: Proteção robusta de todo o tráfego HTTP na rede local com suporte a TLS 1.3/1.2 moderno (aceleração por hardware via AES-NI / ChaCha20-Poly1305), multiplexação com HTTP/2, geração instantânea e automática de certificados autoassinados em memória (zero fricção na LAN) ou suporte a certificados/chaves customizados via flags CLI (`--tls`/`-s`, `--tls-cert`, `--tls-key`), garantindo confidencialidade absoluta contra sniffing e captura de pacotes na LAN sem comprometer a taxa de transferência.
- **Foco no Produto e Excelência Técnica**: Centralização total da identidade no servidor de arquivos de alta performance, posicionando práticas de engenharia (TDD, BDD, Makefile, linters) como pilares de qualidade na camada de desenvolvimento.

## Capabilities

### New Capabilities
- `web-file-server`: Servidor web de alta performance de arquivos e diretórios com navegação interativa, uploads, downloads individuais e em lote (ZIP streaming), camada de criptografia em trânsito (TLS/HTTPS) e proteção estrita de sandbox contra path traversal.

### Modified Capabilities
<!-- Nenhuma especificação existente teve seus requisitos alterados. -->

## Impact

- **Código e Arquitetura**: Extensão da camada de domínio (`internal/core/domain/`), portas (`internal/core/ports/`), serviços (`internal/core/services/`) e adaptadores HTTP (`internal/adapters/handlers/`) para incorporar listagem de arquivos, streaming de downloads, processamento de uploads multipart, validação de segurança de caminhos e geração/gerenciamento de certificados TLS em memória.
- **Interface de Linha de Comando (CLI)**: Atualização dos comandos em `cmd/root.go` e `cmd/serve.go` para aceitar argumento posicional de diretório, flags `--dir`/`-d`, `--tls`/`-s`, `--tls-cert` e `--tls-key`.
- **Frontend / Assets**: Novos templates e componentes HTML (`web/templates/`) e scripts interativos para listagem, navegação, upload com drag-and-drop e acionamento de downloads.
- **Dependências**: Utilização exclusiva de pacotes da biblioteca padrão de Go (`net/http`, `crypto/tls`, `crypto/x509`, `archive/zip`, `io`, `os`, `path/filepath`, `mime/multipart`) mantendo zero overhead externo e binário único autocontido.
- **APIs**: Novas rotas HTTP/HTTPS para navegação (`/files/*`), download direto (`/download/*`), download de pasta em zip (`/zip/*`), upload multipart (`/upload/*`) e API JSON correspondente.
