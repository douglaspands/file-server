## Purpose

Prover servidores dedicados de transferência segura de arquivos nos protocolos FTP (com suporte a FTPS/TLS) e SFTP (sobre SSH), com gerenciamento via CLI, isolamento estrito de sandbox contra path traversal, autenticação flexível e proteção criptográfica contra packet sniffing e scans na rede local.

## ADDED Requirements

### Requirement: Inicialização e Controle do Servidor SFTP via CLI
O sistema SHALL fornecer o subcomando `sftp` na CLI (`file-server sftp [diretório]`) permitindo iniciar o servidor SFTP seguro sobre SSH, suportando configuração de diretório raiz (posicional ou `--dir`/`-d`), porta (`--port`/`-p`, padrão 2222), endereço de bind (`--host`, padrão `0.0.0.0`), autenticação por usuário/senha (`--user`, `--pass`) e chave pública SSH (`--auth-key`), chave privada do host SSH (`--host-key`) ou geração automática de chave de host em memória, e modo somente leitura (`--read-only`).
*(Visão PO: Permite que ferramentas de automação, scripts e clientes SFTP transfiram arquivos com segurança total na LAN com a mesma facilidade de uso do servidor web. Visão QA: Valida parsing de flags, resolução e fallback de diretório, geração de chave host em memória, encerramento gracioso e rejeição de acessos não autorizados).*

#### Scenario: Inicialização padrão do SFTP no diretório atual
- **GIVEN** que o usuário está em um diretório no terminal
- **WHEN** executar o comando `file-server sftp` sem argumentos de diretório
- **THEN** o sistema deve iniciar o servidor SFTP adotando o diretório de trabalho atual (`.`) como raiz do compartilhamento na porta padrão 2222

#### Scenario: Inicialização do SFTP com diretório customizado, caminhos relativos ou expansão de til (~)
- **GIVEN** que o usuário informa `~/Documentos` ou `../dados` como argumento posicional ou na flag `-d`/`--dir`
- **WHEN** o comando `file-server sftp` for executado
- **THEN** o sistema deve expandir e resolver o caminho canônico absoluto correspondente e validar a existência do diretório antes de iniciar o serviço

#### Scenario: Inicialização com credenciais customizadas de usuário e senha
- **GIVEN** que o usuário define `--user admin --pass Segredo123`
- **WHEN** o servidor SFTP for inicializado
- **THEN** o sistema deve aceitar conexões SFTP autenticadas exclusivamente com o usuário e senha configurados

#### Scenario: Inicialização com autenticação por chave pública SSH
- **GIVEN** que o usuário informa o caminho de uma chave pública autorizada via `--auth-key ~/.ssh/id_ed25519.pub`
- **WHEN** um cliente SFTP tentar conectar apresentando a chave privada correspondente
- **THEN** o servidor SFTP deve autenticar o cliente com sucesso com base na chave pública

#### Scenario: Geração automática e transparente de chave de host SSH
- **GIVEN** que nenhuma chave de host for informada via `--host-key`
- **WHEN** o servidor SFTP for iniciado
- **THEN** o sistema deve gerar dinamicamente em memória uma chave privada de host SSH (RSA 2048 ou Ed25519) para o handshake seguro

#### Scenario: Inicialização com chave privada de host SSH persistida
- **GIVEN** que o usuário possui um arquivo de chave privada de host `host_key.pem`
- **WHEN** executar `file-server sftp --host-key host_key.pem`
- **THEN** o servidor deve utilizar essa chave de host para identificar o servidor perante os clientes SSH/SFTP

#### Scenario: Modo somente leitura no SFTP
- **GIVEN** que a flag `--read-only` está ativada
- **WHEN** um cliente SFTP tentar enviar arquivos (upload), criar diretórios (`mkdir`), renomear ou remover arquivos (`rm`)
- **THEN** o servidor SFTP deve rejeitar a operação com código de permissão negada (`SSH_FX_PERMISSION_DENIED`), permitindo apenas leitura e listagem

#### Scenario: Exibição de banner de inicialização do SFTP
- **GIVEN** que o servidor SFTP foi inicializado
- **WHEN** o banner de inicialização for exibido no terminal
- **THEN** o sistema deve imprimir a versão da aplicação, diretório raiz compartilhado, porta, protocolo SFTP (SSH), credenciais ativas e URLs/endereços de acesso local e de rede (LAN)

