#!/usr/bin/env bash
set -euo pipefail

MIN_COVERAGE=80.0
COVERAGE_FILE="coverage.out"
HTML_REPORT="coverage.html"

echo "🧪 Executando testes unitários e gerando perfil de cobertura..."
go test -coverprofile="${COVERAGE_FILE}" -covermode=atomic ./...

TOTAL_COVERAGE=$(go tool cover -func="${COVERAGE_FILE}" | grep total: | awk '{print $3}' | tr -d '%')
echo "📊 Cobertura Total: ${TOTAL_COVERAGE}% (Meta Mínima: ${MIN_COVERAGE}%)"

# Gera relatório HTML opcional
go tool cover -html="${COVERAGE_FILE}" -o "${HTML_REPORT}"
echo "📄 Relatório HTML gerado em: ${HTML_REPORT}"

# Validação com precisão decimal usando awk
IS_SUFFICIENT=$(awk -v total="${TOTAL_COVERAGE}" -v min="${MIN_COVERAGE}" 'BEGIN {print (total >= min)}')

if [ "${IS_SUFFICIENT}" -eq 1 ]; then
    echo "✅ Cobertura de testes APROVADA: ${TOTAL_COVERAGE}% >= ${MIN_COVERAGE}%"
    exit 0
else
    echo "❌ Cobertura de testes REPROVADA: ${TOTAL_COVERAGE}% é inferior ao mínimo exigido de ${MIN_COVERAGE}%"
    exit 1
fi
