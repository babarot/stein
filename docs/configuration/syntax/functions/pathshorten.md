# pathshorten(path)

Returns the file path shortened like [:fa-external-link: Vim's one](http://vimdoc.sourceforge.net/htmldoc/eval.html#pathshorten()).

## Type

Arguments | Return values
---|---
string | string

## Usage

```hcl
"${pathshorten("manifests/microservices/x-gateway-jp/development/Service/a.yaml")}"
# => "m/m/x/d/S/a.yaml"
```
