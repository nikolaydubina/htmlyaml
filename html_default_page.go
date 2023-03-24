package htmlyaml

import (
	"bytes"
	_ "embed"
	"io"
)

//go:embed html_default_page.html
var defaultPageTemplate []byte

var DefaultPageMarshaler = PageMarshaler{
	Title:            "htmlyaml",
	Template:         defaultPageTemplate,
	TemplateTitleKey: `{{.Title}}`,
	TemplateYAMLKey:  `{{.HTMLYAML}}`,
	Marshaler:        &DefaultMarshaler,
}

// PageMarshaler encodes YAML via marshaller into HTML page by placing Title and content appropriately.
type PageMarshaler struct {
	Title            string
	Template         []byte
	TemplateTitleKey string
	TemplateYAMLKey  string

	Marshaler interface {
		MarshalTo(w io.Writer, v any) error
	}

	idxTitle    int
	idxhtmlyaml int
}

func (m *PageMarshaler) Marshal(v any) []byte {
	b := bytes.Buffer{}
	m.MarshalTo(&b, v)
	return b.Bytes()
}

func (m *PageMarshaler) parseTemplate() {
	if m.idxTitle == 0 || m.idxhtmlyaml == 0 {
		m.idxTitle = bytes.Index(m.Template, []byte(m.TemplateTitleKey))
		m.idxhtmlyaml = bytes.Index(m.Template, []byte(m.TemplateYAMLKey))
	}
}

func (m *PageMarshaler) MarshalTo(w io.Writer, v any) error {
	m.parseTemplate()

	var s int

	if f := m.idxTitle; f > 0 {
		w.Write(m.Template[s:f])
		w.Write([]byte(m.Title))
		s = f + len(m.TemplateTitleKey)
	}

	if f := m.idxhtmlyaml; f > 0 {
		w.Write(m.Template[s:f])
		m.Marshaler.MarshalTo(w, v)
		s = f + len(m.TemplateYAMLKey)
	}

	w.Write(m.Template[s:])

	return nil
}
