## Context

O projeto adota Clean Architecture (Ports & Adapters) em Go, onde o domínio e a lógica de negócio residem exclusivamente em `internal/core/`, e os adaptadores de entrada (CLI via Cobra e HTTP via net/http) residem em `cmd/` e `internal/adapters/handlers/`. A interface web utiliza renderização no servidor com templates HTML, Tailwind CSS e Alpine.js embutidos no binário compilado (`embed.FS`).

Veja `proposal.md` para a motivação detalhada da mudança e `specs/web-file-server/spec.md` para os requisitos de comportamento.

## Goals / Non-Goals

**Goals:**
- Implementar um serviço desacoplado de gerenciamento e manipulação de arquivos locais (`FileService`).
- Prover validação de sandbox à prova de falhas contra ataques de path traversal e symlinks externos.
- Assegurar altíssima taxa de transferência na LAN através de `http.ServeContent` com suporte completo a HTTP Range Requests.
- Implementar compactação de pastas inteiras em formato `.zip` sob demanda em streaming direto para a resposta HTTP, sem resíduos em disco.
- Criar interface web moderna, intuitiva e responsiva com visualização em lista/grade, breadcrumbs, busca/filtro instantâneo no cliente e upload via Drag & Drop.
- Permitir configuração flexível do diretório raiz via CLI (argumento posicional, flag `--dir`/`-d` ou fallback para o diretório de execução atual).
- Manter suíte abrangente de testes automatizados com cobertura superior a 80%.

**Non-Goals:**
- Sistema de autenticação com múltiplos usuários, controle de permissões por login ou ACLs complexas.
- Edição de texto ou modificação de conteúdo de arquivos diretamente no navegador.
- Persistência de banco de dados ou indexação assíncrona de arquivos pesados.

## Decisions

### Decisão 1: Arquitetura do Domínio e Portas de Entrada
- **Estrutura**:
  - `internal/core/domain/file.go`: Entidades `FileItem`, `DirectoryListing`, `Breadcrumb`, `UploadResult`.
  - `internal/core/ports/file.go`: Interface `FileService` com métodos para listagem segura, obtenção de arquivo para download, streaming de ZIP, processamento de upload e resolução canônica de caminho.
  - `internal/core/services/file.go`: Implementação `LocalFileService` contendo toda a lógica de filesystem e regras de negócio.
  - `internal/adapters/handlers/file_handler.go`: Adaptador HTTP com endpoints para visualização HTML, download de arquivo, download de ZIP e upload multipart.
- **Alternativas consideradas**: Embutir chamadas de sistema de arquivos diretamente nos handlers HTTP. *Rejeitada* para preservar testabilidade unitária e isolamento de domínio.

### Decisão 2: Streaming Direto de ZIP para Resposta HTTP (Zero I/O Residual)
- **Implementação**: A compactação de diretórios utiliza `archive/zip` gravando diretamente no stream `http.ResponseWriter` via `filepath.WalkDir`.
- **Justificativa**: Elimina a necessidade de criar arquivos `.zip` temporários em disco, garantindo performance máxima (sem escrita intermediária em disco) e prevenindo 100% o risco de arquivos residuais no diretório servido ou em `/tmp`.
- **Controle de Cancelamento**: Monitoramento do contexto da requisição (`r.Context().Done()`) para interromper o percurso e a compactação imediatamente caso o cliente desconecte.
- **Alternativas consideradas**: Gerar arquivo `.zip` temporário em pasta temp do SO e servir via `http.ServeFile`. *Rejeitada* pelo custo extra de I/O em disco e risco de deixar arquivos residuais caso o servidor seja finalizado durante a transferência.

### Decisão 3: Sandboxing Estrito, Resolução de Caminhos da CLI e Prevenção de Path Traversal
- **Mecanismo de Resolução e Validação**:
  1. O caminho raiz informado na CLI (via argumento posicional ou `--dir`/`-d`) passa por pré-processamento de expansão de atalhos de usuário: se iniciar com `~` ou `~/`, o prefixo é substituído pelo diretório home do usuário atual (`os.UserHomeDir()`).
  2. O caminho resultante (seja relativo como `./`, `../`, ou expandido do home) é normalizado para seu caminho canônico absoluto (`filepath.Abs` e `filepath.EvalSymlinks`).
  3. Qualquer caminho relativo solicitado pela web é normalizado (`filepath.Clean`) e concatenado ao caminho base raiz.
  4. O caminho resultante na web é validado via `filepath.Rel(baseDir, targetDir)`. Se o resultado iniciar com `..` ou for um caminho absoluto que não esteja contido sob a raiz, o acesso é sumariamente bloqueado.
  5. Para links simbólicos existentes dentro da árvore, é verificado o destino real através de `filepath.EvalSymlinks` para garantir que o alvo pertença à árvore autorizada.
- **Alternativas consideradas**: Apenas checar substrings `..` na URL. *Rejeitada* por ser vulnerável a codificação percentual (URL encode), múltiplos separadores e manipulação de symlinks.

### Decisão 4: Alta Performance de Transferência na Rede Local (LAN)
- **Implementação**:
  - Para arquivos individuais: utilização de `http.ServeContent` passando o ponteiro de arquivo `*os.File` aberto com `Stat()`.
  - Suporte automático a `Content-Length`, `Last-Modified`, cabeçalhos `ETag` e requisições de faixa de bytes (`Range: bytes=...`, HTTP 206 Partial Content).
  - Uploads utilizam `io.CopyBuffer` com pool de buffers otimizados (32KB a 64KB) para cópia direta do stream multipart para o arquivo destino.

