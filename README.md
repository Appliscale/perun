# CFTool [![Build Status](https://travis-ci.org/Appliscale/cftool.svg?branch=master)](https://travis-ci.org/Appliscale/cftool)
A tool for CloudFormation template validation and conversion.

[Documentation](https://godoc.org/github.com/Appliscale/cftool)

## Goal
CFTool was created to support work with CloudFormation templates. CloudFormation works in a way that it runs template online
in AWS infrastructure and fails after first error - in many cases it is related with particular name length (e.g. maximum
length is 64 characters). Instead of doing a round-trip, we would like to detect such cases locally. 

## Working with CFTool
First of all you need to download CFTool to your GO workspace:

`go get github.com/Appliscale/cftool`

Then install the application by going to cftool directory and typing:

`./build.sh`

Configuration file (config.yaml) will be copied to your home directory under the `~/.config/cftool/main.yaml` path.

The application should be compiled to `cftool` binary file to the `bin` directory in your GO workspace.

To validate your template online, with AWS API, just type:

`cftool --mode=validate --template=[path to your template]`

To validate template offline (well, almost offline - AWS CloudFormation Resource Specification still needs be downloaded) use validate_offline mode:

`cftool --mode=validate_offline --template=[path to your template]`

To convert your template from JSON to YAML and form YAML to JSON type:

`cftool --mode=convert --template=[path to your template] --output=[path to place, where you want to save converted file]
--format=[JSON or YAML]`

## Configuration file
You can find example configuration file in the main directory of the repository (`config.yml`).

It is possible to have multiple configuration files in different locations. Configuration files take precedence, according to the standard `Unix` convention. The application will be looking for the configuration file in the following order:
1. CLI argument (`-c, --config=CONFIG`)
2. Current working directory search (`.cftool` file)
3. Current user local config (`~/.config/cftool/main.yaml`)
4. Global system config (`/etc/cftool/main.yaml`)

Configuration file is mandatory. Minimal configuration file includes AWS CloudFormation Resource Specification URLs, listed under `SpecificationURL` key:
```
SpecificationURL:
  us-east-2: "https://dnwj8swjjbsbt.cloudfront.net"
  us-east-1: "https://d1uauaxba7bl26.cloudfront.net"
  us-west-1: "https://d68hl49wbnanq.cloudfront.net"
  ...
```

There are two optional parameters:
* `Profile` (`default` by default)
* `Region` (`us-east-1` by default)

## AWS MFA
If you want to use MFA add `--mfa` flag to the command:

`cftool --mode=validate --template=[path to your template] --mfa`

Application will use `[profile]-long-term` from the `~/.aws/credentials` file (`[profile]` - profile specified in config.yaml file.).

Example - profile `default`:

> \[default-long-term]

> aws_access_key_id = your access key

> aws_secret_access_key = your secret access key

> mfa_serial = the identification number of the MFA device