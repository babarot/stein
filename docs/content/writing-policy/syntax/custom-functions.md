---
title: "Custom Functions"
date: 2017-10-17T15:26:15Z
draft: false
weight: 21

---

*Notes:* This idea basically comes from https://github.com/hashicorp/hcl2/tree/master/ext/userfunc

## What's custom functions?

The custom function feature is like an user-defined functions. You can freely define functions that Stein doesn't provide as [a built-in function](../syntax/interpolation.md).

Of course you can define it freely, so you can customize it by wrapping a built-in function, or you can use it like an alias.

## Why need custom functions?

Stein functions as a versatile testing framework for configuration files such as YAML. Stein therefore doesn't provide a function only to achieve a specific company's use case or purpose. However, there will be many cases in which you want to do it. This custom function feature covers that.

## Usage

This HCL extension allows a calling application to support user-defined functions.

Functions are defined via a specific block type, like this:

```hcl
function "add" {
  params = [a, b]
  result = a + b
}

function "list" {
  params         = []
  variadic_param = items
  result         = items
}
```

Predefined keywords to be used in `function` block is:

- `params`: Arguments for the function.
- `variadic_param`: Variable-length argument list. It can be omitted.
- `result`: Return value. It can take not only just string but also other functions, variables, etc.

## Examples

```hcl
function "remove_ext" {
  params = [file]
  result = replace(basename(file), ext(file), "")
}

# "${remove_ext("/path/to/secret.txt")}" => secret
```

```hcl
variable "shortened_environment" {
  description = "Shortened environment, such as prod, dev"
  type        = "map"

  default = {
    production  = "prod"
    development = "dev"
    laboratory  = "lab"
  }
}

function "shorten_env" {
  params = [env]
  result = lookup(var.shortened_environment, env)
}

# ${shorten_env("development")} => dev
```

## What's next

- [Function Configuration](../policy/functions.md)
