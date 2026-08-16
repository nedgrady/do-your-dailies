# do-your-dailies

## Database (Docker)

```powershell
docker compose up -d postgres
```

### Troubleshooting

- If Docker Desktop is not running, start it and rerun `./validate.ps1`.
- If Docker cannot mount your drive, enable file sharing for the drive in Docker Desktop settings.

## Server

### Environment variables

- `DATABASE_URL` — required, no default. A Postgres DSN (e.g. `host=localhost user=postgres password=postgres dbname=dailies port=5432 sslmode=disable`). The app fails fast at startup if this is unset. Never commit a real value.
- `PORT` — optional, defaults to `8080`.

### Local execution

```powershell
cd server
cp .env.example .env    # one-time; edit if your local DB differs
docker compose up -d postgres    # from repo root
go run .                 # or `air` for hot reload
```

### Docker execution

```powershell
docker build -t dailies-api -f server/Dockerfile server
docker run --rm -p 8080:8080 -e DATABASE_URL="host=host.docker.internal user=postgres password=postgres dbname=dailies port=5432 sslmode=disable" dailies-api
```

### Migrations

Migrations run automatically on every API boot (`migrations.Migrate(database)` in `main.go`), including in the deployed Cloud Run service. Since Cloud Run can start multiple/concurrent instances, `Migrate` serializes itself with a Postgres transaction-scoped advisory lock (`server/internal/migrations/migrate.go`) so only one instance runs the schema DDL at a time — others block briefly, then their own (idempotent) `AutoMigrate` is a no-op.
