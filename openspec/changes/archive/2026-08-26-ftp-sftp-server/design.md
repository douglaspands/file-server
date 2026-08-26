## Context

O `file-server` foi arquitetado seguindo o padrão Clean Architecture / Ports & Adapters em Go, com modularidade em `internal/core` (serviços de domínio e sistema de arquivos) e `internal/adapters` (handlers web e repositórios). A CLI atual utiliza `spf13/cobra` com os comandos `serve` e `version`.

Para suportar transferências automatizadas e seguras na rede local via SFTP e FTP/FTPS, é necessário introduzir novos adaptadores de protocolo em `internal/adapters/` mantendo o desacoplamento das regras de negócio, a proteção estrita de sandbox e a cobertura de testes &ge; 80%.

## Goals / Non-Goals

**Goals:**
- Implementar subcomandos dedicados `file-server sftp [diretório]` e `file-server ftp [diretório]`.
- Assegurar criptografia robusta de ponta a ponta (SSHv2 no SFTP e TLS 1.2/1.3 no FTP/FTPS) evitando que *sniffing* e *scans* exponham dados ou credenciais na rede local.
- Confinamento absoluto (sandbox) no diretório raiz configurado, com tratamento idêntico ao servidor web (expansão de `~`, caminhos relativos, checagem de existência e bloqueio de *path traversal*).
- Autenticação configurável (usuário/senha, chave pública SSH) com fallback para credenciais seguras geradas automaticamente exibidas no banner de inicialização.
- Suporte a modo somente leitura (`--read-only`) para proteção contra modificações acidentais ou não autorizadas.
- Testes unitários e de integração abrangentes com meta de cobertura &ge; 80% e validação via `make test` e `make check`.

**Non-Goals:**
- Não prover acesso a shell interativo SSH (o servidor SSH deve atender estritamente ao subsistema `sftp`).
- Não prover controle de permissões complexas multi-tenant com múltiplos usuários concorrentes tendo diretórios-raízes distintos (um usuário/sessão por execução de servidor com escopo no diretório compartilhado).
- Não suportar protocolos legados e inseguros como TFTP ou Telnet.

## Decisions

### 1. Bibliotecas de Protocolo e Desacoplamento
- **SFTP / SSH**: Utilização de `golang.org/x/crypto/ssh` combinada com `github.com/pkg/sftp`.
  - *Justificativa*: São os padrões consolidados da comunidade Go para SSH e SFTP, com suporte robusto a cifras modernas (AES-GCM, ChaCha20-Poly1305, Ed25519) e hooks completos para controle de filesystem e handlers de permissão.
  - *Alternativa considerada*: Implementar SFTP do zero ou usar CGo / libssh. Descartado devido à complexidade desnecessária e quebra da portabilidade de compilação pura em Go.
- **FTP / FTPS**: Utilização de `github.com/fclairamb/ftpserverlib`.
  - *Justificativa*: Biblioteca pura em Go, altamente modular, com suporte nativo a FTPS (TLS explícito e implícito), comandos RFC padrão, faixa passiva configurável e interface limpa de driver de sistema de arquivos (`ftpserverlib.MainDriver` / `ClientDriver`).
  - *Alternativa considerada*: `goftp/server` (projeto legado e com menor flexibilidade para TLS customizado) ou implementação manual completa do protocolo RFC 959 (alta complexidade de controle de portas de dados e modo passivo).

