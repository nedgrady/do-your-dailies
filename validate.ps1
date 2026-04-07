$ErrorActionPreference = 'Stop'

$rootDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$dockerfilePath = "$rootDir\server\Dockerfile.validate"
$validationImage = 'do-your-dailies-validate:local'

docker build -f $dockerfilePath -t $validationImage $rootDir
docker run --rm -v "${rootDir}:/workspace" $validationImage
