# Copilot Instructions — Shownest

## Project Overview

Go microservices backend for an event/venue booking platform.
Go workspace (`go.work`) at the repo root with a shared library (`pkg/`) and independent services under `services/`.

## Tech Stack

- **Language**: Go (latest stable)
- **HTTP**: Gin
- **Database**: PostgreSQL 18 via `pgx/v5` + `pgxpool`; queries built with `squirrel` (always use `sq.Dollar` placeholder format); rows scanned with `scany/v2 (pgxscan)`
- **Cache / Locking**: Redis via `go-redis/v9`
- **Auth**: JWT (HMAC-SHA256) — access + refresh tokens; Gin middleware `JWTAuth` and `RequireMerchant` from `pkg/middleware`
- **Logging**: Zap via `pkg/logger`
- **Config**: `pkg/config` — `ConfigProvider` interface; local `config.json` per service; AWS Secrets Manager in prod

## Architecture

Each service follows a strict three-layer pattern:

```
Handler → UseCase → Repository
```

- **Handler**: binds/validates the HTTP request, calls the use case, writes the HTTP response
- **UseCase**: all business logic; calls repository or external service clients
- **Repository**: all DB/cache operations; no business logic

Dependency injection is manual, wired in `wire/wire.go` via `InitializeApp()`. Do not introduce a DI framework.

## Service Structure

```
services/<name>/
  config.json          ← local runtime config
  main.go
  wire/wire.go         ← manual DI entry point
  internal/
    config/            ← config struct + loader
    routes/routes.go   ← Gin router setup
    handlers/          ← one file per domain entity
    usecases/usecases.go
    repository/repository.go
    models/models.go   ← DB row structs (db tags)
    dto/
      request/         ← HTTP input structs (json/uri/form tags + binding)
      response/        ← HTTP output structs (json tags)
    client/            ← typed HTTP clients for other services
    jobs/              ← background goroutines (ticker-based)
    utils/constants.go ← service-scoped constants
```

Cross-service calls use typed clients in `internal/client/`. Internal (service-to-service) routes are mounted under `/internal/` and require no auth.

## Error Handling

Always use the `apperrors` package from `pkg/errors`:

```go
// new error
apperrors.New(apperrors.CodeInvalidArgument, "seat already booked")

// wrap an existing error
apperrors.Wrap(apperrors.CodeDBError, "failed to fetch seat", err)
```

Available codes: `CodeInvalidArgument`, `CodeDBNotFound`, `CodeDBError`, `CodeInternal`, `CodeUnauthenticated`, `CodeTokenExpired`, `CodeTokenInvalid`, `CodePermissionDenied`, `CodeAlreadyExists`, `CodeFailedPrecondition`, `CodeUserBlocked`, `CodeResourceExhausted`.

HTTP error responses are written via `pkg/utils.WriteError(c, err)`.

## Naming Conventions

- **Route path params**: always descriptive — `:bookingId`, `:showtimeId`, `:eventId`, `:venueId`, `:hallId`, `:sessionId`. Never use plain `:id`.
- **JSON fields**: `camelCase`
- **DB columns**: `snake_case`
- **Constants**: `PascalCase` in `utils/constants.go`
- **Go identifiers**: follow standard Go conventions (`MustUserID`, `WriteError`, etc.)

## Code Style

- Minimal comments — only add a comment when it explains *why*, not *what*
- No over-engineering: no helper functions for one-time use, no unnecessary abstractions
- No extra error handling for impossible cases — validate only at system boundaries
- Do not add docstrings or type annotations to code you did not change
- Do not refactor or "improve" code outside the scope of the task

## Background Jobs

Background jobs live in `internal/jobs/` and are started from `wire/wire.go` before the HTTP server. Jobs use `time.NewTicker` inside a goroutine that selects on `ctx.Done()` for graceful shutdown. A `StartJobs(ctx, repo)` aggregator function in `jobs/scheduler.go` launches all jobs.

## After Every File Change

Run `gofmt` and verify the build compiles cleanly:

```sh
gofmt -w <changed_file.go>
go build ./...   # from inside the affected service directory
```

Fix any compilation errors before moving on.

## Database

- All migrations live under `migrations/<service>/schema.sql`
- UUIDs generated with `pgcrypto` (`gen_random_uuid()`)
- Soft deletes via `deleted_at TIMESTAMPTZ`
- `set_updated_at()` trigger on all tables with `updated_at`
- Use `squirrel` with `sq.Dollar` for all query construction

## Shared Library (`pkg/`)

Never duplicate anything that already exists in `pkg/`. Key exports:
- `pkg/errors` — `AppError`, all error codes
- `pkg/utils` — `MustUserID`, `WriteError`, UUID helper, column joiner
- `pkg/middleware` — `JWTAuth`, `RequireMerchant`
- `pkg/jwt` — JWT service
- `pkg/db` — pgxpool init
- `pkg/cache` — Redis client init
- `pkg/logger` — Zap logger