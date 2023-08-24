package services

import (
	"unicode"
)

func validstring(str string, validators ...func(r rune) bool) bool {
	for _, r := range str {
		var cond bool
		for _, v := range validators {
			if v(r) {
				cond = true
				break
			}
		}
		if !cond {
			return false
		}
	}
	return true
}

func isPrintable(str string) bool {
	return validstring(str, func(r rune) bool {
		return unicode.In(r, unicode.Latin)
	})
}
