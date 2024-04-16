@ECHO OFF
set /p targetSystem="Enter target system[windows, darwin, linux]: "
set /p targetPlatform="Enter target platform[386, amd64, arm]: "
SET  CGO_ENABLED=0
SET GOOS=%targetSystem%
SET GOARCH=%targetPlatform%
go build cloud-lock-go-gin
