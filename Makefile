.PHONY: test build install lint clean

BINARY_NAME=zaap

test:
	go test -v ./...

build:
	go build -o $(BINARY_NAME) .

install:
	go install .

lint:
	golangci-lint run ./...

clean:
	rm -f $(BINARY_NAME)
