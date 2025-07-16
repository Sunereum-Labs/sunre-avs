# -----------------------------------------------------------------------------
# This Makefile is used for building your AVS application.
#
# It contains basic targets for building the application, installing dependencies,
# and building a Docker container.
#
# Modify each target as needed to suit your application's requirements.
# -----------------------------------------------------------------------------

GO = $(shell which go)
OUT = ./bin

build: deps
	@mkdir -p $(OUT) || true
	@echo "Building binaries..."
	go build -o $(OUT)/performer ./cmd/main.go

deps:
	GOPRIVATE=github.com/Layr-Labs/* go mod tidy

build/container:
	./.hourglass/scripts/buildContainer.sh

test:
	@echo "Running unit tests..."
	go test ./cmd/... -v -count=1

test-short:
	@echo "Running unit tests (short mode)..."
	go test ./cmd/... -v -short

test-coverage:
	@echo "Running tests with coverage..."
	go test ./cmd/... -v -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-bench:
	@echo "Running benchmarks..."
	go test ./cmd/... -bench=. -benchmem -run=^#

test-integration:
	@echo "Running integration tests with devkit..."
	devkit avs devnet start
	@sleep 5
	go test ./cmd/... -v -tags=integration
	devkit avs devnet stop

.PHONY: build deps build/container test test-short test-coverage test-bench test-integration