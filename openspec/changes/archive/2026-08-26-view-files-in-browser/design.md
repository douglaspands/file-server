## Context

Atualmente o File Server provê endpoints para navegação (`/files/*`), download forçado de arquivos (`/download/*` com `Content-Disposition: attachment`), download de pastas em ZIP (`/zip/*`) e upload (`/upload/*`). Na camada web, a tabela de arquivos renderiza links exclusivamente para `/download/*`.

Para viabilizar a visualização direta e reprodução de mídias no navegador (mantendo a opção de download sob demanda), é necessário estender a Clean Architecture com novo endpoint de visualização, suporte de metadados no domínio para formatos suportados e controles visuais aprimorados no template SSR.

## Goals / Non-Goals

**Goals:**
- Implementar endpoint HTTP `/view/*` com cabeçalho `Content-Disposition: inline` e suporte a HTTP Range (206) para visualização e streaming direto.
- Expandir o domínio (`domain.FileItem`) com a flag `IsViewable` e catalogação precisa de formatos com suporte nativo pelos navegadores (PDFs, imagens, vídeos, áudios, textos e códigos).
- Garantir registro e detecção precisa de Content-Type MIME (incluindo charset UTF-8 para arquivos de texto e código).
- Atualizar o template `file_table.html` com ação dedicada de visualização (ícone de olho, abrindo em nova aba com `target="_blank"`) e link direto no nome do arquivo para itens visualizáveis, mantendo o botão de download (`/download/*`) sempre disponível.
- Manter cobertura de testes $\ge 80\%$ com testes unitários e de integração cobrindo o novo endpoint, cabeçalhos e segurança de sandbox.

**Non-Goals:**
- Implementar visualizadores ou conversores proprietários no servidor (ex: conversor de DOCX para PDF ou renderizador de LibreOffice). O foco é suporte nativo do navegador para formatos web padrão.
- Bloquear downloads forçados — a opção de download continua disponível com paridade de acesso.

## Decisions

### 1. Novo Endpoint Dedicado `/view/*` vs Query Parameter
- **Decisão**: Criar o endpoint `/view/{path}` no `FileHandler` e registrá-lo no mux em `RegisterRoutes`.
- **Justificativa**: Mantém simetria e clareza com as rotas já estabelecidas (`/download/*`, `/files/*`, `/zip/*`, `/upload/*`), permitindo que URLs sejam facilmente compartilháveis e previsíveis na rede local.
- **Alternativa Considerada**: Utilizar `?inline=true` em `/download/*`. Rejeitada para manter responsabilidade única nos endpoints e evitar condicionais em cascata no roteamento.

### 2. Header `Content-Disposition: inline` e Suporte a `http.ServeContent`
- **Decisão**: Utilizar `w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))` e delegar o streaming ao `http.ServeContent`.
- **Justificativa**: `http.ServeContent` gerencia nativamente `Content-Length`, `Last-Modified`, caching `ETag`/`If-Modified-Since` e `HTTP Range (206)` — essencial para streaming progressivo de vídeo, áudio e navegação rápida em PDFs extensos sem sobrecarga de memória.

### 3. Classificação de Formatos Visualizáveis no Domínio (`FileItem.IsViewable`)
- **Decisão**: Adicionar a propriedade booleana `IsViewable` em `domain.FileItem`, calculada por uma função de domínio `IsViewableFormat(category domain.FileCategory, ext string) bool`.
- **Formatos Elegíveis**:
  - Imagens (`.png`, `.jpg`, `.jpeg`, `.gif`, `.svg`, `.webp`, `.bmp`, `.ico`, `.tiff`).
  - Documentos renderizáveis (`.pdf`, `.txt`, `.csv`, `.tsv`, `.md`).
  - Áudio (`.mp3`, `.wav`, `.ogg`, `.flac`, `.aac`, `.m4a`).
  - Vídeo (`.mp4`, `.webm`, `.ogg`, `.mov`, `.m4v`).
  - Código-fonte / Texto plano (`.go`, `.js`, `.ts`, `.html`, `.css`, `.json`, `.xml`, `.yaml`, `.py`, `.sh`, `.sql`, etc.).
- **Justificativa**: Permite que o template SSR e clientes da API JSON saibam imediatamente se o arquivo é visualizável sem precisar duplicar lógica de extensões no HTML/JS.

### 4. Ajustes no Template de Tabela de Arquivos (`file_table.html`)
- **Decisão**:
  - Para arquivos com `IsViewable = true`:
    - Nome do arquivo linka para `/view/{{.RelativePath}}` com `target="_blank"` e `rel="noopener noreferrer"`.
    - Coluna de ações exibe botão de "Visualizar" (ícone de olho) + botão de "Baixar" (ícone de download).
  - Para arquivos com `IsViewable = false`:
    - Nome do arquivo linka para `/download/{{.RelativePath}}`.
    - Coluna de ações exibe apenas o botão de "Baixar".

## Risks / Trade-offs

- **[Risco] Vulnerabilidade de MIME Confusion / XSS em arquivos HTML/SVG servidos inline** → **Mitigação**: O endpoint `/view/*` utiliza validação estrita de sandbox (`ResolveAndValidatePath`) impedindo path traversal, e o servidor pode aplicar header de segurança `X-Content-Type-Options: nosniff` quando necessário.
- **[Risco] Tipos MIME não mapeados pelo sistema operacional host** → **Mitigação**: Inicialização explícita de tipos MIME comuns via `mime.AddExtensionType` para garantir que extensões modernas (ex: `.webp`, `.svg`, `.md`, `.ts`) sejam servidas com Content-Type correto.
