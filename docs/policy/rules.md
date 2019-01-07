# Rule Configuration

The most important thing you'll configure with Stein are rules. Rules are a component of your policies. It might be some rule set such as a region to be deployed, naming convention, or some linting. Or it can be a higher level component such as an email provider, DNS record, or database provider.

This page assumes you're familiar with the [configuration syntax]() already.

## Example

A `rule` configuration looks like the following:

```hcl
rule "replicas" {
  description = "Check the number of replicas is sufficient"

  expressions = [
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
- `ignore_cases` (list of bools) - Conditions to ignore execution of the rule. If it contains one or more True condition, this rule will be skipped.
- `expressions` (list of bools) - Conditions for deciding whether this rule passes or fails. In order to pass, all conditions must return True.
- `report` (configuration block) -
    - `level` (string) - Error level. It can take "ERROR" or "WARN" as the level. In case of "ERROR", this rule fails. But in case of "WARN", this rule doesn't fail.
    - `message` (string) - Error message. Let's write the conditions for passing the role here.

## Syntax

The full syntax is:

```hcl
rule NAME {
  description = DESCRIPTION

  [ignore_cases = [bool, bool, ...]]

  expressions = [bool, bool, ...]

  report {
    level = [ERROR|WARN]
    message = MESSAGE
  }
}
```
