# Capability: web-file-server

## Purpose

Prover um servidor de arquivos web de alta performance via rede local com interface moderna e intuitiva, navegação fluida entre diretórios, download de arquivos individuais e pastas compactadas em ZIP via streaming, upload de múltiplos arquivos e isolamento estrito de sandbox contra path traversal.

## Requirements

### Requirement: Configuração de Diretório Raiz via CLI e Fallback Automático para Diretório Atual
O sistema SHALL permitir que o diretório raiz servido seja definido dinamicamente via CLI, assumindo o diretório de execução atual (`cwd`) caso nenhum parâmetro seja informado, ou recebendo o caminho absoluto/relativo como argumento posicional ou através da flag `--dir` / `-d`.
*(Visão PO: Proporciona facilidade instantânea de uso no terminal, permitindo que qualquer usuário inicie o compartilhamento com um único comando na pasta desejada. Visão QA: Garante validação robusta de caminhos válidos, normalização canônica e mensagens de erro descritivas para diretórios inexistentes ou inacessíveis).*

#### Scenario: Inicialização padrão no diretório corrente
- **GIVEN** que o usuário está em um diretório no terminal
- **WHEN** executar o comando `file-server` ou `file-server serve` sem argumentos de diretório
- **THEN** o sistema deve iniciar o servidor web adotando o diretório de trabalho atual (`.`) como a raiz do compartilhamento

#### Scenario: Inicialização com diretório customizado via argumento posicional
- **GIVEN** que o diretório `/meus/arquivos` existe e tem permissão de leitura
- **WHEN** o usuário executar `file-server /meus/arquivos` ou `file-server serve /meus/arquivos`
- **THEN** o sistema deve iniciar o servidor web utilizando `/meus/arquivos` como diretório raiz canônico

#### Scenario: Inicialização com diretório via flag explícita
- **GIVEN** que o diretório `./dados` existe
- **WHEN** o usuário executar `file-server --dir ./dados` ou `file-server serve -d ./dados`
- **THEN** o sistema deve resolver o caminho para o diretório raiz correspondente

#### Scenario: Inicialização com caminho relativo contendo til (~ / home)
- **GIVEN** que o usuário informa `~` ou `~/Downloads` como argumento posicional ou na flag `-d`
- **WHEN** o servidor for inicializado
- **THEN** o sistema deve expandir `~` para o diretório home do usuário atual (`os.UserHomeDir()`) e resolver o caminho absoluto canônico correspondente

#### Scenario: Inicialização com caminhos relativos (../, ./)
- **GIVEN** que o usuário executa o comando a partir de um subdiretório e passa `../` ou `../../pasta`
- **WHEN** o servidor for inicializado
- **THEN** o sistema deve resolver o caminho relativo em relação ao diretório de trabalho atual para o seu caminho canônico absoluto correspondente

#### Scenario: Exibição de URLs de acesso local e da rede (LAN) na inicialização
- **GIVEN** que o servidor é inicializado na porta 8080 (ou customizada) em uma máquina conectada à rede local (ex: IP 192.168.1.50)
- **WHEN** o comando de inicialização for executado
- **THEN** o sistema deve detectar os endereços IP das interfaces de rede ativas e imprimir no terminal as URLs diretas de acesso para localhost e para cada IP de rede local disponível

#### Scenario: Tentativa de inicialização com diretório inválido ou inexistente
- **GIVEN** que o caminho `/diretorio/inexistente` não existe no sistema de arquivos
- **WHEN** o usuário tentar iniciar o servidor apontando para este caminho
- **THEN** o sistema deve abortar a execução imediatamente com código de erro diferente de zero e exibir mensagem informativa no terminal

### Requirement: Navegação Interativa e Listagem de Diretórios e Arquivos
O sistema SHALL fornecer uma interface web e endpoints para visualização e navegação hierárquica de arquivos e subdiretórios a partir da raiz configurada, exibindo metadados essenciais (nome, tipo, tamanho formatado e data de modificação) e trilha de navegação (*breadcrumbs*).
*(Visão PO: Permite que os usuários encontrem rapidamente qualquer arquivo na estrutura de pastas através de uma interface visual clara e intuitiva. Visão QA: Valida a correta renderização de metadados, ordenação (pastas primeiro), tratamento de caracteres especiais em nomes e integridade da navegação em profundidade).*

#### Scenario: Listagem de conteúdo do diretório raiz
- **GIVEN** que o servidor está em execução
- **WHEN** o usuário acessar a rota `/` ou `/files` no navegador
- **THEN** o sistema deve renderizar a lista de arquivos e subdiretórios do diretório raiz com metadados formatados e breadcrumb indicando a raiz

