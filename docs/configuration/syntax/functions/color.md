# color(text, color)

Returns a string colorized by the color name.

## Type

Arguments | Return values
---|---
string, string | string

## Usage

```hcl
"${color("hello!", "white")}"
# => "\x1b[37mhello!\x1b[0m"

"${color("hello!", "red", "BgBlack")}"
# => "\x1b[31m\x1b[40mhello!\x1b[0m"
```
