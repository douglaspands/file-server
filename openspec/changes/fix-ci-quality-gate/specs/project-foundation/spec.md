## MODIFIED Requirements

### Requirement: Governança de Git, GitHub Actions, Integração Contínua e Arquivamento
O sistema SHALL aplicar convenções de Conventional Commits, ganchos de pre-commit, workflows do GitHub Actions com toolchain Go e linters perfeitamente alinhados para validação contínua, e procedimento mandatório de commit Git na conclusão/arquivamento de qualquer especificação OpenSpec.
*(Visão PO: Garante histórico de mudanças limpo, rastreabilidade de valor de negócio, esteira de CI sem falsos negativos e persistência garantida do código após a conclusão de uma spec. Visão QA: Funciona como Quality Gate automatizado determinístico, garantindo que o compilador, o linter estrito e a suíte de testes executem de forma idêntica localmente e no CI).*

#### Scenario: Validação automática em Pull Request via GitHub Actions
- **WHEN** um Pull Request ou push for enviado ao repositório GitHub
- **THEN** o workflow de CI deve sincronizar a versão do Go a partir do `go.mod`, executar o `golangci-lint` sem conflitos de versão do compilador, auditar vulnerabilidades via `govulncheck`, rodar os testes automatizados com validação de cobertura >= 80% e verificar a integridade do OpenSpec

#### Scenario: Compatibilidade e integridade do linter no Quality Gate
- **WHEN** a etapa de linting (`Run GolangCI-Lint`) for executada no pipeline de CI
- **THEN** a versão de compilação do `golangci-lint` deve ser estritamente igual ou superior à versão alvo declarada em `go.mod`, executando todas as checagens estáticas com sucesso e código de saída zero

#### Scenario: Validação local antes do commit (Pre-commit hook)
- **WHEN** um desenvolvedor tentar efetuar um commit local
- **THEN** os hooks configurados devem verificar a formatação do código e a mensagem de commit de acordo com a especificação Conventional Commits

#### Scenario: Commit Git obrigatório no arquivamento da especificação
- **WHEN** o comando de arquivamento de especificação (`openspec archive` ou `/openspec-archive-change`) for executado após a conclusão da implementação
- **THEN** as orientações operacionais (`operations.archive.guidance`) devem instruir o agente a criar o commit com Conventional Commits contendo as alterações e a sincronização da spec principal
