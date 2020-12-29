.PHONY: build image fmt test clean

build:
	go build -o dist/vault-init ./cmd/vault-init/...

image:
	docker build -t vault-init:latest .

fmt:
	go fmt ./...

test:
	go test -race -v ./...

clean:
	go clean ./...
	rm -rf ./dist
