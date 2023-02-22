BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

HOME_PATH := $(shell pwd)
PATH := $(PATH)

DOCKER_COMPOSE_FILE="docker-compose.yaml"
DOCKER_COMPOSE_DEV_FILE="./deployments/docker-compose-dev.yaml"
DOCKER_ENV_FILE="./deployments/.env.dev"
DOCKER_CMD="docker-compose --file=./docker-compose.yaml --env-file ./deployments/.env.dev "
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
	docker-compose --file=$(DOCKER_COMPOSE_FILE) --env-file $(DOCKER_ENV_FILE) up -d app-http app-scheduler app-sender

up-no-cache:
	docker-compose --file=$(DOCKER_COMPOSE_FILE) --env-file $(DOCKER_ENV_FILE) build --no-cache pgs rabbitmq app-http app-scheduler app-sender
	docker-compose --file=$(DOCKER_COMPOSE_FILE) --env-file $(DOCKER_ENV_FILE) up -d app-http app-scheduler app-sender

down:
	docker-compose --file=$(DOCKER_COMPOSE_FILE) --env-file $(DOCKER_ENV_FILE) down

recreate:
	docker-compose --file=$(DOCKER_COMPOSE_FILE) rm -f
	docker-compose --file=$(DOCKER_COMPOSE_FILE) pull
	docker-compose --file=$(DOCKER_COMPOSE_FILE) build --no-cache
	docker-compose --file=$(DOCKER_COMPOSE_FILE) up --build -d

ps:
	docker-compose --file=$(DOCKER_COMPOSE_FILE) --env-file $(DOCKER_ENV_FILE) ps

ps-db:
	docker-compose --file=$(DOCKER_COMPOSE_FILE) --env-file $(DOCKER_ENV_FILE) exec pgs psql -U calendar -d calendar -c "SELECT * FROM events;"
	docker-compose --file=$(DOCKER_COMPOSE_FILE) --env-file $(DOCKER_ENV_FILE) exec pgs psql -U calendar -d calendar -c "SELECT * FROM users;"

generate:
	docker-compose --file=$(DOCKER_COMPOSE_FILE) exec otus-env go generate

go-fix:
	docker-compose --file=$(DOCKER_COMPOSE_DEV_FILE) run ubuntu gci write /app
	docker-compose --file=$(DOCKER_COMPOSE_DEV_FILE) run ubuntu gofumpt -l -w /app

integration-tests:
	set -e ;\
	test_status_code=0 ;\
	docker-compose  --file=$(DOCKER_COMPOSE_FILE) --env-file $(DOCKER_ENV_FILE) run integration-tests go test || test_status_code=$$? ;\
	docker-compose  --file=$(DOCKER_COMPOSE_FILE) --env-file $(DOCKER_ENV_FILE) down ;\
	exit $$test_status_code ;

integration-tests-no-cache:
	set -e ;\
	make up-no-cache
	docker-compose --file=$(DOCKER_COMPOSE_FILE) --env-file $(DOCKER_ENV_FILE) build --no-cache integration-tests
	test_status_code=0 ;\
	docker-compose  --file=$(DOCKER_COMPOSE_FILE) --env-file $(DOCKER_ENV_FILE) run integration-tests go test || test_status_code=$$? ;\
	docker-compose  --file=$(DOCKER_COMPOSE_FILE) --env-file $(DOCKER_ENV_FILE) down ;\
	exit $$test_status_code ;
clear:
	docker system prune --volumes
	docker image prune
	docker container prune

.PHONY: build run build-img run-img version test lint