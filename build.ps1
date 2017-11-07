param([switch]$SKIP_TESTS, [switch]$SKIP_ANALYZE)

$ConfigDir = "$env:APPDATA\.config\cftool"

$ConfigDirExists = Test-Path $ConfigDir
if ($ConfigDirExists -eq $False) {
    New-Item -ItemType directory -Path $ConfigDir >$null 2>&1
}
Copy-Item config.yaml "$ConfigDir\config.yaml"

Write-Host "$ConfigDir\config.yaml:"
Get-Content "$ConfigDir\config.yaml"

go get -t -d .\...
if ($? -ne $True) {
    Exit 1
}

if ($SKIP_TESTS -eq $False) {
    Invoke-Expression "go test github.com/Appliscale/cftool/... -cover"
    if ($? -ne $True) {
        Exit 1
    }
}

if ($SKIP_ANALYZE -eq $False) {
    go tool vet .\
}

Invoke-Expression "go install github.com/Appliscale/cftool"
if ($? -ne $True) {
    Exit 1
}