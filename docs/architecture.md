# Folder Architecture

This project is organized as a modular monolith. The main dependency flow should stay:

```text
GraphQL resolver -> service/application layer -> repository -> PostgreSQL
```

Resolvers handle GraphQL transport concerns only. Business rules and authorization live in services. SQL and persistence concerns live in repositories.

## Tree

```text
.
|-- cmd/
|   `-- api/
|-- internal/
|   |-- app/
|   |-- auth/
|   |-- authorization/
|   |-- comment/
|   |-- config/
|   |-- database/
|   |-- dataloader/
|   |-- graphql/
|   |   |-- mapper/
|   |   |-- model/
|   |   |-- resolver/
|   |   `-- schema/
|   |-- issue/
|   |-- pagination/
|   |-- platform/
|   |   |-- clock/
|   |   `-- logger/
|   |-- project/
|   |-- server/
|   |-- user/
|   `-- validation/
|-- migrations/
|-- scripts/
|-- test/
|   `-- integration/
`-- docs/
```

## Top-Level Folders

### `cmd/api`

Application entrypoint. This should contain the `main` package that loads configuration, opens infrastructure connections, wires dependencies, creates the HTTP/GraphQL server, and starts the process.

Keep business logic out of `cmd/api`. It is composition code.

### `internal`

Private application code. Packages inside `internal` are not importable by external Go modules, which is a good fit for application internals.

### `migrations`

PostgreSQL schema migrations. This is where tables, constraints, indexes, and reversible migration files should live.

Expected tables include:

- `users`
- `projects`
- `project_members`
- `issues`
- `comments`
- `issue_history`

Indexes should be based on real access patterns: project membership lookup, issues by project, issues by assignee, comments by issue, history by issue, and cursor pagination fields.

### `scripts`

Small developer scripts for repeatable local tasks such as running migrations, generating gqlgen code, starting local dependencies, or seeding development data.

Avoid putting application logic here.

### `test/integration`

Integration tests that need a real GraphQL server, PostgreSQL, migrations, or cross-package behavior. Unit tests should stay near the package they test.

### `docs`

Human-facing project documentation. Keep architectural notes, API usage examples, and important decisions here.

## Internal Packages

### `internal/app`

Application wiring layer. This package can assemble services, repositories, loaders, subscription broker, and server dependencies into one application struct.

Use it to keep `cmd/api` small.

### `internal/config`

Configuration loading and validation. It should read environment variables or config files and expose typed config structs.

Examples:

- HTTP port
- database DSN
- JWT secret and token lifetime
- log level
- GraphQL playground enablement

### `internal/server`

HTTP server setup: routes, middleware chain, GraphQL handler mounting, WebSocket transport setup, health checks, graceful shutdown.

This package should not contain business rules.

### `internal/graphql/schema`

GraphQL schema files and gqlgen configuration-adjacent schema organization.

The schema should be business-oriented, not a direct mirror of database tables.

### `internal/graphql/resolver`

gqlgen resolver implementations. Resolvers should be thin:

- read GraphQL arguments
- read authenticated user identity from context
- call services
- map domain results to GraphQL models
- return safe GraphQL errors

Do not put authorization or business workflows directly in resolvers.

### `internal/graphql/model`

GraphQL-facing model types when they are manually defined or separated from generated code.

Do not use these as domain entities unless the project intentionally decides that the coupling is acceptable for a specific simple case.

### `internal/graphql/mapper`

Mapping between domain/service models and GraphQL models.

This keeps gqlgen-generated or GraphQL-specific shapes from leaking into business logic.

### `internal/auth`

Authentication infrastructure:

- password hashing
- JWT creation and verification
- auth middleware
- authenticated user context helpers

JWT parsing belongs in middleware, not in every resolver.

### `internal/authorization`

Centralized permission rules. This package should answer questions like:

- Can this user access this project?
- Can this user manage project members?
- Can this user assign an issue?
- Can this user read a sensitive field?

Permissions must be resource-aware. A role in one project does not grant access to another project.

### `internal/user`

User domain module:

- domain model
- repository
- service
- registration/login support
- profile retrieval

Passwords must be hashed and never exposed through GraphQL.

### `internal/project`

Project and membership domain module:

- project model
- project repository
- project service
- membership model
- role management
- owner/admin/member rules

Project creation should automatically create an owner membership for the creator.

### `internal/issue`

Issue domain module:

- issue model
- repository
- service
- status and priority changes
- assignment rules
- issue history creation

Operations that update an issue and insert history should use one database transaction.

### `internal/comment`

Comment domain module:

- comment model
- repository
- service
- create/update/delete rules

Comments must enforce project access through the related issue.

### `internal/dataloader`

Request-scoped GraphQL DataLoaders for batching and caching nested resolver lookups.

Typical loaders:

- user by ID
- project by ID
- comments by issue ID
- members by project ID

Do not use global DataLoader caches.

### `internal/subscription`

GraphQL subscription event handling. Start with an in-memory Pub/Sub abstraction, but keep it replaceable by Redis Pub/Sub, NATS, or another broker later.

Subscription authorization should happen when the subscription is established.

### `internal/database`

PostgreSQL connection setup, transaction helpers, and low-level database abstractions shared by repositories.

Repositories should receive database handles from here instead of opening their own connections.

### `internal/pagination`

Shared cursor pagination helpers:

- opaque cursor encoding/decoding
- stable ordering helpers
- connection/pageInfo helpers
- pagination limit validation

Cursors should be stable and should not rely only on non-unique sortable values.

### `internal/validation`

Reusable validation helpers for application inputs:

- email format
- password requirements
- title/body limits
- pagination limits

Business-specific validation can live in services when it depends on repositories or domain rules.

### `internal/platform/logger`

Logging setup and thin wrappers around the selected logging package.

Keep logs structured where practical. Avoid permanently enabling verbose SQL/debug logs in production config.

### `internal/platform/clock`

Small time abstraction for testable timestamp behavior.

This helps services create deterministic `createdAt` and `updatedAt` values in tests without global time mocking.

## Testing Placement

Use package-local unit tests for isolated service, authorization, pagination, and validation behavior.

Use `test/integration` for tests that need PostgreSQL, migrations, GraphQL execution, subscriptions, or realistic request context with authentication and DataLoaders.

## Boundary Rules

- Resolvers call services, not repositories directly.
- Services own business rules, validation coordination, authorization decisions, and transactions.
- Repositories own SQL and persistence mapping.
- Domain packages should not depend on `internal/graphql`.
- GraphQL models should not become the source of truth for business logic.
- DataLoaders are request-scoped and used by nested resolvers.
- Filtering, sorting, and pagination happen in PostgreSQL.
- Multi-step state changes that must stay consistent use transactions.
