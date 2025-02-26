.PHONY: build test coverage lint clean install all windows linux darwin

# Default target
all: build

# Install dependencies
install:
	go mod download
	go install github.com/golangci/golint/cmd/golangci-lint@latest

# Build the application for current platform
build:
	go build -v -o screenshot-sorter$(if $(filter windows,$(OS)),.exe,) .

# Platform-specific builds
windows:
	cmd /V:ON /C "set GOOS=windows&& set GOARCH=amd64&& go build -v -o screenshot-sorter.exe ."

linux:
	bash -c "GOOS=linux GOARCH=amd64 go build -v -o screenshot-sorter ."

darwin:
	bash -c "GOOS=darwin GOARCH=amd64 go build -v -o screenshot-sorter ."

# Run all tests
test:
	cmd /V:ON /C "set SCREENSHOT_SORTER_TEST_USE_MODTIME=1&& go test -v ./..."

# Run tests with coverage (without race detector)
coverage:
	cmd /V:ON /C "set SCREENSHOT_SORTER_TEST_USE_MODTIME=1&& go test -v -coverprofile=coverage.txt -covermode=count ./..."
	go tool cover -html=coverage.txt -o coverage.html

# Run tests with coverage and race detector (requires GCC)
coverage-race:
	cmd /V:ON /C "set CGO_ENABLED=1&& set SCREENSHOT_SORTER_TEST_USE_MODTIME=1&& go test -v -race -coverprofile=coverage.txt -covermode=atomic ./..."
	go tool cover -html=coverage.txt -o coverage.html

# Run linter
lint:
	golangci-lint run

# Clean build artifacts
clean:
	del /F /Q screenshot-sorter.exe 2>nul || exit 0
	del /F /Q coverage.txt coverage.html 2>nul || exit 0