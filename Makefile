.PHONY: info

APP_NAME_PRODUCER := producer
APP_NAME_CONSUMER := consumer
BUILD_DIR := bin
#//VERSION := $(shell git describe --tags --always --dirty)
VERSION := 1.0.0
COMMIT_HASH := $(shell git rev-parse HEAD)
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

GO := go
GOLINT := golangci-lint


IMPORT_PATH := producer-consumer/pkg/build

info:
	@echo "Version : ${VERSION}"
	@echo "Commit hash : ${COMMIT_HASH}"
	@echo "Build time : ${BUILD_TIME}"

LDFLAGS := -X '${IMPORT_PATH}.Version=$(VERSION)' -X '${IMPORT_PATH}.CommitHash=$(COMMIT_HASH)' -X '${IMPORT_PATH}.BuildTime=$(BUILD_TIME)' -s -w


all: clean proto-generate dep generate fmt lint test_all

dep:
	go mod tidy
	go mod download

dep-update:
	go get -t -u ./...

dep-all: dep-update dep

build-producer: $(BUILD_DIR)/$(APP_NAME_PRODUCER)
$(BUILD_DIR)/$(APP_NAME_PRODUCER):
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME_PRODUCER) ./cmd/$(APP_NAME_PRODUCER)/main.go


build-consumer:
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME_CONSUMER) ./cmd/$(APP_NAME_CONSUMER)/main.go


clean:
	@rm -rf $(BUILD_DIR)

lint: dep
	$(GOLINT) run --timeout=5m -c .golangci.yml

test_all: test

test:
	go test -cover -race -count=1 -timeout=60s ./...

coverage-html: coverage-html-maker

coverage:  coverage-maker

coverage-html-maker:
	go test -tags=unit,integration -coverpkg=./... -covermode atomic -coverprofile=/tmp/coverage.out ./...
	cat /tmp/coverage.out | grep -v "postgres/db" | grep -v ".pb.go" | grep -v "mock" > /tmp/coverage_cleaned.out
	mv /tmp/coverage_cleaned.out /tmp/coverage.out
	go tool cover -html=/tmp/coverage.out

coverage-maker:
	go test -tags=unit,integration -coverpkg=./... -covermode atomic -coverprofile=/tmp/coverage.out ./...
	cat /tmp/coverage.out | grep -v ".pb.go" | grep -v "mock" > /tmp/coverage_cleaned.out
	mv /tmp/coverage_cleaned.out /tmp/coverage.out
	go tool cover -func=/tmp/coverage.out

fmt:
	go fmt producer-consumer/...

sqlc-generate:
	docker build -f sql/Dockerfile.sqlc -t sqlc-generator .
	docker run --rm -v `pwd`:/app sqlc-generator

proto-generate:
	docker run --platform linux/amd64 --rm -v `pwd`:/defs namely/protoc-all:1.51_2 -i common -f pkg/proto/task/common.proto -l go -o internal/handler/grpc/gen/task

generate:
	go generate ./...

migrate-create: ## Create migration file with name
	migrate create -ext sql -dir sql/migrations -seq -digits 10 $(name)

migrate-up: ## Run migrations
	$(MIGRATE) up

migrate-down: ## Rollback migrations
	$(MIGRATE) down

migrate-fix: ## Fix migrations
	$(MIGRATE) force $(v)

tools:
	go install github.com/vektra/mockery/v2@v2.46.0
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1