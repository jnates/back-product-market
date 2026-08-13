.PHONY: modd
modd:
	go mod download
	go mod tidy

.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: lint
lint:
	$(subst \,/,$(GOPATH))/bin/golangci-lint.exe -v run

lint-install:
	curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(GOPATH)/bin v2.12.2


lint-install-apple-silicon:
	brew install golangci-lint
	brew upgrade golangci-lint

code-format-install:
	go install golang.org/x/tools/cmd/goimports@latest

.PHONY: code-format
code-format:
	goimports -l -w .
	gofmt -l -w .

mockery-install:
	go install github.com/vektra/mockery/v2@latest


generate-mocks:
	mockery --all --recursive --case=underscore --output=./mocks

.PHONY: generate-mocks

install-swag:
	go install -v github.com/swaggo/swag/cmd/swag@latest

.PHONY: swag
swag:
	swag init -g cmd/api/main.go --parseDependency --parseInternal
	cp ./docs/swagger.json ./docs/product-market-openapi.json
	cp ./docs/swagger.yaml ./docs/product-market-openapi.yaml

.PHONY: build
build:
	go build ./...

.PHONY: build-apple-silicon
apple-silicon:
	go build -tags dynamic ./...
