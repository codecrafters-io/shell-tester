package internal

import (
	"bytes"
	"fmt"
	"unicode"
	"unicode/utf8"
)

// removeNonPrintable removes all non-printable characters from a byte slice.
func removeNonPrintable(data []byte) []byte {
	var buffer bytes.Buffer

	for len(data) > 0 {
		r, size := utf8.DecodeRune(data)
		if r == utf8.RuneError && size == 1 {
			// Invalid UTF-8 encoding, skip this byte
			data = data[size:]
			continue
		}

		if unicode.IsPrint(r) {
			buffer.WriteRune(r)
		}

		data = data[size:]
	}

	return buffer.Bytes()
}

func printAllChars(data []byte) {
	for _, r := range string(data) {
		fmt.Printf("%q ", r)
	}
	fmt.Println()
}