### 2. Arquitetura de Camadas e Sandboxing
```
  [ CLI (cmd/sftp.go, cmd/ftp.go) ]
               │
               ▼
  [ Adapters (internal/adapters/sftp, internal/adapters/ftp) ]
               │
               ▼
  [ Core Domain & Services (internal/core/services/file_service.go, sandbox) ]
               │
               ▼
  [ Sistema de Arquivos Local (LocalFsRepository) ]
```
- **Driver de Sandbox Compartilhado**: Os adaptadores SFTP e FTP traduzem requisições de arquivos para o `FileService` ou `LocalFsRepository` já existente em `internal/core`, garantindo que toda resolução de caminho passe por `filepath.Clean` e validação com `strings.HasPrefix(canonicalPath, canonicalRoot)`.
- **Mapeamento de Erros**: Erros de acesso fora da raiz ou restrições de `--read-only` são convertidos para os códigos de erro padrão dos respectivos protocolos (`sftp.ErrSSHFxPermissionDenied` ou FTP `550 Permission Denied`).

### 3. Gestão de Chaves e Credenciais
- **SFTP Chave de Host**: Se `--host-key` não for informado, o servidor gera automaticamente em memória uma chave privada Ed25519 / RSA temporária usando `crypto/rand` a cada inicialização, sem requerer arquivos no disco.
- **FTPS Certificados**: Se `--tls` for ativado sem certificados customizados (`--tls-cert`, `--tls-key`), reutiliza `services.CreateSelfSignedTLSConfig()` para criar certificados X.509 em memória para os IPs e hostnames locais.
- **Credenciais de Acesso**: Se `--user` e `--pass` não forem informados pelo usuário, o sistema gera dinamicamente um usuário padrão (`fileserver`) e uma senha aleatória de 12 caracteres alfanuméricos com alta entropia, logando no banner inicial.

### 4. CLI e Experiência do Usuário
- Criação dos comandos `cmd/sftp.go` (`file-server sftp [diretório]`) e `cmd/ftp.go` (`file-server ftp [diretório]`).
- Reutilização das funções `ResolveDirectory`, `ExpandHomeDir` e `GetLANIPAddresses` para garantir paridade total com o comando `file-server serve`.
- Banners formatados no terminal indicando com clareza o protocolo seguro, a porta, o diretório servido e as URLs/comandos de conexão rápida para o cliente (ex: `sftp -P 2222 user@<IP>`).

## Risks / Trade-offs

- **[Risco: Conflito de portas padrão em ambientes de teste ou sem privilégio de root]**
  → *Mitigação*: Portas padrão altas definidas para 2222 (SFTP) e 2121 (FTP), evitando a necessidade de `sudo`/root e permitindo execução simultânea ao lado de outros daemons SSH/FTP do sistema.
- **[Risco: Vulnerabilidades de Path Traversal em clientes com caminhos relativos complexos]**
  → *Mitigação*: Implementação de wrapper de sistema de arquivos com verificação estrita de limite de pasta raiz canônica e chroot virtual para cada operação de I/O.
- **[Risco: Interceptação em redes locais inseguras]**
  → *Mitigação*: O servidor SFTP opera obrigatoriamente sob canal seguro SSHv2; o servidor FTP suporta e incentiva FTPS com TLS autoassinado imediato com aviso em log caso executado em modo plaintext.

## Test Strategy

- **Testes Unitários**:
  - Testes de flags, parsing de opções e validação de argumentos em `cmd/sftp_test.go` e `cmd/ftp_test.go`.
  - Testes de geração de chaves em memória, certificados TLS e geradores de senhas seguras em `internal/core/services/`.
  - Testes de validação de sandbox e prevenção de path traversal nos drivers de SFTP e FTP.
- **Testes de Integração**:
  - Inicialização de servidor SFTP em porta randômica (`127.0.0.1:0`), conexão via cliente `pkg/sftp` / `golang.org/x/crypto/ssh`, autenticação por senha/chave, upload de arquivo, download de arquivo, tentativa de leitura fora da raiz e encerramento limpo via contexto.
  - Inicialização de servidor FTP/FTPS em porta randômica, conexão via cliente FTP, autenticação, transferência de arquivo em modo passivo e encerramento.
- **Verificação de Cobertura**:
  - Manter cobertura global e por pacote &ge; 80% validada continuamente via `make test` e `make check`.
