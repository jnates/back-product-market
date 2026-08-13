#!/usr/bin/env bash

set -o pipefail

# Rango de diff (ej: "origin/development...HEAD")
diff_range="${1:-}"

# Buscar archivos .go agregados/modificados,
# excluyendo rutas que contengan '/mocks' y archivos *_test.go
go_files_to_test=$(git diff --name-status "$diff_range" | \
  grep -E '^(A|M|R)' | \
  awk '{ print $NF }' | \
  grep -E '\.go$' | \
  # Excluir carpetas mocks
grep -Ev '(^|/)mocks(/|$)' | \
  # Excluir tests
grep -Ev '_test\.go$'
)

exit_code=0

for file in $go_files_to_test; do
    # Doble comprobación para asegurarse
    [[ "$file" == *"/mocks/"* ]] && continue
    [[ "$file" == *_test.go ]] && continue

    echo "Checking format for $file..."

    # Revisión de formato con goimports
    file_diff=$(goimports -d "$file")
    if [ -n "$file_diff" ]; then
        echo "❌ File $file is not well formatted"
        echo "$file_diff"
        exit_code=1
    fi
done

exit $exit_code
