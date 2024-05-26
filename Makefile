deps-dl:
	@echo "Installing all dependencies"
	@go mod download all

deps-upd:
	@echo "Updating all dependencies"
	@go get -u ./...
	@go mod download all
	@go mod tidy

wire-dl:
	@echo "Installing wire"
	@go install github.com/google/wire/cmd/wire@latest

wire:
	@echo "Running wire"
	@wire ./...

swaggo-dl:
	@echo "Installing swag"
	@go install github.com/swaggo/swag/cmd/swag@latest

swaggo:
	@echo "Generating swagger docs"
	@swag init -g internal/router/router.go -o api/docs

mockgen-dl:
	@echo "Installing mockgen"
	@go install github.com/golang/mock/mockgen@latest

mockgen:
	@echo "Deleting prev mocks"
	@rm -rf ./test/mock
	@echo "Generating mocks"
	@go generate ./...

mockery-dl:
	@echo "Installing mockery"
	@go install github.com/vektra/mockery/v2@latest

mockery:
	@echo "Deleting prev mocks"
	@rm -rf ./test/mocks
	@echo "Generating mocks"
	@go run github.com/vektra/mockery/v2@latest --keeptree --dir=internal/app --output=test/mocks --all --case=snake

all-dl: deps-dl swaggo-dl mockgen-dl mockery-dl lint-dl
	@echo "Downloaded and install all prerequisites"

build-all: swaggo mockgen mockery
	@echo "Building app"
	@go build -v ./cmd/app/main.go

build:
	@echo "Formatting"
	@go fmt ./...
	@echo "Building app"
	@go build -v ./cmd/app/main.go

tests-short:
	@echo "Running short tests"
	@go test -v -short -cover ./...

tests-all:
	@echo "Running all tests"
	@go test -v -cover ./...

lint-dl:
	@echo "Installing golangci-lint"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint:
	@echo "Running lint with config"
	@golangci-lint run -v ./...

lint-full:
	@echo "Running lint enable all"
	@golangci-lint run -v --no-config --timeout 5m --enable-all --disable gci --allow-parallel-runners --skip-dirs test --out-format colored-line-number ./...> .reports/lint.txt

govuln:
	@echo "Installing govulncheck"
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@echo "Running govulncheck"
	@govulncheck ./...

gokart:
	@echo "Installing gokart"
	@go install github.com/praetorian-inc/gokart@latest
	@echo "Running gokart"
	@gokart scan .
