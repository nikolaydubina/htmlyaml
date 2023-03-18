package htmlyaml

import (
	"bytes"
	"errors"
	"io"
	"sort"
	"strconv"
	"strings"
)

// Marshaler converts YAML stored as Go `any` object represented into HTML.
// Rendering is customized by providing renderers for specific YAML elements.
// This facilitates CSS styling, CSS animations, and JavaScript events.
// YAML element renderers receive JSON Path and value of element.
// Should be used only for types: bool, float64, string, []any, map[string]any, nil.
// You can get allowed input easily with yaml.Unmarshal to any.
// Safe for repeated use.
// Not safe for concurrent use.
// TODO: string quotation check for whitespace inside strings
// TODO: root keys start same depth when it is dict
type Marshaler struct {
	Null       func(key string) string
	Bool       func(key string, v bool) string
	String     func(key string, v string) string
	Number     func(key string, v float64, s string) string
	ArrayDash  string
	ArrayEmpty string
	MapColon   string
	MapEmpty   string
	MapKey     func(key string, v string) string
	Row        func(s string, padding int) string

	*rowWriter
	depth int
	key   string
	err   []error
}

// Marshaler converts YAML stored as Go `any` object represented into HTML.
func (s *Marshaler) Marshal(v any) []byte {
	b := bytes.Buffer{}
	s.MarshalTo(&b, v)
	return b.Bytes()
}

// MarshalTo converts YAML stored as Go `any` object represented into HTML.
func (s *Marshaler) MarshalTo(w io.Writer, v any) error {
	s.depth = 0
	s.key = "$"
	s.rowWriter = &rowWriter{
		b:   strings.Builder{},
		w:   w,
		row: s.Row,
	}
	s.marshal(v)
	s.flush(s.depth)
	s.err = append(s.err, s.rowWriter.err...)
	return errors.Join(s.err...)
}

func (s *Marshaler) marshal(v any) {
	if v == nil {
		s.encodeNull()
	}
	switch q := v.(type) {
	case bool:
		s.encodeBool(q)
	case string:
		s.encodeString(q)
	case float64:
		s.encodeFloat64(q)
	case map[string]any:
		s.encodeMap(q)
	case []any:
		s.encodeArray(q)
	default:
		s.err = append(s.err, errors.New("skip unsupported type at key("+s.key+")"))
	}
}

func (s *Marshaler) encodeNull() {
	s.write(s.Null(s.key))
	s.flush(s.depth)
}

func (s *Marshaler) encodeBool(v bool) {
	s.write(s.Bool(s.key, v))
	s.flush(s.depth)
}

func (s *Marshaler) encodeString(v string) {
	s.write(s.String(s.key, v))
	s.flush(s.depth)
}

func (s *Marshaler) encodeFloat64(v float64) {
	s.write(s.Number(s.key, v, strconv.FormatFloat(v, 'f', -1, 64)))
	s.flush(s.depth)
}

func (s *Marshaler) encodeArray(v []any) {
	if len(v) == 0 {
		s.write(s.ArrayEmpty)
		s.flush(s.depth)
		return
	}

	// write array
	k, d := s.key, s.depth
	defer func() { s.key, s.depth = k, d }()
	s.flush(s.depth)

	s.depth = d + 1
	for i, q := range v {
		s.key = k + "[" + strconv.Itoa(i) + "]"

		s.write(s.ArrayDash)
		s.marshal(q)
	}
}

// TODO: anonymous map starts same line first key
func (s *Marshaler) encodeMap(v map[string]any) {
	if len(v) == 0 {
		s.write(s.MapEmpty)
		s.flush(s.depth)
		return
	}

	// extract and sort the keys
	type kv struct {
		k string
		v any
	}
	sv := make([]kv, 0, len(v))
	for k, v := range v {
		sv = append(sv, kv{k: k, v: v})
	}
	sort.Slice(sv, func(i, j int) bool { return sv[i].k < sv[j].k })

	// write map
	k, d := s.key, s.depth
	defer func() { s.key, s.depth = k, d }()

	s.flush(s.depth)

	s.depth = d + 1
	for _, kv := range sv {
		s.key = k + "." + kv.k

		// key
		s.write(s.MapKey(s.key, kv.k))
		s.write(s.MapColon)

		// value
		s.marshal(kv.v)
	}
}
