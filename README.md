stein
=====

[![][release-badge]][release-link] [![][license-badge]][license-link]

[release-badge]: https://img.shields.io/github/release/b4b4r07/stein.svg?style=popout
[release-link]:  https://github.com/b4b4r07/stein/releases
[license-badge]: https://img.shields.io/github/license/b4b4r07/stein.svg?style=popout
[license-link]:  ./LICENSE

Stein is a linter for config files with a customizable rule set.
Supported config file types are JSON, YAML and HCL for now.

The basic design of this tool are heavily inspired by [HashiCorp Sentinel](https://www.hashicorp.com/sentinel) and its lots of implementations come from [Terraform](https://www.terraform.io/).

## Motivation

As the motivation of this tool, the factor which accounts for the most of them is the [Policy as Code](docs/policy-as-code.md).

Thanks to [Infrastructure as Code](https://en.wikipedia.org/wiki/Infrastructure_as_code), the number of cases that the configurations of its infrastructure are described as code is increasing day by day.
Then, it became necessary to set the lint or policy for the config files.
As an example: the namespace of Kubernetes to be deployed, the number of replicas of Pods, the naming convention of a namespace, etc.

This tool makes it possible to describe those requests as code (called as the [rules](docs/policy/rules.md)).

## Documentations

- [Concepts](docs/concepts.md)
- [Commands (CLI)](docs/commands.md)
- [Writing Policy](docs/writing-policy.md)

## License

MIT

## Author

b4b4r07
