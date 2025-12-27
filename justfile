#go-cleanstack justfile

set shell := ["bash", '-c']

[doc('Start the server in development mode')]
dev:
    go run . serve

[doc('API Code Generation')]
[group('generate')]
generate-api:
    cd internal/infra/api && buf generate

[doc('unit tests')]
[group('test')]
test:
    gotestsum -- ./...

[doc('integration tests')]
[group('test')]
test-int:
    gotestsum -- -tags=integration ./tests/integration/...

[doc('end-to-end tests')]
[group('test')]
test-e2e:
    gotestsum -- -tags=e2e ./tests/e2e/...

[doc('all tests')]
[group('test')]
test-all:
    gotestsum -- -tags=integration,e2e ./...

[doc('coverage tests')]
[group('test')]
test-cover:
    gotestsum -- -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html

[doc('Migrations')]
[group('migrate')]
migrate-up:
    go run ./main.go migrate up

[group('migrate')]
migrate-down:
    go run ./main.go migrate down

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

[group('docker')]
up:
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
