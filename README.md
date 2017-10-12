# CFTool
A tool for CloudFormation template validation and conversion.

## Goal
CFTool was created to support work with CloudFormation templates. CloudFormation works in a way that it runs template online
in AWS infrastructure and fails after first error - in many cases it is related with particular name length (e.g. maximum
length is 64 characters). Instead of doing a round-trip, we would like to detect such cases locally. 

## Working with CFTool
First of all you need to download CFTool to your GO workspace:

`go get github.com/Appliscale/cftool`

Then install the application:

`go install github.com/Appliscale/cftool`

The application should be compiled to `cftool` binary file to the `bin` directory in your GO workspace.

To validate your template just type:

`cftool -mode=validate -file=[path to your template] -region=[region, e.g. eu-central-1]`

To convert your template from JSON to YAML and form YAML to JSON type:

`cftool -mode=convert -file=[path to your template] -output=[path to place, where you want to save converted file] 
-format=[JSON or YAML]`
