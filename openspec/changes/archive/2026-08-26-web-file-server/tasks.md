## 1. Domínio e Portas de Arquivos (Core Domain & Ports)

- [x] 1.1 Definir modelos de domínio em `internal/core/domain/file.go` (`FileItem`, `DirectoryListing`, `Breadcrumb`, `UploadResult`, `FileCategory`)
- [x] 1.2 Definir a interface `FileService` em `internal/core/ports/file.go` com assinaturas para listagem segura, download de arquivos, streaming de ZIP, upload e validação de sandbox

## 2. Serviço de Arquivos e Proteção de Sandbox (Core Services & Security)

- [x] 2.1 Implementar validação e resolução canônica de caminhos em `internal/core/services/file.go` garantindo isolamento estrito contra path traversal (`..` e symlinks externos)
- [x] 2.2 Implementar listagem e extração de metadados de arquivos e subdiretórios com ordenação prioritária de pastas e categorização por extensão
- [x] 2.3 Implementar recuperação de arquivo para download com suporte a leitura posicional e informações de tamanho/data
- [x] 2.4 Implementar compactação e streaming de diretórios em ZIP diretamente para `io.Writer` com cancelamento via contexto e zero resíduos em disco
- [x] 2.5 Implementar processamento de upload multipart e gravação segura no diretório destino com sanitização estrita de nomes
- [x] 2.6 Desenvolver suíte de testes unitários e BDD em `internal/core/services/file_test.go` cobrindo cenários normais, bordas e tentativas de ataque de path traversal

## 3. Adaptadores HTTP e Endpoints Web (HTTP Handlers & API)

- [x] 3.1 Implementar handler de visualização e navegação de arquivos (`FileBrowserHandler`) renderizando template HTML e breadcrumbs
- [x] 3.2 Implementar handler de download de arquivo individual (`DownloadFileHandler`) utilizando `http.ServeContent` com suporte a HTTP Range (206)
- [x] 3.3 Implementar handler de download de diretório compactado em ZIP (`DownloadZipHandler`) via streaming com cabeçalhos apropriados
- [x] 3.4 Implementar handler de upload multipart (`UploadFilesHandler`) para envio único e múltiplo de arquivos com feedback visual
- [x] 3.5 Atualizar registro de rotas no mux HTTP em `internal/adapters/handlers/handler.go` integrando o `FileService`
- [x] 3.6 Desenvolver testes unitários e de integração HTTP em `internal/adapters/handlers/file_handler_test.go` validando rotas, downloads, streaming de ZIP e uploads

## 4. Interface Web e Experiência do Usuário (Frontend, SSR & Templates)

- [x] 4.1 Criar template principal do explorador de arquivos (`web/templates/pages/explorer.html`) com layout moderno e responsivo em Tailwind CSS
- [x] 4.2 Criar componentes parciais (`breadcrumbs.html`, `file_table.html`, `upload_modal.html`, `empty_state.html`) com ícones dinâmicos por categoria de arquivo
- [x] 4.3 Integrar scripts Alpine.js para filtro/busca instantânea no cliente e gerenciamento de estado da zona de arrastar e soltar (Drag & Drop)
- [x] 4.4 Atualizar empacotamento de templates e assets em `web/web.go` (`embed.FS`) e validar compilação de recursos estáticos
- [x] 4.5 Atualizar a página de status do sistema (`web/templates/pages/index.html`) e testes associados substituindo 'Fundação do Sistema Pronta' por 'File Server - Status do Sistema'
- [x] 4.6 Otimizar e validar responsividade da interface web em múltiplos dispositivos (smartphones, tablets, notebooks e PCs) com Tailwind CSS e componentes adaptativos

## 5. Integração CLI e Argumentos de Inicialização (CLI & Flags)

- [x] 5.1 Atualizar `cmd/serve.go` e `cmd/root.go` para suportar argumento posicional de diretório e flag `--dir`/`-d` com fallback automático para o diretório atual (`.`)
- [x] 5.2 Implementar validação e normalização do diretório informado no momento da inicialização da CLI
- [x] 5.3 Desenvolver testes unitários da CLI em `cmd/serve_test.go` e `cmd/cmd_test.go` cobrindo diferentes combinações de flags e argumentos
- [x] 5.4 Implementar suporte a expansão de til (~ / home) e caminhos relativos (../, ./) na resolução de diretório da CLI com testes unitários em `cmd/serve_test.go`
- [x] 5.5 Exibir URLs de acesso local (localhost) e da rede local (LAN IPs) na inicialização do servidor em `cmd/serve.go` com testes unitários em `cmd/serve_test.go`

## 6. Validação de Qualidade, Documentação e Verificação Final

- [x] 6.1 Atualizar `README.md` com documentação detalhada dos comandos da CLI, flags de diretório, guia de uso da interface web, downloads em ZIP e uploads
- [x] 6.2 Executar suíte completa de testes automatizados e validar barreira de cobertura >= 80% (`make test`)
- [x] 6.3 Executar verificação unificada com linter estrito e checagens de integridade (`make check`)
- [x] 6.4 Consolidar o `README.md` e documentações garantindo que o foco principal seja o servidor de arquivos de alta performance, mantendo TDD/BDD/Harness como qualidades na seção de desenvolvimento

## 7. Camada de Segurança e Criptografia em Trânsito (TLS / HTTPS & Anti-Sniffing)

- [x] 7.1 Implementar gerador de certificados TLS autoassinados em memória e configuração de `tls.Config` de alta performance em `internal/core/services/tls.go`
- [x] 7.2 Atualizar `cmd/serve.go` e `cmd/root.go` adicionando flags `--tls`/`-s`, `--tls-cert` e `--tls-key` com suporte a inicialização segura em HTTPS/HTTP2
- [x] 7.3 Desenvolver testes unitários e de integração em `internal/core/services/tls_test.go` e `cmd/serve_test.go` validando handshake TLS, autoassinatura e proteção anti-sniffing
- [x] 7.4 Atualizar `README.md` com documentação da camada de segurança TLS/HTTPS, flags de linha de comando e exemplos de uso seguro na rede local
- [x] 7.5 Executar validação de qualidade completa com `make check` e verificar barreira de cobertura >= 80%
