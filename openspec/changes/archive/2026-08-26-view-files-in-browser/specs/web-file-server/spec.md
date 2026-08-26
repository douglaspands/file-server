## ADDED Requirements

### Requirement: Visualização Direta (Inline Preview) de Arquivos no Navegador
O sistema SHALL disponibilizar endpoint `/view/{path}` para visualização direta (*inline preview*) e streaming de arquivos com cabeçalho `Content-Disposition: inline` e `Content-Type` detectado automaticamente (MIME type adequado), permitindo que formatos suportados pelo navegador (documentos PDF, imagens, vídeos, áudios, textos e códigos-fonte) sejam abertos e renderizados diretamente na interface do navegador em nova aba, enquanto formatos não renderizáveis ou não suportados pelo navegador realizem o download normalmente, mantendo em paralelo botões e links dedicados para forçar o download direto (`/download/{path}`) de qualquer arquivo.
*(Visão PO: Proporciona agilidade imediata aos usuários para inspecionar fotos, documentos e mídias sem poluir o armazenamento local com downloads desnecessários, mantendo total flexibilidade para baixar quando desejado. Visão QA: Garante envio correto dos headers HTTP Content-Disposition inline vs attachment, MIME types precisos, streaming sem consumo excessivo de memória, Range requests (206) e isolamento rigoroso contra path traversal).*

#### Scenario: Visualização direta de documento PDF em nova aba
- **GIVEN** que o arquivo `documento.pdf` existe no diretório servido
- **WHEN** o usuário clicar no nome do arquivo ou no botão "Visualizar" na interface web, ou requisitar `GET /view/documento.pdf`
- **THEN** o servidor deve responder com status HTTP 200 (ou 206), cabeçalho `Content-Type: application/pdf` e cabeçalho `Content-Disposition: inline; filename="documento.pdf"`, permitindo a renderização nativa pelo navegador

#### Scenario: Visualização direta de imagens no navegador
- **GIVEN** que existe uma imagem `foto.png` no diretório
- **WHEN** o usuário acessar a URL `/view/foto.png` ou clicar na ação de visualização
- **THEN** o servidor deve responder com `Content-Type: image/png` e `Content-Disposition: inline; filename="foto.png"`, exibindo a imagem diretamente

#### Scenario: Streaming e reprodução inline de áudio e vídeo com HTTP Range
- **GIVEN** que existe o arquivo `video.mp4` ou `musica.mp3` no diretório
- **WHEN** o cliente acessar `/view/video.mp4` com cabeçalho `Range: bytes=0-1048575`
- **THEN** o servidor deve responder com status HTTP 206 (Partial Content), cabeçalho `Content-Range`, `Content-Disposition: inline; filename="video.mp4"` e `Content-Type: video/mp4`, permitindo reprodução e busca temporal (seek) no player nativo do navegador

#### Scenario: Visualização inline de arquivos de texto e código-fonte
- **GIVEN** que existe um arquivo `script.py` ou `documento.txt` no diretório
- **WHEN** o usuário acessar `/view/script.py` ou clicar em visualizar
- **THEN** o servidor deve responder com `Content-Disposition: inline; filename="script.py"` e `Content-Type: text/plain; charset=utf-8` (ou tipo específico de texto), permitindo a leitura direta no navegador

#### Scenario: Preservação de botão de download direto independente
- **GIVEN** que o arquivo `relatorio.pdf` está listado na tabela de arquivos
- **WHEN** o usuário clicar no botão dedicado "Baixar" com link `/download/relatorio.pdf`
- **THEN** o servidor deve responder com cabeçalho `Content-Disposition: attachment; filename="relatorio.pdf"`, forçando a caixa de diálogo de download do navegador

#### Scenario: Fallback transparente para arquivos não renderizáveis
- **GIVEN** que existe um arquivo binário ou executável `instalador.bin` ou arquivo compactado `dados.tar.gz`
- **WHEN** o usuário clicar no nome do arquivo na interface web
- **THEN** o link deve apontar para `/download/instalador.bin` ou o endpoint `/view/instalador.bin` deve ser processado pelo navegador efetuando o download padrão

#### Scenario: Bloqueio estrito de Path Traversal no endpoint de visualização
- **WHEN** uma requisição for enviada para `/view/../../etc/passwd` ou manipulação de caminho fora da raiz
- **THEN** o sistema deve bloquear o acesso respondendo com status HTTP 403 (Forbidden) e não expor qualquer dado fora da sandbox raiz
