---
title: "Writing Stein rules"
date: 2019-01-17T15:26:15Z
draft: false
weight: 52

---

Let's say you want to create a lint policy for the next YAML file.

```yaml
apiVersion: v1
metadata:
  name: my-service
  # namespace: echo
spec:
  selector:
    app: MyApp
  ports:
  - protocol: TCP
    port: 80
    targetPort: 9376
```

This is Kubernetes YAML of [Service](https://kubernetes.io/docs/concepts/services-networking/service/) manifest.
The field `metadata.namespace` in Service can be omitted.
However, let's say you want to define it explicitly and force the owner to specify this.
In such a case, [Stein's Rule]({{< ref "/configuration/policy/rules" >}}) is useful.
A rule is simple block which can be represented by simple DSL schema by using HCL.

The rule suitable for this case is as follows.

```hcl
rule "namespace_specification" {
  description = "Check namespace name is not empty"

  conditions = [
    "${jsonpath("metadata.namespace") != ""}",
  ]

  report {
    level   = "ERROR"
    message = "Namespace is not specified"
  }
}
```

The most important attributes in rule block is `conditions` list.

This list is a collections of boolean values.
If this list contains one or more ***false*** values, this rule will fail.
The failed rule will output an error message according to the report block.

By the way, `jsonpath` is provided as a built-in function.
The available functions are here: [Interpolation Syntax]({{< ref "/configuration/syntax/interpolation" >}}).
