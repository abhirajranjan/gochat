package jwtHandler

import (
	"strings"
	"unicode"
)

var SpecialCharacters string = "@!#^$&*()_-=+?"

func IsAlphanum(s string) bool {
	for _, i := range s {
		if !(unicode.IsDigit(i) || unicode.IsLetter(i)) {
			return false
		}
	}
	return true
}

func IsAlphanumWithSpecialChar(s string) bool {
	for _, i := range s {
		if !(unicode.IsDigit(i) || unicode.IsLetter(i) || strings.ContainsRune(SpecialCharacters, i)) {
			return false
		}
	}
	return true
}
