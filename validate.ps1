param(
	[switch]$RebuildImage,
	[switch]$NoQuiet,
	[switch]$Quick,
	[switch]$Verbose
)

$ErrorActionPreference = 'Stop'

$rootDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$dockerfilePath = "$rootDir\server\Dockerfile.validate"
$validationImage = 'do-your-dailies-validate:local'

docker image inspect $validationImage *> $null
$imageExists = $LASTEXITCODE -eq 0

if ($RebuildImage -or -not $imageExists) {
	$buildArgs = @('build', '--quiet', '-f', $dockerfilePath, '-t', $validationImage, $rootDir)
	if ($NoQuiet) {
		$buildArgs = @('build', '-f', $dockerfilePath, '-t', $validationImage, $rootDir)
	}
	docker @buildArgs
}

$runArgs = @('run', '--rm', '-v', "${rootDir}:/workspace")
if ($Quick) {
	$runArgs += @('-e', 'VALIDATE_QUICK=1')
}
$runArgs += $validationImage

$summaryPattern = '==>|ALLOWLISTED FAIL|^The mutation score|^FAIL\b|^ok\b|All validation checks passed'

if ($Verbose) {
	docker @runArgs
} else {
	docker @runArgs 2>&1 | Select-String -Pattern $summaryPattern | ForEach-Object { $_.Line }
}

exit $LASTEXITCODE
