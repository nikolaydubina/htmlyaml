package htmlyaml

import "strings"

func NullHTML(k string) string { return `<div class="yaml-lang yaml-value yaml-null">null</div>` }

func BoolHTML(k string, v bool) string {
	x := "false"
	if v {
		x = "true"
	}
	return `<div class="yaml-lang yaml-value yaml-bool">` + x + `</div>`
}

func StringHTML(k string, v string) string {
	return `<div class="yaml-value yaml-string">"` + v + `"</div>`
}

func NumberHTML(k string, v float64, s string) string {
	return `<div class="yaml-value yaml-number">` + s + `</div>`
}

var DefaultArrayDashHTML = `<div class="yaml-lang">&nbsp-&nbsp</div>`

var DefaultMapColonHTML = `<div class="yaml-lang">:</div>`

func DefaultMapKeyHTML(key string, v string) string {
	return `<div class="yaml-key yaml-string">"` + v + `"</div>`
}

type DefaultRowHTML struct {
	Padding int
}

func (s DefaultRowHTML) Marshal(v string, depth int) string {
	p := `<div class="yaml-container-padding">` + strings.Repeat("&nbsp", s.Padding*2*depth) + `</div>`
	return `<div class="yaml-container-row">` + p + v + `</div>`
}

// DefaultMarshaler adds basic HTML div classes for further styling.
var DefaultMarshaler = Marshaler{
	Null:      NullHTML,
	Bool:      BoolHTML,
	String:    StringHTML,
	Number:    NumberHTML,
	ArrayDash: DefaultArrayDashHTML,
	MapKey:    DefaultMapKeyHTML,
	MapColon:  DefaultMapColonHTML,
	Row:       DefaultRowHTML{Padding: 4}.Marshal,
}
