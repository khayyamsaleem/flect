package flect

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Ident struct {
	original string
	parts    []string
}

func (i Ident) String() string {
	return i.original
}

func New(s string) Ident {
	i := Ident{
		original: s,
		parts:    toParts(s),
	}

	return i
}

var splitRx = regexp.MustCompile("[^\\p{L}]")

func toParts(s string) []string {
	parts := []string{}
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return parts
	}
	if _, ok := baseAcronyms[strings.ToUpper(s)]; ok {
		return []string{strings.ToUpper(s)}
	}
	var prev rune
	var x string
	for _, c := range s {
		cs := string(c)
		// fmt.Println("### cs ->", cs)
		// fmt.Println("### unicode.IsControl(c) ->", unicode.IsControl(c))
		// fmt.Println("### unicode.IsDigit(c) ->", unicode.IsDigit(c))
		// fmt.Println("### unicode.IsGraphic(c) ->", unicode.IsGraphic(c))
		// fmt.Println("### unicode.IsLetter(c) ->", unicode.IsLetter(c))
		// fmt.Println("### unicode.IsLower(c) ->", unicode.IsLower(c))
		// fmt.Println("### unicode.IsMark(c) ->", unicode.IsMark(c))
		// fmt.Println("### unicode.IsPrint(c) ->", unicode.IsPrint(c))
		// fmt.Println("### unicode.IsPunct(c) ->", unicode.IsPunct(c))
		// fmt.Println("### unicode.IsSpace(c) ->", unicode.IsSpace(c))
		// fmt.Println("### unicode.IsTitle(c) ->", unicode.IsTitle(c))
		// fmt.Println("### unicode.IsUpper(c) ->", unicode.IsUpper(c))
		if !utf8.ValidRune(c) {
			continue
		}

		if isSpace(c) {
			parts = xappend(parts, x)
			x = cs
			prev = c
			continue
		}
		if unicode.IsUpper(c) && !unicode.IsUpper(prev) {
			parts = xappend(parts, x)
			x = cs
			prev = c
			continue
		}
		if unicode.IsLetter(c) || unicode.IsDigit(c) || unicode.IsPunct(c) || c == '`' {
			prev = c
			x += cs
			continue
		}
		parts = xappend(parts, x)
		x = ""
		prev = c
	}
	parts = xappend(parts, x)

	return parts
}

var spaces = []rune{'_', ' ', ':', '-', '/'}

func isSpace(c rune) bool {
	for _, r := range spaces {
		if r == c {
			return true
		}
	}
	return unicode.IsSpace(c)
}

func xappend(a []string, ss ...string) []string {
	for _, s := range ss {
		s = strings.TrimSpace(s)
		for _, x := range spaces {
			s = strings.Trim(s, string(x))
		}
		if _, ok := baseAcronyms[strings.ToUpper(s)]; ok {
			s = strings.ToUpper(s)
		}
		if s != "" {
			a = append(a, s)
		}
	}
	return a
}