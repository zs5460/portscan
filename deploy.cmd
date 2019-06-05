@echo off
cls
rmdir /q /s dest
mkdir dest

SET CGO_ENABLED=1
SET GOOS=windows
SET GOARCH=amd64
go build -o dest/portscan.exe

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o dest/portscan



