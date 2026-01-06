#go-cleanstack justfile

set shell := ["bash", '-uc']
set dotenv-load := true

GOTESTSUM_OPTIONS := "--format testdox"

[doc('Start the server in development mode')]
dev:
    go run . serve

[doc('API Code Generation')]
[group('generate')]
generate-api:
    cd internal/app/user/infra/api && buf generate

[doc('unit tests')]
[group('test')]
test:
    go tool gotestsum {{ GOTESTSUM_OPTIONS }} -- -v ./...

[doc('integration tests')]
[group('test')]
test-int:
    go tool gotestsum {{ GOTESTSUM_OPTIONS }} -- -v -tags=integration ./tests/integration/...

[doc('end-to-end tests')]
[group('test')]
test-e2e:
    go tool gotestsum {{ GOTESTSUM_OPTIONS }} -- -v -tags=e2e ./tests/e2e/...

[doc('all tests')]
[group('test')]
test-all: test
    go tool gotestsum {{ GOTESTSUM_OPTIONS }} -- -v -tags=integration,e2e ./...

[doc('coverage tests')]
[group('test')]
test-cover:
    go tool gotestsum {{ GOTESTSUM_OPTIONS }} -- -v -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html

[doc('Migrations Up')]
[group('database')]
migrate-up:
    go run . migrate up

[doc('Migrations Down')]
[group('database')]
migrate-down:
    go run . migrate down

[doc('Migrations Create (interactive)')]
[group('database')]
migrate-create:
    go run . migrate create

[doc('Connect to the Dockerized Database')]
[group('database')]
db-connect:
    docker exec -it cleanstack-db-${APP_ENV} psql -U user -d cleanstack_${APP_ENV}

[doc('Linting')]
[group('lint')]
lint:
    golangci-lint run

[group('lint')]
lint-fix:
    golangci-lint run --fix

[doc('Build the application')]
build:
    go build -o bin/cleanstack ./main.go

[doc('Start docker-compose for current APP_ENV')]
[group('docker')]
up:
    #!/usr/bin/env bash
    set -euo pipefail

    if [ -z "${APP_ENV:-}" ]; then
        echo "Error: APP_ENV is not set. Run './configure' first or 'source .envrc'"
        exit 1
    fi

    # Set ports based on APP_ENV
    case "$APP_ENV" in
        development)
            export APP_PORT=4224
            export DB_PORT=5435
            ;;
        staging)
            export APP_PORT=4225
            export DB_PORT=5436
            ;;
        production)
            export APP_PORT=4226
            export DB_PORT=5437
            ;;
        *)
            echo "Error: APP_ENV must be one of: development, staging, production"
            exit 1
            ;;
    esac

    echo "Starting docker-compose for APP_ENV=$APP_ENV (app:$APP_PORT, db:$DB_PORT)"
    docker-compose up -d

[group('docker')]
down:
    docker-compose down

[group('docker')]
logs:
    docker-compose logs -f

[doc('Clean')]
clean:
    rm -rf bin/ coverage.out coverage.html
