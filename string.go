package htmlyaml

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func hasWhiteSpace(s string) bool {
	for _, c := range s {
		if unicode.IsSpace(c) {
			return true
		}
	}
	return false
}

func tryEscapeString(s string) string {
	s = strings.TrimSpace(s)
	if hasWhiteSpace(s) {
		fmt.Printf("%#v\n", s)
		return strconv.Quote(s)
	}
	return s
}
