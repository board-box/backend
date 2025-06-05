.PHONY: build run stop db-up db-down generate migration-create migration-up bin-deps

-include .env

##############################
# General
##############################

CURDIR := $(shell pwd)
LOCAL_BIN:=$(CURDIR)/bin

GOLANGCI_BIN:=$(LOCAL_BIN)/golangci-lint
SMART_IMPORTS := ${LOCAL_BIN}/smartimports

export GOPROXY=direct
export GOBIN=$(LOCAL_BIN)

# build app
build:
	go mod download && CGO_ENABLED=0
	go build -o ./bin/main$(shell go env GOEXE) ./cmd/main.go

# run all in docker
run: build
	docker-compose -f $(CURDIR)/build/docker/docker-compose.yml up --force-recreate --build -d
	make migration-up

# stop in docker
stop:
	docker-compose -f $(CURDIR)/build/docker/docker-compose.yml down


##############################
# Database
##############################

POSTGRES_USER ?= user
POSTGRES_PASSWORD ?= password
POSTGRES_DB ?= boardbox
POSTGRES_HOST ?= localhost
POSTGRES_PORT ?= 5432
POSTGRES_SSLMODE ?= disable

POSTGRES_SETUP := user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) dbname=$(POSTGRES_DB) host=$(POSTGRES_HOST) port=$(POSTGRES_PORT) sslmode=$(POSTGRES_SSLMODE)

db-up:
	docker compose -f build/docker/docker-compose.yml up -d --build postgres

db-down:
	docker compose -f build/docker/docker-compose.yml down postgres


# Migrations

INTERNAL_PKG_PATH=$(CURDIR)/internal
MIGRATION_FOLDER=$(CURDIR)/migrations
migration_name = collection_game_table

migration-create:
	bin/goose -dir "$(MIGRATION_FOLDER)" create "$(migration_name)" sql

migration-up:
	bin/goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" up

##############################
# Generate
##############################

# generate
generate:
	$(LOCAL_BIN)/swag init -g cmd/main.go --parseDependency --parseInternal

bin-deps:
	$(info Installing binary dependencies...)
	GOBIN=$(LOCAL_BIN) go install github.com/swaggo/swag/cmd/swag@latest
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest


##############################
# CI-CD
##############################

.PHONY: test
test:
	go test -v ./...

# precommit jobs
.PHONY: precommit
precommit: format lint

.PHONY: install-lint
install-lint:
ifeq ($(wildcard $(GOLANGCI_BIN)),)
	$(info Downloading golangci-lint v$(GOLANGCI_TAG))
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
endif

# run diff lint like in pipeline
.PHONY: lint
lint: install-lint
	$(info Running lint...)
	$(GOLANGCI_BIN) run --new-from-rev=origin/master --config=build/linter/.golangci.yaml ./...

# run full lint like in pipeline
.PHONY: lint-full

lint-full: install-lint
	$(GOLANGCI_BIN) run --config=build/linter/.golangci.yaml ./...

.PHONY: format
format:
	$(info Running goimports...)
	test -f ${SMART_IMPORTS} || GOBIN=${LOCAL_BIN} go install github.com/pav5000/smartimports/cmd/smartimports@latest
	${SMART_IMPORTS} -exclude pkg/,internal/pb  -local 'gitlab.ozon.dev'
