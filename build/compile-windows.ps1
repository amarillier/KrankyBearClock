# Build the Fyne app for Windows (amd64) natively on the Windows build host. cgo + gcc (MSYS2 ucrt64) are required, so it cannot be
# cross-compiled from macOS. Run it from anywhere: it moves to the project root (the folder holding go.mod).
# Names come from build\app.env; the version is the VERSION file, compiled in with go:embed. The icon and version information of the exe
# are the tracked rsrc_windows_amd64.syso (regenerated on the Mac by build/winres.sh), so nothing icon-related is needed here.
# Output: bin\<KB_WINDOWS_EXE>  (a windowsgui build: no console window)
$ErrorActionPreference = 'Continue'

function Log([string]$msg) { Write-Host $msg }

$root = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
Set-Location $root

# build\app.env is KEY="value" lines shared with the bash scripts
$app = @{}
Get-Content (Join-Path $root 'build\app.env') | ForEach-Object {
    if ($_ -match '^\s*(KB_\w+)="([^"]*)"') { $app[$Matches[1]] = $Matches[2] }
}
$KB_WINDOWS_EXE = $app['KB_WINDOWS_EXE']
if (-not $KB_WINDOWS_EXE) { Log 'ERROR: KB_WINDOWS_EXE is not set in build\app.env.'; exit 1 }
$version = (Get-Content (Join-Path $root 'VERSION') -ErrorAction Stop | Select-Object -First 1).Trim()

Log ($app['KB_PROJECT_TITLE'] + " - Windows compile ($version, native Windows host)")
Log '============================================================='
Log ('Working directory: ' + (Get-Location).Path)
Log ''

if (-not (Test-Path 'bin')) { New-Item -ItemType Directory -Path 'bin' -Force | Out-Null }
Remove-Item -Path ('bin\' + $KB_WINDOWS_EXE) -Force -ErrorAction SilentlyContinue

if (-not (Get-Command go -ErrorAction SilentlyContinue)) { Log 'ERROR: go not found in PATH.'; exit 1 }
if (-not (Get-Command gcc -ErrorAction SilentlyContinue)) {
    $msys = 'C:\msys64\ucrt64\bin'
    if (Test-Path (Join-Path $msys 'gcc.exe')) { $env:PATH = "$msys;$env:PATH"; Log "Added $msys to PATH" }
    else { Log 'ERROR: gcc not found (install MSYS2 and its ucrt64 gcc toolchain).'; exit 1 }
}

if (-not (Test-Path 'go.mod')) {
    Log 'ERROR: go.mod not found in the project root (wrong folder?).'
    exit 1
}
$goModText = Get-Content -Path 'go.mod' -Raw -ErrorAction SilentlyContinue
if (-not $goModText -or $goModText -notmatch 'github\.com/go-gl/gl' -or $goModText -notmatch 'fyne\.io/fyne/v2') {
    Log 'ERROR: go.mod looks stale or wrong (missing fyne / go-gl lines).'
    Log 'If you use build/sync2windows.sh, the compile path must be the tree that receives the sync, not an old copy.'
    exit 1
}

# Fyne is pinned in go.mod (kept in sync from macOS): no `go get @latest` here. To bump Fyne intentionally, run
# `go get fyne.io/fyne/v2@<version>` on the Mac, rebuild and test, then sync.
$env:GOWORK = 'off'
$gf = $env:GOFLAGS
if ($gf -and $gf -match '-mod=vendor') {
    Log 'NOTE: Replacing -mod=vendor in GOFLAGS with -mod=mod (the sync leaves out vendor/).'
    $gf = $gf -replace '-mod=vendor', '-mod=mod'
}
# No VCS stamping, so builds work when .git is absent (the source is synced without it).
$env:GOFLAGS = if ($gf) { "$gf -buildvcs=false" } else { '-buildvcs=false' }

$env:GOOS = 'windows'
$env:GOARCH = 'amd64'
$env:CGO_ENABLED = '1'

Log ('Diagnostics: ' + (go version))
Log 'Diagnostics: go env GOMOD GOWORK GOFLAGS GOOS GOARCH CGO_ENABLED'
& go env GOMOD GOWORK GOFLAGS GOOS GOARCH CGO_ENABLED
Log ''

go mod download
if ($LASTEXITCODE -ne 0) { Log 'ERROR: go mod download failed.'; exit 1 }

Log 'Caching OpenGL modules (go-gl)...'
go mod download github.com/go-gl/gl github.com/go-gl/glfw/v3.4/glfw
if ($LASTEXITCODE -ne 0) { Log 'ERROR: go mod download github.com/go-gl/gl failed.'; exit 1 }

Remove-Item Env:\CC -ErrorAction SilentlyContinue
if (Get-Command x86_64-w64-mingw32-gcc -ErrorAction SilentlyContinue) { $env:CC = 'x86_64-w64-mingw32-gcc' }
Get-ChildItem -Path $root -Filter '*386.syso' -Recurse -ErrorAction SilentlyContinue | Remove-Item -Force -ErrorAction SilentlyContinue

Log ('Building Windows GUI binary (bin\' + $KB_WINDOWS_EXE + ')...')
go build -mod=mod -ldflags '-s -w -H windowsgui' -trimpath -o ('bin\' + $KB_WINDOWS_EXE) .
if ($LASTEXITCODE -ne 0) { Log 'ERROR: Windows build failed.'; exit 1 }
Log ('OK: bin\' + $KB_WINDOWS_EXE)

Log ''
Log '============================================================='
Log 'Compile complete.'
Get-ChildItem ('bin\' + $KB_WINDOWS_EXE) | Format-Table Name, Length -AutoSize
