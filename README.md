# Miniflux at rss.pinescore.com

DATETIME of last agent review: 27/11/2025 14:55 GMT

Self-hosted Miniflux RSS reader backed by PostgreSQL, running behind Apache/Virtualmin at `https://rss.pinescore.com`.

## Quick Start
- Prereqs: Go 1.24+ (we use 1.25.x from `/usr/local/go`) and PostgreSQL 15+.
- Setup: `make miniflux` (builds the `miniflux` binary from the repo root).
- Run (local dev): `DATABASE_URL=postgres://miniflux:<password>@localhost/miniflux?sslmode=disable RUN_MIGRATIONS=1 go run main.go`.

## Development
- Test: `make test` (runs Go tests with race detection).
- Lint/format: `make lint` (go vet + gofmt + golangci-lint if installed).
- Useful: `make integration-test` (starts a temporary instance and runs API integration tests against PostgreSQL).

## Architecture
- `main.go` and `internal/cli` — CLI entrypoint that parses flags, loads configuration, and starts the HTTP server and background workers.
- `internal/http` — HTTP handlers, routing, authentication, and UI endpoints.
- `internal/database` and `internal/storage` — PostgreSQL access and persistence helpers.
- `internal/reader`, `internal/worker`, `internal/mediaproxy` — feed fetching, background jobs, and media proxying.
- `client` and `internal/template/templates` — front-end assets and HTML templates bundled into the binary.

## Configuration
- `DATABASE_URL` (required) — PostgreSQL connection string (`postgres://miniflux:<password>@localhost/miniflux?sslmode=disable` on this host).
- `RUN_MIGRATIONS` (recommended on first run) — when `1`, applies database migrations automatically.
- `LISTEN_ADDR` — bind address for the HTTP server (production uses `127.0.0.1:8180` behind Apache).
- `BASE_URL` — external URL Miniflux should use when generating links (`https://rss.pinescore.com` here).
- `CREATE_ADMIN`, `ADMIN_USERNAME`, `ADMIN_PASSWORD` — optional one-time bootstrap of the initial admin account.

## Troubleshooting
- HTTP server fails with “bind: address already in use” → update `LISTEN_ADDR` to a free port or stop the conflicting service, then restart Miniflux.
- Login keeps failing for the admin user → reset the password via the UI and ensure `CREATE_ADMIN=0` so the environment is not overwriting it.
- Service will not start under systemd → run `sudo systemctl status miniflux` and `sudo journalctl -u miniflux` to inspect logs and check for database or config errors.
- Database-related errors (e.g. “database does not exist” or auth failures) → confirm `DATABASE_URL` matches an existing PostgreSQL database and user.

## Deployment
- This host runs Miniflux as a `miniflux` systemd service behind Apache/Virtualmin; see `ops/deployment.md` for PostgreSQL, systemd, and vhost configuration used on `rss.pinescore.com`.

## Links
- Miniflux upstream documentation — https://miniflux.app/docs/
