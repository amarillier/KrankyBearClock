# prepare-deps.ps1 - Prepare Go dependencies (Windows)
# Downloads modules, pre-fetches Fyne/OpenGL deps, tidies and verifies go.mod (module cache; sync excludes vendor/).

$ErrorActionPreference = "Stop"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "KrankyBear Template - Dependency Setup" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Check if Go is installed
$goCmd = Get-Command go -ErrorAction SilentlyContinue
if (-not $goCmd) {
    Write-Host "Error: Go is not installed or not in PATH" -ForegroundColor Red
    Write-Host "Please install Go from https://golang.org/dl/" -ForegroundColor Yellow
    exit 1
}

# Display Go version
$goVersion = & go version
Write-Host "✓ Found Go: $goVersion" -ForegroundColor Green
Write-Host ""

# Check if go.mod exists
if (-not (Test-Path "go.mod")) {
    Write-Host "Error: go.mod not found in current directory" -ForegroundColor Red
    exit 1
}

$env:GOWORK = "off"
$gf = $env:GOFLAGS
if ($gf -and $gf -match "-mod=vendor") {
    Write-Host "NOTE: Replacing -mod=vendor in GOFLAGS with -mod=mod (sync2windows excludes vendor/)." -ForegroundColor Yellow
    $gf = $gf -replace "-mod=vendor", "-mod=mod"
    $env:GOFLAGS = $gf
}

Write-Host "Step 1: Downloading all dependencies..." -ForegroundColor Yellow
Write-Host "Running: go mod download"
try {
    & go mod download
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Dependencies downloaded successfully" -ForegroundColor Green
    } else {
        Write-Host "✗ Failed to download dependencies" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "✗ Failed to download dependencies: $_" -ForegroundColor Red
    exit 1
}
Write-Host "Running: go mod download github.com/go-gl/gl (Fyne / Windows desktop)"
try {
    & go mod download github.com/go-gl/gl github.com/go-gl/glfw/v3.4/glfw
    if ($LASTEXITCODE -ne 0) {
        Write-Host "[FAIL] go-gl download failed" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "[FAIL] go-gl download failed: $_" -ForegroundColor Red
    exit 1
}
Write-Host ""

Write-Host "Step 2: Tidying module dependencies..." -ForegroundColor Yellow
Write-Host "Running: go mod tidy"
try {
    & go mod tidy
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Module dependencies tidied" -ForegroundColor Green
    } else {
        Write-Host "✗ Failed to tidy dependencies" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "✗ Failed to tidy dependencies: $_" -ForegroundColor Red
    exit 1
}
Write-Host ""

Write-Host "Step 3: Verifying module dependencies..." -ForegroundColor Yellow
Write-Host "Running: go mod verify"
try {
    & go mod verify
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Module dependencies verified" -ForegroundColor Green
    } else {
        Write-Host "⚠ Module verification had issues (this may be normal)" -ForegroundColor Yellow
    }
} catch {
    Write-Host "⚠ Module verification had issues (this may be normal)" -ForegroundColor Yellow
}
Write-Host ""


# Count dependencies (with error handling)
$directDeps = 0
$indirectDeps = 0
$vendorCount = 0

try {
    $requireLines = Select-String -Path "go.mod" -Pattern "^\s+\S+" -ErrorAction SilentlyContinue
    if ($requireLines) {
        $directDeps = @($requireLines | Where-Object { $_.Line -notmatch "// indirect" }).Count
    }
} catch {
    $directDeps = 0
}

try {
    $indirectLines = Select-String -Path "go.mod" -Pattern "// indirect" -ErrorAction SilentlyContinue
    if ($indirectLines) {
        $indirectDeps = @($indirectLines).Count
    }
} catch {
    $indirectDeps = 0
}

try {
    if (Test-Path "vendor") {
        $vendorDirs = Get-ChildItem -Path "vendor" -Directory -Recurse -ErrorAction SilentlyContinue
        if ($vendorDirs) {
            $vendorCount = @($vendorDirs | Measure-Object).Count
        }
    }
} catch {
    $vendorCount = 0
}

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Dependency Setup Complete!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Summary:"
Write-Host "  • Direct dependencies: $directDeps"
Write-Host "  • Indirect dependencies: $indirectDeps"
Write-Host "  • Vendor packages: $vendorCount"
Write-Host ""
Write-Host "You can now build the application with:" -ForegroundColor Green
Write-Host "  go build -trimpath -ldflags="-s -w" -o <name>.exe"
Write-Host ""

# Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning. - Winston Churchill, November 10, 1942
