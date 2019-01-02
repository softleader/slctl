@echo off

set BUILD=_build
set DIST=_dist
set BINARY=%1
set VERSION=%2
set COMMIT=%3
set LDFLAGS="-X main.version=%VERSION% -X main.commit=%COMMIT%"
set MAIN=./cmd/%BINARY%
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64
set CHOCO_SERVER=http://ci.softleader.com.tw:8081/repository/choco/
set CHOCO_USER=choco:choco

if not exist %BUILD% mkdir %BUILD%
if not exist %DIST% mkdir %DIST%
go build -o %BUILD%/%BINARY%.exe -ldflags %LDFLAGS% -a -tags netgo %MAIN%
copy README.md %BUILD%
copy LICENSE %BUILD%
copy .nuspec %BUILD%
choco pack --version %VERSION% --outputdirectory %DIST% %BUILD%/.nuspec
curl -X PUT -F "file=@%DIST%/%BINARY%.%VERSION%.nupkg" %CHOCO_SERVER% -u %CHOCO_USER% -v