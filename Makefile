PROJECT := petsera
OUT_DIR := ./build
OUT := $(OUT_DIR)/$(PROJECT)
DIST_DIR := ./dist
MAIN_PKG := ./cmd/$(PROJECT)

GOBIN := $(shell go env GOPATH)/bin

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

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor

.PHONY: run
run:
	@./scripts/dev-server.sh
