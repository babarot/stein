# jsonpath(query)

This function returns the value corresponding to given query. The query should be followed JSONPATH. However, as of now, this jsonpath function uses [tidwall/gjson](https://github.com/tidwall/gjson) package as JSONPATH internally. So basically for now, please refer to [its godoc page](https://godoc.org/github.com/tidwall/gjson).

## Type

Arguments | Return values
---|---
string | string / number / list / map

## Usage

Let's say you run some queries with `jsonpath` function in your rule file against the following [Kubernetes Deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) manifest file.

!!! Info "A example manifest file for Kubernetes Deployment"

    These config files listed with three format are the same. Technically Kubernetes allows to accept only JSON and YAML as manifest file. A HCL code listed here is just a example for explaining HCL is the compatible for JSON and YAML.

    ``` json tab="JSON"
    {
      "apiVersion": "extensions/v1beta1",
      "kind": "Deployment",
      "metadata": {
        "name": "nginx"
      },
      "spec": {
        "replicas": 2,
        "template": {
          "metadata": {
            "labels": {
              "run": "nginx"
            }
          },
          "spec": {
            "containers": [
              {
                "name": "nginx",
                "image": "nginx:1.11",
                "ports": [
                  {
                    "containerPort": 80
                  }
                ]
              }
            ]
          }
        }
      }
    }
    ```

    ```yaml tab="YAML"
    apiVersion: extensions/v1beta1
    kind: Deployment
    metadata:
      name: nginx
    spec:
      replicas: 2
      template:
        metadata:
          labels:
            run: nginx
        spec:
          containers:
          - name: nginx
            image: nginx:1.11
            ports:
            - containerPort: 80
    ```

    ```terraform tab="HCL"
    "apiVersion" = "extensions/v1beta1"

    "kind" = "Deployment"

    "metadata" = {
      "name" = "nginx"
    }

    "spec" = {
      "replicas" = 2

      "template" "metadata" "labels" {
        "run" = "nginx"
      }

      "template" "spec" {
        "containers" = {
          "image" = "nginx:1.11"

          "name" = "nginx"

          "ports" = {
            "containerPort" = 80
          }
        }
      }
    }
    ```

First, let's say to use this query against above manifest file.

```hcl
jsonpath("spec.replicas")
```


It will return `2` (type is number).
