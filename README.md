# API PRODUCT-MARKET

## Prerequisites

Before you begin, ensure you have met the following requirements:
* You have installed the latest version of [Go](https://go.dev/dl/)
* You have `make` (GNU Make) available on your `PATH`. On Windows it is not installed by default:
```bash
 winget install ezwinports.make
```

## Installation

* Install dependencies:
```bash
$ make modd
```

### Install linter:
```bash
 make lint-install
```
##### or build apple-silicon
```bash
 make lint-install-apple-silicon
```

##### command to validate the version
```bash
 golangci-lint --version
```

### Formatter of code
#### Install formatter:
```bash
 make code-format-install
```
#### Execute formatter:
```bash
 make code-format
```

Generate all mocks with expected feat include, see official documentation. This line would go to the beginning of the test file

## API Documentation (Swagger)

Install the `swag` CLI once, then (re)generate the docs after adding or changing endpoints:
```bash
 make install-swag
 make swag
```

This regenerates `docs/docs.go`, `docs/swagger.json` and `docs/swagger.yaml` from the `@`-annotations in the handlers and `cmd/main.go`. With the app running, the Swagger UI is served at:
```
http://localhost:<API_PORT>/api/market/docs/index.html
```

## Execute Test
```bash
 make coverage
```

## Run application

prior to execution you must include the following environment variables
```
DB_HOST=
DB_USER=
DB_PASSWORD=
DB_NAME=back_anime
DB_PORT=
DB_SSL_MODE=disable
API_PORT=
LOGGER_DEBUG=
SECRET_KEY=
```

### Execute go build
```bash
 make build
```

##### or build apple-silicon
```bash
 make build-apple-silicon
```
