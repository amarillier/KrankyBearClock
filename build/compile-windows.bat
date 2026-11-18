@echo off
REM Wrapper for compile-windows.ps1: starting PowerShell through a batch file behaves better over ssh.
cd /d "%~dp0.."
powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%~dp0compile-windows.ps1" %*
