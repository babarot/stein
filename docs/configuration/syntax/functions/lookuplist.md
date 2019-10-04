# lookuplist(map, key)

Returns a list matched by the key in the given map.

Like the Terraform's [`lookup`](https://www.terraform.io/docs/configuration/interpolation.html#lookup-map-key-default-) but this is only for returning a list.

## Type

Arguments | Return values
---|---
map, string | list(string)

## Usage

```hcl
variable "colors" {
  type = "map"

  default = {
    "red" = [
      "burgundy",
      "terracotta",
      "scarlet",
    ]
    "blue" = [
      "heliotrope",
      "cerulean blue",
      "turquoise blue",
    ]
  }
}
```

```hcl
"${lookuplist(var.colors, "red")}"
# => ["burgundy", "terracotta", "scarlet"]

"${contains(lookuplist(var.colors, "red"), "scarlet")}"
# => true
```
