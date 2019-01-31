# Perun [![Release](https://img.shields.io/github/release/Appliscale/perun.svg?style=flat-square)](https://github.com/Appliscale/perun/releases/latest) [![Build_Status](https://travis-ci.org/Appliscale/perun.svg?branch=master)](https://travis-ci.org/Appliscale/perun) [![License](https://img.shields.io/badge/License-Apache%202.0-orange.svg)](https://github.com/Appliscale/perun/blob/master/LICENSE.md) [![Go_Report_Card](https://goreportcard.com/badge/github.com/Appliscale/perun?style=flat-square&fuckgithubcache=1)](https://goreportcard.com/report/github.com/Appliscale/perun) [![GoDoc](https://godoc.org/github.com/Appliscale/perun?status.svg)](https://godoc.org/github.com/Appliscale/perun)

<p align="center">
<img src="perun_logo.png" alt="Perun logo" width="400">
</p>

A command-line validation tool for *AWS Cloud Formation* that allows to conquer the cloud faster!

## Goal

Perun was created to improve work experience with CloudFormation. The idea came from the team constantly using AWS CloudFormation - it runs a template online in AWS infrastructure and fails after first error - which in many cases is trivial (e.g. maximum name length is 64 characters). Instead of doing a round-trip, we would like to detect such cases locally. 
## Building and Installation

### OSX
#### Homebrew:
```bash
$ brew install Appliscale/tap/perun
```
#### From binaries:
* Go to Perun’s releases https://github.com/Appliscale/perun/releases
* Find and download perun-darwin-amd64.tar.gz
* Unpack the archive

### Debian
#### Dpkg package manager:
* Go to https://github.com/Appliscale/perun-dpkg
* Download perun.deb
* Install:
```bash
$ dpkg -i perun.deb
```
#### From binaries:
* Go to Perun’s releases https://github.com/Appliscale/perun/releases
* Find and download perun-linux-amd64.tar.gz
* Unpack:
```bash
$ tar xvzf perun-linux-amd64.tar.gz
```

### Linux
#### Rpm package manager:
* Go to: https://github.com/Appliscale/rpmbuild/tree/master/RPMS/x86_64
* Download perun-linux-amd64-1.2.0-1.x86_64.rpm
* Install:
```bash
$ rpm -ivh perun-linux-amd64-1.2.0-1.x86_64.rpm
```

#### From binaries:
* Go to Perun’s releases https://github.com/Appliscale/perun/releases
* Find and download perun-linux-amd64.tar.gz
* Unpack:
```bash
tar xvzf perun-linux-amd64.tar.gz
```

### Building from sources

First of all you need to download Perun to your GO workspace:

```bash
$GOPATH $ go get github.com/Appliscale/perun
$GOPATH $ cd perun
```

Then build and install configuration for the application inside perun directory by executing:

```bash
perun $ make
```

After this, application will be compiled as a `perun` binary inside `bin` directory in your `$GOPATH/perun` workspace.


## Working with Perun

### Commands

#### Validation
To validate your template, just type:

```bash
~ $ perun validate <PATH TO YOUR TEMPLATE>
```
Your template will be then validated using both our validation mechanism and AWS API
(*aws validation*).

#### Configuration
To create your own configuration file use `configure` mode:

```bash
~ $ perun configure
```

Then type path and name of new configuration file.

#### Stack Parameters
Bored of writing JSON parameter files? Perun allows you to interactively create parameters file
for a given template. You can either pass the parameters interactively or as a command-line argument.

##### Command Line Parameter way:
```bash
~ $ perun create-parameters <PATH TO YOUR TEMPLATE> <OUTPUT PARAMETER FILE> --parameter=MyParameter1:<PARAMETER VALUE>
```

The greatest thing is that you can mix those in any way you want. Perun will validate the
given parameters from command line. If everything is OK, it will just create the parameters file.
If anything is missing or invalid, it will let you know and ask for it interactively.

#### Working with stacks

Perun allows to create and destroy stacks.

Cloud Formation templates can be in JSON or YAML format.

Example JSON template which describe S3 Bucket:

```json
{
    "Resources" : {
        "HelloPerun" : {
            "Type" : "AWS::S3::Bucket"
        }
    }
}
```

Before you create stack Perun will validate it by default :wink:. You can disable it with flag `--no-validate`.

To create new stack you have to type:

```bash
~ $ perun create-stack <NAME OF YOUR STACK>  <PATH TO YOUR TEMPLATE>
```

To destroy stack just type:

```bash
~ $ perun delete-stack <NAME OF YOUR STACK>
```

You can use option ``--progress`` to show the stack creation/deletion progress in the console, but
note, that this requires setting up a remote sink.

##### Remote sink

To setup remote sink type:

```bash
~ $ perun setup-remote-sink
```

This will create an sns topic and sqs queue with permissions for the sns topic to publish on the sqs
queue. Using above services may produce some cost:
According to the AWS SQS and SNS pricing:

