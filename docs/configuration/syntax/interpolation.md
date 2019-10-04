# Interpolation Syntax

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

# Built-in Functions

Stein ships with built-in functions. Functions are called with the syntax `name(arg, arg2, ...)`. For example, to read a file: `${file("path.txt")}`.

Stein supports all Terraform's built-in functions listed in [this page](https://www.terraform.io/docs/configuration/interpolation.html#built-in-functions).

In addition to these functions, it also comes with the original built-in functions to make it even easier to write rules.

For more details, please see also ==Built-in Functions== in Navigation bar on left.

# Custom Functions

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
