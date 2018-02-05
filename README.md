# Perun [![Build Status](https://travis-ci.org/Appliscale/perun.svg?branch=master)](https://travis-ci.org/Appliscale/perun) [![GoDoc](https://godoc.org/github.com/Appliscale/perun?status.svg)](https://godoc.org/github.com/Appliscale/perun)

A swiss army knife for *AWS CloudFormation* templates - validation, conversion, generators and other various stuff.

## Goal

Perun was created to support work with CloudFormation templates. CloudFormation works in a way that it runs template online in AWS infrastructure and fails after first error - in many cases it is related with particular name length (e.g. maximum length is 64 characters). Instead of doing a round-trip, we would like to detect such cases locally.

## Building and Installation

### OSX

If you are using *homebrew* just run:

```bash
brew install Appliscale/tap/perun
```

### Building from sources

First of all you need to download Perun to your GO workspace:

```bash
$GOPATH $ go get github.com/Appliscale/perun
$GOPATH $ cd perun
```

Then build and install configuration for the application inside perun directory by executing:

```bash
perun $ make config-install
perun $ make all
```

With first command a default configuration file (`defaults/main.yaml`) will be copied to your home directory under the `~/.config/perun/main.yaml` path. After second command application will be compiled as a `perun` binary inside `bin` directory in your `$GOPATH/perun` workspace.

## Working with Perun

### Commands

To validate your template with AWS API (*online validation*), just type:

```bash
~ $ perun --mode=validate --template=<PATH TO YOUR TEMPLATE>
```

To validate your template offline (*well*, almost offline :wink: - *AWS CloudFormation Resource Specification* still needs to be downloaded for a fresh installation) use `validate_offline` mode:

```bash
~ $ perun --mode=validate_offline --template=<PATH TO YOUR TEMPLATE>
```

To convert your template between JSON and YAML formats you have to type:

```bash
~ $ perun --mode=convert
           --template=<PATH TO YOUR INCOMING TEMPLATE>
           --output=<PATH FOR A CONVERTED FILE, INCLUDING FILE NAME>
           --format=<JSON or YAML>
```
To create your own configuration file use mode configure:

```bash
~ $ perun --mode=configure
```
Then type path and name of new configuration file.

### Configuration file

You can find an example configuration file in the main directory of the repository in file `defaults/main.yml`.

perun supports multiple configuration files for different locations. Configuration files take precedence, according to the typical `UNIX` convention. The application will be looking for the configuration file in the following order:

1. CLI argument (`-c=<CONFIG FILE>, --config=<CONFIG FILE>`).
2. Current working directory (`.perun` file).
3. Current user local configuration (`~/.config/perun/main.yaml`).
4. System global configuration (`/etc/perun/main.yaml`).

Having a configuration file is mandatory. Minimal configuration file requires only *AWS CloudFormation Resource Specification* URLs, listed under `SpecificationURL` key:

```yaml
SpecificationURL:
  us-east-2: "https://dnwj8swjjbsbt.cloudfront.net"
  ...
```

There are two optional parameters:

* `Profile` (`default` taken by default, when no value found inside configuration files).
* `Region` (`us-east-1` taken by default, when no value found inside configuration files).

### Supporting  MFA

If you account is using *MFA* (which we strongly recommend to enable) you should add `--mfa` flag to the each executed command.

```bash
~ $ perun --mode=validate --template=<PATH TO YOUR TEMPLATE> --mfa
```

In that case application will use `[profile]-long-term` from the `~/.aws/credentials` file (`[profile]` is a placeholder filled with adequate value taken from configuration files).

Example profile you need to setup - in this case `default`:

```ini
[default-long-term]
aws_access_key_id = <YOUR ACCESS KEY>
aws_secret_access_key = <YOUR SECRET ACCESS KEY>
mfa_serial = <IDENTIFICATION NUMBER FOR MFA DEVICE>
```

## License

[Apache License 2.0](LICENSE)

## Maintainers

- [Piotr Figwer](https://github.com/pfigwer)
- [Wojciech Gawroński](https://github.com/afronski)
- [Kacper Patro](https://github.com/morfeush22)
- [Paweł Pikuła](https://github.com/ppikula)

## Contributors

- *Here is a place for you!*
