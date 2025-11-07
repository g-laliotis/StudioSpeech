# Makefile for StudioSpeech
# Professional build automation and development workflow

.DEFAULT_GOAL := help
.PHONY: help build run test test-coverage test-verbose bench fmt vet clean install version deps check-deps

# Build configuration
BINARY_NAME := ttscli
BIN_DIR := bin
CMD_DIR := ./cmd/ttscli
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Go build flags
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)"
BUILD_FLAGS := -trimpath $(LDFLAGS)

# Colors for output
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
MAGENTA := \033[35m
CYAN := \033[36m
WHITE := \033[37m
RESET := \033[0m

## help: Show this help message
help:
	@echo "$(CYAN)StudioSpeech Build System$(RESET)"
	@echo "$(CYAN)========================$(RESET)"
	@echo ""
	@echo "$(YELLOW)Available commands:$(RESET)"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
	@echo ""
	@echo "$(YELLOW)Examples:$(RESET)"
	@echo "  $(GREEN)make build$(RESET)        - Build optimized binary"
	@echo "  $(GREEN)make test$(RESET)         - Run all tests"
	@echo "  $(GREEN)make run \"file.txt\"$(RESET)    - Convert file.txt to file.mp3"
	@echo "  $(GREEN)make run-greek \"file.txt\"$(RESET) - Convert to Greek speech"
	@echo "  $(GREEN)make run-male \"file.txt\"$(RESET)  - Convert to male voice"
	@echo "  $(GREEN)make clean$(RESET)        - Clean build artifacts"

## build: Build the CLI binary with version info
build: check-deps
	@echo "$(BLUE)Building $(BINARY_NAME) v$(VERSION)...$(RESET)"
	@mkdir -p $(BIN_DIR)
	@go build $(BUILD_FLAGS) -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "$(GREEN)✓ Built $(BIN_DIR)/$(BINARY_NAME)$(RESET)"

## run: Run synthesis on sample data or specified file
run: build
	@if [ "$(filter-out $@,$(MAKECMDGOALS))" ]; then \
		INPUT_FILE="$(filter-out $@,$(MAKECMDGOALS))"; \
		OUTPUT_FILE="$${INPUT_FILE%.*}.mp3"; \
		echo "$(BLUE)Converting $$INPUT_FILE to $$OUTPUT_FILE...$(RESET)"; \
		./$(BIN_DIR)/$(BINARY_NAME) synth --in "$$INPUT_FILE" --out "$$OUTPUT_FILE"; \
		echo "$(GREEN)✓ Generated $$OUTPUT_FILE$(RESET)"; \
	else \
		echo "$(BLUE)Running synthesis on sample data...$(RESET)"; \
		if [ -f "testdata/samples/sample.txt" ]; then \
			./$(BIN_DIR)/$(BINARY_NAME) synth --in testdata/samples/sample.txt --out output.mp3; \
			echo "$(GREEN)✓ Generated output.mp3$(RESET)"; \
		else \
			echo "$(YELLOW)⚠ Sample file not found, creating test file...$(RESET)"; \
			mkdir -p testdata/samples; \
			echo "Hello world. This is a test sentence with proper punctuation!" > testdata/samples/sample.txt; \
			./$(BIN_DIR)/$(BINARY_NAME) synth --in testdata/samples/sample.txt --out output.mp3; \
			echo "$(GREEN)✓ Generated output.mp3$(RESET)"; \
		fi; \
	fi

# Allow make run "filename" to work
%:
	@:

## test: Run all tests
test:
	@echo "$(BLUE)Running tests...$(RESET)"
	@go test -race ./...
	@echo "$(GREEN)✓ All tests passed$(RESET)"

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "$(BLUE)Running tests with coverage...$(RESET)"
	@go test -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report generated: coverage.html$(RESET)"

## test-verbose: Run tests with verbose output
test-verbose:
	@echo "$(BLUE)Running tests (verbose)...$(RESET)"
	@go test -v -race ./...

## bench: Run benchmark tests
bench:
	@echo "$(BLUE)Running benchmarks...$(RESET)"
	@go test -bench=. -benchmem ./...
	@echo "$(GREEN)✓ Benchmarks completed$(RESET)"

## fmt: Format Go code
fmt:
	@echo "$(BLUE)Formatting code...$(RESET)"
	@go fmt ./...
	@echo "$(GREEN)✓ Code formatted$(RESET)"

## vet: Run static analysis
vet:
	@echo "$(BLUE)Running static analysis...$(RESET)"
	@go vet ./...
	@echo "$(GREEN)✓ Static analysis passed$(RESET)"

## clean: Remove build artifacts and temporary files
clean:
	@echo "$(BLUE)Cleaning build artifacts...$(RESET)"
	@rm -rf $(BIN_DIR)
	@rm -f $(BINARY_NAME)
	@rm -f output.mp3 output.wav
	@rm -f coverage.out coverage.html
	@rm -rf /tmp/studiospeech_*
	@echo "$(GREEN)✓ Cleaned build artifacts$(RESET)"

## clear-cache: Clear Go build and test cache
clear-cache:
	@echo "$(BLUE)Clearing Go cache...$(RESET)"
	@go clean -cache -testcache -modcache
	@echo "$(GREEN)✓ Go cache cleared$(RESET)"

## install: Install system dependencies (macOS)
install:
	@echo "$(BLUE)Installing system dependencies...$(RESET)"
	@if command -v brew >/dev/null 2>&1; then \
		echo "$(YELLOW)Installing FFmpeg via Homebrew...$(RESET)"; \
		brew install ffmpeg; \
		echo "$(GREEN)✓ FFmpeg installed$(RESET)"; \
	else \
		echo "$(RED)✗ Homebrew not found. Please install FFmpeg manually.$(RESET)"; \
		exit 1; \
	fi

## deps: Download and verify dependencies
deps:
	@echo "$(BLUE)Downloading dependencies...$(RESET)"
	@go mod download
	@go mod verify
	@echo "$(GREEN)✓ Dependencies verified$(RESET)"

## check-deps: Check if required dependencies are available
check-deps:
	@echo "$(BLUE)Checking dependencies...$(RESET)"
	@go version >/dev/null 2>&1 || (echo "$(RED)✗ Go not installed$(RESET)" && exit 1)
	@echo "$(GREEN)✓ Go $(shell go version | cut -d' ' -f3)$(RESET)"

## version: Show version information
version:
	@echo "$(CYAN)StudioSpeech Build Information$(RESET)"
	@echo "$(CYAN)=============================$(RESET)"
	@echo "Version:    $(VERSION)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Go Version: $(shell go version | cut -d' ' -f3)"

## check: Run all quality checks (fmt, vet, test)
check: fmt vet test
	@echo "$(GREEN)✓ All quality checks passed$(RESET)"

## ci: Run CI pipeline (used by GitHub Actions)
ci: check-deps fmt vet test
	@echo "$(GREEN)✓ CI pipeline completed$(RESET)"

## release: Prepare release build
release: clean check-deps
	@echo "$(BLUE)Building release version...$(RESET)"
	@mkdir -p $(BIN_DIR)
	@CGO_ENABLED=0 go build $(BUILD_FLAGS) -a -installsuffix cgo -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "$(GREEN)✓ Release build completed: $(BIN_DIR)/$(BINARY_NAME)$(RESET)"

## dev: Development workflow (fmt, vet, test, build)
dev: fmt vet test build
	@echo "$(GREEN)✓ Development workflow completed$(RESET)"

## run-greek: Convert file to Greek speech
run-greek: build
	@if [ "$(filter-out $@,$(MAKECMDGOALS))" ]; then \
		INPUT_FILE="$(filter-out $@,$(MAKECMDGOALS))"; \
		OUTPUT_FILE="$${INPUT_FILE%.*}.mp3"; \
		echo "$(BLUE)Converting $$INPUT_FILE to Greek speech...$(RESET)"; \
		./$(BIN_DIR)/$(BINARY_NAME) synth --in "$$INPUT_FILE" --lang el-GR --gender female --out "$$OUTPUT_FILE"; \
		echo "$(GREEN)✓ Generated $$OUTPUT_FILE with Greek voice$(RESET)"; \
	else \
		echo "$(RED)Usage: make run-greek \"filename.txt\"$(RESET)"; \
	fi

## run-male: Convert file to male voice
run-male: build
	@if [ "$(filter-out $@,$(MAKECMDGOALS))" ]; then \
		INPUT_FILE="$(filter-out $@,$(MAKECMDGOALS))"; \
		OUTPUT_FILE="$${INPUT_FILE%.*}.mp3"; \
		echo "$(BLUE)Converting $$INPUT_FILE to male voice...$(RESET)"; \
		./$(BIN_DIR)/$(BINARY_NAME) synth --in "$$INPUT_FILE" --gender male --out "$$OUTPUT_FILE"; \
		echo "$(GREEN)✓ Generated $$OUTPUT_FILE with male voice$(RESET)"; \
	else \
		echo "$(RED)Usage: make run-male \"filename.txt\"$(RESET)"; \
	fi

## run-female: Convert file to female voice
run-female: build
	@if [ "$(filter-out $@,$(MAKECMDGOALS))" ]; then \
		INPUT_FILE="$(filter-out $@,$(MAKECMDGOALS))"; \
		OUTPUT_FILE="$${INPUT_FILE%.*}.mp3"; \
		echo "$(BLUE)Converting $$INPUT_FILE to female voice...$(RESET)"; \
		./$(BIN_DIR)/$(BINARY_NAME) synth --in "$$INPUT_FILE" --gender female --out "$$OUTPUT_FILE"; \
		echo "$(GREEN)✓ Generated $$OUTPUT_FILE with female voice$(RESET)"; \
	else \
		echo "$(RED)Usage: make run-female \"filename.txt\"$(RESET)"; \
	fi