# Stein CLI Commands

Stein is a command-line interface (CLI) for developing and testing policies.
Having a standard workflow to develop policies is critical for our mission of policy as code.
The CLI takes a subcommand to execute.
The complete list of subcommands is below.

The Stein CLI is a well-behaved command line application.
In erroneous cases, a non-zero exit status will be returned.
It also responds to -h and --help as you'd expect.
To view a list of the available commands at any time, just run `stein` with no arguments.

## Command: `apply`

The `stein apply` command is used to execute a policy locally for development purposes.

```
Usage: stein apply [options] POLICY
```

This command executes the policy file at the path specified by POLICY.

The output will indicate whether the policy passed or failed. The exit code also reflects the status of the policy: 0 is pass, 1 is fail, 2 is undefined (fail, but because the result was undefined), and 2 is a runtime error.

A configuration file can be specified with `-config` to define available import plugins, mock data, and global values. This is used to simulate a policy embedded within an application. The documentation for this configuration file is below.

The command-line flags are all optional. The list of available flags are:

- `-policy=file[,file,dir,...]` - Path to HCL file path or a directory path located in HCL files. You can specify multiple paths (directory or just HCL file) with a comma. The `STEIN_POLICY` variable is the environment variable version of this flag.

See also [How policies are loaded by `stein` - Policy](policy.md#how-policies-are-loaded-by-stein).

## Command: `fmt`

The `stein fmt` command formats a policy source to a canonical format.

```
Usage: stein fmt [options] FILE ...
```

This command formats all the specified policy files to a canonical format.

By default, policy files are overwritten in place. This behavior can be changed with the `-write` flag. If a specified FILE is - then stdin is read and the output is always written to stdout.

The command-line flags are all optional. The list of available flags are:

- `-write=true` - Write formatted policy to the named source file. If false, output will go to stdout. If multiple files are specified, the output will be concatenated directly.
- `-check=false` - Don't format, only check if formatting is necessary. Files that require formatting are printed, and a non-zero exit code is returned if changes are required.

```console
$ stein fmt -check test.hcl
   config {
-
     report {
       format = "{{.Level}}: {{.Rule}}: {{.Message}}"
-      style = "console"
-      color = true
+      style  = "console"
+      color  = true
     }
-
   }
```
