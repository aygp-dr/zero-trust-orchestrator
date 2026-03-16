.PHONY: build run test clean

build:
	go build -o bin/zero-trust-orchestrator .

run: build
	./bin/zero-trust-orchestrator

test:
	go test ./...

clean:
	rm -rf bin/
