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
// You can get allowed input easily with yaml.Unmarshal or json.Unmarshal to any.
// Since HTML automatically removes whitespace, to make indentation YAML conformant,
// spaces are wrapped in their own div element.
// Safe for repeated use.
// Not safe for concurrent use.
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
	depth        int
	isParentList bool
	isRoot       bool
	key          string
	err          []error
}

// Marshaler converts YAML stored as Go `any` object represented into HTML.
func (s *Marshaler) Marshal(v any) []byte {
	b := bytes.Buffer{}
	s.MarshalTo(&b, v)
	return b.Bytes()
}

// MarshalTo converts YAML stored as Go `any` object represented into HTML.
func (s *Marshaler) MarshalTo(w io.Writer, v any) error {
	s.isRoot = true
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

// go-yaml decodes YAML numbers into different numeric Go types,
// as opposed to encoding/json decoding always to float64.
func (s *Marshaler) marshal(v any) {
	if v == nil {
		s.write(s.Null(s.key))
	}
	switch q := v.(type) {
	case bool:
		s.write(s.Bool(s.key, q))
	case string:
		s.write(s.String(s.key, tryEscapeString(q)))
	case int:
		s.write(s.Number(s.key, float64(q), strconv.Itoa(q)))
	case int64:
		s.write(s.Number(s.key, float64(q), strconv.Itoa(int(q))))
	case uint:
		s.write(s.Number(s.key, float64(q), strconv.Itoa(int(q))))
	case uint64:
		s.write(s.Number(s.key, float64(q), strconv.Itoa(int(q))))
	case float32:
		s.write(s.Number(s.key, float64(q), strconv.FormatFloat(float64(q), 'f', -1, 32)))
	case float64:
		s.write(s.Number(s.key, q, strconv.FormatFloat(q, 'f', -1, 64)))
	case map[string]any:
		s.encodeMap(q)
		return
	case []any:
		s.encodeArray(q)
		return
	default:
		s.err = append(s.err, errors.New("skip unsupported type at key("+s.key+")"))
	}
	s.flush(s.depth)
}

func (s *Marshaler) encodeArray(v []any) {
	if len(v) == 0 {
		s.write(s.ArrayEmpty)
		s.flush(s.depth)
		return
	}

	// write array
	k, d, pl := s.key, s.depth, s.isParentList
	defer func() { s.key, s.depth, s.isParentList = k, d, pl }()

	s.flush(s.depth)

	if s.isRoot {
		s.isRoot = false
		s.depth = 0
	} else {
		s.depth = d + 1
	}

	s.isParentList = true
	for i, q := range v {
		s.key = k + "[" + strconv.Itoa(i) + "]"
		s.write(s.ArrayDash)
		s.marshal(q)
	}
}

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
	k, d, r := s.key, s.depth, s.isRoot
	defer func() { s.key, s.depth = k, d }()

	if !s.isParentList {
		s.flush(s.depth)
	}

	if s.isRoot {
		s.isRoot = false
	}

	for i, kv := range sv {
		if !r && ((s.isParentList && i > 0) || !s.isParentList) {
			s.depth = d + 1
		}

		// key
		s.key = k + "." + kv.k
		s.write(s.MapKey(s.key, kv.k))
		s.write(s.MapColon)

		// value
		s.marshal(kv.v)
	}
}
