PLUGIN_NAME := docker-tree
VERSION := 0.0.3

build: test
	go build -ldflags="-X 'main.version=${VERSION}'" -o ${PLUGIN_NAME} ./cmd/

test:
	go vet ./...
	go test -v ./...
