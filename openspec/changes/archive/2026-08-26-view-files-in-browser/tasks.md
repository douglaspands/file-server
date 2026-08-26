## 1. Modelo de Domínio e Suporte a Formatos Visualizáveis

- [x] 1.1 Expandir `domain.FileItem` adicionando o campo `IsViewable bool` com tag JSON
- [x] 1.2 Implementar a função `IsViewableFormat(category domain.FileCategory, ext string) bool` no pacote `domain` cobrindo PDFs, imagens, vídeos, áudios, textos e códigos-fonte
- [x] 1.3 Atualizar a construção de itens em `LocalFileService.ListDirectory` para preencher o atributo `IsViewable`
- [x] 1.4 Adicionar e atualizar testes unitários em `core/domain/file_test.go` e `core/services/file_test.go`

## 2. Handler HTTP e Roteamento para Visualização Inline

- [x] 2.1 Implementar `ViewFileHandler` em `internal/adapters/handlers/file_handler.go` configurando `Content-Disposition: inline` e servindo arquivos via `http.ServeContent` com suporte a HTTP Range (206)
- [x] 2.2 Registrar rotas `/view/` e `/view` no multiplexador em `internal/adapters/handlers/handler.go`
- [x] 2.3 Garantir mapeamento de extensões e Content-Type MIME adicionais (ex: `.md`, `.webp`, `.svg`, `.ts`)
- [x] 2.4 Escrever testes unitários e de integração em `internal/adapters/handlers/file_handler_test.go` cobrindo sucesso, HTTP Range, headers de inline, arquivos inexistentes e bloqueio estrito de path traversal

## 3. Atualização da Interface Web e Templates SSR

- [x] 3.1 Atualizar `web/templates/partials/file_table.html` para que arquivos com `IsViewable = true` tenham links abrindo `/view/*` com `target="_blank"` e `rel="noopener noreferrer"`
- [x] 3.2 Adicionar botão de ação com ícone de visualização ("Visualizar") na coluna de ações para arquivos visualizáveis
- [x] 3.3 Preservar o botão dedicado "Baixar" (`/download/*`) e o comportamento de download para arquivos não visualizáveis
- [x] 3.4 Validar renderização e testes de template em `web/web_test.go`

## 4. Validação de Qualidade e Cobertura de Testes

- [x] 4.1 Executar `make test` e verificar cobertura de código $\ge 80\%$
- [x] 4.2 Executar `make check` (linting e vet) garantindo conformidade total do projeto
