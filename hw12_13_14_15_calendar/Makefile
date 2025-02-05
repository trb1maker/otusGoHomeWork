BIN := "./bin/calendar"
SCHEDULER := "./bin/calendar_scheduler"
SENDER := "./bin/calendar_sender"

DOCKER_IMG="calendar:develop"

POSTGRES_HOST="192.168.0.103"
POSTGRES_USER="otus"
POSTGRES_PASSWORD="20240518_otus"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/scheduler
	go build -v -o $(SENDER) -ldflags "$(LDFLAGS)" ./cmd/sender

run: build
	$(BIN) -config ./configs/calendar_config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

migrations_up:
	GOOSE_MIGRATION_DIR=./migrations goose postgres "host=$(POSTGRES_HOST) port=5432 dbname=otus user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD)" up

migrations_down:
	GOOSE_MIGRATION_DIR=./migrations goose postgres "host=$(POSTGRES_HOST) port=5432 dbname=otus user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD)" down

test:
	go test -v -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.55.2

lint: install-lint-deps
	golangci-lint run ./...

generate:
	go generate ./...

.PHONY: build run build-img run-img version test lint
