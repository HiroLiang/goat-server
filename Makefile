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

.PHONY: build clean run
build:
	GOOS=linux GOARCH=amd64 go build -o bin/goat-api ./cmd/api

clean:
	rm -rf bin

run:
	go run ./cmd/api

swag:
	swag init -g cmd/api/main.go