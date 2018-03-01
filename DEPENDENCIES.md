# Dependencies and Rationale

## MIT

- https://github.com/alecthomas/kingpin
  - Command-line parser. It is type-safe and allows to have short version of commands (e.g. `--config`, `-c`).
- https://github.com/ghodss/yaml
  - *Go* lacks of YAML support out of the box, so we need this one.
- https://github.com/asaskevich/govalidator
  - Additional validators and sanitizers, like `isCIDR()` or `isIP()`.
- https://github.com/mitchellh/mapstructure
  - Library for decoding generic map to go structures. We need this one for type-aware validators implementation.
- https://github.com/stretchr/testify
  - Test framework and assertions.

## Apache 2.0

- https://github.com/aws/aws-sdk-go
  - AWS API.
- https://github.com/go-ini/ini
  - For handling AWS credential files.
- https://github.com/awslabs/goformation
  - Library for working with AWS CloudFormation templates (capable of resolving the intrinsic functions).