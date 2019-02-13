---
title: "Rule Configuration"
date: 2017-10-17T15:26:15Z
draft: false
weight: 20

---

The most important thing you'll configure with Stein are rules. Rules are a component of your policies. It might be some rule set such as a region to be deployed, naming convention, or some linting. Or it can be a higher level component such as an email provider, DNS record, or database provider.

This page assumes you're familiar with the [configuration syntax]() already.

## Example

A `rule` configuration looks like the following:

```hcl
rule "replicas" {
  description = "Check the number of replicas is sufficient"

  conditions = [
    "${jsonpath(".spec.replicas") > 3}",
  ]

  report {
    level   = "ERROR"
    message = "Too few replicas"
  }
}
```

## Description

The `rule` block creates a rule set of the given *NAME* (first parameter). The name must be unique.

Within the block (the `{ }`) is configuration for the rule.

### Meta-parameters

There are **meta-parameters** available to all rules:

- `description` (string) - A human-friendly description for the rule. This is primarily for documentation for users using your Stein configuration. When a module is published in Terraform Registry, the given description is shown as part of the documentation.
- `depends_on` (list of strings) - Other rules which this rule depends on. This rule will be skipped if the dependency rules has failed. The rule name which will be described in "depends_on" list should follow as "rule.xxx".
- `precondition` (configuration block; optional) -
    - `cases` (list of bools) - Conditions to determine whether the rule should be executed. This rule will only be executed if all preconditions return true.
- `conditions` (list of bools) - Conditions for deciding whether this rule passes or fails. In order to pass, all conditions must return True.
- `report` (configuration block) -
    - `level` (string) - Error level. It can take "ERROR" or "WARN" as the level. In case of "ERROR", this rule fails. But in case of "WARN", this rule doesn't fail.
    - `message` (string) - Error message. Let's write the conditions for passing the role here.

## Syntax

The full syntax is:

```hcl
rule NAME {
  description = DESCRIPTION

  [depends_on = [NAME, ...]]

  [PRECONDITION]

  conditions = [CONDITION, ...]

  REPORT
}
```

where PRECONDITION is:

```hcl
precondition {
  cases = [CONDITION, ...]
}
```

where REPORT is:

```hcl
report {
  level = [ERROR|WARN]
  message = MESSAGE
}
```
