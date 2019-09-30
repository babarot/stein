---
title: "Stein CLI Commands"
date: 2019-01-17T15:26:15Z
draft: false
weight: 60
item: "Commands (CLI)"

---

Stein is controlled via a very easy to use command-line interface (CLI). Stein is only a single command-line application: `stein`. This application then takes a subcommand such as "apply" or "plan". The complete list of subcommands is in the navigation to the left.

The Stein CLI is a well-behaved command line application. In erroneous cases, a non-zero exit status will be returned. It also responds to `-h` and `--help` as you'd expect. To view a list of the available commands at any time, just run `stein` with no arguments.

To view a list of the available commands at any time, just run stein with no arguments:

```console
$ stein
Usage: stein [--version] [--help] <command> [<args>]

Available commands are:
    apply    Applies a policy to arbitrary config files.
    fmt      Formats a policy source to a canonical format.

```

To get help for any specific command, pass the `-h` flag to the relevant subcommand. For example, to see help about the apply subcommand:

```console
$ stein apply -h
Usage of apply:
  Applies a policy to arbitrary config files.

Options:
  -policy string
        path to the policy files or the directory where policy files are located

```

## Shell Tab-completion

TBD

## Debug stein

Stein has detailed logs which can be enabled by setting the `STEIN_LOG` environment variable to any value. This will cause detailed logs to appear on stderr.

You can set `STEIN_LOG` to one of the log levels `TRACE`, `DEBUG`, `INFO`, `WARN` or `ERROR` to change the verbosity of the logs. `TRACE` is the most verbose and it is the default if `STEIN_LOG` is set to something other than a log level name.

To persist logged output you can set `STEIN_LOG_PATH` in order to force the log to always be appended to a specific file when logging is enabled. Note that even when `STEIN_LOG_PATH` is set, `STEIN_LOG` must be set in order for any logging to be enabled.

If you find a bug with Stein, please include the detailed log by using a service such as gist.
