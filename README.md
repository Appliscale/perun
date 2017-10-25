# CFTool [![Build Status](https://travis-ci.org/Appliscale/cftool.svg?branch=master)](https://travis-ci.org/Appliscale/cftool)
A tool for CloudFormation template validation and conversion.

## Goal
CFTool was created to support work with CloudFormation templates. CloudFormation works in a way that it runs template online
in AWS infrastructure and fails after first error - in many cases it is related with particular name length (e.g. maximum
length is 64 characters). Instead of doing a round-trip, we would like to detect such cases locally. 

## Working with CFTool
First of all you need to download CFTool to your GO workspace:

`go get github.com/Appliscale/cftool`

Then install the application by going to cftool directory and typing:

`./build.sh`

Configuration file (config.yaml) will be copied to `/etc/.Appliscale/cftool`

The application should be compiled to `cftool` binary file to the `bin` directory in your GO workspace.

To validate your template online, with AWS API, just type:

`cftool -mode=validate -file=[path to your template] -region=[region, e.g. eu-central-1]`

To validate template offline (well, almost offline - AWS CloudFormation Resource Specification still needs be downloaded) use validate_offline mode:

`cftool -mode=validate_offline -file=[path to your template] -region=[region]`

To convert your template from JSON to YAML and form YAML to JSON type:

`cftool -mode=convert -file=[path to your template] -output=[path to place, where you want to save converted file] 
-format=[JSON or YAML]`
