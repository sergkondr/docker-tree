DOCKER_CONFIG ?= $(HOME)/.docker
APP_NAME := docker-tree
VERSION := 0.1.0

.PHONY: build build-all test install

build:
	go build -ldflags="-X 'main.version=${VERSION}'" -o ${APP_NAME} ./cmd/

build-all: test
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.version=${VERSION}'" -o ${APP_NAME}-linux-amd64 ./cmd/
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-X 'main.version=${VERSION}'" -o ${APP_NAME}-linux-arm64 ./cmd/
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-X 'main.version=${VERSION}'" -o ${APP_NAME}-darwin-amd64 ./cmd/
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-X 'main.version=${VERSION}'" -o ${APP_NAME}-darwin-arm64 ./cmd/

test:
	go vet ./...
	go test -v ./...

install: build
	mkdir -p $(DOCKER_CONFIG)/cli-plugins
	install -m 0755 $(APP_NAME) $(DOCKER_CONFIG)/cli-plugins/$(APP_NAME)
