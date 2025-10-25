Goat Application Server
===

## Introduction

---

Goat-Chat is a personal practice project of a Tauri desktop application with the Ollama agent. 
I use GoLang to be its server side.

## Planing

---

- User authentication
- Chat control (websocket, LLM redirect, chat group switcher...)
- For practice. Other features are on the road... Maybe to send an order to resp berry pi IOT?

## Commands

---

- Api
```shell
# test locally
make run

# build
make build

# delete build
make clean
```

- Swagger
```shell
# install swag
go install github.com/swaggo/swag/cmd/swag@latest

# check swag version
swag --version

# generate swagger docs
make swag
```