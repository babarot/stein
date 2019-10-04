# glob(pattern)

Return an array of filenames or directories that matches the specified pattern.

## Type

Arguments | Return type
---|---
string | list(string)

## Usage

```hcl
"${glob("*.txt")}"
```

The output of the code above could be:

```hcl
["a.txt", "b.txt", "c.txt"]
```
