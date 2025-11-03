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

## For developers

---

- Build up environment

```shell
# 1. Install Go: version 1.25.0
# mac
brew install go
# Windows : https://go.dev/dl/ 

# 2. Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# 3. Change the @host 'api.hiroliang.com' path at cmd/api/main.go to localhost:8080

# 4. Install dependencies
go mod tidy

# 5. Generate swagger docs
make swag

# 6. Build datasource at first time (use docker or podman)

# Create network
docker network create goat-net

# (Option) Create volume if you want to keep data persist in host machine (you can change the path)
mkdir -p ~/podman-data/redis

# Run Redis (-v ~/podman-data/redis:/data \ if folder is created)
docker run -d \
  --name redis \
  --network goat-net \
  -p 6379:6379 \
  -v ~/podman-data/redis:/data \
  docker.io/library/redis:8 \
  redis-server --appendonly yes --requirepass "1234"
  
# (Option) Create volume if you want to keep data persist in host machine (you can change the path)
mkdir -p ~/podman-data/postgres 

# Run Postgres (-v ~/podman-data/postgres:/data \ if folder is created)
docker run -d \
  --name postgres \
  --network goat-net \
  -e POSTGRES_USER=root \
  -e POSTGRES_PASSWORD=1234 \
  -e POSTGRES_DB=goat \
  -p 5432:5432 \
  docker.io/library/postgres:18
  
# 7.Add required documents below

# 8. Run server
make run
```

- Deploy locally

```shell
# 1. Build image
docker build -t goat-server:latest .

# 2. Run container
docker run -d \
  --name goat-server \
  --network goat-net \
  -p 8080:8080 \
  goat-server:latest

```

- Documents

```yaml
# ./config/config.yaml (relative to project root)
databases:
  mysql: # not used
    driver: mysql
    dsn: "root:1234@tcp(localhost:3306)/goat?charset=utf8mb4&parseTime=True&loc=Local"
    config:
      max_open_conns: 20
      max_idle_conns: 10
      conn_max_lifetime: 3600 # second
      conn_max_idle_time: 600 # second
  postgres:
    driver: pgx
    dsn: "postgres://root:1234@localhost:5432/goat?sslmode=disable"
    config:
      max_open_conns: 30
      max_idle_conns: 15
      conn_max_lifetime: 3600 # second
      conn_max_idle_time: 600 # second
redis:
  addr: "localhost:6379"
  password: "1234"
  db: 0
```

```dotenv
# .env (relative to project root)
APP_ENV=dev

SERVER_PORT=8080
```