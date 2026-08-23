APP_ENV ?= development
HTTP_ADDR ?= :8080
DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/qhira?sslmode=disable
JWT_SECRET ?= change-me-in-development
LOG_LEVEL ?= info

.PHONY: run test fmt gqlgen migrate-up migrate-down migrate-status goose-create

run:
	APP_ENV=$(APP_ENV) HTTP_ADDR=$(HTTP_ADDR) DATABASE_URL=$(DATABASE_URL) JWT_SECRET=$(JWT_SECRET) LOG_LEVEL=$(LOG_LEVEL) go run ./cmd/api

test:
	go test ./...

fmt:
	go fmt ./...

gqlgen:
	go run github.com/99designs/gqlgen generate

migrate-up:
	goose -dir migrations postgres "$(DATABASE_URL)" up

migrate-down:
	goose -dir migrations postgres "$(DATABASE_URL)" down

migrate-status:
	goose -dir migrations postgres "$(DATABASE_URL)" status

goose-create:
	goose -dir migrations create $(name) sql
