## Why

Atualmente, na interface web do File Server, todos os links e botões de ação para arquivos apontam exclusivamente para o endpoint `/download/...` com o cabeçalho `Content-Disposition: attachment`, forçando o navegador a sempre baixar o arquivo para o disco local em vez de exibi-lo. Usuários que precisam apenas consultar rapidamente documentos em PDF, imagens (PNG, JPG, SVG, WebP), áudios, vídeos ou arquivos de texto/código são obrigados a efetuar o download e abrir em aplicativos externos.

Esta mudança introduz a capacidade de visualização direta e reprodução (*inline preview*) no próprio navegador para formatos suportados nativamente (como PDFs, imagens, mídias e textos), preservando concomitantemente a opção de download forçado via botão dedicado.

## What Changes

- **Novo Endpoint de Visualização Inline (`/view/*`)**: Endpoint HTTP dedicado que serve arquivos com o cabeçalho `Content-Disposition: inline; filename="..."` e `Content-Type` detectado com precisão, com suporte completo a Range requests (HTTP 206) para streaming de mídia e visualização de documentos.
- **Detecção de Tipos Visualizáveis no Domínio**: Adição de método/propriedade no domínio (`FileItem.IsViewable` ou `CanPreview`) para identificar formatos suportados pelo navegador (PDFs, imagens, vídeos, áudios, textos e códigos-fonte).
- **Interface Web com Ação de Visualização e Download**:
  - Para arquivos visualizáveis, o clique no nome do arquivo e o botão de ação rápida "Visualizar" abrem o arquivo diretamente no navegador (em nova aba com `target="_blank"`).
  - Adição de botão de ação de visualização com ícone descritivo (olho / preview) na tabela de arquivos.
  - Manutenção integral do botão de download direto (`/download/*`) com `Content-Disposition: attachment` para todos os arquivos.
  - Para arquivos não visualizáveis (ex: arquivos binários, executáveis, arquivos compactados), o clique no nome preserva o comportamento de download direto.

## Capabilities

### Modified Capabilities
- `web-file-server`: Adiciona o requisito de visualização direta (inline preview) de arquivos no navegador com endpoint `/view/*`, metadados de suporte no modelo de domínio e controles visuais dedicados na interface web.

## Impact

- **Código Go**:
  - `internal/core/domain/file.go`: Inclusão de helper/campo `IsViewable` no struct `FileItem` e mapeamento de extensões com suporte a visualização inline pelo navegador.
  - `internal/adapters/handlers/file_handler.go`: Implementação do `ViewFileHandler` que atende à rota `/view/*` com `Content-Disposition: inline` e headers de MIME adequados.
  - `internal/adapters/handlers/handler.go`: Registro das rotas `/view/` e `/view` no multiplexador HTTP.
- **Templates Web**:
  - `web/templates/partials/file_table.html`: Atualização dos links e botões de ação para incorporar o botão "Visualizar" e direcionar o clique de arquivos visualizáveis para a rota `/view/*`.
- **Testes**:
  - Testes unitários e de integração em Go para `ViewFileHandler`, validação de headers `Content-Disposition: inline`, Range requests e segurança de sandbox / path traversal na rota `/view/*`.
