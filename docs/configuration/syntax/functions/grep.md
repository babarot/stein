# grep(pattern, text)

Returns the text block matched with the given pattern.

## Type

Arguments | Return values
---|---
string | string

## Usage

```
My life didn't please me,
so I created my life.
- Coco Chanel
```

```hcl
"${grep(file("text.txt"), "life")}"
# => "My life didn't please me,\nso I created my life."
```
