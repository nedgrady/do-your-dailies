$ErrorActionPreference = 'Stop'

$rootDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$serverDir = Join-Path $rootDir 'server'

if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
  Write-Error 'go is not installed or not on PATH'
}

Write-Host "==> validating Go module in $serverDir"
Set-Location $serverDir

Write-Host '==> go mod tidy (verify go.mod/go.sum are consistent)'
go mod tidy

Write-Host '==> go fmt ./...'
go fmt ./...

Write-Host '==> go vet ./...'
go vet ./...

Write-Host '==> go test ./...'
go test ./...

if (-not (Get-Command wsl -ErrorAction SilentlyContinue)) {
  Write-Error 'WSL is required for mutation testing. Install WSL, then rerun validate.ps1.'
}

$drive = $serverDir.Substring(0, 1).ToLower()
$pathWithoutDrive = $serverDir.Substring(2).Replace('\', '/')
$wslServerDir = "/mnt/$drive$pathWithoutDrive"

$wslCommandTemplate = @'
cd "__WSL_SERVER_DIR__"; if ! command -v go >/dev/null 2>&1; then echo 'go is not installed in WSL. Install Go in your WSL distro first.'; exit 1; fi; if command -v go-mutesting >/dev/null 2>&1; then go-mutesting ./internal/api; elif [ -x "$(go env GOPATH)/bin/go-mutesting" ]; then "$(go env GOPATH)/bin/go-mutesting" ./internal/api; else echo 'go-mutesting not found in WSL.'; echo 'Install with: go install github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest'; exit 1; fi
'@

$wslCommand = $wslCommandTemplate.Replace('__WSL_SERVER_DIR__', $wslServerDir).Trim()

Write-Host '==> go-mutesting ./internal/api (WSL)'
$wslArgs = @()

$availableDistros = @()
try {
  $availableDistros = @(wsl -l -q 2>$null | ForEach-Object { $_.Trim() } | Where-Object { $_ })
} catch {
  $availableDistros = @()
}

$selectedDistro = $env:WSL_DISTRO_NAME
if (-not $selectedDistro -and $availableDistros -contains 'Ubuntu') {
  $selectedDistro = 'Ubuntu'
}

if ($selectedDistro) {
  $wslArgs += @('-d', $selectedDistro)
}

if ($env:WSL_USER) {
  $wslArgs += @('-u', $env:WSL_USER)
}

$wslArgs += @('bash', '-lc', $wslCommand)

$previousErrorActionPreference = $ErrorActionPreference
$ErrorActionPreference = 'Continue'
$wslOutput = & wsl @wslArgs 2>&1
$wslExitCode = $LASTEXITCODE
$ErrorActionPreference = $previousErrorActionPreference
if ($wslOutput) {
  $wslOutput | ForEach-Object { Write-Host $_ }
}

if ($wslExitCode -ne 0) {
  $distroHint = if ($selectedDistro) { " (distro: $selectedDistro)" } else { '' }
  Write-Error "WSL mutation testing failed$distroHint. Set WSL_DISTRO_NAME/WSL_USER if needed, then rerun validate.ps1."
}

Write-Host '==> validate passed'
