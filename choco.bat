@echo off

set BUILD=_build
set DIST=_dist
set BINARY=slctl
set VERSION=%1
set COMMIT=%2
set LDFLAGS="-X main.version=%VERSION% -X main.commit=%COMMIT%"
set BINARY=slctl
set MAIN=./cmd/slctl
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64
set CHOCO_SERVER=http://ci.softleader.com.tw:8081/repository/choco/
set CHOCO_USER=choco:choco

if not exist %BUILD% mkdir %BUILD%
if not exist %DIST% mkdir %DIST%
go build -o %BUILD%/%BINARY% .exe -ldflags %LDFLAGS% -a -tags netgo %MAIN%
cp README.md %BUILD% && cp LICENSE %BUILD% && cp .nuspec %BUILD%
choco pack --version %VERSION% --outputdirectory %DIST% %BUILD%/.nuspec
curl -X PUT -F "file=@%DIST%/slctl.%VERSION%.nupkg" %CHOCO_SERVER% -u %CHOCO_USER% -v