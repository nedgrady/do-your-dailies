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

Write-Host '==> validate passed'
