package htmlyaml_test

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/nikolaydubina/htmlyaml"
)

//go:embed testdata/example.json
var exampleJSON []byte

//go:embed testdata/example.html
var exampleHTML string

func TestDefaultMarshaler(t *testing.T) {
	var v any
	json.Unmarshal(exampleJSON, &v)

	h := htmlyaml.DefaultMarshaler.Marshal(v)

	os.WriteFile("testdata/example.out.html", h, 0666)
	if strings.TrimSpace(exampleHTML) != strings.TrimSpace(string(h)) {
		t.Errorf("wrong output: %s", string(h))
	}
}

func TestDefaultMarshaler_Repeated(t *testing.T) {
	var v any
	json.Unmarshal(exampleJSON, &v)

	s := htmlyaml.DefaultMarshaler

	for i := 0; i < 10; i++ {
		h := s.Marshal(v)
		if strings.TrimSpace(exampleHTML) != strings.TrimSpace(string(h)) {
			t.Errorf("%d: wrong output: %s", i, string(h))
		}
	}
}

func TestMarshaler_JSONPath(t *testing.T) {
	var v any
	json.Unmarshal(exampleJSON, &v)

	j := htmlyaml.NewJSONPathCollector()
	s := htmlyaml.DefaultMarshaler
	s.Null = j.Null
	s.Bool = j.Bool
	s.String = j.String
	s.Number = j.Number
	s.MapKey = j.MapKey

	s.Marshal(v)

	type kv struct {
		k string
		v any
	}
	kvs := make([]kv, 0, len(j.Keys))
	for k, v := range j.Keys {
		kvs = append(kvs, kv{k: k, v: v})
	}
	sort.Slice(kvs, func(i, j int) bool { return kvs[i].k < kvs[j].k })
	for _, k := range kvs {
		t.Log(k)
	}

	exp := []kv{
		{k: "$.bookings", v: "bookings"},
		{k: "$.bookings.monday", v: true},
		{k: "$.bookings.tuesday", v: false},
		{k: "$.box-colors", v: "box-colors"},
		{k: "$.box-colors[0]", v: "\"red with blue stripes\""},
		{k: "$.box-colors[1]", v: "green"},
		{k: "$.box-sizes", v: "box-sizes"},
		{k: "$.box-sizes[0]", v: "10"},
		{k: "$.box-sizes[1]", v: "11"},
		{k: "$.box-sizes[2]", v: "12"},
		{k: "$.box-with-boxes", v: "box-with-boxes"},
		{k: "$.box-with-boxes[1][2][0]", v: "2"},
		{k: "$.cakes", v: "cakes"},
		{k: "$.cakes.chocolate-cake", v: "chocolate-cake"},
		{k: "$.cakes.strawberry-cake", v: "strawberry-cake"},
		{k: "$.cakes.strawberry-cake.color", v: "white"},
		{k: "$.cakes.strawberry-cake.ingredients", v: "ingredients"},
		{k: "$.cakes.strawberry-cake.ingredients[0]", v: "cream"},
		{k: "$.cakes.strawberry-cake.ingredients[1]", v: "strawberry"},
		{k: "$.cakes.strawberry-cake.size", v: "10"},
		{k: "$.drinks", v: "drinks"},
		{k: "$.drinks[0].name", v: "soda"},
		{k: "$.drinks[0].price", v: "10.23"},
		{k: "$.drinks[1].name", v: "tea"},
		{k: "$.drinks[1].price", v: "1.12"},
		{k: "$.fruits", v: "fruits"},
		{k: "$.fruits[0]", v: "null"},
		{k: "$.fruits[1]", v: "null"},
		{k: "$.ice-cream", v: "null"},
		{k: "$.tables", v: "tables"},
	}
	if len(exp) != len(kvs) {
		t.Fatalf("len exp(%d) != %d", len(exp), len(kvs))
	}
	for i := range kvs {
		var errs []error
		if exp[i] != kvs[i] {
			errs = append(errs, fmt.Errorf("exp(%#v) != got(%#v)", exp[i], kvs[i]))
		}
		if err := errors.Join(errs...); err != nil {
			t.Error(errors.Join(errs...))
		}
	}
}
