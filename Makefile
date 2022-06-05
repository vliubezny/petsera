PROJECT := petsera
OUT_DIR := ./build
OUT := $(OUT_DIR)/$(PROJECT)
DIST_DIR := ./dist
MAIN_PKG := ./cmd/$(PROJECT)

GCP_PROJECT_ID := petsera
DOCKER_REGISTRY := europe-north1-docker.pkg.dev/$(GCP_PROJECT_ID)/petsera/app
TAG ?= latest

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

.PHONY: image
image:
	docker build -t $(DOCKER_REGISTRY):$(TAG) -f scripts/Dockerfile .

.PHONY: push
push:
	docker push $(DOCKER_REGISTRY):$(TAG)

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

.PHONY: compose-up
compose-up:
	@env $$(cat ./cloud.local.env | grep -Ev '^#' | xargs) docker-compose -f scripts/docker-compose.yml up -d

.PHONY: compose-down
compose-down:
	@docker-compose -f scripts/docker-compose.yml down

.PHONY: install-migrate
install-migrate:
	@echo INSTALLING $(MIGRATE_NAME) for MacOS
	@echo 'Check https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md for details'
	brew install golang-migrate
	@echo DONE

.PHONY: new-migration
new-migration:
	migrate create -ext sql -dir scripts/migrations/postgres -seq $(NAME)

.PHONY: migrate
migrate:
	@echo Running migration for Cloud SQL
	@set -a && . ./cloud.local.env \
	&& migrate -database "postgres://$$PETSERA_DB_USER:$$PETSERA_DB_PASSWORD@localhost:5432/$$PETSERA_DB_NAME?sslmode=disable" \
		-path scripts/migrations/postgres up

.PHONY: deploy
deploy:
	helm upgrade --install -f ./charts/prod-values.yaml \
		--set config.mapsAPIKey=$$(gcloud secrets versions access latest --secret=maps-api-key-prod) \
		--set config.db.password=$$(gcloud secrets versions access latest --secret=petsera-db-password) \
		--set image.pullPolicy=Always --set image.tag=$(TAG) \
		petsera ./charts/petsera
