# API PRODUCT-MARKET

## Prerequisites

Before you begin, ensure you have met the following requirements:
* You have installed the latest version of [Go](https://go.dev/dl/)

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

### Generate Swag (Swagger):
Swag converts Go annotations to Swagger Documentation 2.0. We've created a variety of plugins for popular Go web frameworks. This allows you to quickly integrate with an existing Go project (using Swagger UI).

**Link:** https://github.com/swaggo/swag


Install Swagger:
```bash
 make lint-swag
```

##### command generate Swagger doc
```bash
 make swag
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

### Generate Mock Interface
This is an automatic mock generator using mockery, the first thing we must do is go to the path of the file that we want to autogenerate:

**Link:** https://github.com/vektra/mockery

Install the library
```bash
 make mockery-install
```
#### Command:
``//go: mockery --outpkg=_mocks --all --case=underscore --output=banner``

Generate all mocks with expected feat include, see official documentation. This line would go to the beginning of the test file

## Execute Test
```bash
 make go-test
```

## Run application

prior to execution you must include the following environment variables
```
DB_HOST=
DB_DRIVER=
DB_USER=
DB_PASSWORD=
DB_NAME=back_anime
DB_PORT=
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
