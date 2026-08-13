#!/usr/bin/env bash

if tty -s; then
    base_dir=$(dirname $0)
    . $base_dir/use_colors.sh
fi

# Si GOPATH no está definido, usar el valor por defecto de go env
if [[ -z "$GOPATH" ]]; then
    export GOPATH=$(go env GOPATH)
fi

cd /tmp

if [[ ! -x $GOPATH/bin/errcheck ]] || [[ $FORCE_UPDATE == "yes" ]]; then
    green "Installing errcheck"
    go install -v github.com/kisielk/errcheck@v1.10.0
fi

if [[ ! -x $GOPATH/bin/go-errorlint ]] || [[ $FORCE_UPDATE == "yes" ]]; then
    green "Installing go err lint"
    go install -v github.com/polyfloyd/go-errorlint@v1.8.0
fi

if [[ ! -x $GOPATH/bin/goimports ]] || [[ $FORCE_UPDATE == "yes" ]]; then
    green "Installing goimports"
    go install -v golang.org/x/tools/cmd/goimports@v0.42.0
fi

if [[ ! -x $GOPATH/bin/gosec ]] || [[ $FORCE_UPDATE == "yes" ]]; then
    green "Installing gosec"
    go install -v github.com/securego/gosec/v2/cmd/gosec@v2.22.0
fi
