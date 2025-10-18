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

.PHONY: build clean run
build:
	go build -o bin/api ./cmd/api

clean:
	rm -rf bin

run:
	go run ./cmd/api