- SNS:
  - notifications to the SQS queue are free
- SQS:
  - The first 1 million monthly requests are free.
  - After that: 0.40$ per million requests after Free Tier (Monthly)
  - Typical stack creation uses around a hundred requests
  
More information about pricing can be found [here](https://aws.amazon.com/sqs/pricing/).

To destroy remote sink just type:

```bash
~ $ perun destroy-remote-sink
```

#### Cost estimation

```bash
~ $ perun estimate-cost <PATH TO YOUR TEMPLATE>
```
To estimate template's cost run the command above with path to file. Perun resolves parameters located in the template and checks if it’s correct. Then you get url to Simple Monthly Calculator which will be filled with data from the template.

#### Protecting Stack

You can protect your stack by using Stack Policy file. It's JSON file where you describe which action is allowed or denied. This example allows to all Update Actions.

```json
{
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": "*",
      "Action": "Update:*",
      "Resource": "*"
    }
  ]
}
```

To apply your Stack Policy file you have to type:

```bash
~ $ perun set-stack-policy <NAME OF YOUR STACK>  <PATH TO YOUR TEMPLATE>
```

Perun has some default flags:

- ``--block`` - Block all Update actions in stack.

- ``--unblock`` - Unblock all Update actions in stack.

- ``--disable-stack-termination`` - Protect stack from being deleted.

- ``--enable-stack-termination`` - Allow to destroy stack.

You use flag instead of template.

```bash
~ $ perun set-stack-policy <NAME OF YOUR STACK> <FLAG>
```

### Configuration files

Perun will help you in setting up all the needed configuration files on you first run - no previous setup required.

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

There are 6 other parameters:

* `DefaultProfile` (`default` taken by default, when no value found inside configuration files).
* `DefautRegion` (`us-east-1` taken by default, when no value found inside configuration files).
* `DefaultDurationForMFA`: (`3600` taken by default, when no value found inside configuration files).
* `DefaultDecisionForMFA`: (`false` taken by default, when no value found inside configuration files).
* `DefaultVerbosity`: (`INFO` taken by default, when no value found inside configuration files).
* `DefaultTemporaryFilesDirectory`: (`.` taken by default, when no value found inside configuration files).

### Supporting  MFA

If you account is using *MFA* (which we strongly recommend to enable) you should add `--mfa` flag to the each executed command or set `DefaultDecisionForMFA` to `true` in the configuration file.

```bash
~ $ perun validate <PATH TO YOUR TEMPLATE> --mfa
```

In that case application will use `[profile]-long-term` from the `~/.aws/credentials` file (`[profile]` is a placeholder filled with adequate value taken from configuration files).

Example profile you need to setup - in this case `default`:

```ini
[default-long-term]
aws_access_key_id = <YOUR ACCESS KEY>
aws_secret_access_key = <YOUR SECRET ACCESS KEY>
mfa_serial = <IDENTIFICATION NUMBER FOR MFA DEVICE>
```

You do not need to use Perun for validation, you can just use it to obtain security credentials and use them in AWS CLI. To do this type:

```bash
~ $ perun mfa
```

### Capabilities

If your template includes resources that can affect permissions in your AWS account,
you must explicitly acknowledge its capabilities by adding `--capabilities=CAPABILITY` flag.

Valid values are `CAPABILITY_IAM` and `CAPABILITY_NAMED_IAM`.
You can specify both of them by adding `--capabilities=CAPABILITY_IAM --capabilities=CAPABILITY_NAMED_IAM`.

### Inconsistencies between official documentation and Resource Specification

Perun uses Resource Specification provided by AWS - using this we can determine if fields are required etc. Unfortunately, during the development process, we found inconsistencies between documentation and Resource Specification. These variances give rise to a mechanism that allows patching those exceptions in place via configuration. In a few words, inconsistency is the variation between information which we get from these sources.

To specify inconsistencies edit `~/.config/perun/specification_inconsistency.yaml` file.

Example configuration file:

```yaml
  SpecificationInconsistency:
    AWS::CloudFront::Distribution.DistributionConfig:
      DefaultCacheBehavior:
        - Required
```

## License

[Apache License 2.0](LICENSE)

## Maintainers

- [Maksymilian Wojczuk](https://github.com/maxiwoj)
- [Piotr Figwer](https://github.com/pfigwer)
- [Sylwia Gargula](https://github.com/SylwiaGargula)
- [Mateusz Piwowarczyk](https://github.com/piwowarc)

## Contributors

- [Wojciech Gawroński](https://github.com/afronski) (originator)
- [Jakub Lamparski](https://github.com/jlampar)
- [Aleksander Mamla](https://github.com/amamla)
- [Kacper Patro](https://github.com/morfeush22)
- [Paweł Pikuła](https://github.com/ppikula)
- [Michał Połcik](https://github.com/mwpolcik)
- [Tomasz Raus](https://github.com/rusty-2)
