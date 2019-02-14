---
title: "Config Configuration"
item: "Config"
date: 2019-01-17T15:26:15Z
draft: false
weight: 43

---

The `config` is a block that can describe settings related to stein lint. Basically stein configuration is based on "Smart default" concept. It means that it has been set up sufficiently from the beginning. Moreover, this means that you can use it without having to define this block and no need to change the setting. However, depending on the item, you may want to customize it. Therefore, you can change the setting according to the `config` block accordingly.

## Example

A `config` configuration looks like the following:

```hcl
config {
  report {
    format = "${format("{{.Level}}:  {{.Rule}}  %s", color("{{.Message}}", "white"))}"
    style  = "console"
    color  = true
  }
}
```

## Description

Only one `config` block can be defined.

Within the block (the `{ }`) is configuration for the config block.

### Meta-parameters

There are **meta-parameters** available to config block:

- `report` (configuration block) -
    - `format` (string) - Report message format. It's shown in lint message. In format, it can be described with [Go template](https://golang.org/pkg/text/template/). `{{.Level}}` is converted to the lint level, `{{.Rule}}` is converted to the rule name, and `{{.Message}}` is converted to the lint message.
    - `style` (string) - Report style. It can take "console", "inline" now.
    - `color` (bool) - Whether to color output.

If config block isn't defined, the following configuration is used by default.

```hcl
config {
  report {
    format = "${format("[{{.Level}}]  {{.Rule}}  {{.Message}}")}"
    style  = "console"
    color  = true
  }
}
```

## Syntax

The full syntax is:

```hcl
config {
  [REPORT]
}
```

where REPORT is:

```hcl
report {
  format = FORMAT
  style  = [console|inline]
  color  = bool
}
```