#### Scenario: Navegação para subdiretório e exibição de breadcrumbs
- **GIVEN** que existe a pasta `documentos/projetos` dentro da raiz
- **WHEN** o usuário clicar na pasta ou acessar `/files/documentos/projetos`
- **THEN** o sistema deve exibir os arquivos desse subdiretório e uma barra de navegação (*breadcrumbs*) permitindo voltar para `documentos` ou para a raiz com um clique

#### Scenario: Listagem de diretório vazio
- **GIVEN** que um subdiretório não contém arquivos nem pastas
- **WHEN** o usuário navegar para este subdiretório
- **THEN** a interface deve exibir um estado visual amigável indicando que a pasta está vazia, mantendo ativas as opções de upload e retorno

#### Scenario: Distinção visual e ordenação de itens
- **WHEN** uma listagem de pasta for renderizada
- **THEN** as pastas devem ser listadas prioritariamente antes dos arquivos, acompanhadas por ícones distintos representativos para diretórios e tipos comuns de arquivos (áudio, vídeo, imagem, código, pdf, compactados)

### Requirement: Transferência e Download de Alta Performance de Arquivos Individuais
O sistema SHALL disponibilizar download de arquivos individuais com alta taxa de transferência na rede local, suporte a streaming direto (`http.ServeContent`), requisições parciais (header `Range` para pausa e retomada de downloads) e baixo consumo de memória RAM.
*(Visão PO: Assegura downloads ultrarrápidos de arquivos pesados (vídeos, ISOs, backups) na LAN sem travamentos. Visão QA: Valida integridade do payload binário, headers MIME corretos, código HTTP 206 para requisições Range e código HTTP 404 para arquivos inexistentes).*

#### Scenario: Download direto de arquivo com sucesso
- **GIVEN** que o arquivo `instalador.iso` de 2GB existe no diretório
- **WHEN** o usuário solicitar o download via interface web ou requisição GET direta
- **THEN** o servidor deve transmitir o arquivo via streaming direto com Content-Type apropriado e header `Content-Disposition: attachment; filename="instalador.iso"`

#### Scenario: Suporte a requisição parcial (HTTP Range Request)
- **GIVEN** que um cliente solicita os bytes `0-1048575` de um arquivo
- **WHEN** o cabeçalho `Range: bytes=0-1048575` for enviado na requisição HTTP
- **THEN** o servidor deve responder com status HTTP 206 (Partial Content) contendo exatamente o intervalo de bytes solicitado e cabeçalho `Content-Range`

#### Scenario: Tentativa de download de arquivo inexistente
- **WHEN** uma requisição de download for feita para um caminho de arquivo inexistente
- **THEN** o servidor deve responder com status HTTP 404 (Not Found)

### Requirement: Download de Diretórios Compactados em ZIP sob Demanda
O sistema SHALL permitir o download de diretórios inteiros empacotados em arquivo `.zip` gerado sob demanda, realizando a compactação em streaming direto para a resposta HTTP (`io.Writer`) sem deixar arquivos residuais ou lixo no diretório servido e garantindo limpeza estrita de buffers ou arquivos temporários após a conclusão ou interrupção do download.
*(Visão PO: Facilita o resgate de projetos ou conjuntos completos de fotos/pastas em um único clique sem exigir compactação manual prévia no servidor. Visão QA: Valida que nenhum arquivo `.zip` temporário ou resíduo permaneça no disco após download completo ou cancelado, e que a estrutura interna do arquivo compactado reproduza a hierarquia original).*

#### Scenario: Download bem-sucedido de pasta compactada em ZIP
- **GIVEN** que a pasta `fotos_viagem` contém múltiplos arquivos e subpastas
- **WHEN** o usuário acionar a opção "Baixar Pasta (.zip)"
- **THEN** o servidor deve gerar o arquivo `fotos_viagem.zip` em streaming com header `Content-Disposition: attachment; filename="fotos_viagem.zip"`, entregando toda a árvore de arquivos compactada

#### Scenario: Ausência total de resíduos no diretório servido e limpeza temporária
- **GIVEN** que um download de diretório em ZIP foi concluído ou abortado pelo cliente
- **WHEN** o sistema de arquivos local for inspecionado
- **THEN** nenhum arquivo temporário, lock ou resíduo deve existir dentro da pasta servida ou na área temporária do sistema

#### Scenario: Download em ZIP de pasta vazia
- **GIVEN** que um subdiretório está completamente vazio
- **WHEN** o usuário solicitar o download em ZIP dessa pasta
- **THEN** o servidor deve retornar um arquivo `.zip` válido sem erros, contendo um arquivo vazio ou estrutura correspondente

