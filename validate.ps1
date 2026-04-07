param(
	[switch]$RebuildImage,
	[switch]$NoQuiet
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

docker run --rm -v "${rootDir}:/workspace" $validationImage
