Push-Location .\server 
go generate ./...

Pop-Location
Push-Location .\client
npx orval