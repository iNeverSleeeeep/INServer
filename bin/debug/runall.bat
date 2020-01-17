@echo off
for /f "tokens=1 delims=." %%i in ('dir /a-d /b ^| findstr \-') do (   
    START "%%i" %~dp0%%i.bat
    if "%%i" equ "0-center" timeout 2
)