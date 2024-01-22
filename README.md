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

## Execute Test
```bash
 make coverage
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
