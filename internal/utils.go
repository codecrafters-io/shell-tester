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

func removeControlSequence(data []byte) []byte {
	PROMPT_START := '$'

	for startIdx, r := range string(data) {
		// Standard escape codes are prefixed with Escape (27)
		if r == 27 {
			// remove from here upto PROMPT_START
			for endIdx, r2 := range string(data[startIdx:]) {
				if r2 == PROMPT_START {
					// Remove from start_idx to end_idx-1
					data = append(data[:startIdx], data[endIdx:]...)
					break
				}
			}
		}
	}

	return data
}

func printAllChars(data []byte) {
	for _, r := range string(data) {
		fmt.Printf("%d ", r)
	}
	fmt.Println()
}
