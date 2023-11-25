RELEASE_PKG := main
BUILD_VERSION ?= $(shell git rev-parse --short=8 HEAD || echo unknown)
LDFLAGS := -X '$(RELEASE_PKG).Version=$(BUILD_VERSION)'
LDFLAGS += -X '$(RELEASE_PKG).DBUser=db-user'
LDFLAGS += -X '$(RELEASE_PKG).DBPass=db-pass'
LDFLAGS += -X '$(RELEASE_PKG).DBHost=localhost:5432'
LDFLAGS += -X '$(RELEASE_PKG).DBName=db-name'


.PHONY: build
build:
	go build -v -ldflags "-s $(LDFLAGS)"

.PHONY: test
test:
	go test ./...

.PHONY: up
up:
	docker compose -f docker-compose.yaml up -d --build

.PHONY: down
down:
	docker compose -f docker-compose.yaml down

