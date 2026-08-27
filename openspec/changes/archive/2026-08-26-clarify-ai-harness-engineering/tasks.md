## 1. Atualização do Prompt Mestre de Fundação (docs/foundation-spec-prompt.md)

- [x] 1.1 Atualizar o Pilar 6 no prompt mestre de `docs/foundation-spec-prompt.md`, dissociando explicitamente o conceito de Harness Engineering do Makefile (coberto pelo Pilar 4) e consolidando-o como arcabouço e scaffolding de Engenharia de IA (boas práticas, governança de contexto, guardrails e operação do Antigravity CLI).
- [x] 1.2 Atualizar o checklist dos artefatos e referências textuais em `docs/foundation-spec-prompt.md` para alinhar as definições de Harness, Loop e Graph Engineering.

## 2. Atualização das Regras de Governança de IA (.agent/rules/)

- [x] 2.1 Refatorar `.agent/rules/agent_harness_engineering.md` para redefinir a seção 1 (Harness Engineering) em torno do scaffolding de regras, guardrails operacionais, matriz de permissões (`.agent/settings.json`) e limites de sandbox para o Antigravity CLI.
- [x] 2.2 Harmonizar as diretrizes operacionais de ferramentas nativas e economia de tokens mantendo coerência com a nova definição de Harness de IA.

## 3. Atualização das Especificações e Documentação do Projeto

- [x] 3.1 Sincronizar o requisito `Requirement: Estratégias do Antigravity CLI (Harness, Loop e Graph Engineering) e Otimização de Tokens` e seu cenário BDD em `openspec/specs/project-foundation/spec.md`.
- [x] 3.2 Harmonizar referências textuais em `README.md`, `AGENTS.md` e `GEMINI.md` para eliminar conflitos terminológicos entre automação de compilação e harness de IA.

## 4. Validação e Quality Gate

- [x] 4.1 Validar a integridade de todas as especificações e esquemas OpenSpec executando `openspec validate --all`.
- [x] 4.2 Executar a esteira local de verificação de qualidade (`make check`) com suíte de testes automatizados e validação de cobertura (>= 80%).
