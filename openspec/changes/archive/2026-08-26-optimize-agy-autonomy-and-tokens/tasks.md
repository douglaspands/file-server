## 1. Configuração de Permissões e Autonomia no Repositório

- [x] 1.1 Atualizar `.agent/settings.json` com permissões automáticas (`allow`) expandidas para comandos e ferramentas de ciclo de vida (`go`, `make`, `git`, `openspec`, `golangci-lint`, `govulncheck`, scripts).
- [x] 1.2 Ajustar filtros de segurança (`ask` e `deny`) mantendo proteção para operações destrutivas ou remotas sem interromper o fluxo local.

## 2. Diretrizes de Ferramentas Nativas e Economia de Tokens

- [x] 2.1 Criar a regra `.agent/rules/agent_tooling_autonomy.md` estabelecendo a prioridade obrigatória de ferramentas nativas (`write_to_file`, `replace_file_content`, `view_file`, `grep_search`, `find_by_name`, `list_dir`) e a proibição estrita de comandos bash para manipulação de arquivos.
- [x] 2.2 Atualizar `.agent/rules/agent_harness_engineering.md` incorporando diretrizes de economia ativa de tokens, saídas concisas e leitura/edição cirúrgica de arquivos.
- [x] 2.3 Atualizar `AGENTS.md` e `GEMINI.md` com as diretrizes consolidadas de autonomia operacional, prioridade de ferramentas nativas e economia de tokens.

## 3. Correção de Configurações do OpenSpec

- [x] 3.1 Ajustar a estrutura de `openspec/config.yaml` corrigindo o formato das seções de orientação (`guidance`) para eliminar avisos de validação.

## 4. Validação e Qualidade

- [x] 4.1 Validar a integridade de todas as especificações e do repositório através de `openspec validate --all`.
- [x] 4.2 Executar a suíte de testes e verificação completa (`make check`) garantindo cobertura >= 80% e ausência de regressões.
