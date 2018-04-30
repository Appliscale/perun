# Perun [![Build Status](https://travis-ci.org/Appliscale/perun.svg?branch=master)](https://travis-ci.org/Appliscale/perun) [![GoDoc](https://godoc.org/github.com/Appliscale/perun?status.svg)](https://godoc.org/github.com/Appliscale/perun)

<p align="center">
<img src="perun_logo.png" alt="Perun logo">
</p>

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

#### Validation
To validate your template with AWS API (*online validation*), just type:

```bash
~ $ perun validate <PATH TO YOUR TEMPLATE>
```

To validate your template offline (*well*, almost offline :wink: - *AWS CloudFormation Resource Specification* still needs to be downloaded for a fresh installation) use `validate_offline` mode:

```bash
~ $ perun validate_offline <PATH TO YOUR TEMPLATE>
```

#### Conversion
To convert your template between JSON and YAML formats you have to type:

```bash
~ $ perun convert
           <PATH TO YOUR INCOMING TEMPLATE>
           <PATH FOR A CONVERTED FILE, INCLUDING FILE NAME>
           <JSON or YAML>
```
#### Configuration
To create your own configuration file use `configure` mode:

```bash
~ $ perun configure
```
Then type path and name of new configuration file.

#### Stack Creation
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
  
To destroy remote sink just type:

```bash
~ $ perun destroy-remote-sink
``` 

#### Protecting Stack

You can protect your stack by using Stack Policy file. It's JSON file where you describe which action is allowed or denied.
This example allows to all Update actions.

```ini
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

or
```bash 
~ $ perun set-stack-policy --stack=<NAME OF YOUR STACK> --template=<PATH TO YOUR TEMPLATE>
```

Perun has some default flags:

``--block``
Block all Update actions in stack.

``--unblock``
Unblock all Update actions in stack.

``--disable-stack-termination``
Protect stack from being deleted.

``--enable-stack-termination`` 
Allow to destroy stack.

You use flag instead of template.

```bash
~ $ perun set-stack-policy <NAME OF YOUR STACK> <FLAG>
```

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

### Working with stacks

Perun allows to create and destroy stacks.

To create stack it uses your template. It can be JSON or YAML format.

Example JSON template which describe S3 Bucket:

```ini
{
    "Resources" : {
        "HelloPerun" : {
            "Type" : "AWS::S3::Bucket"
        }
    }
}
```

If you want to destroy stack just type its name.
Before you create stack you should validate it with perun :wink:.

### Capabilities

If your template includes resources that can affect permissions in your AWS account, 
you must explicitly acknowledge its capabilities by adding `--capabilities=CAPABILITY` flag.

Valid values are `CAPABILITY_IAM` and `CAPABILITY_NAMED_IAM`.
You can specify both of them by adding `--capabilities=CAPABILITY_IAM --capabilities=CAPABILITY_NAMED_IAM`.

## License

[Apache License 2.0](LICENSE)

## Maintainers

- [Piotr Figwer](https://github.com/pfigwer)
- [Sylwia Gargula](https://github.com/SylwiaGargula)
- [Wojciech Gawroński](https://github.com/afronski)
- [Mateusz Piwowarczyk](https://github.com/piwowarc)

## Contributors

- [Jakub Lamparski](https://github.com/jlampar)
- [Aleksander Mamla](https://github.com/amamla)
- [Kacper Patro](https://github.com/morfeush22)
- [Paweł Pikuła](https://github.com/ppikula)
- [Michał Połcik](https://github.com/mwpolcik)
- [Tomasz Raus](https://github.com/rusty-2)
- [Maksymilian Wojczuk](https://github.com/maxiwoj)
