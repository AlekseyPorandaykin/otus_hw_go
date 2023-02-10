BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

HOME_PATH := $(shell pwd)
PATH := $(PATH)

DOCKER_DIR="./deployments/docker-compose.yaml"
MIGRATE_SQL := $(shell cat < ./migrations/specification.sql;)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/config.toml

run-scheduler: build
	$(BIN) scheduler --config=./configs/scheduler.toml

run-sender: build
	$(BIN) sender --config=./configs/sender.toml

run-grpc: build
	$(BIN) grpc_server --config=./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

up:
	docker-compose --file=$(HOME_PATH)/deployments/docker-compose.yaml up -d

down:
	docker-compose --file=$(DOCKER_DIR) down

recreate:
	docker-compose --file=$(DOCKER_DIR) rm -f
	docker-compose --file=$(DOCKER_DIR) pull
	docker-compose --file=$(DOCKER_DIR) up --build -d

ps:
	docker-compose --file=$(DOCKER_DIR) ps

ps-db:
	docker-compose --file=$(DOCKER_DIR) exec pgs psql -U calendar -d calendar -c "SELECT * FROM events;"
	docker-compose --file=$(DOCKER_DIR) exec pgs psql -U calendar -d calendar -c "SELECT * FROM users;"

migrate:
	docker-compose --file=$(DOCKER_DIR) exec pgs psql -U calendar -d calendar -c "$(MIGRATE_SQL)"

generate:
	docker-compose --file=$(DOCKER_DIR) exec go16 go generate
go-fix:
	docker-compose --file=$(DOCKER_DIR) exec ubuntu gci write /app
	docker-compose --file=$(DOCKER_DIR) exec ubuntu gofumpt -l -w /app
.PHONY: build run build-img run-img version test lint