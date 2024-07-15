VERSION := 0.0.1

build: test
	go build -ldflags="-X 'main.version=$(VERSION)'" -o docker-tree ./cmd/

test:
	go vet ./...
	go test -v ./...
