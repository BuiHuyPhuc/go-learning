APP_NAME := server

run:
  go run ./cmd/${APP_NAME}/

GOOSE_DRIVER ?= mysql
GOOSE_DBSTRING ?= "root:root123@tcp(127.0.0.1:33306)/shopdevgo"
GOOSE_MIGRATION_DIR ?= sqlc/schema

dc_up:
  docker-compose up -d
dc_down:
  docker-compose down

upse:
  @GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir $(GOOSE_MIGRATION_DIR) up
  # goose -dir sqlc/schema mysql "root:root123@tcp(127.0.0.1:33306)/shopdevgo" up
downse:
  @GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir $(GOOSE_MIGRATION_DIR) down
  # goose -dir sqlc/schema mysql "root:root123@tcp(127.0.0.1:33306)/shopdevgo" down
resetse:
  @GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir $(GOOSE_MIGRATION_DIR) reset
  # goose -dir sqlc/schema mysql "root:root123@tcp(127.0.0.1:33306)/shopdevgo" reset

sqlgen:
  sqlc generate

.PHONY: run upse downse resetse dc_up dc_down

.PHONY: air