### Decisão 5: Interface Web Moderna com SSR, Tailwind CSS e Alpine.js (Multi-Dispositivo)
- **Design & Layout Responsivo**:
  - Estilo visual profissional com tema moderno, paleta refinada, tipografia clara e ícones SVG intuitivos para diferentes tipos de arquivos (pastas, imagens, vídeos, áudios, código, documentos, arquivos compactados).
  - **Estratégia de Breakpoints e Fatores de Forma**:
    - **Smartphones (`< 640px` / `sm`)**: Tap targets táteis confortáveis (&ge; 44px), breadcrumbs com rolagem horizontal suave sem quebra de layout, ocultação de colunas secundárias de data/tipo na tabela para priorizar nome, tamanho e download rápido, modal de upload adaptado à largura da tela.
    - **Tablets (`640px - 1024px` / `md` e `lg`)**: Grid e tabela balanceados com espaçamento otimizado para toque e mouse, exibição de metadados principais e barra de busca integrada.
    - **Notebooks e PCs (`>= 1024px` / `xl` e `2xl`)**: Visualização completa de tabela, área ampla de arrastar e soltar (Drag & Drop) com overlay interativo, largura centralizada (`max-w-7xl`) para evitar distorções em telas ultrawide.
  - Breadcrumbs interativos no topo para navegação rápida entre níveis de diretórios com rolagem responsiva.
  - Tabela responsiva com estados visuais elegantes e ordenação prioritária de pastas.
  - Zona de arrastar e soltar (*Drag & Drop*) com feedback visual animado e fallback direto via botão para seletor de arquivos em touchscreens.
  - Campo de busca rápida no cliente via Alpine.js filtrando arquivos instantaneamente conforme digitação sem overhead no servidor.
  - Página de Status do Sistema (`/status`) alinhada com o nome e propósito do projeto ("File Server - Status do Sistema"), apresentando diagnóstico e métricas em harmonia com o produto principal.

### Decisão 6: Camada de Transporte Seguro (TLS 1.3 / HTTPS & Geração de Certificados Autoassinados em Memória)
- **Implementação**:
  - Módulo `internal/core/services/tls.go` com função `GenerateSelfSignedCertificate(hosts ...string) (tls.Certificate, error)` para gerar dinamicamente em memória um certificado X.509 RSA/ECDSA válido para `localhost`, `127.0.0.1` e IPs locais.
  - Configuração do `*tls.Config`:
    - `MinVersion: tls.VersionTLS12` (priorizando TLS 1.3).
    - Cifras modernas de alta performance aceleradas por hardware (`TLS_AES_128_GCM_SHA256`, `TLS_AES_256_GCM_SHA384`, `TLS_CHACHA20_POLY1305_SHA256`).
    - Multiplexação nativa com HTTP/2 habilitada.
  - Flags da CLI: `--tls` / `-s` (ativa HTTPS automático com certificado autoassinado), `--tls-cert` e `--tls-key` (para certificados customizados).

### Decisão 7: Identificação e Exibição de URLs de Acesso Local e LAN na Inicialização
- **Implementação**:
  - Na rotina de inicialização em `cmd/serve.go` (`RunServerWithOptions`), o servidor detecta os IPs de rede através de `services.GetLANIPAddresses()`.
  - É exibido no terminal um banner claro com:
    - URL de Acesso Local (`http://127.0.0.1:<port>` ou `https://...`)
    - URLs de Acesso na Rede Local (`http://<ip-da-lan>:<port>` ou `https://...`)
    - Informações de versão, raiz compartilhada e status de criptografia TLS.
- **Justificativa**: Facilita enormemente o acesso a partir de outros dispositivos na mesma rede local (smartphones, tablets, outros computadores), sem obrigar o usuário a abrir outro terminal para rodar `ifconfig` / `ip addr`.

## Risks / Trade-offs

- **[Risco: Alerta de certificado não confiável no navegador ao usar autoassinado]** → *Mitigação*: Comportamento padrão esperado em ambientes LAN; o servidor exibirá aviso claro no terminal instruindo como aceitar o certificado local para tráfego 100% criptografado.
- **[Risco: Upload de arquivos gigantescos exaurindo recursos de memória]** → *Mitigação*: Processamento em streaming contínuo através de `multipart.Reader` gravando diretamente no disco sem carregar o corpo completo em memória RAM.
- **[Risco: Loops infinitos ao compactar diretórios com links simbólicos cíclicos]** → *Mitigação*: Rastreamento de nós visitados (`visitedPaths` map) e limitação de profundidade de recursão na geração de arquivos `.zip`.
- **[Risco: Sobrescrita de arquivos com nomes duplicados no upload]** → *Mitigação*: Estratégia configurável de versionamento de nome (`arquivo (1).ext`) ou substituição controlada com sanitização estrita via `filepath.Base`.

## Estratégia de Testes

1. **Testes Unitários (TDD/BDD)**:
   - `internal/core/services/file_test.go`: Testes de listagem, ordenação, permissões, ZIP streaming e uploads.
   - `internal/core/services/tls_test.go`: Testes de geração de certificados autoassinados, validação de pares de chaves e configuração de `tls.Config`.
   - Testes exaustivos de segurança e sandbox: tentativas de path traversal com `../`, caminhos absolutos, caracteres nulos, symlinks externos.
2. **Testes de Integração HTTP & TLS**:
   - `internal/adapters/handlers/file_handler_test.go`: Testes de endpoints HTTP/HTTPS, downloads parciais (206), uploads multipart e respostas de erro.
3. **Testes de CLI**:
   - `cmd/root_test.go` e `cmd/serve_test.go`: Testes de inicialização com `--tls`, `--tls-cert`, `--tls-key`, argumento posicional, flag `--dir`/`-d`.
4. **Validação de Qualidade**:
   - Verificação de cobertura total `>= 80%` via `make test` e aprovação no linter estrito via `make check`.
