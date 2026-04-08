# do-your-dailies

## Database (Docker)

```powershell
docker compose up -d postgres
```

## Mutation testing

Validation runs inside a Linux Docker container so Windows host setup stays simple.

### Run all validation checks

From PowerShell in the repository root:

```powershell
./validate.ps1
```

This script builds `server/Dockerfile.validate` and runs these checks in the container.

### Run mutation test only

From PowerShell in the repository root:

```powershell
docker build -f .\server\Dockerfile.validate -t do-your-dailies-validate:local .
docker run --rm -v "${PWD}:/workspace" -w /workspace/server do-your-dailies-validate:local bash -lc "go-mutesting ./internal/api ; go-mutesting ./internal/docs"
```

### Troubleshooting

- If Docker Desktop is not running, start it and rerun `./validate.ps1`.
- If Docker cannot mount your drive, enable file sharing for the drive in Docker Desktop settings.
