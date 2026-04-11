<!--
Sync Impact Report
Version change: 1.0.0 -> 1.1.0
Modified principles:
- V. Entrega Incremental com Qualidade -> V. Entrega Incremental com Qualidade e Documentacao
Added sections:
- Nenhuma
Removed sections:
- Nenhuma
Templates requiring updates:
- ✅ .specify/templates/plan-template.md
- ✅ .specify/templates/tasks-template.md
- ✅ .specify/templates/spec-template.md (validada; sem mudancas necessarias)
- ✅ .github/prompts/speckit.constitution.prompt.md (validada; sem mudancas necessarias)
Follow-up TODOs:
- Nenhum (README.md na raiz criado e validado).
-->

# Pic2Video Constitution

## Core Principles

### I. Clean Code Inegociavel
Todo codigo novo ou alterado MUST ser legivel, coeso e orientado a responsabilidade unica.
Funcoes e classes MUST ter nomes explicitos e tamanho reduzido; duplicacao MUST ser removida
de forma ativa. Revisoes MUST rejeitar implementacoes que escondam complexidade acidental.
Rationale: codigo limpo reduz custo de manutencao e evita regressao em projetos pequenos.

### II. Simplicidade e Escopo Pequeno
Cada entrega MUST resolver um problema pequeno e claro, sem antecipar funcionalidades
futuras sem demanda real. Solucoes MUST preferir a alternativa mais simples que atenda
os requisitos atuais, e toda complexidade extra MUST ser justificada no plano.
Rationale: simplicidade acelera entrega e reduz risco de arquitetura desnecessaria.

### III. Testes Unitarios Sempre
Toda mudanca de comportamento MUST incluir testes unitarios cobrindo caminho feliz,
validacoes de entrada e falhas relevantes. O fluxo MUST seguir red-green-refactor:
escrever teste, observar falha, implementar, e refatorar com seguranca.
Rationale: testes unitarios garantem confiabilidade local e permitem evolucao segura.

### IV. Testes E2E Sempre
Toda funcionalidade entregue ao usuario MUST incluir pelo menos um teste end-to-end
que valide a jornada principal em condicoes reais do sistema. Mudancas em fluxos
criticos MUST atualizar a suite E2E correspondente antes do merge.
Rationale: testes E2E previnem regressao de integracao e confirmam valor percebido.

### V. Entrega Incremental com Qualidade e Documentacao
Cada historia MUST ser implementavel e validavel de forma independente, preservando
MVP funcional ao final de cada incremento. Merge so pode ocorrer quando lint, testes
unitarios e testes E2E estiverem verdes no pipeline local/CI. Toda mudanca que altere
uso, fluxo operacional, instalacao, build, ou execucao MUST atualizar README.md na raiz.
Rationale: incrementos pequenos com gates objetivos reduzem risco e retrabalho, e
documentacao obrigatoria reduz friccao operacional.

## Padroes Tecnicos Obrigatorios

- O projeto MUST manter arquitetura enxuta com separacao minima entre dominio,
	aplicacao e infraestrutura quando aplicavel.
- Dependencias MUST ser poucas e justificadas; preferir bibliotecas maduras e
	evitar frameworks pesados para problemas simples.
- Convencoes de nomenclatura, testes e estrutura de pastas MUST ser consistentes
	em todo o repositorio.
- README.md na raiz MUST existir e refletir o estado atual de uso, build e testes.

## Fluxo de Trabalho e Quality Gates

- Especificacao MUST definir escopo pequeno, criterios de sucesso mensuraveis e
	plano de testes unitarios + E2E para cada historia.
- Plano de implementacao MUST passar no Constitution Check antes da execucao.
- Lista de tarefas MUST incluir tarefas explicitas de testes unitarios e E2E,
	com execucao antes da implementacao de codigo funcional.
- Lista de tarefas MUST incluir tarefa de atualizacao de README.md quando houver
	impacto em uso, setup, build, run, flags, ou troubleshooting.
- Apos implementacao bem-sucedida, o fluxo MUST executar `make build-all`
	e registrar o resultado antes de considerar a entrega concluida.
- Pull requests MUST documentar quais principios foram atendidos e quais testes
	comprovam conformidade, incluindo confirmacao explicita de revisao do README.md.

## Governance

Esta Constituicao prevalece sobre praticas locais conflitantes. Mudancas nesta
Constituicao MUST ser propostas via pull request com justificativa, impacto em
templates e plano de migracao quando aplicavel.

Politica de versao constitucional (SemVer):
- MAJOR: remocao ou redefinicao incompativel de principios ou regras obrigatorias.
- MINOR: adicao de principio, secao, ou expansao material de obrigacoes.
- PATCH: clarificacoes editoriais sem alterar obrigacoes normativas.

Revisao de conformidade MUST ocorrer em tres pontos: durante especificacao,
durante planejamento e durante revisao de PR final. Nao conformidades MUST bloquear
aprovacao ate adequacao ou excecao formal documentada.

**Version**: 1.2.0 | **Ratified**: 2026-04-09 | **Last Amended**: 2026-04-09
