.PHONY: generate
.DEFAULT_GOAL := build

# generate
generate:
	go mod tidy
	go get github.com/google/wire/cmd/wire@latest
	go generate ./...
	go mod tidy

.PHONY: all
# generate all
all:
	make generate;

all_tests:
	go test -v ./internal/v1/http/handler/... ./internal/v1/biz/... -bench=. -cover  -coverprofile=coverage.out -benchmem -cpu=1,2,3,4 -timeout=500ms


bench_tests:
	go test -v ./internal/v1/http/handler/... ./internal/v1/biz/... -bench=. -benchmem -cpu=1,2,3,4 -timeout=100ms

unit_tests:
	go test -v ./internal/v1/http/handler/... ./internal/v1/biz/...

coverage_tests:
	go test -v ./internal/v1/http/handler/... ./internal/v1/biz/... -cover  -coverprofile=coverage.out

fmt:
	gofumpt -l -w .;gci write ./

devtools:
	@echo "Installing devtools"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install mvdan.cc/gofumpt@latest
	go install go.uber.org/mock/mockgen@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/daixiang0/gci@v0.11.2

check:
	golangci-lint run --config .golangci.yml

build:
	docker build -t buildf .

