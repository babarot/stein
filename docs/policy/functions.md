# Function Configuration

It is recommended that you read the [Custom Functions](../syntax/custom-functions.md) page prior to reading this section of the documentation. The page will explain what the custom functions are and how to use them. On the other hands, this documentation will guide you the basics of writing custom functions and introducing it into your policies efficiently.

## Example

A `function` configuration looks like the following:

```hcl
function "get_service_name" {
  params = [file]
  result = basename(dirname(dirname(dirname(file))))
}

function "get_env" {
  params = [file]
  result = basename(dirname(dirname(file)))
}

function "get_service_id_with_env" {
  params = [file]
  result = format("%s-%s", get_service_name(file), lookup(var.shortened_environment, get_env(file)))
}
```

## Description

The `function` block creates an user-defined function of the given *NAME* (first parameter). The name must be unique.

Within the block (the `{ }`) is configuration for the function.

### Meta-parameters

There are **meta-parameters** available to all rules:

- `params` (list of strings) - Parameters for the function. Like arguments. It can be referenced within the function. The variable name for `params` can specify arbitrary string.
- `variadic_param` (list of strings) - Variable arguments for the function.
- `result` (any) - Return value of the function. It can take just string of course, but also take variables, built-in functions and custom functions even.

## Syntax

The full syntax is:

```hcl
rule NAME {
  params = [ARG, ...]

  [variadic_param = [ARG, ...]]

  result = RETURN-VALUE
}
```
