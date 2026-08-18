@echo off
setlocal
cd /d "%~dp0\..\.."

call scripts\windows\build.bat
if errorlevel 1 exit /b %errorlevel%

call scripts\windows\run.bat
