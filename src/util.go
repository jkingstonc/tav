package src

import (
	"strings"
	"unicode"
)

func LineInString(source string, line int) string {
	lineCounter := 0
	found := false
	sb := strings.Builder{}
	for i := 0; i < len(source); i++ {
		if !found && lineCounter == line {
			found = true
		}
		if source[i] == '\n' || source[i] == '\r' {
			// if we have reached the end of the current line
			if found == true {
				return sb.String()
			}
			lineCounter++
		}
		if found {
			sb.WriteByte(source[i])
		}
	}
	return ""
}

func IsChar(r rune) bool {
	return r != ' ' && unicode.IsLetter(r)
}

func IsNum(r rune) bool {
	return unicode.IsDigit(r)
}