CALENDAR_BIN := "./bin/calendar"
SCHEDULER_BIN := "./bin/scheduler"
SENDER_BIN := "./bin/sender"
CALENDAR_DOCKER_IMG="calendar:develop"
SCHEDULER_DOCKER_IMG="calendar_scheduler:develop"
SENDER_DOCKER_IMG="calendar_sender:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(CALENDAR_BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(SCHEDULER_BIN) -ldflags "$(LDFLAGS)" ./cmd/scheduler
	go build -v -o $(SENDER_BIN) -ldflags "$(LDFLAGS)" ./cmd/sender

run: build
	$(CALENDAR_BIN) -config ./configs/calendar.dev.toml
	$(SCHEDULER_BIN) -config ./configs/scheduler.dev.toml
	$(SENDER_BIN) -config ./configs/sender.dev.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(CALENDAR_DOCKER_IMG) \
		-f build/calendar/Dockerfile .
	
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(SCHEDULER_DOCKER_IMG) \
		-f build/scheduler/Dockerfile .

	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(SENDER_DOCKER_IMG) \
		-f build/sender/Dockerfile .

generate-pb:
	mkdir -p internal/app/delivery/grpc/pb
	protoc \
		-I /usr/local/include \
		--proto_path=api/ \
		--go_out=internal/app/delivery/grpc/pb \
		--go-grpc_out=internal/app/delivery/grpc/pb \
		api/*.proto

run-img: build-img
	docker run $(CALENDAR_DOCKER_IMG)
	docker run $(SCHEDULER_DOCKER_IMG)
	docker run $(SENDER_DOCKER_IMG)

create-network:
	docker network create calendar || true

up: create-network
	docker-compose -f ./deployments/docker-compose.yaml up -d

down:
	docker-compose -f ./deployments/docker-compose.yaml down

version: build
	$(CALENDAR_BIN) version

test:
	go test -race -v ./internal/...

integration-tests: create-network
	docker-compose -f ./deployments/docker-compose.yaml -f ./deployments/docker-compose.test.yaml up --build -d
	docker-compose -f ./deployments/docker-compose.yaml -f ./deployments/docker-compose.test.yaml run integration_tests
	docker-compose -f ./deployments/docker-compose.yaml -f ./deployments/docker-compose.test.yaml down

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint
