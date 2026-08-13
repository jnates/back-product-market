#!/usr/bin/env bash
set -o pipefail

# Ruta del script
SCRIPT_DIR="./scripts"
[[ -f "$SCRIPT_DIR/use_colors.sh" ]] && . "$SCRIPT_DIR/use_colors.sh"

# Entornos y variables de entorno
export AWS_XRAY_SDK_DISABLED=true
export LOGS_EXTERNAL_ENABLED=false

# Elegir grep (si existe ggrep en macOS)
if command -v ggrep >/dev/null; then
  GREP_COMMAND="ggrep"
else
  GREP_COMMAND="grep"
fi

HAVE_BC=false
if command -v bc >/dev/null 2>&1; then
    HAVE_BC=true
    green "✅ Comando 'bc' detectado. Se usará para comparaciones de cobertura."
else
    yellow "⚠️ Comando 'bc' no encontrado. Se usará 'awk' como alternativa."
fi

# Validar comandos necesarios
tools=(go awk "$GREP_COMMAND" gofmt)
for cmd in "${tools[@]}"; do
  command -v "$cmd" >/dev/null || { echo "❌ Missing required command: $cmd"; exit 1; }
done

# Rango de diff (ej: "origin/development...HEAD")
diff_range="${1:-}"

# Obtener paquetes afectados, excluyendo carpetas 'mocks', 'interfaces', 'models', 'enums' y 'proto'
pkgs_to_test=$(git diff --name-status "$diff_range" \
  | $GREP_COMMAND -E '^(A|M|R)' \
  | awk '{ print $NF }' \
  | $GREP_COMMAND '\.go$' \
  | xargs -r -n1 dirname \
  | sort -u \
  | $GREP_COMMAND -Ev '(^|/)mocks(/|$)' \
  | $GREP_COMMAND -Ev '(^|/)(interfaces|models|model|enums|constant|constants|docs|configs|proto)(/|$)'
)

exit_code=0
for pkg in $pkgs_to_test; do
  [[ ! -d "$pkg" ]] && continue
  green "📦 Package to test: $pkg"

  # Archivos .go en el paquete, excluyendo main.go
  gofiles=$(find "$pkg" -maxdepth 1 -type f -name '*.go' ! -name 'main.go')
  [[ -z "$gofiles" ]] && continue

  coverage_dir=".coverage/${pkg}"
  coverage_file="$coverage_dir/coverage.out"
  mkdir -p "$coverage_dir"

  # Configuración de prueba
  timeout="160s"
  min_cov=80

  # Ejecutar tests con coverage
  go test -shuffle=on -cover -coverprofile="$coverage_file" -timeout="$timeout" -short "./$pkg"
  [[ ! -f "$coverage_file" ]] && continue

  # Extraer métricas de cobertura, omitiendo total, main, init y SC
  coverage_output=$(go tool cover -func="$coverage_file" | \
    $GREP_COMMAND -Ev '^(total:|init|SC)'
  )
  if [[ -z "$coverage_output" ]]; then
    # Sin información de cobertura: listar funciones sin tests
    funcs=$($GREP_COMMAND -hE '^func\s' $gofiles | $GREP_COMMAND -Ev '^func (Test|Benchmark)')
    while IFS= read -r fn; do
      red "🚨 Please define tests for: $fn"
      exit_code=1
    done <<< "$funcs"
  else
    # Revisar cobertura de cada función
    while IFS= read -r line; do
      # Obtener valor de cobertura (último campo) y quitar '%'
      p=$(echo "$line" | awk '{ print $NF }' | tr -d '%')
      # Comparar con mínimo
      is_below_coverage=0 # 0 para falso (no está por debajo), 1 para verdadero

      if [ "$HAVE_BC" = true ]; then
        # Usar 'bc' si está disponible
        if (( $(echo "$p < $min_cov" | bc) )); then
          is_below_coverage=1
        fi
      else
        # Usar 'awk' como alternativa
        if [ "$(awk -v p="$p" -v min="$min_cov" 'BEGIN { print (p < min) }')" -eq 1 ]; then
          is_below_coverage=1
        fi
      fi

      # Evaluar el resultado
      if [ "$is_below_coverage" -eq 1 ]; then
        red "❌ $line"
        exit_code=1
      else
        green "✅ $line"
      fi
    done <<< "$coverage_output"
  fi
done

exit $exit_code