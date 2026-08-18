@echo off
setlocal
cd /d "%~dp0\..\.."

set "PS1_FILE=%TEMP%\spinner.ps1"
echo param($cmd, $msg) > "%PS1_FILE%"
echo $p = Start-Process cmd -ArgumentList "/c $cmd" -PassThru -WindowStyle Hidden >> "%PS1_FILE%"
echo $spinner = @('^|','/','-','\') >> "%PS1_FILE%"
echo $i = 0 >> "%PS1_FILE%"
echo Write-Host -NoNewline "$msg...  " >> "%PS1_FILE%"
echo while (-not $p.HasExited) { >> "%PS1_FILE%"
echo     Write-Host -NoNewline "`b$($spinner[$i])" >> "%PS1_FILE%"
echo     $i = ($i + 1) %% 4 >> "%PS1_FILE%"
echo     Start-Sleep -Milliseconds 100 >> "%PS1_FILE%"
echo } >> "%PS1_FILE%"
echo if ($p.ExitCode -eq 0) { Write-Host "`b`b Done!" -ForegroundColor Green } >> "%PS1_FILE%"
echo else { Write-Host "`b`b Failed!" -ForegroundColor Red; exit $p.ExitCode } >> "%PS1_FILE%"

echo 🚀 Starting DeepPhotos Build Process...
type nul > build.log

powershell -ExecutionPolicy Bypass -File "%PS1_FILE%" "cd frontend && npm install >> ..\build.log 2>&1" "Installing Frontend Dependencies"
if errorlevel 1 exit /b %errorlevel%

powershell -ExecutionPolicy Bypass -File "%PS1_FILE%" "cd frontend && npm run build >> ..\build.log 2>&1" "Building Frontend (SvelteKit)"
if errorlevel 1 exit /b %errorlevel%

powershell -ExecutionPolicy Bypass -File "%PS1_FILE%" "cd backend && go build -o server.exe ./cmd/server >> ..\build.log 2>&1" "Compiling Backend (Go)"
if errorlevel 1 exit /b %errorlevel%

del "%PS1_FILE%"
echo ✅ DeepPhotos Build Complete! (See build.log for details)
