#!/usr/bin/env pwsh

$GGHInstallDir = "$env:USERPROFILE\.ggh"
$GGHCliName = "ggh.exe"
$GGHCliPath = "${GGHInstallDir}\${GGHCliName}"
Write-Host $GGHInstallDir
$githubHeader = @{}
if (-not ((Get-CimInstance Win32_ComputerSystem)).SystemType -match "x64-based PC")
{
    Write-Host "Automatic GGH install is available for x64 PC's only`n"
    return 1
}

if ((Get-ExecutionPolicy) -gt 'RemoteSigned' -or (Get-ExecutionPolicy) -eq 'ByPass')
{
    Write-Host "PowerShell requires an execution policy of 'RemoteSigned'."
    Write-Host "To make this change please run:"
    Write-Host "'Set-ExecutionPolicy RemoteSigned -scope CurrentUser'"
    break
}

[Net.ServicePointManager]::SecurityProtocol = "tls12, tls11, tls"

Write-Host "Installing GGH" -ForegroundColor DarkCyan
Write-Host "Creating the directory in $GGHInstallDir" -ForegroundColor Green
New-Item -ErrorAction Ignore -Path $GGHInstallDir -ItemType "directory"
if (!(Test-Path $GGHInstallDir -PathType Container))
{
    throw "Could not create $GGHInstallDir"
}

$githubHeader.Accept = "application/octet-stream"
Invoke-WebRequest -Headers $githubHeader -Uri "https://github.com/byawitz/ggh/releases/latest/download/ggh_windows_x86_64.exe" -OutFile $GGHCliPath

if (!(Test-Path $GGHCliPath -PathType Leaf))
{
    throw "Failed to download GGH"
}

Write-Host "Attempting to add $GGHInstallDir to User Path Environment variable..."
$UserPathEnvionmentVar = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($UserPathEnvionmentVar -like "*$GGHInstallDir*")
{
    Write-Host "GGH already in the path, skipping..."  -ForegroundColor Cyan
}
else
{
    [System.Environment]::SetEnvironmentVariable("PATH", $UserPathEnvionmentVar + ";$GGHInstallDir", "User")
    $UserPathEnvionmentVar = [Environment]::GetEnvironmentVariable("PATH", "User")
    Write-Host "Added $GGHInstallDir to User Path" -ForegroundColor Green
}

Write-Host "`r`GGH was installed successfully to $GGHInstallDir"  -ForegroundColor Green
Write-Host "`r`Restart the terminal to start use GGH"