### Requirement: Inicialização e Controle do Servidor FTP / FTPS via CLI
O sistema SHALL fornecer o subcomando `ftp` na CLI (`file-server ftp [diretório]`) permitindo iniciar o servidor FTP com suporte a conexões criptografadas (FTPS explícito via `AUTH TLS` e implícito), porta configurável (`--port`/`-p`, padrão 2121), faixa de portas de modo passivo (`--passive-ports`), autenticação (`--user`, `--pass`), certificados TLS autoassinados ou customizados (`--tls`/`-s`, `--tls-cert`, `--tls-key`), e modo somente leitura (`--read-only`).
*(Visão PO: Garante compatibilidade com clientes clássicos de FTP oferecendo camada robusta de criptografia TLS para evitar sniffing de senhas e dados. Visão QA: Valida handshake FTPS, comandos FTP padrão (USER, PASS, LIST, RETR, STOR, DELE, MKD, RMD, PWD, CWD, EPSV, PASV), modo passivo e bloqueio de tráfego inseguro quando TLS for exigido).*

#### Scenario: Inicialização padrão do FTP no diretório atual
- **GIVEN** que o usuário está no terminal
- **WHEN** executar `file-server ftp` sem argumentos
- **THEN** o servidor FTP deve iniciar escutando na porta 2121 servindo o diretório de trabalho atual

#### Scenario: Inicialização do FTP com TLS autoassinado e proteção criptográfica
- **GIVEN** que a flag `--tls` ou `-s` está ativada na inicialização do FTP
- **WHEN** o servidor for iniciado
- **THEN** o sistema deve gerar automaticamente um certificado TLS em memória e exigir negociação TLS (`AUTH TLS`) antes de permitir a transmissão de credenciais ou dados

#### Scenario: Inicialização do FTP com certificados TLS customizados
- **GIVEN** que os arquivos `cert.pem` e `key.pem` são informados via `--tls-cert` e `--tls-key`
- **WHEN** o servidor FTP for iniciado
- **THEN** o servidor deve utilizar os certificados fornecidos para estabelecer conexões seguras FTPS

#### Scenario: Configuração de credenciais de usuário para o FTP
- **GIVEN** que o usuário define `--user transfer --pass MinhaSenha123`
- **WHEN** um cliente FTP enviar comandos `USER transfer` e `PASS MinhaSenha123`
- **THEN** o servidor FTP deve autenticar a sessão com sucesso

#### Scenario: Modo somente leitura no FTP
- **GIVEN** que a flag `--read-only` está ativada
- **WHEN** um cliente FTP tentar enviar comandos de gravação ou exclusão (`STOR`, `DELE`, `RMD`, `MKD`, `RNFR/RNTO`)
- **THEN** o servidor deve responder com código de erro 550 (Permissão negada)

#### Scenario: Exibição de banner de inicialização do FTP
- **GIVEN** que o servidor FTP foi inicializado
- **WHEN** o banner de inicialização for exibido no terminal
- **THEN** o sistema deve imprimir a versão da aplicação, diretório raiz compartilhado, porta, indicação de TLS/FTPS, credenciais ativas e endereços de acesso local e de rede (LAN)

### Requirement: Criptografia de Ponta a Ponta e Proteção contra Packet Sniffing na Rede Local
O sistema SHALL assegurar que todas as transferências realizadas via SFTP e FTPS utilizem protocolos criptográficos consolidados (SSHv2 e TLS 1.2/1.3 com cifras seguras como ChaCha20-Poly1305, AES-GCM e RSA/Ed25519), impedindo que sniffing de pacotes e ferramentas de varredura (scans) na rede local capturem credenciais, comandos ou arquivos em texto plano.
*(Visão PO: Assegura sigilo e integridade corporativa dos arquivos que trafegam na rede local. Visão QA: Valida que payloads de comandos e dados trafegam cifrados e que conexões sem negociação de segurança válida sejam rejeitadas).*

#### Scenario: Criptografia integral do canal de controle e dados no SFTP
- **GIVEN** que um cliente conecta ao servidor SFTP
- **WHEN** qualquer comando, autenticação ou transferência de blocos de arquivo for executada
- **THEN** todo o tráfego deve ser cifrado através do túnel SSHv2 negociado, sem exposição de credenciais ou dados em texto claro na rede

#### Scenario: Criptografia do canal de controle e dados no FTPS
- **GIVEN** que o servidor FTP está operando com TLS ativado
- **WHEN** o cliente efetuar login e transferir arquivos via canal de dados (modo passivo/ativo protegido por `PROT P`)
- **THEN** os fluxos de controle e de dados devem ser integralmente cifrados via TLS

