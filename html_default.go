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
	return `<div class="yaml-value yaml-string">` + v + `</div>`
}

func NumberHTML(k string, v float64, s string) string {
	return `<div class="yaml-value yaml-number">` + s + `</div>`
}

var (
	DefaultArrayDashHTML  = `<div class="yaml-lang">&nbsp;-&nbsp;</div>`
	DefaultArrayEmptyHTML = `<div class="yaml-lang">&nbsp;[]</div>`
	DefaultMapColonHTML   = `<div class="yaml-lang">:&nbsp;</div>`
	DefaultMapEmptyHTML   = `<div class="yaml-lang">&nbsp;{}</div>`
)

func DefaultMapKeyHTML(key string, v string) string {
	return `<div class="yaml-key yaml-string">` + v + `</div>`
}

type DefaultRowHTML struct {
	Padding int
}

func (s DefaultRowHTML) Marshal(v string, depth int) string {
	p := `<span class="yaml-container-padding">` + strings.Repeat("&nbsp;", s.Padding*depth) + `</span>`
	return `<div class="yaml-container-row">` + p + v + `</div>`
}

// DefaultMarshaler adds basic HTML div classes for further styling.
var DefaultMarshaler = Marshaler{
	Null:       NullHTML,
	Bool:       BoolHTML,
	String:     StringHTML,
	Number:     NumberHTML,
	MapKey:     DefaultMapKeyHTML,
	ArrayDash:  DefaultArrayDashHTML,
	ArrayEmpty: DefaultArrayEmptyHTML,
	MapColon:   DefaultMapColonHTML,
	MapEmpty:   DefaultMapEmptyHTML,
	Row:        DefaultRowHTML{Padding: 2}.Marshal,
}
