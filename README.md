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

Migrations are **not** run automatically on API startup (Cloud Run can start multiple/concurrent instances, and GORM's `AutoMigrate` has no locking). There is currently no dedicated migrate command — run them manually when the schema changes by temporarily uncommenting the `migrations.Migrate(database)` call in `main.go`, running the app once against the target `DATABASE_URL`, then re-commenting it. A dedicated `cmd/migrate` binary is a planned follow-up.
