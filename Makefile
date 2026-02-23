.PHONY: test build lint clean

BINARY_NAME=zaap

test:
	go test -v ./...

build:
	go build -o $(BINARY_NAME) .

lint:
	golangci-lint run ./...

clean:
	rm -f $(BINARY_NAME)
