# wc(text, [l, c, w])

Returns the counted number of text as options (**l** ines, **w** ords, **c** hars).
Default option is **l** ines. Same as UNIX's one.

## Type

Arguments | Return values
---|---
string, (string...) | number

## Usage

```hcl
"${wc("foo\nbar baz")}"
# => 1

"${wc("foo\nbar baz", "c")}"
# => 11

"${wc("foo\nbar baz", "w")}"
# => 3
```

