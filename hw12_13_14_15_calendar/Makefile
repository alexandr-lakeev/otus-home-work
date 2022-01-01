BIN := "./bin/calendar"
CALENDAR_DOCKER_IMG="calendar:develop"
SCHEDULER_DOCKER_IMG="calendar_scheduler:develop"
SENDER_DOCKER_IMG="calendar_sender:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/calendar.dev.toml

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

create-network:
	docker network create calendar || true

run-compose: create-network
	docker-compose -f ./deployments/docker-compose.yaml up

version: build
	$(BIN) version

test:
	go test -race -v ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint
