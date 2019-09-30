---
title: "Command: fmt"
item: "fmt"
date: 2019-01-17T15:26:15Z
draft: false
weight: 62

---

The `stein fmt` command formats a policy source to a canonical format.

```
Usage: stein fmt [options] FILE ...
```

This command formats all the specified policy files to a canonical format.

By default, policy files are overwritten in place. This behavior can be changed with the `-write` flag. If a specified FILE is - then stdin is read and the output is always written to stdout.

The command-line flags are all optional. The list of available flags are:

- [`-write=true`](#write-true) - Write formatted policy to the named source file. If false, output will go to stdout. If multiple files are specified, the output will be concatenated directly.
- [`-check=false`](#check-false) - Don't format, only check if formatting is necessary. Files that require formatting are printed, and a non-zero exit code is returned if changes are required.

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
