# Build variables
BINARY_NAME = nxpx
BUILD_DIR = build
VERSION ?= $(shell git tag --points-at HEAD | tail -n 1)
BUILD_DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
COMMIT_SHA = $(shell git rev-parse --short HEAD)
LDFLAGS = -ldflags "-w -X main.version=${VERSION} -X main.buildDate=${BUILD_DATE} -X main.commit=${COMMIT_SHA}"

# Docker variables
DOCKER_IMAGE ?= nextpax/np
DOCKER_TAG ?= dev
GOOS ?= linux

.PHONY: env
env: ## Show env configuration
	go run ./cmd/${BINARY_NAME} --help

.PHONY: docker
docker: ## Build a Docker image
	docker build -f ./Dockerfile --rm -t ${DOCKER_IMAGE}:${DOCKER_TAG} .

.PHONY: dcup
dcup: ## Local docker-compose up
	VERSION=${DOCKER_TAG} docker-compose  -p $(BINARY_NAME) -f docker-compose.yml up -d --build

.PHONY: import
import: ## Local docker-compose up
	docker-compose -p $(BINARY_NAME) -f docker-compose.yml exec -d db "/var/fixtures/import.sh"

.PHONY: dcdown
dcdown: ## Local docker-compose down
	VERSION=${DOCKER_TAG} docker-compose  -p $(BINARY_NAME) -f docker-compose.yml down -v --remove-orphans --rmi=local

.PHONY: dep
dep: ## Install dependencies
	$(eval PACKAGE := $(shell go list -m))
	@go mod tidy
	@go mod download
	@go mod vendor

.PHONY: test
test: dep ## Run unit tests
	@go test -v ./...

.PHONY: build
build: test ## Build a binary executable file
	GOOS=${GOOS} GOARCH=amd64 CGO_ENABLED=0 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} ${PACKAGE}/cmd/${BINARY_NAME}

.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# Variable outputting/exporting rules
var-%: ; @echo $($*)
varexport-%: ; @echo $*=$($*)
