APP_NAME := server

run:
  go run ./cmd/${APP_NAME}/

GOOSE_DRIVER ?= mysql
GOOSE_DBSTRING ?= "root:root123@tcp(127.0.0.1:33306)/shopdevgo"
GOOSE_MIGRATION_DIR ?= sqlc/schema

kill:
  docker compose kill
up:
  docker compose up -d
down:
  docker compose down

upse:
  @GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir $(GOOSE_MIGRATION_DIR) up
  # goose -dir sqlc/schema mysql "root:root123@tcp(127.0.0.1:33306)/shopdevgo" up
downse:
  @GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir $(GOOSE_MIGRATION_DIR) down
  # goose -dir sqlc/schema mysql "root:root123@tcp(127.0.0.1:33306)/shopdevgo" down
resetse:
  @GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir $(GOOSE_MIGRATION_DIR) reset
  # goose -dir sqlc/schema mysql "root:root123@tcp(127.0.0.1:33306)/shopdevgo" reset

.PHONY: run upse downse resetse

.PHONY: air
