@echo off
:loop
cls

SET filename=quickmenu.exe

FOR /F "usebackq" %%A IN ('%filename%') DO SET /A beforeSize=%%~zA

: Build https://golang.org/cmd/go/
go build -v -ldflags="-w -s -H windowsgui" -o %filename%

FOR /F "usebackq" %%A IN ('%filename%') DO SET /A size=%%~zA
SET /A diffSize = %size% - %beforeSize%
SET /A size=(%size%/1024)+1
IF %diffSize% EQU 0 (
    ECHO %size%kb
) ELSE (
    IF %diffSize% GTR 0 (
        ECHO %size% kb [+%diffSize% b]
    ) ELSE (
        ECHO %size% kb [%diffSize% b]
    )
)

: IF %ERRORLEVEL%==0 start /wait %filename%

pause
goto loop
