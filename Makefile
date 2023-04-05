.PHONY: modd
modd:
	go mod download
	go mod tidy

.PHONY: go-test
go-test:
	go test -coverprofile=coverage.out -v ./...
	go tool cover -func coverage.out

go-test-report:
	go tool cover -html=coverage.out

.PHONY: lint
lint:
	golangci-lint -v run

lint-install:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.51.2

lint-swag:
	go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: swag
swag:
	swag init -g cmd/api/main.go --tags Enrollment,Health
	copy .\docs\swagger.json .\docs\api-data-enrollment-openapi.json
	copy .\docs\swagger.yaml .\docs\api-data-enrollment-openapi.yaml

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

.PHONY: build
build:
	go build ./...

.PHONY: build-apple-silicon
apple-silicon:
	go build -tags dynamic ./...
