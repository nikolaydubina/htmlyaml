package htmlyaml_test

import (
	_ "embed"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/nikolaydubina/htmlyaml"
)

//go:embed testdata/example-array-root.json
var exampleArrayRootJSON []byte

//go:embed testdata/example-array-root.html
var exampleArrayRootPageHTML string

//go:embed testdata/example-page.html
var examplePageHTML string

//go:embed testdata/example-page-color.html
var examplePageColorHTML string

func TestMarshalHTML(t *testing.T) {
	var v any
	json.Unmarshal(exampleJSON, &v)

	h := htmlyaml.DefaultPageMarshaler.Marshal(v)

	os.WriteFile("testdata/example-page.out.html", h, 0666)
	if strings.TrimSpace(examplePageHTML) != strings.TrimSpace(string(h)) {
		t.Errorf("wrong output: %s", string(h))
	}
}

func TestMarshalHTML_ArrayRoot(t *testing.T) {
	var v any
	json.Unmarshal(exampleArrayRootJSON, &v)

	h := htmlyaml.DefaultPageMarshaler.Marshal(v)

	os.WriteFile("testdata/example-array-root.out.html", h, 0666)
	if strings.TrimSpace(exampleArrayRootPageHTML) != strings.TrimSpace(string(h)) {
		t.Errorf("wrong output: %s", string(h))
	}
}

func TestMarshalHTML_Color(t *testing.T) {
	var v any
	json.Unmarshal(exampleJSON, &v)

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

	htmlPage := m.Marshal(v)

	os.WriteFile("testdata/example-page-color.out.html", []byte(htmlPage), 0666)
	if strings.TrimSpace(examplePageColorHTML) != strings.TrimSpace(string(htmlPage)) {
		t.Errorf("wrong output: %s", string(htmlPage))
	}
}
