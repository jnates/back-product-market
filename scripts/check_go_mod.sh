#!/usr/bin/env bash
echo "Executing check go mod..."
# Definir colores si están disponibles
SCRIPT_DIR="./scripts"
if [[ -f "$SCRIPT_DIR/use_colors.sh" ]]; then
  . "$SCRIPT_DIR/use_colors.sh"
else
  purple() { echo "$@"; }
  green() { echo "$@"; }
  yellow() { echo "$@"; }
  red() { echo "$@"; }
  printcolor() { cat; }
fi

exit_code=0

go mod tidy
git diff --exit-code -- go.mod go.sum | fileToRedMsg

if [ $? -ne 0 ]; then
  exit_code=1
fi

set +o pipefail

exit $exit_code