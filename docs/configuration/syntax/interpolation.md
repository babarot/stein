---
title: "Interpolation Syntax"
item: "Interpolation"
date: 2019-01-17T15:26:15Z
draft: false
weight: 51

---

Embedded within strings in Terraform, whether you're using the Terraform syntax or JSON syntax, you can interpolate other values. These interpolations are wrapped in `${}`, such as `${var.foo}`.

The interpolation syntax is powerful and allows you to reference variables, attributes of resources, call functions, etc.

You can perform simple math in interpolations, allowing you to write conditions such as `${count.index + 1}`. And you can also use conditionals to determine a value based on some logic.

You can escape interpolation with double dollar signs: `$${foo}` will be rendered as a literal `${foo}`.

## Available Variables

There are a variety of available variable references you can use.

### User string variables

Use the `var.prefix` followed by the variable name. For example, `${var.foo}` will interpolate the `foo` variable value.

### User map variables

The syntax is `var.MAP["KEY"]`. For example, `${var.amis["us-east-1"]}` would get the value of the `us-east-1` key within the amis map variable.

### User list variables

The syntax is `"${var.LIST}"`. For example, `"${var.subnets}"` would get the value of the subnets list, as a list. You can also return list elements by index: `${var.subnets[idx]}`.

### Path information

*WIP*

The syntax is `path.TYPE`. `TYPE` can be `file`, `dir`, or `policies`. cwd will interpolate the current working directory. module will interpolate the path to the current module. root will interpolate the path of the root module. In general, you probably want the path.module variable.

```hcl
"${path.file}"
# => manifests/microservices/x-gateway-jp/development/Service/a.yaml
# [Notes]
#  this variable is an alias of `filename` variable

"${path.dir}"
# => manifests/microservices/x-gateway-jp/development/Service
```

### Predefined variables

- `filename`: Filename to be applied policy (alias of `path.policy`)

### Environment variables information

The syntax is `env.ENV`. `ENV` can be `USER`, `HOME`, etc. These values comes from `env` command output.

```hcl
"${env.HOME}"
# => /home/username

"${env.EDITOR}"
# => vim
```

## Conditionals

Interpolations may contain conditionals to branch on the final value.

```hcl
"${var.user == "john" ? var.member : env.anonymous}"

# => var.member (if var.user is john)
# => var.anonymous (if var.user is not john)
```

The conditional syntax is the well-known ternary operation:

```
CONDITION ? TRUEVAL : FALSEVAL
```

The condition can be any valid interpolation syntax, such as variable access, a function call, or even another conditional. The true and false value can also be any valid interpolation syntax. The returned types by the true and false side must be the same.

The supported operators are:

- Equality: `==` and `!=`
- Numerical comparison: `>`, `<`, `>=`, `<=`
- Boolean logic: `&&`, `||`, unary `!`

## Built-in Functions

Stein ships with built-in functions. Functions are called with the syntax `name(arg, arg2, ...)`. For example, to read a file: `${file("path.txt")}`.

Stein supports all Terraform's built-in functions listed in [this page](https://www.terraform.io/docs/configuration/interpolation.html#built-in-functions).

In addition to these functions, it also comes with the original built-in functions to make it even easier to write rules.

### `glob(pattern)`

Returns the files matched by given pattern

Types:

- input args: `string`
- return values: `list of string`

Usage:

```hcl
"${glob("*.txt")}"
# => ["a.txt", "b.txt", "c.txt"]
```

### `pathshorten(path)`

Returns the file path shortened like [Vim's one](http://vimdoc.sourceforge.net/htmldoc/eval.html#pathshorten()).

Types:

- input args: `string`
- return values: `string`

Usage:

```hcl
"${pathshorten("manifests/microservices/x-gateway-jp/development/Service/a.yaml")}"
# => "m/m/x/d/S/a.yaml"
```

### `ext(path)`

Returns the file extensions.

Types:

- input args: `string`
- return values: `string`

Usage:

```hcl
"${ext("a.txt")}"
# => ".txt"
```

### `wc(text, ["l", "w", "c"])`

Returns the counted number of text as options (**l**ines, **w**ords, **c**hars). Default option is **l**ines.

Types:

- input args: `string`, `string`... (Optional)
- return values: `number`

Usage:

```hcl
"${wc("foo\nbar baz")}"
# => 1

"${wc("foo\nbar baz", "c")}"
# => 11

"${wc("foo\nbar baz", "w")}"
# => 3
```

### `grep(text, pattern)`

Returns the text block matched with the given pattern.

Types:

- input args: `string`
- return values: `string`

Usage:

```
My life didn't please me,
so I created my life.
- Coco Chanel
```

```hcl
"${grep(file("text.txt"), "life")}"
# => "My life didn't please me,\nso I created my life."
```

### `lookuplist(map, key)`

Returns a list matched by the key in the given map.

Like the Terraform's [`lookup`](https://www.terraform.io/docs/configuration/interpolation.html#lookup-map-key-default-) but this is only for returning a list.

Types:

- input args: `map`, `string`
- return values: `list of string`

Usage:

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

### `match(pattern, text)`

Returns a true if the text is matched with the pattern.

Types:

- input args: `string`, `string`
- return values: `boolean`

Usage:

```hcl
"${match("abcdef", "^a")}"
# => true
```

### `color(str, attrs...)`

Returns a string colorized by the color name.

Types:

- input args: `string`, `string`...
- return values: `string`

Usage:

```hcl
"${color("hello!", "white")}"
# => "\x1b[37mhello!\x1b[0m"

"${color("hello!", "red", "BgBlack")}"
# => "\x1b[31m\x1b[40mhello!\x1b[0m"
```

### `exist(path)`

Returns true if path exists

Types:

- input args: `string`
- return values: `boolean`

Usage:

```hcl
"${exist("/path/to/whatever")}"
# => true (if exists)
```

### `jsonpath(query, default...)`

WIP

## Custom Functions

While supporting some useful built-in functions, Stein allows to create user-defined functions.

```hcl
function "add" {
  params = [a, b]
  result = a + b
}
```

```hcl
"${add(1, 3)}"
# => 4
```

For more details, please see also [Custom Functions](custom-functions.md)

## Math

Almost the same as [Terraform Math](https://www.terraform.io/docs/configuration/interpolation.html#math) mechanism.
