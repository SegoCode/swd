@ECHO off

:: --------------------------------------------------------------------------------
::
:: Author: SegoCode
::
:: Script function:
::	Create release binary for golang script on different OS
::
:: Note:
::  You need golang.org/x/sys/unix dependence in your system if
::  you want to compile for linux, run 'go get golang.org/x/sys/unix'
::
::  For compiling whiout go.mod file set GO111MODULE to off, this
::  isnt recommended, if you dont have go.mod run 'go mod init main'
::  and 'go mod tidy' for create the go.mod file
::
::  Be sure there is no binary from any previous build in the root directory
::
:: --------------------------------------------------------------------------------

:: --------------------------------------------------------------------------------
:: Change the values ​​shown below they must match the name of the file to compile
SET gofile=swd.go 
SET gofilename=swd
:: --------------------------------------------------------------------------------

CD %~dp0

ECHO Compiling for: linux-amd64. . .
SET GOOS=linux
SET GOARCH=amd64

go build -o "%gofilename%-%GOOS%-%GOARCH%" %gofile%

ECHO Compiling for: linux-386. . .
SET GOOS=linux
SET GOARCH=386

go build -o "%gofilename%-%GOOS%-%GOARCH%" %gofile%

ECHO Compiling for: linux-arm. . .
SET GOOS=linux
SET GOARCH=arm

go build -o "%gofilename%-%GOOS%-%GOARCH%" %gofile%

ECHO Compiling for: windows-386. . .
SET GOOS=windows
SET GOARCH=386

go build -o "%gofilename%-%GOOS%-%GOARCH%.exe" %gofile%

ECHO Compiling for: windows-amd64. . .
SET GOOS=windows
SET GOARCH=amd64

go build -o "%gofilename%-%GOOS%-%GOARCH%.exe" %gofile% 

ECHO DONE!
pause > nul