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

## Installation

```console
$ go get github.com/b4b4r07/stein
```

or

```bash
VERSION="X.X.X"
OS="darwin" # or "linux"
wget "https://github.com/b4b4r07/stein/releases/download/v${VERSION}/stein_${OS}_amd64.zip"
unzip -n "stein_${OS}_amd64.zip"
```

## Try to use!

Example is here: [`test.sh`](./scripts/test.sh) (`make run`). It also responds to `-h` and `--help` as you'd expect.
To view a list of the available commands at any time, just run `stein` with no arguments.

```console
$ stein apply _examples/spinnaker/*/*/*
_examples/spinnaker/x-echo-jp/development/deploy-to-dev-v2.yaml (Block 1)
  No violated rules

_examples/spinnaker/x-echo-jp/development/deploy-to-dev-v2.yaml (Block 2)
  [ERROR]  rule.namespace_name            Namespace name "x-echo-jp-prod" is invalid

=====================
1 error(s), 0 warn(s)
```

## Documentations

- [Concepts](docs/concepts.md)
- [Commands (CLI)](docs/commands.md)
- [Writing Policy](docs/writing-policy.md)

## License

MIT

## Author

b4b4r07
