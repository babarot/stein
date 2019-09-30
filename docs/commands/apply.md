---
title: "Command: apply"
item: "apply"
date: 2019-01-17T15:26:15Z
draft: false
weight: 61

---

The `stein apply` command is used to execute a policy locally for development purposes.

```
Usage: stein apply [options] POLICY
```

This command executes the policy file at the path specified by POLICY.

The output will indicate whether the policy passed or failed. The exit code also reflects the status of the policy: 0 is pass, 1 is fail, 2 is undefined (fail, but because the result was undefined), and 2 is a runtime error.

A configuration file can be specified with `-policy` to define available import plugins, mock data, and global values. This is used to simulate a policy embedded within an application. The documentation for this configuration file is below.

The command-line flags are all optional. The list of available flags are:

- [`-policy=file[,file,dir,...]`](#policy-file) - Path to HCL file path or a directory path located in HCL files. You can specify multiple paths (directory or just HCL file) with a comma. The `STEIN_POLICY` variable is the environment variable version of this flag.

See also [How policies are loaded by Stein](../configuration/load.md)
