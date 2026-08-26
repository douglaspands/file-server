## 1. Dependências e Infraestrutura de Autenticação e Chaves

- [x] 1.1 Adicionar dependências Go para SFTP (`golang.org/x/crypto/ssh`, `github.com/pkg/sftp`) e FTP/FTPS (`github.com/fclairamb/ftpserverlib`) no `go.mod` e `go.sum`
- [x] 1.2 Implementar gerador de chaves SSH em memória e utilitários de credenciais temporárias seguras em `internal/core/services/`
- [x] 1.3 Criar testes unitários para geradores de chaves e utilitários de credenciais em `internal/core/services/` com cobertura >= 80%

## 2. Implementação do Adaptador do Servidor SFTP

- [x] 2.1 Criar o handler de sistema de arquivos com sandbox estrito e controle de somente leitura para SFTP em `internal/adapters/sftp/`
- [x] 2.2 Implementar o servidor SSH/SFTP com suporte a autenticação por senha, chave pública e chave privada de host em `internal/adapters/sftp/`
- [x] 2.3 Implementar encerramento gracioso via contexto e canal de erro no adaptador SFTP
- [x] 2.4 Criar testes unitários e de integração para o servidor SFTP (autenticação, operações de arquivo, sandbox e flag read-only) com cobertura >= 80%

## 3. Implementação do Adaptador do Servidor FTP / FTPS

- [x] 3.1 Criar o driver de sistema de arquivos com sandbox estrito e modo somente leitura para FTP em `internal/adapters/ftp/`
- [x] 3.2 Implementar o adaptador de servidor FTP/FTPS com suporte a conexões seguras TLS (autoassinado e certificados customizados) em `internal/adapters/ftp/`
- [x] 3.3 Implementar encerramento gracioso via contexto e controle de conexões ativas no adaptador FTP
- [x] 3.4 Criar testes unitários e de integração para o servidor FTP/FTPS (handshake FTPS, modo passivo, autenticação, sandbox e flag read-only) com cobertura >= 80%

## 4. Integração dos Comandos CLI e Banners Informativos

- [x] 4.1 Implementar o comando `cmd/sftp.go` (`file-server sftp`) com suporte a resolução de diretório, expansão de til, flags de configuração e ciclo de vida gracioso
- [x] 4.2 Implementar o comando `cmd/ftp.go` (`file-server ftp`) com suporte a resolução de diretório, flags de configuração TLS, modo passivo e ciclo de vida gracioso
- [x] 4.3 Implementar banners de inicialização formatados exibindo endereços IP locais e da rede (LAN) e credenciais ativas para SFTP e FTP
- [x] 4.4 Criar testes unitários para a CLI em `cmd/sftp_test.go` e `cmd/ftp_test.go` com cobertura >= 80%

## 5. Documentação e Validação Final

- [x] 5.1 Atualizar a documentação do projeto (`README.md`) com os novos comandos CLI `sftp` e `ftp`, exemplos de conexão via clientes populares e guias de segurança
- [x] 5.2 Executar a suíte completa de testes e análise estática via `make test` e `make check`, garantindo cobertura &ge; 80% e zero advertências de lint
