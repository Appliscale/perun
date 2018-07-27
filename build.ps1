param([string]$go)

function Invoke-AndFailOnError {
    param([string]$cmd)

    Invoke-Expression $cmd
    if ($LastExitCode -ne 0) {
        Exit $LastExitCode
    }
}

$ErrorActionPreference = "Stop"

$ConfigDir = "$env:LOCALAPPDATA\.config\perun"
$PoliciesDir = "$env:LOCALAPPDATA\.config\perun\stack-policies"

$ConfigDirExists = Test-Path $ConfigDir
if ($ConfigDirExists -eq $False) {
    New-Item -ItemType directory -Path $ConfigDir >$null 2>&1
}

$PoliciesDirExists = Test-Path $PoliciesDir
if ($PoliciesDirExists -eq $False) {
    New-Item -ItemType directory -Path $PoliciesDir >$null 2>&1
}

Copy-Item defaults\main.yaml "$ConfigDir\main.yaml"
Copy-Item defaults\style.yaml "$ConfigDir\style.yaml"
Copy-Item defaults\specification_inconsistency.yaml "$ConfigDir\specification_inconsistency.yaml"
Copy-Item defaults\blocked.json "$PoliciesDir\blocked.json"
Copy-Item defaults\unblocked.json "$PoliciesDir\unblocked.json"

Invoke-AndFailOnError "$go get -t .\..."
Invoke-AndFailOnError "$go install .\..."
Invoke-AndFailOnError "$go build ."
Invoke-AndFailOnError "$go fmt .\..."
Invoke-AndFailOnError "$go vet .\..."

$Mockgen = "$env:GOPATH\bin\mockgen"
Invoke-AndFailOnError "$Mockgen -source '.\awsapi\cloudformation.go' -destination '.\stack\mocks\mock_aws_api.go' -package mocks CloudFormationAPI"

Invoke-AndFailOnError "$go test -cover .\..."
