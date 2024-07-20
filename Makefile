VERSION := 0.0.2

build: #test
	go build -ldflags="-X 'main.version=$(VERSION)'" -o docker-tree ./cmd/

test:
	go vet ./...
	go test -v ./...
