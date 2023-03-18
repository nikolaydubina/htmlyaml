package htmlyaml

import (
	"io"
	"strings"
)

// rowWriter accumulates items written to row and flushes it as a wrapped row on flush calls.
// flush has to be called eventually.
type rowWriter struct {
	b   strings.Builder
	w   io.Writer
	row func(s string, padding int) string
	err []error
}

func (s *rowWriter) write(v string) {
	_, err := s.b.WriteString(v)
	s.err = append(s.err, err)
}

func (s *rowWriter) flush(depth int) {
	v := s.row(s.b.String()+"\n", depth)
	_, err := io.WriteString(s.w, v)
	s.err = append(s.err, err)
	s.b.Reset()
}
