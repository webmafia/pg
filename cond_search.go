package pg

import (
	"strings"

	"github.com/webmafia/fast"
)

type SearchOptions struct {
	Dictionary   string
	Preprocessor func(string) string
}

func (opt *SearchOptions) setDefaults() {
	if opt.Dictionary == "" {
		opt.Dictionary = "simple"
	}
}

func Search(col any, val string, options ...SearchOptions) QueryEncoder {
	var opt SearchOptions

	if len(options) > 0 {
		opt = options[0]
	}

	opt.setDefaults()

	if opt.Preprocessor != nil {
		val = opt.Preprocessor(val)
	}

	return Cond(func(buf *fast.StringBuffer, queryArgs *[]any) {
		writeAnyIdentifier(buf, col)
		buf.WriteString(" @@ to_tsquery('")
		buf.WriteString(opt.Dictionary)
		buf.WriteString("', ")
		writeAny(buf, queryArgs, val)
		buf.WriteByte(')')
	})
}

func PrefixSearch(str string) string {
	var buf strings.Builder
	// Pre-allocate to avoid more than one allocation,
	// Estimation might not be perfect due to escaping, but it avoids gross underestimation.
	buf.Grow(len(str) * 2)

	// State tracking variables
	inWord := false // Tracks if we're currently processing within a word
	for i := 0; i < len(str); i++ {
		c := str[i]

		// Check if character is a space
		if c == ' ' {
			if inWord {
				// End of a word, so add suffix and reset inWord flag
				buf.WriteString(":*")
				inWord = false
			}
			// Skip extra spaces
			continue
		}

		// If starting a new word, add a space if it's not the first word
		if !inWord && buf.Len() > 0 {
			buf.WriteByte(' ')
		}
		inWord = true // We are now in a word

		// Check for characters that need escaping
		switch c {
		case '\\', '*', '&', '|', ':':
			buf.WriteByte('\\') // Prefix with a backslash
			fallthrough
		default:
			buf.WriteByte(c) // Add the current character to the buffer
		}
	}

	// If the last character processed was part of a word, add the suffix
	if inWord {
		buf.WriteString(":*")
	}

	return buf.String()
}
