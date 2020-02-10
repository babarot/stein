stein
=====

[![test](https://github.com/b4b4r07/stein/workflows/test/badge.svg)][test]
[![release](https://github.com/b4b4r07/stein/workflows/release/badge.svg)][release]
[![docs](https://github.com/b4b4r07/stein/workflows/docs/badge.svg)][docs]

[test]: https://github.com/b4b4r07/stein/actions?query=workflow%3Atest
[release]: https://github.com/b4b4r07/stein/actions?query=workflow%3Arelease
[docs]: https://github.com/b4b4r07/stein/actions?query=workflow%3Adocs

[![][release-badge]][release-link] [![][license-badge]][license-link] [![][report-badge]][report-link] [![][go-version-badge]][go-version-link] [![][website-badge]][website-link]

[release-badge]: https://img.shields.io/github/release/b4b4r07/stein.svg?style=popout
[release-link]:  https://github.com/b4b4r07/stein/releases

[license-badge]: https://img.shields.io/github/license/b4b4r07/stein.svg?style=popout
[license-link]:  https://b4b4r07.mit-license.org

[report-badge]: https://goreportcard.com/badge/github.com/b4b4r07/stein
[report-link]:  https://goreportcard.com/report/github.com/b4b4r07/stein

[go-version-badge]: https://img.shields.io/github/go-mod/go-version/b4b4r07/stein
[go-version-link]:  https://github.com/b4b4r07/stein/blob/master/go.mod

[website-badge]: https://img.shields.io/website?down_color=lightgrey&down_message=down&up_color=green&up_message=up&url=https%3A%2F%2Fbabarot.me%2Fstein
[website-link]:  https://babarot.me/stein

Stein is a linter for config files with a customizable rule set.
Supported config file types are JSON, YAML and HCL for now.

The basic design of this tool are heavily inspired by [HashiCorp Sentinel](https://www.hashicorp.com/sentinel) and its lots of implementations come from [Terraform](https://www.terraform.io/).

![](https://user-images.githubusercontent.com/4442708/66107167-8a83f800-e5fa-11e9-9719-f7f03624ee46.png)

## Motivation

As the motivation of this tool, the factor which accounts for the most of them is the [Policy as Code](https://b4b4r07.github.io/stein/concepts/policy-as-code/).

Thanks to [Infrastructure as Code](https://en.wikipedia.org/wiki/Infrastructure_as_code), the number of cases that the configurations of its infrastructure are described as code is increasing day by day.
Then, it became necessary to set the lint or policy for the config files.
As an example: the namespace of Kubernetes to be deployed, the number of replicas of Pods, the naming convention of a namespace, etc.

This tool makes it possible to describe those requests as code (called as the [rules](https://b4b4r07.github.io/stein/configuration/policy/rules/)).

## Documentations

[Stein Documentations][website-link]
