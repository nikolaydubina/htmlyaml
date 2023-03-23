package htmlyaml

import "strings"

func DefaultNull(k string) string { return `<div class="yaml-lang yaml-value yaml-null">null</div>` }

func DefaultBool(k string, v bool) string {
	x := "false"
	if v {
		x = "true"
	}
	return `<div class="yaml-lang yaml-value yaml-bool">` + x + `</div>`
}

func DefaultString(k string, v string) string {
	return `<div class="yaml-value yaml-string">` + v + `</div>`
}

func DefaultNumber(k string, v float64, s string) string {
	return `<div class="yaml-value yaml-number">` + s + `</div>`
}

var (
	DefaultArrayDash    = `<div class="yaml-lang">-&nbsp;</div>`
	DefaultArrayEmpty   = `<div class="yaml-lang">&nbsp;[]</div>`
	DefaultMapColon     = `<div class="yaml-lang">:&nbsp;</div>`
	DefaultMapEmpty     = `<div class="yaml-lang">&nbsp;{}</div>`
	DefaultPaddingSpace = `<span class="yaml-padding-space">&nbsp;</span>`
)

func DefaultMapKey(key string, v string) string {
	return `<div class="yaml-key yaml-string">` + v + `</div>`
}

type DefaultRow struct {
	Padding int
}

func (s DefaultRow) Marshal(v string, depth int) string {
	p := `<div class="yaml-container-padding">` + strings.Repeat(DefaultPaddingSpace, s.Padding*depth) + `</div>`
	return `<div class="yaml-container-row">` + p + v + `</div>`
}

// DefaultMarshaler adds basic HTML div classes for further styling.
var DefaultMarshaler = Marshaler{
	Null:       DefaultNull,
	Bool:       DefaultBool,
	String:     DefaultString,
	Number:     DefaultNumber,
	MapKey:     DefaultMapKey,
	ArrayDash:  DefaultArrayDash,
	ArrayEmpty: DefaultArrayEmpty,
	MapColon:   DefaultMapColon,
	MapEmpty:   DefaultMapEmpty,
	Row:        DefaultRow{Padding: 2}.Marshal,
}
