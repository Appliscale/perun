# Dependencies

MIT:

* https://github.com/alecthomas/kingpin/tree/v2.2.5 - Command-line parser. It is type-safe and allows to have short version of commands (e.g. —mode, -m).

* https://github.com/ghodss/yaml - Go lacks of YAML support, so we need this one.

* https://github.com/asaskevich/govalidator - A lot of validators and sanitizers, like isCIRD() or isIP().

* https://github.com/mitchellh/mapstructure - Library for decoding generic map to go structures. We need this one for type-aware validators implementation.


Apache 2.0:

* https://github.com/aws/aws-sdk-go - AWS API.

* https://github.com/go-ini/ini - for handling AWS credential files.