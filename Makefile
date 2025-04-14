.PHONY: build run test clean lint build-and-run

clean:
	rm -f chirpy

lint:
	golangci-lint run ./...

build:
	go build -o chirpy main.go

run:
	./chirpy

test:
	go test -v ./tests/...

build-and-run: clean lint test build run

start:
	find . -name '*.go' | entr -r make build-and-run
