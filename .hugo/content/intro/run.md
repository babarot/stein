---
title: "Running Stein"
date: 2019-01-17T15:26:15Z
draft: false
weight: 53

---

## Apply stein rules

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

For more details about this behavior, see also [How policies are loaded by Stein]({{< ref "/configuration/load" >}}).

## Debug stein

Stein has detailed logs which can be enabled by setting the `STEIN_LOG` environment variable to any value. This will cause detailed logs to appear on stderr.

You can set `STEIN_LOG` to one of the log levels `TRACE`, `DEBUG`, `INFO`, `WARN` or `ERROR` to change the verbosity of the logs. `TRACE` is the most verbose and it is the default if `STEIN_LOG` is set to something other than a log level name.

To persist logged output you can set `STEIN_LOG_PATH` in order to force the log to always be appended to a specific file when logging is enabled. Note that even when `STEIN_LOG_PATH` is set, `STEIN_LOG` must be set in order for any logging to be enabled.

If you find a bug with Stein, please include the detailed log by using a service such as gist.

## What's next

- [Writing Policy](writing-policy.md)
