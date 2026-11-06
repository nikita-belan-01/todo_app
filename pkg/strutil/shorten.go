package strutil

import "unicode/utf8"

func Shorten(value string, maxLen int) string {
	if len(value) <= maxLen {
		return value
	}

	cut := maxLen
	for cut > 0 && !utf8.RuneStart(value[cut]) {
		cut--
	}

	return value[:cut] + "..."
}
