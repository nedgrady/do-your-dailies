# do-your-dailies

## Mutation testing

Use WSL on Windows for mutation testing.

Native PowerShell is not supported for this toolchain. Run mutation tests inside Ubuntu WSL.

### Windows setup (WSL)

From PowerShell, install Ubuntu WSL once:

```powershell
wsl --install -d Ubuntu
```

Then open Ubuntu and install Go (example version):

```bash
cd /tmp
curl -LO https://go.dev/dl/go1.25.5.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.5.linux-amd64.tar.gz
echo 'export PATH=/usr/local/go/bin:$PATH' >> ~/.profile
source ~/.profile
go version
```

Set local toolchain mode so Go does not try to auto-download another version:

```bash
go env -w GOTOOLCHAIN=local
```

### Run mutation tests

In Ubuntu WSL:

```bash
cd /mnt/c/Code/do-your-dailies/server
go test ./...
go install github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest
$(go env GOPATH)/bin/go-mutesting ./internal/api
```

`./validate.ps1` now includes this mutation step and runs it through WSL.

### Troubleshooting

- If you see toolchain not available, your go.mod version is newer than your installed Go version. Install a matching Go version or lower the go directive in server/go.mod.
- If go-mutesting is not found, run it with full path: $(go env GOPATH)/bin/go-mutesting.
- If you want validate.ps1 to use a specific distro, set WSL_DISTRO_NAME before running it.
