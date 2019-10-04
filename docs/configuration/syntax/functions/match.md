# match(text, text)

Returns a true if the text is matched with the pattern.

## Type

Arguments | Return values
---|---
string, string | boolean

## Usage

```hcl
"${match("abcdef", "^a")}"
# => true
```
