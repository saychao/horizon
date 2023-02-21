//Package utf8 contains utilities for working with utf8 data.
package utf8

import (
	"bytes"
	"strings"
	"unicode/utf8"
)

// Scrub ensures that a given string is valid utf-8, replacing any invalid byte
// sequences with the utf-8 replacement character.
func Scrub(in string) string {
	// drop nil bytes
	in = strings.Replace(in, "\u0000", "", -1)

	// drop line breaks
	// (originating from withdraw external_details hotfix, where mysterious line break appeared in tx_hex)
	in = strings.Replace(in, "\n", "", -1)

	// First check validity using the stdlib, returning if the string is already
	// valid
	if utf8.ValidString(in) {
		return in
	}

	left := []byte(in)
	var result bytes.Buffer

	for len(left) > 0 {
		r, n := utf8.DecodeRune(left)

		_, err := result.WriteRune(r)
		if err != nil {
			panic(err)
		}

		left = left[n:]
	}

	return result.String()
}
