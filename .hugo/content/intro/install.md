---
title: "Installing Stein"
date: 2019-01-17T15:26:15Z
draft: false
weight: 51

---

Installing Stein is simple. There are two approaches to installing Stein:

1. Using a precompiled binary
2. Installing from source

Downloading a precompiled binary is easiest, and we provide downloads over TLS along with SHA256 sums to verify the binary. We also distribute a PGP signature with the SHA256 sums that can be verified.

## Precompiled Binaries

To install the precompiled binary, download the appropriate package for your system. Stein is currently packaged as a zip file. We do not have any near term plans to provide system packages.

[Releases · b4b4r07/stein](https://github.com/b4b4r07/stein/releases)

Once the zip is downloaded, unzip it into any directory. The stein binary inside is all that is necessary to run Stein (or `stein.exe` for Windows). Any additional files, if any, aren't required to run Stein.

Copy the binary to anywhere on your system. If you intend to access it from the command-line, make sure to place it somewhere on your `PATH`.

## Compiling from Source

To compile from source, you will need [Go](https://golang.org/) installed and configured properly (including a `GOPATH` environment variable set), as well as a copy of [git](https://www.git-scm.com/) in your `PATH`.

1. Clone the Stein repository from GitHub into your `GOPATH`:

    ```console
    $ mkdir -p $GOPATH/src/github.com/b4b4r07 && cd $_
    $ git clone https://github.com/b4b4r07/stein.git
    $ cd stein
    ```

2. Build Stein for your current system and put the binary in ./bin/ (relative to the git checkout). The make dev target is just a shortcut that builds stein for only your local build environment (no cross-compiled targets).

    ```console
    $ make build
    ```

## Verifying the Installation

To verify Stein is properly installed, run `stein -h` on your system. You should see help output. If you are executing it from the command line, make sure it is on your `PATH` or you may get an error about Stein not being found.

```console
$ stein -h
```
