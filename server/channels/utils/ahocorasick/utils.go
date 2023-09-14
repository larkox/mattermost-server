package ahocorasick

import "unicode"

func mapToSlice[K comparable](mymap map[K]struct{}) []K {
	keys := make([]K, len(mymap))

	i := 0
	for k := range mymap {
		keys[i] = k
		i++
	}

	return keys
}

func isAlphanumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r)
}

func isWordSeparator(r rune) bool {
	return unicode.IsMark(r) || unicode.IsSpace(r) || unicode.IsPunct(r)
}
