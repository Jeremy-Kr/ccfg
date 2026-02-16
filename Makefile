.PHONY: build test lint clean run

build:
	go build -o ccfg ./cmd/ccfg

run: build
	./ccfg

test:
	go test ./...

lint:
	golangci-lint run

clean:
	rm -f ccfg
