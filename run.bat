@echo off
:loop
cls

gocritic check -enableAll -disable='#experimental,#opinionated,#commentedOutCode' ./...

go build
IF %ERRORLEVEL% EQU 0 quickmenu.exe

pause
goto loop