#### Scenario: Rejeição de algoritmos ou cifras inseguras
- **GIVEN** que um cliente tenta negociar conexões com cifras obsoletas (como DES, RC4 ou SSHv1)
- **WHEN** o handshake de conexão for iniciado
- **THEN** o servidor deve recusar a negociação e abortar a conexão imediatamente

### Requirement: Isolamento Estrito de Sandbox e Prevenção de Path Traversal no FTP e SFTP
O sistema SHALL garantir que todas as operações executadas via sessões FTP e SFTP (listagem de diretórios, navegação com `cd`/`cwd`, leitura de arquivos, download, upload, exclusão e criação de pastas) fiquem estritamente restritas ao diretório raiz configurado na inicialização da aplicação, tratando a raiz como `/` virtual e bloqueando qualquer tentativa de escape para fora dos limites autorizados (Sandboxing absoluto).
*(Visão PO: Garante a proteção e o confinamento total do sistema operacional hospedeiro contra acessos não autorizados a arquivos confidenciais do servidor. Visão QA: Testa injeções de caminhos relativos como `../../`, sequências de escape, manipulação de symlinks externos e tentativa de escrita em locais proibidos).*

#### Scenario: Bloqueio de navegação acima da raiz configurada (chroot virtual)
- **GIVEN** que a raiz configurada é `/home/usuario/compartilhado`
- **WHEN** o cliente SFTP ou FTP enviar comando para navegar para `..` ou `/etc` a partir da raiz
- **THEN** o servidor deve confinar a navegação à raiz virtual `/`, não permitindo atingir pastas superiores no host

#### Scenario: Bloqueio de leitura e download de arquivos fora da raiz
- **WHEN** um cliente tentar requisitar leitura de arquivo com caminhos relativos manipulados (ex: `../../../../etc/shadow`)
- **THEN** o sistema deve sanitizar e validar o caminho canônico, recusando o acesso com erro de arquivo não encontrado ou permissão negada

#### Scenario: Bloqueio de gravação e upload fora da raiz
- **WHEN** um cliente tentar fazer upload com nome de arquivo contendo sequências de traversal (ex: `../../malware.bin`)
- **THEN** o sistema deve restringir a gravação estritamente dentro da pasta de destino sob a raiz sanitizada

#### Scenario: Bloqueio de links simbólicos externos
- **GIVEN** que existe um symlink dentro do diretório raiz apontando para um local fora da raiz
- **WHEN** o cliente tentar ler, navegar ou baixar o destino desse symlink
- **THEN** o sistema deve resolver o destino canônico e bloquear o acesso se o alvo estiver fora da raiz

### Requirement: Autenticação Segura e Gerenciamento de Sessões
O sistema SHALL validar credenciais fornecidas por clientes FTP e SFTP contra a configuração definida, gerando credenciais seguras e aleatórias quando nenhuma for definida explicitamente pelo usuário na inicialização, registrando tentativas de acesso e encerrando sessões inativas ou expiradas.
*(Visão PO: Evita que servidores iniciados sem configuração prévia fiquem abertos sem controle na rede. Visão QA: Valida rejeição de senhas incorretas, suporte a senhas aleatórias e isolamento entre sessões de clientes simultâneos).*

#### Scenario: Autenticação bem-sucedida com credenciais configuradas
- **GIVEN** que o usuário e a senha foram especificados via linha de comando
- **WHEN** um cliente FTP ou SFTP submeter as credenciais correspondentes
- **THEN** o servidor deve autenticar a conexão e conceder acesso ao sistema de arquivos servido

#### Scenario: Geração automática de senha temporária segura quando não fornecida
- **GIVEN** que o usuário não informou `--user` ou `--pass` nem `--auth-key`
- **WHEN** o servidor for inicializado
- **THEN** o sistema deve gerar automaticamente um usuário padrão (ex: `fileserver`) e uma senha temporária aleatória segura, exibindo-os com destaque no banner do terminal

#### Scenario: Rejeição de autenticação com credenciais inválidas
- **WHEN** um cliente tentar autenticar com credenciais incorretas
- **THEN** o servidor deve recusar a conexão e registrar a tentativa de autenticação sem vazar informações internas

#### Scenario: Isolamento e suporte a múltiplos clientes simultâneos
- **GIVEN** múltiplos clientes conectados simultaneamente ao servidor FTP ou SFTP
- **WHEN** os clientes executarem operações concorrentes de leitura e escrita
- **THEN** cada sessão de cliente deve ser executada de forma isolada e segura sem interferência ou corrupção de estado
