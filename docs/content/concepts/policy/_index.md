---
title: "Policy"
date: 2017-10-17T15:26:15Z
draft: false
weight: 20

---

A policy is the collection of Stein rules written in HCL.

## What's a policy?

## Why needs policies?

Nowadays, thanks to the infiltration of the concept of "Infrastructure as Code", many infrastructure settings are coded in a configuration file language such as YAML.

For YAML files that we have to maintain like Kubernetes manifest files, as we continue to maintain them in the meantime, we will want to unify the writing style with policies like a style guide.

## How policies are loaded by `stein`

To understand how `stein` loads policy files and recognizes them is very important for writing and applying policies to the files effectively.
`stein apply` requires always one or more arguments only.
It assumes the config file paths such as YAML, JSON and so on.

The path may have a hierarchical structure.
In Stein, when a path with a hierarchical structure is given as arguments, `stein` recognizes the HCL file in `.policy` directory placed in the path included in that path as a policy to be applied.

Let's see a concrete example.

```
_examples
|-- .policy/
|   |-- config.hcl
|   |-- functions.hcl
|   |-- rules.hcl
|   `-- variables.hcl
|-- manifests/
|   |-- .policy/
|   |   |-- functions.hcl
|   |   `-- rules.hcl
|   `-- microservices/
|       |-- x-echo-jp/
|       |   `-- development/
|       |       |-- Deployment/
|       |       |   |-- redis-master.yaml
|       |       |   |-- test.yaml
|       |       |   `-- test.yml
|       |       |-- PodDisruptionBudget/
|       |       |   `-- pdb.yaml
|       |       `-- Service/
|       |           `-- service.yaml
|       `-- x-gateway-jp/
|           `-- development/
|               `-- Deployment/
|                   `-- test.yaml
`-- spinnaker/
    |-- .policy/
    |   `-- functions.hcl
    `-- x-echo-jp/
        `-- development/
            `-- deploy-to-dev-v2.yaml
```

There are some Kubernetes YAML with hierarchical structure and some policies here.

In this case, `stein` recognizes these HCL files as the policy to be applied to the arguments if `_examples/manifests/microservices/x-echo-jp/development/Deployment/test.yaml` is given as arguments of `stein`:

- `_examples/.policy/*.hcl`
- `_examples/manifests/.policy/*.hcl`

This is because given argument file contains `_examples/` and `_examples/manifests`.

That is, all YAML files located in `_examples/manifests/` is applied with `_examples/.policy/*.hcl` and `_examples/manifests/.policy/*.hcl`.

On the other hand, all YAML files located in `_examples/spinnaker/` is applied with `_examples/.policy/*.hcl` and `_examples/spinnaker/.policy/*.hcl`.

So, you can control the policy to apply by appropriately creating the directory and placing the YAML files and `.policy` directory there.

In addition, if you want to apply policies placed in places that have no relation to given arguments, you can control by environment variable or `apply` flag.

```bash
export STEIN_POLICY=/path/to/policy
stein apply deployment.yaml

# or

stein apply -policy /path/to/policy deployment.yaml
```

Also `STEIN_POLICY` (`-policy`) can take multiple values separated by a comma, also can take directories and files:

```bash
STEIN_POLICY=root-policy/,another-policy/special.hcl
# -> these files are applied, besides ".policy/*.hcl" included in given arguments
#    root-policy/*.hcl
#    another-policy/special.hcl
```

## What's next

- [Rules](rules.md)
- [Variables](variables.md)
- [Functions](functions.md)
- [Config](config.md)
