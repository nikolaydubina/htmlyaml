## 🐹 htmlyaml: rich rendering of YAML as HTML in Go

[![codecov](https://codecov.io/gh/nikolaydubina/htmlyaml/branch/master/graph/badge.svg?token=yXmNdIDn8O)](https://codecov.io/gh/nikolaydubina/htmlyaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/nikolaydubina/htmlyaml)](https://goreportcard.com/report/github.com/nikolaydubina/htmlyaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/nikolaydubina/htmlyaml.svg)](https://pkg.go.dev/github.com/nikolaydubina/htmlyaml)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/nikolaydubina/htmlyaml/badge)](https://securityscorecards.dev/viewer/?uri=github.com/nikolaydubina/htmlyaml)

* pure Go
* no Javascript
* no dependencies
* no `reflect`
* no `fmt`
* no `text/template`, `html/template`
* 300 LOC
* customizable rendering
* JSON Path for elements access

![](./doc/example-color.png)

```go
// YAML has to be any
var v any
yaml.Unmarshal(exampleYAML, &v)

htmlPage := htmlyaml.DefaultPageMarshaler.Marshal(v)
```

```go
// YAML has to be any
var v any
yaml.Unmarshal(exampleYAML, &v)

// customize how to render HTML elements
s := htmlyaml.DefaultMarshaler
s.Number = func(k string, v float64, s string) string {
    if k == "$.cakes.strawberry-cake.size" {
        return `<div class="yaml-value yaml-number" style="color:red;">` + s + `</div>`
    }
    if v > 10 {
        return `<div class="yaml-value yaml-number" style="color:blue;">` + s + `</div>`
    }
    return `<div class="yaml-value yaml-number">` + s + `</div>`
}

m := htmlyaml.DefaultPageMarshaler
m.Marshaler = &s

// write HTML page
htmlPage := m.Marshal(v)
```

## Related Work

- https://github.com/alecthomas/chroma can render many languages into HTML with syntax highligting
