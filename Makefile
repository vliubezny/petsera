PROJECT := petsera
OUT_DIR := ./build
OUT := $(OUT_DIR)/$(PROJECT)
DIST_DIR := ./dist
MAIN_PKG := ./cmd/$(PROJECT)

GOBIN := $(shell go env GOPATH)/bin

MIGRATE_NAME := migrate
MIGRATE_VERSION := v4.14.1

default: build

.PHONY: build
build:
	CGO_ENABLED=0 go build -mod=vendor -o $(OUT) $(MAIN_PKG)

.PHONY: build-ui
build-ui:
	@(cd ui && npm run build)

.PHONY: linux
linux: export GOOS := linux
linux: export GOARCH := amd64
linux: LINUX_OUT := $(OUT)-$(GOOS)-$(GOARCH)
linux:
	@echo BUILDING $(LINUX_OUT)
	CGO_ENABLED=0 go build -mod=vendor -o $(LINUX_OUT) $(MAIN_PKG)
	@echo DONE

.PHONY: clean
clean:
	rm -rf $(OUT_DIR) $(DIST_DIR)
	go clean -testcache

.PHONY: test
test: GO_TEST_FLAGS := -race
test:
	go test -v -mod=vendor $(GO_TEST_FLAGS) $(GO_TEST_TAGS) ./...

.PHONY: fulltest
fulltest: GO_TEST_TAGS := --tags=integration
fulltest: test

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor

.PHONY: run
run:
	@./scripts/dev-server.sh

.PHONY: sandbox-up
sandbox-up:
	@docker-compose -f scripts/sandbox.yml up -d

.PHONY: sandbox-down
sandbox-down:
	@docker-compose -f scripts/sandbox.yml down

.PHONY: install-migrate
install-migrate:
	@echo INSTALLING $(MIGRATE_NAME) for MacOS
	@echo 'Check https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md for details'
	brew install golang-migrate
	@echo DONE

.PHONY: new-migration
new-migration:
	migrate create -ext sql -dir scripts/migrations/postgres -seq $(NAME)