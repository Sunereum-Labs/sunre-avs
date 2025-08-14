# SunRe AVS Makefile

GO = $(shell which go)
OUT = ./bin
DEVKIT = devkit

# Primary DevKit commands
build: ## Build using DevKit
	@echo "Building AVS with DevKit..."
	$(DEVKIT) avs build

test: ## Test using DevKit
	@echo "Testing AVS with DevKit..."
	$(DEVKIT) avs test

deploy: ## Deploy using DevKit
	@echo "Deploying AVS with DevKit..."
	$(DEVKIT) avs deploy

# DevKit-specific commands
init: ## Initialize DevKit project
	@echo "Initializing DevKit project..."
	$(DEVKIT) avs init --template hourglass

validate: ## Validate DevKit configuration
	@echo "Validating DevKit configuration..."
	$(DEVKIT) avs validate

manifest: ## Generate DevKit manifest
	@echo "Generating DevKit manifest..."
	$(DEVKIT) avs manifest generate

# Development environment
devnet: ## Start local devnet
	$(DEVKIT) avs devnet start

devnet-stop: ## Stop local devnet
	$(DEVKIT) avs devnet stop

# Go-specific builds (when needed outside DevKit)
build-go: deps
	@mkdir -p $(OUT) || true
	@echo "Building SunRe performer..."
	go build -o $(OUT)/performer ./cmd/main.go

deps:
	GOPRIVATE=github.com/Layr-Labs/* go mod tidy

# Container build using DevKit patterns
build-container:
	./.hourglass/scripts/buildContainer.sh

# Testing
test-go:
	@echo "Running unit tests..."
	@go test ./cmd/... -v -count=1

test-integration:
	@echo "Running integration tests..."
	$(DEVKIT) avs devnet start
	@sleep 5
	go test ./cmd/... -v -tags=integration
	$(DEVKIT) avs devnet stop

test-quick: ## Run quick tests
	@./scripts/test.sh

demo: ## Run full demo
	@./scripts/demo.sh

# Task operations
create-task: ## Create weather verification task
	$(DEVKIT) avs call

status: ## Check AVS status
	$(DEVKIT) avs status

logs: ## View logs
	$(DEVKIT) avs logs

clean: ## Clean build artifacts
	rm -rf $(OUT) coverage.out coverage.html

help: ## Show this help
	@echo "SunRe AVS - Weather Insurance Verification"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

.PHONY: build test deploy devnet devnet-stop build-go deps build-container test-go test-integration create-task status logs clean help
.DEFAULT_GOAL := help