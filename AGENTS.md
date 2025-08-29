# Repository Guidelines

## Project Structure & Module Organization
- `cmd/api`: HTTP service entrypoint (`main.go`), handlers, routing.
- `internal`: shared packages (e.g., `db`, `env`, `log`, `repo`).
- `cmd/migrate`: DB migrations (`migrations/`) and seeds (`seed/`).
- `docs`: OpenAPI/Swagger files (`swagger.yaml`, `swagger.json`).
- `src`: playground/examples (not part of the service binary).
- `bin`: optional build outputs.

## Build, Run, and Docs
- Run API: `go run ./cmd/api`
- Build binary: `go build -o bin/api ./cmd/api && ./bin/api`
- Generate Swagger: `make gen-docs` (requires `swag` installed).

## Migrations & Seed
- Create migration: `make migrate-create name=create_table`
- Apply migrations: `make migrate-up` | rollback: `make migrate-down`
- Current version: `make migrate-version` | force: `make migrate-force version=N`
- Seed data: `make seed`
- Configure DB via `DB_ADDR`; avoid committing secrets. Example: `export DB_ADDR=postgres://user:pass@host:port/db?sslmode=require`.

## Coding Style & Naming
- Go 1.x, `gofmt`-formatted. Format and vet: `go fmt ./... && go vet ./...`.
- Package names: short, lower-case (`repo`, `env`).
- Files: descriptive (`posts.go`, `users.go`).
- Identifiers: Exported `CamelCase`, unexported `lowerCamel`.
- Handlers return JSON consistently via helpers in `cmd/api`.

## Testing Guidelines
- Framework: standard `testing` package.
- Add tests next to code as `*_test.go` (e.g., `repo/posts_test.go`).
- Names: `TestXxx(t *testing.T)` with table tests where useful.
- Run all tests: `go test ./...` (add coverage flags as needed).

## Commit & Pull Request Guidelines
- Use Conventional Commits: `feat:`, `fix:`, `refactor:`, `docs:`, etc. (matches history).
- PRs: clear description, link issues, list migrations and env changes, and include API/Swagger updates when endpoints change.
- Include run steps for reviewers (commands, env vars). Add screenshots for API responses where helpful.

## Security & Configuration
- Prefer env vars: `ADDR`, `DB_ADDR`, `DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS`, `DB_MAX_IDLE_TIME`, `API_URL`.
- Do not commit secrets; use `.env` locally and secret managers in CI/Prod.
