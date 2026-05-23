include .env

export PROJECT_ROOT=$(shell pwd)


.PHONY: env-up env-down migrate-up migrate-up

env-up:
	@docker compose up -d avito_postgres

env-down:
	@docker compose down 

DB_DSN=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

migrate-up:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_DSN) goose -dir ./migrations up

migrate-down:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_DSN) goose -dir ./migrations down