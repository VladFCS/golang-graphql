# Go GraphQL Issue Tracker

Production-style educational backend for a Jira-like issue tracker built with Go, gqlgen, PostgreSQL, JWT authentication, DataLoader, and GraphQL subscriptions.

Start with the folder map in [docs/architecture.md](docs/architecture.md).

## Local Setup

Copy the example environment and adjust values for your machine:

```sh
cp .env.example .env
```

The current bootstrap expects:

- `APP_ENV`: application environment, defaults to `development`
- `HTTP_ADDR`: HTTP listen address, defaults to `:8080`
- `DATABASE_URL`: PostgreSQL connection string
- `JWT_SECRET`: secret used later for JWT signing
- `LOG_LEVEL`: `debug`, `info`, `warn`, or `error`

## Commands

Run the API:

```sh
make run
```

Run tests:

```sh
make test
```

Format code:

```sh
make fmt
```

Generate GraphQL code after creating or updating schema files:

```sh
make gqlgen
```

Run goose migrations:

```sh
make migrate-up
```

Rollback one migration:

```sh
make migrate-down
```

Check migration status:

```sh
make migrate-status
```

Create a new goose migration:

```sh
make goose-create name=create_issues
```

## Health Check

After starting the API, verify the server:

```sh
curl http://localhost:8080/healthz
```

Expected response:

```json
{"status":"ok"}
```
