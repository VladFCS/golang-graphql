APP_ENV ?= development
HTTP_ADDR ?= :8080
DATABASE_URL ?= postgres://postgres:postgres@localhost:5433/qhira?sslmode=disable
JWT_SECRET ?= change-me-in-development
LOG_LEVEL ?= info
GOOSE ?= go run github.com/pressly/goose/v3/cmd/goose

.PHONY: run test fmt gqlgen sqlc db-up db-down db-logs migrate-up migrate-down migrate-status goose-create

run:
	APP_ENV=$(APP_ENV) HTTP_ADDR=$(HTTP_ADDR) DATABASE_URL=$(DATABASE_URL) JWT_SECRET=$(JWT_SECRET) LOG_LEVEL=$(LOG_LEVEL) go run ./cmd/api

test:
	go test ./...

fmt:
	go fmt ./...

gqlgen:
	go run github.com/99designs/gqlgen generate

sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc generate

db-up:
	docker compose up -d postgres

db-down:
	docker compose down

db-logs:
	docker compose logs -f postgres

migrate-up:
	$(GOOSE) -dir migrations postgres "$(DATABASE_URL)" up

migrate-down:
	$(GOOSE) -dir migrations postgres "$(DATABASE_URL)" down

migrate-status:
	$(GOOSE) -dir migrations postgres "$(DATABASE_URL)" status

goose-create:
	$(GOOSE) -dir migrations create $(name) sql
