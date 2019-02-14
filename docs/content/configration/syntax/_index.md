---
title: "Syntax"
date: 2017-10-17T15:26:15Z
draft: false
weight: 50

---

The syntax of Stein configurations is called [HashiCorp Configuration Language (HCL)](https://github.com/hashicorp/hcl). It is meant to strike a balance between human readable and editable as well as being machine-friendly. For machine-friendliness, Stein can also read JSON configurations. For general Stein configurations, however, we recommend using the HCL Stein syntax.

## Stein Syntax

Almost all syntax comes from HCL and Terraform one. This tool is heavily inspired by those.

So for HCL syntax or some, please see the [official documentation](https://www.terraform.io/docs/configuration/syntax.html)

## Interpolation Syntax

The interpolation syntax is powerful mechanism to make it easy to define the custom rule set.
This design and implementaion comes from Terraform one.

- [Interpolation Syntax](interpolation.md)
