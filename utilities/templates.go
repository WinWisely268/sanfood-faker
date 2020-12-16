package utilities

import (
	"regexp"
	"strings"
)

var (
	configNumberSequence    = regexp.MustCompile(`([a-zA-Z])(\d+)([a-zA-Z]?)`)
	configNumberReplacement = []byte("$1 $2 $3")
	upperCaseAcronym        = map[string]bool{
		"ID": true,
	}
)

func ToCaseSnake(s string) string {
	return toUpperCaseDelimited(s, '_', 0, false)
}

func addWordBoundariesToNumbers(s string) string {
	b := []byte(s)
	b = configNumberSequence.ReplaceAll(b, configNumberReplacement)
	return string(b)
}

func toUpperCaseDelimited(s string, delimiter uint8, ignore uint8, upper bool) string {
	s = addWordBoundariesToNumbers(s)
	s = strings.Trim(s, " ")
	n := ""
	for i, v := range s {
		// treat acronyms as words, eg for JSONData -> JSON is a whole word
		nextCaseIsChanged := false
		if i+1 < len(s) {
			next := s[i+1]
			vIsCap := v >= 'A' && v <= 'Z'
			vIsLow := v >= 'a' && v <= 'z'
			nextIsCap := next >= 'A' && next <= 'Z'
			nextIsLow := next >= 'a' && next <= 'z'
			if (vIsCap && nextIsLow) || (vIsLow && nextIsCap) {
				nextCaseIsChanged = true
			}
			if ignore > 0 && i-1 >= 0 && s[i-1] == ignore && nextCaseIsChanged {
				nextCaseIsChanged = false
			}
		}

		if i > 0 && n[len(n)-1] != delimiter && nextCaseIsChanged {
			// add underscore if next letter case type is changed
			if v >= 'A' && v <= 'Z' {
				n += string(delimiter) + string(v)
			} else if v >= 'a' && v <= 'z' {
				n += string(v) + string(delimiter)
			}
		} else if v == ' ' || v == '_' || v == '-' {
			// replace spaces/underscores with delimiters
			if uint8(v) == ignore {
				n += string(v)
			} else {
				n += string(delimiter)
			}
		} else {
			n = n + string(v)
		}
	}

	if upper {
		n = strings.ToUpper(n)
	} else {
		n = strings.ToLower(n)
	}
	return n
}

// Converts a string to CamelCase
func toCamelInitCase(s string, initCase bool) string {
	s = addWordBoundariesToNumbers(s)
	s = strings.Trim(s, " ")
	n := ""
	capNext := initCase
	for _, v := range s {
		if v >= 'A' && v <= 'Z' {
			n += string(v)
		}
		if v >= '0' && v <= '9' {
			n += string(v)
		}
		if v >= 'a' && v <= 'z' {
			if capNext {
				n += strings.ToUpper(string(v))
			} else {
				n += string(v)
			}
		}
		if v == '_' || v == ' ' || v == '-' || v == '.' {
			capNext = true
		} else {
			capNext = false
		}
	}
	return n
}

// toCamelCase converts a string to lowerCamelCase
func ToCamelCase(s string) string {
	if s == "" {
		return s
	}
	if upperCaseAcronym[s] {
		s = strings.ToLower(s)
	}
	if r := rune(s[0]); r >= 'A' && r <= 'Z' {
		s = strings.ToLower(string(r)) + s[1:]
	}
	return toCamelInitCase(s, false)
}
