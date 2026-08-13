#!/usr/bin/env bash
set -o pipefail

diff=""
if [ "$1" != "" ]
then
    diff="$1"
fi

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


go_files_to_test=$(git diff --name-status "$diff" |  grep -E '^(A|M|R)'  | awk '{ print $NF }' | grep '\.go$')

# Extract unique directories from the changed Go files
packages=$(echo "$go_files_to_test" | xargs -n1 dirname | sort -u)

exit_code=0

if [ -z "$packages" ]; then
    echo "No Go files to check"
    exit 0
fi

for pkg in $packages
do
    echo "Checking security for package ./$pkg..."
    if ! gosec "./$pkg" 2>&1 | printcolor; then
        exit_code=1
    fi
done

set +o pipefail

exit $exit_code