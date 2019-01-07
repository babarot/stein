# Variable Configuration

Input variables serve as parameters for a Terraform module.

When used in the root module of a configuration, variables can be set from CLI arguments and environment variables. For [child modules](), they allow values to pass from parent to child.

Input variable usage is introduced in the Getting Started guide section [Input Variables]().

This page assumes you're familiar with the [configuration syntax]() already.

## Example

Input variables can be defined as follows:

```hcl
variable "key" {
  type = "string"
}

variable "images" {
  type = "map"

  default = {
    us-east-1 = "image-1234"
    us-west-2 = "image-4567"
  }
}

variable "zones" {
  default = ["us-east-1a", "us-east-1b"]
}
```

## Description

The `variable` block configures a single input variable for a Terraform module. Each block declares a single variable.

The name given in the block header is used to assign a value to the variable via the CLI and to reference the variable elsewhere in the configuration.

Within the block body (between `{ }`) is configuration for the variable, which accepts the following arguments:

*WIP*
