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
DATABASE_MIGRATION := postgres://postgres:12345@host.docker.internal:5432/testdb?sslmode=disable
MIGRATE := docker run -v ./sql/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database ${DATABASE_MIGRATION}
CPU_PPROF_PRODUCER := cpu-${APP_NAME_PRODUCER}.pprof
CPU_PPROF_CONSUMER := cpu-${APP_NAME_CONSUMER}.pprof
HEAP_PPROF_PRODUCER := heap-${APP_NAME_PRODUCER}.pprof
HEAP_PPROF_CONSUMER := heap-${APP_NAME_CONSUMER}.pprof
ALLOC_PPROF_PRODUCER := alloc-${APP_NAME_PRODUCER}.pprof
ALLOC_PPROF_CONSUMER := alloc-${APP_NAME_CONSUMER}.pprof


IMPORT_PATH := github.com/WsDev69/producer-consumer/pkg/build

info:
	@echo "Version : ${VERSION}"
	@echo "Commit hash : ${COMMIT_HASH}"
	@echo "Build time : ${BUILD_TIME}"

LDFLAGS := -X '${IMPORT_PATH}.Version=$(VERSION)' -X '${IMPORT_PATH}.CommitHash=$(COMMIT_HASH)' -X '${IMPORT_PATH}.BuildTime=$(BUILD_TIME)' -s -w
PGO_MEM := -pgo=mem.prof
PGO_CPU := -pgo=cpu.prof


all: clean proto-generate dep generate fmt lint test_all

dep:
	go mod tidy
	go mod download

dep-update:
	go get -t -u ./...

dep-all: dep-update dep

build-producer:
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
	cat /tmp/coverage.out | grep -v ".pb.go" | grep -v "mock" > /tmp/coverage_cleaned.out
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
	$(MIGRATE) down $(v)

migrate-fix: ## Fix migrations
	$(MIGRATE) force $(v)

show-profile-consumer:
	@rm -rf ${CPU_PPROF_CONSUMER}
	@curl -o ${CPU_PPROF_CONSUMER} 'http://localhost:1377/debug/pprof/profile?seconds=30'
	go tool pprof -http=:8081 ${CPU_PPROF_CONSUMER}

show-profile-producer:
	@rm -rf ${CPU_PPROF_PRODUCER}
	@curl -o ${CPU_PPROF_PRODUCER} 'http://localhost:1378/debug/pprof/profile?seconds=30'
	go tool pprof -http=:8082 ${CPU_PPROF_PRODUCER}

show-heap-consumer:
	@rm -rf ${CPU_PPROF_CONSUMER}
	@curl -o ${CPU_PPROF_CONSUMER} 'http://localhost:1377/debug/pprof/heap'
	go tool pprof -http=:8081 ${CPU_PPROF_CONSUMER}

show-heap-producer:
	@rm -rf ${HEAP_PPROF_PRODUCER}
	@curl -o ${HEAP_PPROF_PRODUCER} 'http://localhost:1378/debug/pprof/heap'
	go tool pprof -http=:8082 ${HEAP_PPROF_PRODUCER}

show-alloc-consumer:
	@rm -rf ${ALLOC_PPROF_CONSUMER}
	@curl -o ${ALLOC_PPROF_CONSUMER} 'http://localhost:1377/debug/pprof/allocs'
	go tool pprof -http=:8081 ${ALLOC_PPROF_CONSUMER}

show-alloc-producer:
	@rm -rf ${ALLOC_PPROF_PRODUCER}
	@curl -o ${ALLOC_PPROF_PRODUCER} 'http://localhost:1378/debug/pprof/allocs'
	go tool pprof -http=:8082 ${ALLOC_PPROF_PRODUCER}

tools:
	go install github.com/vektra/mockery/v2@v2.46.0
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0