### Requirement: Upload Simples e Múltiplo de Arquivos no Diretório Navegado
O sistema SHALL disponibilizar mecanismo de upload de arquivos únicos ou múltiplos via formulário multipart/form-data e interface drag-and-drop, gravando os arquivos recebidos diretamente no diretório que o usuário está visualizando no momento.
*(Visão PO: Proporciona experiência fluida para envio de documentos e arquivos diretamente para a pasta correta sem fricção. Visão QA: Valida integridade dos arquivos salvos, suporte a uploads simultâneos, prevenção de sobrescrita acidental descontrolada e limites de segurança).*

#### Scenario: Upload de arquivo único com sucesso
- **GIVEN** que o usuário está navegando no diretório `documentos/trabalho`
- **WHEN** o usuário selecionar o arquivo `relatorio.pdf` e enviar pelo formulário ou arrastar para a área de upload
- **THEN** o arquivo `relatorio.pdf` deve ser salvo no diretório `documentos/trabalho` com permissões adequadas e a listagem de arquivos deve ser atualizada

#### Scenario: Upload de múltiplos arquivos simultaneamente
- **GIVEN** que o usuário seleciona 5 fotos simultâneas para envio
- **WHEN** a requisição de upload em lote for enviada
- **THEN** todos os 5 arquivos devem ser gravados com integridade no diretório atual e o usuário deve receber confirmação de sucesso

#### Scenario: Tratamento de erro em caso de diretório sem permissão de escrita
- **GIVEN** que o diretório de destino é somente leitura no sistema operacional
- **WHEN** o usuário tentar efetuar um upload
- **THEN** o servidor deve retornar status de erro explicativo (ex: HTTP 403 ou 500) informando a impossibilidade de gravação sem corromper o estado do servidor

### Requirement: Proteção Estrita de Sandbox e Prevenção de Path Traversal
O sistema SHALL proibir terminantemente qualquer tentativa de acesso, leitura, escrita, download ou navegação que ultrapasse os limites do diretório raiz configurado na inicialização da aplicação (Sandboxing absoluto).
*(Visão PO: Garante a segurança absoluta do host onde o servidor é executado, impedindo que usuários na rede acessem arquivos confidenciais do sistema operacional ou de outras pastas não autorizadas. Visão QA: Testa vetores de ataque como `../`, `..%2f`, caminhos absolutos forçados, bypass de normalização e links simbólicos que apontem para fora do diretório raiz).*

#### Scenario: Bloqueio de navegação com sequências de escape (Path Traversal)
- **WHEN** um usuário ou atacante tentar acessar caminhos como `/files/../../etc/passwd` ou `/files/..%2f..%2f`
- **THEN** o sistema deve sanitizar e validar o caminho contra a raiz canônica, rejeitando a requisição com HTTP 403 (Forbidden) ou HTTP 400 (Bad Request)

#### Scenario: Bloqueio de download fora do diretório raiz
- **WHEN** uma requisição de download direto ou ZIP for enviada com caminho relativo ou manipulado para arquivo fora da raiz
- **THEN** o sistema deve negar imediatamente a operação com HTTP 403 (Forbidden) e não expor qualquer dado do arquivo

#### Scenario: Bloqueio de upload fora dos limites da raiz
- **WHEN** uma requisição de upload tentar injetar arquivos em diretórios superiores à raiz através de nomes manipulados (ex: `../../malicioso.sh`)
- **THEN** o sistema deve sanitizar o nome base do arquivo (`filepath.Base`), restringindo a gravação estritamente dentro do diretório de destino validado sob a raiz

#### Scenario: Bloqueio de links simbólicos externos à raiz
- **GIVEN** que existe um symlink dentro do diretório raiz apontando para `/etc` ou outra pasta fora da raiz
- **WHEN** for realizada uma requisição de navegação ou download através desse symlink
- **THEN** o sistema deve resolver o caminho canônico de destino e bloquear o acesso se estiver localizado fora dos limites da pasta raiz configurada

### Requirement: Interface Web Intuitiva, Responsiva e Agradável em Múltiplos Dispositivos
O sistema SHALL apresentar uma interface gráfica web moderna, polida e estritamente responsiva (construída com HTML5 semântico, Tailwind CSS, HTMX e Alpine.js), adaptando-se perfeitamente a diferentes fatores de forma — incluindo smartphones, tablets, notebooks e computadores de mesa (PCs) —, oferecendo experiência fluida com busca/filtro em tempo real no cliente e área intuitiva de arrastar e soltar (drag & drop) com alternativa acessível para dispositivos sensíveis ao toque.
*(Visão PO: Garante uma experiência de usuário profissional e sem atritos em qualquer tela (do celular do colaborador até o monitor ultrawide do desenvolvedor na rede local). Visão QA: Valida legibilidade, tap targets adequados em telas móveis, overflow horizontal controlado, responsividade de modais e adaptação dinâmica de tabelas em múltiplos breakpoints).*

