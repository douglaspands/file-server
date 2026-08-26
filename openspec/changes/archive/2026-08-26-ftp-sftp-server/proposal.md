## Why

Atualmente, o `file-server` disponibiliza uma interface web HTTP/HTTPS e API para transferência de arquivos. No entanto, para fluxos automatizados de transferência, pipelines, scripts de sincronização, integrações de terminal e clientes dedicados de transferência (como FileZilla, Cyberduck, scripts em shell, rsync/scp/sftp), há uma forte demanda pelo suporte nativo a protocolos padronizados de transferência de arquivos em rede: **SFTP** (SSH File Transfer Protocol) e **FTP/FTPS** (File Transfer Protocol com suporte a TLS).

Em ambientes de rede local (LAN), a transferência em texto plano expõe credenciais e dados a *packet sniffing* e varreduras de rede (*scans*). Esta proposta introduz os modos de servidor FTP e SFTP no `file-server`, garantindo tráfego totalmente criptografado e seguro contra interceptações, mantendo o isolamento estrito de sandbox contra *path traversal* e a mesma facilidade de configuração via CLI (parâmetros de diretório raiz, porta, bind address, certificados TLS e chaves SSH).

## What Changes

- **Novos Comandos CLI (`file-server ftp` e `file-server sftp`)**:
  - Comando `file-server ftp [diretório]`: Inicia o servidor FTP com suporte nativo a criptografia TLS (FTPS explícito/implícito), suporte a modo passivo configurável, autenticação de usuário e controle de diretório raiz com proteção de sandbox.
  - Comando `file-server sftp [diretório]`: Inicia o servidor SFTP operando sobre o subsistema SSH com criptografia de ponta a ponta (cifras modernas), autenticação por usuário/senha ou chave pública SSH, geração automática de chaves de host em memória ou uso de chaves existentes, e isolamento de diretório.
- **Resolução e Sandbox de Diretório Compartilhado**:
  - Ambos os servidores herdam as mesmas regras do servidor web: aceitam argumento posicional de diretório ou flag `-d`/`--dir`, resolução de caminhos relativos, expansão de til (`~`) para o diretório home do usuário e validação estrita de sandbox impedindo *path traversal* (`../`, symlinks externos).
- **Proteção Criptográfica contra Sniffing e Scans**:
  - SFTP: Criptografia de canal SSH obrigatória com algoritmos seguros (ChaCha20-Poly1305, AES-GCM, Ed25519/RSA).
  - FTP/FTPS: Suporte a TLS 1.3 / 1.2 (autoassinado dinâmico ou certificados customizados) para proteger controle e dados contra interceptação.
- **Autenticação Flexível e Segura**:
  - Flags de CLI e variáveis de ambiente para definir credenciais de acesso (`--user`, `--pass`, `--auth-key`, `--read-only`).
  - Geração segura de credenciais temporárias ou chaves quando não informadas explicitamente, exibidas no banner de inicialização do terminal.
- **Banners Informativos no Terminal**:
  - Exibição na inicialização com endereços IP locais e da rede LAN, porta, protocolo seguro ativo e credenciais de acesso configuradas.

## Capabilities

### New Capabilities
- `ftp-sftp-server`: Servidores FTP/FTPS e SFTP seguros com inicialização via CLI, controle de autenticação, criptografia em trânsito e isolamento estrito de sandbox.

### Modified Capabilities
<!-- Nenhuma especificação existente de web-file-server ou project-foundation tem seus requisitos alterados. -->

## Impact

- **CLI (`cmd/`)**: Adição de `cmd/ftp.go`, `cmd/sftp.go` e testes correspondentes em `cmd/`.
- **Core & Adapters (`internal/`)**:
  - Implementação dos adaptadores de servidor FTP (`internal/adapters/ftp/`) e SFTP (`internal/adapters/sftp/`).
  - Integração com `internal/core/services/file_service.go` ou driver de sistema de arquivos seguro que reutiliza as regras de sandbox.
  - Reutilização dos utilitários de rede (`services/network.go`) e TLS (`services/tls.go`).
- **Dependências (`go.mod`)**: Adição de bibliotecas Go consolidadas e seguras para SFTP/SSH (`golang.org/x/crypto/ssh`, `github.com/pkg/sftp`) e FTP/FTPS (`github.com/fclairamb/ftpserverlib` ou similar em Go padrão/puro).
- **Testes & Cobertura**: Testes unitários e de integração com cobertura &ge; 80% e validação via `make test` e `make check`.
