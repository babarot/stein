---
title: "Getting Started"
date: 2017-10-17T15:26:15Z
draft: false
weight: 20

---

Stein is a lint tool to test config files such as YAML by using a customizable rule set called as a policy.
This page covers how to start to test your config files with stein and how to write Stein policies what you want to test.

## Install stein

At first, let's install stein command by the following command. Stein is written in Go. So just running `go get`.

```console
$ go get github.com/b4b4r07/stein
$ stein --version
```

Also you can grab the package from [GitHub Releases](https://github.com/b4b4r07/stein/releases) if you want to use stable version.

## Write rules

Let's say you want to create a lint policy for the next YAML file.

```yaml
apiVersion: v1
metadata:
  name: my-service
  # namespace: echo
spec:
  selector:
    app: MyApp
  ports:
  - protocol: TCP
    port: 80
    targetPort: 9376
```

This is Kubernetes YAML of [Service](https://kubernetes.io/docs/concepts/services-networking/service/) manifest.
The field `metadata.namespace` in Service can be omitted.
However, let's say you want to define it explicitly and force the owner to specify this.
In such a case, [Stein's Rule](./policy/rules.md) is useful.
A rule is simple block which can be represented by simple DSL schema by using HCL.

The rule suitable for this case is as follows.

```hcl
rule "namespace_specification" {
  description = "Check namespace name is not empty"

  conditions = [
    "${jsonpath("metadata.namespace") != ""}",
  ]

  report {
    level   = "ERROR"
    message = "Namespace is not specified"
  }
}
```

The most important attributes in rule block is `conditions` list.

This list is a collections of boolean values.
If this list contains one or more ***false*** values, this rule will fail.
The failed rule will output an error message according to the report block.

By the way, `jsonpath` is provided as a built-in function.
The available functions are here: [Interpolation Syntax](syntax/interpolation.md).

## Run stein

After writing up your rules, let's run stein command.

The Stein CLI is a well-behaved command line application.
In erroneous cases, a non-zero exit status will be returned.
It also responds to `-h` and `--help` as you'd expect.
To view a list of the available commands at any time, just run stein with no arguments.

To apply the rule to that YAML file and run the test you can do with the [`apply`](commands.md#command-apply) subcommand.

```console
$ stein apply -policy rules.hcl service.yaml
service.yaml
  [ERROR]  rule.namespace_specification  Namespace is not specified

=====================
1 error(s), 0 warn(s)
```

You can show the error message with exit code `1`.

The location (a file path directly or a directory path which is located policies) of policy files can be specified with `-policy` flag.
Otherwise, you can tell stein the location of policies with `STEIN_POLICY` environment variable.

Moreover, stein automatically checks `.policy` directory whether policies written in HCL are located or not when running.
So you can put it on `.policy` directory like the following:

```console
$ tree .
service.yaml
.policy/
`-- rules.hcl
```

For more details about this behavior, see also [How policies are loaded by `stein` - Policy](policy.md#how-policies-are-loaded-by-stein).

## Debug stein

Stein has detailed logs which can be enabled by setting the `STEIN_LOG` environment variable to any value. This will cause detailed logs to appear on stderr.

You can set `STEIN_LOG` to one of the log levels `TRACE`, `DEBUG`, `INFO`, `WARN` or `ERROR` to change the verbosity of the logs. `TRACE` is the most verbose and it is the default if `STEIN_LOG` is set to something other than a log level name.

To persist logged output you can set `STEIN_LOG_PATH` in order to force the log to always be appended to a specific file when logging is enabled. Note that even when `STEIN_LOG_PATH` is set, `STEIN_LOG` must be set in order for any logging to be enabled.

If you find a bug with Stein, please include the detailed log by using a service such as gist.

## What's next

- [Writing Policy](writing-policy.md)
