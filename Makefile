# StudioSpeech Makefile

.PHONY: build test clean install deps check

# Build the CLI binary
build:
	@echo "Building ttscli..."
	@mkdir -p bin
	@go build -o bin/ttscli ./cmd/ttscli

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p bin
	@GOOS=darwin GOARCH=amd64 go build -o bin/ttscli-darwin-amd64 ./cmd/ttscli
	@GOOS=darwin GOARCH=arm64 go build -o bin/ttscli-darwin-arm64 ./cmd/ttscli
	@GOOS=windows GOARCH=amd64 go build -o bin/ttscli-windows-amd64.exe ./cmd/ttscli
	@GOOS=linux GOARCH=amd64 go build -o bin/ttscli-linux-amd64 ./cmd/ttscli

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf tmp/

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

# Check system requirements (piper, ffmpeg)
check:
	@echo "Checking system requirements..."
	@command -v piper >/dev/null 2>&1 || echo "WARNING: piper not found in PATH"
	@command -v ffmpeg >/dev/null 2>&1 || echo "WARNING: ffmpeg not found in PATH"
	@go version

# Install the binary to GOPATH/bin
install: build
	@echo "Installing ttscli..."
	@cp bin/ttscli $(GOPATH)/bin/ttscli