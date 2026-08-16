@echo off
setlocal

call build.bat
if errorlevel 1 exit /b %errorlevel%

call run.bat
