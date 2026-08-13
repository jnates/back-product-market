#!/usr/bin/env bash
set -o pipefail

# Ruta base del script
SCRIPT_DIR="./scripts"
[[ -f "$SCRIPT_DIR/use_colors.sh" ]] && . "$SCRIPT_DIR/use_colors.sh"

diff_range="${1:-}"
exit_code=0

purple "🔍 Executing check go error..."

# Buscar paquetes con archivos .go cambiados (excluyendo cualquier carpeta 'mocks')
pkgs_to_test=$(git diff --name-status "$diff_range" | \
  grep -E '^(A|M|R)' | \
  awk '{ print $NF }' | \
  grep '\.go$' | \
  xargs -r -n1 dirname | \
  sort -u | \
  # <-- filtramos rutas que contengan '/mocks' o terminen/empiecen con 'mocks'
  grep -Ev '(^|/)(mocks)(/|$)'
)

for pkg in $pkgs_to_test; do
  # Alternativa: si prefieres mantener el filtrado en el bucle
  # [[ "$pkg" == *"/mocks"* ]] && continue

  [[ ! -d "$pkg" ]] && continue

  green "📦 Package to test: $pkg"

  yellow "🧪 Checking errcheck..."
  errcheck -asserts "./$pkg" | fileToRedMsg
  [[ $? -ne 0 ]] && exit_code=1

  yellow "🕵️ Checking go-errorlint..."
  if find "$pkg" -name '*.go' | grep -q .; then
  go-errorlint ./$pkg/...  | fileToRedMsg
    [[ $? -ne 0 ]] && exit_code=1
  fi
done

set +o pipefail
exit $exit_code