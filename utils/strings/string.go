package strings

import "strings"

func IsBlank(buffer string) bool {
	bufferTimed := strings.TrimSpace(buffer)
	if len(bufferTimed) == 0 {
		return true
	}
	return false
}

func IsNotBlank(buffer string) bool {
	return !IsBlank(buffer)
}
