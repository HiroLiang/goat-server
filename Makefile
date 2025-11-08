# ---- config ----
APP       := api
CMD_DIR   := ./cmd/api
BIN       := bin/$(APP)
GO        := go

# ---- targets ----
.PHONY: help
help:
	@echo "Commends："
	@echo "  make build    # Compile -> $(BIN)"
	@echo "  make clean    # clean compiled output"
	@echo "  make run      # Start goat-api locally"
	@echo "  make swag     # Build swagger docs in ./docs/"
	@echo "  make test     # Run unit & BDD test"

.PHONY: setup build clean run

setup: setup_swag setup_godog
	@echo "Installing dependencies..."

setup_swag:
	@echo "Installing Swag CLI..."
	go install github.com/swaggo/swag/cmd/swag@latest

setup_godog:
	@echo "Installing Godog CLI..."
	go install github.com/cucumber/godog/cmd/godog@latest

build:
	go build -o bin/goat-api ./cmd/api

clean:
	rm -rf bin

run:
	go run ./cmd/api

swag:
	swag init -g cmd/api/main.go -o docs

.PHONY: test unit bdd

test: unit bdd

unit:
	@echo ""
	@echo "Running unit tests..."
	go test ./internal/... -v

bdd:
	@echo ""
	@echo "Running BDD feature tests..."
	go test -v ./features