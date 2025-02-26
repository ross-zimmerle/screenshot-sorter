.PHONY: build test coverage lint clean

# Build the application
build:
	go build -v .

# Run all tests
test:
	go test -v ./...

# Run tests with coverage
coverage:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -html=coverage.txt -o coverage.html

# Run linter
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -f screenshot-sorter
	rm -f coverage.txt
	rm -f coverage.html