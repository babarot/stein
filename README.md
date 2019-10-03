stein
=====

[![][circleci-badge]][circleci-link] [![][release-badge]][release-link] [![][license-badge]][license-link] [![][report-budge]][report-link]

[circleci-link]:  https://circleci.com/gh/b4b4r07/stein/tree/master
[circleci-badge]: https://circleci.com/gh/b4b4r07/stein/tree/master.svg?style=svg
[release-badge]: https://img.shields.io/github/release/b4b4r07/stein.svg?style=popout
[release-link]:  https://github.com/b4b4r07/stein/releases
[license-badge]: https://img.shields.io/github/license/b4b4r07/stein.svg?style=popout
[license-link]:  ./LICENSE
[report-budge]: https://goreportcard.com/badge/github.com/b4b4r07/stein
[report-link]:  https://goreportcard.com/report/github.com/b4b4r07/stein

Stein is a linter for config files with a customizable rule set.
Supported config file types are JSON, YAML and HCL for now.

The basic design of this tool are heavily inspired by [HashiCorp Sentinel](https://www.hashicorp.com/sentinel) and its lots of implementations come from [Terraform](https://www.terraform.io/).

## Motivation

As the motivation of this tool, the factor which accounts for the most of them is the [Policy as Code](https://b4b4r07.github.io/stein/concepts/policy-as-code/).

Thanks to [Infrastructure as Code](https://en.wikipedia.org/wiki/Infrastructure_as_code), the number of cases that the configurations of its infrastructure are described as code is increasing day by day.
Then, it became necessary to set the lint or policy for the config files.
As an example: the namespace of Kubernetes to be deployed, the number of replicas of Pods, the naming convention of a namespace, etc.

This tool makes it possible to describe those requests as code (called as the [rules](https://b4b4r07.github.io/stein/configuration/policy/rules/)).

## Try to use!

After [installing stein](https://b4b4r07.github.io/stein/intro/install/), run `make run` to quick start!

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

[Stein Documentations](https://b4b4r07.github.io/stein/)

## License

MIT

## Author

b4b4r07
