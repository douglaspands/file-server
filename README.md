# ⚡ File Server

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go&logoColor=white)](https://golang.org)
[![CI Quality Gate](https://img.shields.io/badge/CI-Passing-success?style=flat&logo=githubactions&logoColor=white)](.github/workflows/ci.yml)
[![Test Coverage](https://img.shields.io/badge/Coverage-%E2%89%A5%2080%25-brightgreen?style=flat)](scripts/coverage.sh)
[![Architecture](https://img.shields.io/badge/Architecture-Clean%20%2F%20Ports%20%26%20Adapters-informational?style=flat)](#-arquitetura-e-estrutura-de-pacotes)
[![Platforms](https://img.shields.io/badge/Platforms-Linux%20%7C%20Windows-blueviolet?style=flat)](#-compilação-multiplataforma)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat)](LICENSE)

> Servidor web, SFTP e FTP/FTPS de alta performance para rede local (LAN) com interface moderna, streaming direto com suporte a HTTP Range (206), downloads de pastas compactadas em ZIP via streaming sob demanda, uploads multipart/drag-and-drop, isolamento estrito de sandbox contra path traversal e criptografia em trânsito (SSHv2 / TLS / HTTPS) para evitar que scans e sniffs interceptem o tráfego. Empacotado em um **único binário executável 100% autocontido e portátil** (`go:embed`).

---

## 📖 Visão Geral

O **File Server** é uma solução leve, portátil e de alta performance desenvolvida em Go para compartilhamento instantâneo de arquivos e pastas no terminal ou rede local via Web (HTTP/HTTPS), SFTP (SSH) e FTP/FTPS.

### ✨ Principais Recursos do Produto

- 🚀 **Inicialização Instantânea**: Execute `file-server` em qualquer pasta e comece a compartilhar imediatamente. Fallback automático para o diretório atual (`.`) ou configuração de raiz via argumento posicional ou flag `--dir`/`-d`.
- 🔒 **Criptografia em Trânsito (TLS / SSHv2 / FTPS)**: Suporte nativo a conexões 100% criptografadas para navegação web (HTTPS), SFTP (SSHv2 com Ed25519/RSA e ChaCha20/AES) e FTPS (TLS 1.3/1.2), evitando que *packet sniffing* e *scans* na rede local interceptem credenciais ou dados.
- 🛡️ **Sandboxing Rigoroso & Segurança**: Validação de fronteira canônica e bloqueio absoluto contra ataques de *path traversal* (`../`, links simbólicos externos e caminhos absolutos forçados) em todos os protocolos (Web, SFTP e FTP).
- 🌐 **Interface Web Moderna & Intuitiva**: Interface visual desenvolvida com **Tailwind CSS**, **Alpine.js** e **HTMX**, com ícones dinâmicos por categoria de arquivo e navegação por *breadcrumbs*.
- 🔑 **SFTP e FTP/FTPS Nativos**: Transfira arquivos utilizando clientes de terminal (`sftp`, `lftp`, scripts de automação, rsync/scp) ou clientes gráficos (FileZilla, Cyberduck, WinSCP) com suporte a autenticação por senha ou chave pública SSH.
- ⚡ **Alta Performance em LAN**: Download direto via streaming com `http.ServeContent`, suporte completo a requisições parciais (header `Range: bytes=...`, HTTP 206) e modo somente leitura (`--read-only`).
- 📦 **Download de Pastas em ZIP sob Demanda**: Compactação em streaming direto para a resposta HTTP (`io.Writer`) com cancelamento via contexto e **zero resíduos ou arquivos temporários deixados em disco**.
- 📤 **Upload Simples e em Lote (Drag & Drop)**: Envio de arquivos únicos ou múltiplos diretamente para o diretório visualizado.
- 📦 **Binário Único Autocontido**: Todos os templates HTML, estilos CSS e scripts JS são compilados dentro do executável via `go:embed`. Transporte apenas um arquivo para qualquer servidor Linux ou Windows.

---

## 🚀 Guia de Uso da CLI

### 1. Obtenção e Compilação

Você pode compilar localmente com um único comando:

```bash
# Compilar o binário para seu sistema operacional
make build

# O binário será gerado em bin/file-server
./bin/file-server --help
```

---

### 2. Catálogo Completo de Comandos da CLI

| Comando / Flag | Descrição | Exemplo de Uso |
| :--- | :--- | :--- |
| `file-server [dir]` | Inicia o servidor Web na pasta informada (ou na pasta atual caso omitida) | `./bin/file-server` ou `./bin/file-server /meus/arquivos` |
| `file-server serve [dir]` | Subcomando explícito para iniciar o servidor web HTTP/HTTPS | `./bin/file-server serve ./dados -p 8080` |
| `file-server sftp [dir]` | Inicia o servidor SFTP seguro sobre SSH com chaves e senha | `./bin/file-server sftp ./dados -p 2222` |
| `file-server ftp [dir]` | Inicia o servidor FTP com suporte a FTPS (TLS) e modo passivo | `./bin/file-server ftp ./dados -p 2121 --tls` |
| `--dir` (`-d`) | Define o caminho do diretório raiz a ser compartilhado | `./bin/file-server -d /var/public` |
| `--port` (`-p`) | Define a porta TCP de escuta (Web: `8080`, SFTP: `2222`, FTP: `2121`) | `./bin/file-server sftp -p 2222` |
| `--host` | Define o endereço IP/host de escuta (default: `0.0.0.0`) | `./bin/file-server --host 127.0.0.1` |
| `--user` (`-u`) | Define o usuário para autenticação SFTP ou FTP (default: `fileserver`) | `./bin/file-server sftp -u admin -P segredo` |
| `--pass` (`-P`) | Define a senha de autenticação (gerada automaticamente caso omitida) | `./bin/file-server ftp -P MinhaSenha123` |
| `--auth-key` (`-k`) | Caminho da chave pública SSH autorizada para login SFTP | `./bin/file-server sftp -k ~/.ssh/id_ed25519.pub` |
| `--host-key` | Caminho da chave privada de host SSH (gerada em memória caso omitida) | `./bin/file-server sftp --host-key host.pem` |
| `--tls` (`-s`) | Ativa criptografia segura TLS/HTTPS ou FTPS (certificado autoassinado automático) | `./bin/file-server -s` ou `./bin/file-server ftp --tls` |
| `--tls-cert` | Caminho do arquivo PEM com certificado público TLS customizado | `./bin/file-server serve --tls-cert cert.pem --tls-key key.pem` |
| `--tls-key` | Caminho do arquivo PEM com chave privada TLS customizada | `./bin/file-server ftp --tls-cert cert.pem --tls-key key.pem` |
| `--passive-ports` | Faixa de portas para o modo passivo FTP (ex: `50000-50100`) | `./bin/file-server ftp --passive-ports 50000-50100` |
| `--read-only` (`-r`) | Ativa modo somente leitura (bloqueia uploads, edições e remoções) | `./bin/file-server sftp -r` ou `./bin/file-server ftp -r` |
| `file-server version` | Exibe versão semântica, commit, data de build e plataforma | `./bin/file-server version` |
| `file-server --help` (`-h`) | Exibe a ajuda interativa completa da CLI | `./bin/file-server --help` |

---

### 3. Exemplos Práticos de Inicialização

#### Servidor SFTP Seguro (Recomendado para Redes Locais)
```bash
# Inicia servidor SFTP na pasta ~/Documentos na porta 2222 com senha temporária segura
./bin/file-server sftp ~/Documentos -p 2222

# Ou configurando usuário e senha explícitos:
./bin/file-server sftp /dados -u admin -P SegredoForte123

# Conectar via terminal:
sftp -P 2222 admin@<ip-do-servidor>

# Ou via FileZilla / Cyberduck:
# Protocolo: SFTP (SSH File Transfer Protocol)
# Host: <ip-do-servidor>, Porta: 2222
```

#### Servidor FTP com Criptografia FTPS / TLS
```bash
# Inicia servidor FTP com FTPS ativado na porta 2121
./bin/file-server ftp /dados -p 2121 --tls -u transfer -P Senha123

# Conectar via FileZilla com "Usar FTP explícito sobre TLS se disponível"
```

#### Compartilhar a Pasta Atual Instantaneamente via Web (HTTP)
```bash
# Inicia na porta 8080 servindo o diretório corrente
./bin/file-server

# Acesse no navegador:
# http://localhost:8080 ou http://<ip-da-maquina>:8080
```

#### Iniciar Servidor Web Seguro com TLS/HTTPS Autoassinado
```bash
# Inicia com HTTPS na porta 8443 e certificado gerado em memória
./bin/file-server -s -p 8443

# Acesse no navegador:
# https://localhost:8443 ou https://<ip-da-maquina>:8443
```

---

## 🌐 Guia de Uso da Interface Web

A interface gráfica do **File Server** foi projetada para oferecer uma experiência ágil tanto no desktop quanto em dispositivos móveis:

1. **Navegação Hierárquica por Breadcrumbs**:
   - Clique em qualquer pasta na lista para entrar nela.
   - Utilize a barra de navegação (*breadcrumbs*) no topo para retornar instantaneamente a qualquer nível superior ou para a raiz do compartilhamento.

2. **Filtro de Busca em Tempo Real**:
   - Digite na caixa de busca (*"Filtrar arquivos nesta pasta..."*) para filtrar instantaneamente os itens visíveis no navegador via Alpine.js sem recarregar a página ou fazer requisições extras ao servidor.

3. **Download de Arquivos e Pastas**:
   - **Download Individual**: Clique no nome do arquivo ou no ícone de download para baixar diretamente com suporte a pausas e retomadas (*HTTP Range requests*).
   - **Download de Pasta em ZIP**: Clique no botão *"Baixar Pasta (.zip)"* no topo ou no ícone correspondente na linha de qualquer pasta para baixar a árvore inteira compactada em `.zip` em streaming contínuo.

4. **Upload Simples e Múltiplo (Drag & Drop)**:
   - **Botão Enviar**: Clique em *"Enviar Arquivos"*, selecione um ou mais arquivos e confirme o envio.
   - **Arrastar e Soltar**: Arraste arquivos de qualquer lugar do seu computador e solte sobre a janela do navegador para iniciar o envio diretamente para a pasta atualmente aberta.

---

## 📡 Catálogo de Rotas e Endpoints da API Web

| Método | Rota | Descrição |
| :--- | :--- | :--- |
| `GET` | `/` ou `/files/*` | Interface gráfica do explorador de arquivos (HTML/SSR) |
| `GET` | `/api/files/*` | Metadados e listagem do diretório em JSON |
| `GET` | `/download/*` | Download de arquivo individual com suporte a HTTP Range (206) |
| `GET` | `/zip/*` | Download de diretório compactado em `.zip` via streaming |
| `POST`| `/upload/*` | Upload multipart de arquivos para o diretório de destino |
| `GET` | `/status` | Painel de diagnóstico e integridade dos serviços do servidor |
| `GET` | `/partials/health` | Fragmento HTML do card de saúde para requisições HTMX |
| `GET` | `/api/health` | Status de saúde e métricas de uptime em JSON |
| `GET` | `/api/version` | Metadados estruturados de versão e build em JSON |

---

## 🛠️ Guia para Desenvolvedores e Padrões de Qualidade

Esta seção descreve os padrões de engenharia de software empregados no desenvolvimento do **File Server**.

### 📋 Pré-requisitos

- **Go**: Versão `1.25` ou superior instalada ([golang.org](https://golang.org)).
- **Make**: Utilitário Make (`/usr/bin/make`) para automação de tarefas.
- **Git**: Controle de versão.

---

### 🏗️ Arquitetura e Estrutura de Pacotes

O projeto adota rigorosamente os princípios de **Clean Architecture** e **Ports & Adapters**, isolando todo o domínio de negócio em `internal/`:

```
file-server/
├── cmd/                          # Camada de entrada da CLI (Cobra)
│   ├── root.go                   # Comando base e inicialização padrão
│   ├── version.go                # Subcomando 'version'
│   ├── serve.go                  # Subcomando 'serve' (Web HTTP/HTTPS)
│   ├── sftp.go                   # Subcomando 'sftp' (SFTP sobre SSH)
│   └── ftp.go                    # Subcomando 'ftp' (FTP/FTPS com TLS)
├── internal/                     # Domínio isolado e privado (internal/)
│   ├── core/
│   │   ├── domain/               # Entidades de domínio (FileItem, DirectoryListing, Breadcrumb)
│   │   ├── ports/                # Interfaces e contratos (FileService, HealthService)
│   │   └── services/             # Serviços centrais (LocalFileService, TLS, Auth & SSH Host Keys)
│   ├── adapters/
│   │   ├── handlers/             # Controladores HTTP, templates e downloads
│   │   ├── sftp/                 # Adaptador do servidor SFTP (SSH listener e FSHandler)
│   │   └── ftp/                  # Adaptador do servidor FTP/FTPS (Driver e Settings)
│   ├── testutils/                # Fixtures e mocks para testes automatizados
│   └── version/                  # Metadados de compilação injetados via ldflags
├── web/                          # Recursos de frontend
│   ├── templates/                # Templates Go HTML (layouts, pages, partials)
│   ├── static/                   # CSS, scripts JS e assets estáticos
│   └── web.go                    # Empacotamento embutido (go:embed embed.FS)
├── scripts/                      # Harness de automação (scripts de cobertura e setup)
├── openspec/                     # Especificações normativas do projeto
└── Makefile                      # Interface universal de comandos
```

---

### 🕹️ Harness de Automação (Makefile)

```bash
# Exibir o menu interativo com todos os comandos disponíveis
make help
```

| Alvo | Finalidade |
| :--- | :--- |
| `make help` | Exibe o menu com todos os comandos disponíveis e descrições |
| `make setup` | Instala e checa ferramentas locais (`golangci-lint`, `air`, `govulncheck`) |
| `make dev` | Inicia o servidor em modo de desenvolvimento com live-reload (Air) |
| `make run` | Executa a aplicação diretamente do código fonte |
| `make fmt` | Formata todo o código Go e templates |
| `make lint` | Executa análise estática com `golangci-lint` (ou `go vet`) |
| `make vulncheck` | Executa auditoria de vulnerabilidades de segurança com `govulncheck` |
| `make test` | Executa toda a suíte de testes e valida a barreira de cobertura (&ge; 80%) |
| `make test-unit` | Executa testes unitários rápidos |
| `make check` | Executa o pipeline completo de qualidade local (`fmt` + `lint` + `test`) |
| `make build` | Compila o binário de produção otimizado com assets embutidos |
| `make build-all` | Realiza compilação cruzada para Linux e Windows (`amd64` e `arm64`) |
| `make clean` | Limpa diretórios de build (`bin/`, `dist/`), relatórios e temporários |

---

### 🧪 Práticas de Testes e Qualidade (TDD & BDD)

O desenvolvimento segue **TDD (Test-Driven Development)** e especificações em estilo **BDD (Behavior-Driven Development)** estruturadas em subtestes com `testing` + `testify`:

```go
func TestLocalFileService_GetFile(t *testing.T) {
    t.Run("sucesso ao obter arquivo existente", func(t *testing.T) {
        file, info, err := svc.GetFile(ctx, "sample.txt")
        require.NoError(t, err)
        defer file.Close()
        assert.Equal(t, "sample.txt", info.Name())
    })
}
```

- **Barreira Mínima de Cobertura**: Meta contínua de cobertura de código &ge; 80%.
- **Execução e Relatório**:
  ```bash
  make test
  # Relatório HTML gerado em: coverage.html
  ```

---

### 🌍 Compilação Multiplataforma (Linux & Windows)

Para gerar executáveis portáteis para Linux e Windows (arquiteturas `amd64` e `arm64`):

```bash
make build-all
```

Os binários serão gerados no diretório `dist/`:
- `dist/file-server-linux-amd64`
- `dist/file-server-linux-arm64`
- `dist/file-server-windows-amd64.exe`
- `dist/file-server-windows-arm64.exe`

---

## 📄 Licença

Este projeto está sob a licença [MIT](LICENSE).
