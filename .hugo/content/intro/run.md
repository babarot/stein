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
