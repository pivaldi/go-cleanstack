#go-cleanstack justfile

[doc('Start the server in development mode')]
dev:
    go run main.go serve

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
    go run ./cmd migrate up

[group('migrate')]
migrate-down:
    go run ./cmd migrate down

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

[doc('install/update code automation for development (prettier, pre-commit, goreturns, lintpack, gocritic, golangci-lint)')]
install-dev: install-dep
    curl https://pre-commit.com/install-local.py | python3 -
    go get github.com/sqs/goreturns
    go get github.com/go-lintpack/lintpack/...
    go get github.com/go-critic/go-critic/...
    go get github.com/golangci/golangci-lint/cmd/golangci-lint

[doc('install/update binary dependencies (buf, protoc-gen-go)')]
install-dep:
    go install github.com/direnv/direnv
    go install github.com/bufbuild/buf/cmd/buf@latest
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

[doc('setup/update pre-commit hooks (optional)')]
setup-dev:
    pre-commit install --install-hooks

[doc('Clean')]
clean:
    rm -rf bin/ coverage.out coverage.html