#### Scenario: Experiência interativa de upload via Drag & Drop e alternativa tátil
- **WHEN** o usuário arrasta arquivos sobre a janela do navegador no computador ou notebook
- **THEN** a interface deve exibir uma sobreposição visual destacada indicando a área de soltura e, ao soltar os arquivos, iniciar imediatamente o upload com feedback de progresso
- **AND** em dispositivos móveis (smartphones/tablets), o botão de envio deve acionar o seletor nativo de arquivos do sistema operacional

#### Scenario: Filtro dinâmico de arquivos na visualização atual
- **GIVEN** que uma pasta contém dezenas de arquivos
- **WHEN** o usuário digitar um termo no campo de busca/filtro da interface
- **THEN** a lista exibida deve filtrar instantaneamente os itens no cliente (via Alpine.js), ocultando itens não correspondentes sem disparar requisições extras ao servidor

#### Scenario: Adaptação responsiva para smartphones (telas compactas < 640px)
- **WHEN** a interface for acessada por um smartphone
- **THEN** os botões e links de navegação devem possuir dimensões mínimas de toque confortáveis (tap targets >= 44px), a barra de breadcrumbs deve possuir rolagem horizontal suave sem quebrar a tela, colunas não essenciais da tabela devem ser ocultadas ou condensadas e os modais devem preencher a largura útil da tela

#### Scenario: Adaptação responsiva para tablets (telas médias sensíveis ao toque 640px a 1024px)
- **WHEN** a interface for acessada por um tablet em modo retrato ou paisagem
- **THEN** a interface deve equilibrar o espaçamento visual exibindo metadados essenciais (nome, tamanho, data) com controles táteis confortáveis para download e navegação

#### Scenario: Adaptação responsiva para notebooks e PCs (telas amplas >= 1024px)
- **WHEN** a interface for acessada por um notebook ou computador de mesa
- **THEN** a tabela deve exibir todas as colunas com informações detalhadas, área ampla de drag-and-drop, limites de largura centralizada (max-w-7xl) para evitar estiramento excessivo em monitores ultrawide e suporte completo a atalhos de mouse e teclado

#### Scenario: Exibição coerente da página de status do sistema
- **GIVEN** que o usuário acessa a rota `/status`
- **WHEN** a página de diagnóstico do sistema for renderizada
- **THEN** o título e cabeçalho principal devem exibir "File Server - Status do Sistema" alinhados com o propósito de servidor de arquivos de alta performance, sem menções a etapas de fundação inicial

### Requirement: Criptografia em Trânsito e Proteção contra Packet Sniffing na Rede Local (TLS/HTTPS)
O sistema SHALL disponibilizar suporte a conexões seguras criptografadas via TLS 1.3 e TLS 1.2 com cifras de alta performance (AES-GCM / ChaCha20-Poly1305) e suporte a HTTP/2, permitindo inicialização com certificado TLS autoassinado gerado automaticamente em memória (zero configuração) ou com certificados e chaves privadas informadas pelo usuário através de flags da CLI (`--tls`/`-s`, `--tls-cert`, `--tls-key`).
*(Visão PO: Garante privacidade e sigilo total dos dados na rede local, impedindo que ferramentas de sniffing capturem nomes de arquivos, parâmetros de URL ou conteúdos em trânsito, mantendo velocidade máxima de transferência. Visão QA: Valida handshake TLS 1.3/1.2, rejeição de conexões inseguras quando TLS estiver ativado, geração correta de certificado autoassinado para o host e funcionamento transparente de downloads e uploads em HTTPS).*

#### Scenario: Inicialização com TLS autoassinado automático para rede local
- **GIVEN** que o usuário inicia o servidor com a flag `--tls` ou `-s` sem informar caminhos de certificado
- **WHEN** o servidor for inicializado
- **THEN** o sistema deve gerar dinamicamente em memória um certificado X.509 autoassinado válido para localhost e endereços IP locais e escutar conexões exclusivamente via HTTPS

#### Scenario: Inicialização com certificado e chave TLS customizados
- **GIVEN** que os arquivos `cert.pem` e `key.pem` existem e são válidos
- **WHEN** o usuário executar o servidor passando `--tls-cert cert.pem --tls-key key.pem`
- **THEN** o servidor deve carregar o par de chaves e aceitar conexões HTTPS com o certificado configurado

#### Scenario: Criptografia e integridade total de downloads e uploads
- **GIVEN** que o servidor está rodando com TLS ativado
- **WHEN** um cliente realizar navegação, download de arquivo ou upload de dados via HTTPS
- **THEN** todo o payload, cabeçalhos e metadados devem ser transmitidos criptografados via TLS/HTTP2, tornando ilegível qualquer tentativa de captura de pacotes na rede local
