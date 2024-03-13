package pg

import (
	"fmt"

	"github.com/webmafia/fast"
)

// func (db *DB) formatQuery(sql string, args ...any) {
// 	q := `
// 		SELECT *
// 		FROM %T
// 		WHERE %T

// 	`
// }

func indexQuery(format string) (q indexedQuery) {
	end := len(format)

formatLoop:
	for i := 0; i < end; {
		for i < end && format[i] != '%' {
			i++
		}

		if i >= end {
			// done processing format string
			break
		}

		idx := i

		// Process one character (after %)
		i++

		// Double % means escape
		if format[i] == '%' {
			i++
			continue formatLoop
		}

		if i >= end {
			break
		}

		c := format[i]

		if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') {
			i++
			q.indices = append(q.indices, idx)

			continue formatLoop
		}
	}

	q.query = format

	return
}

type indexedQuery struct {
	query   string
	indices []int
}

func (q indexedQuery) encodeQuery(buf *fast.StringBuffer, args []any) (err error) {
	if expect, got := len(q.indices), len(args); expect != got {
		return fmt.Errorf("invalid number of arguments; expected %d, got %d", expect, got)
	}

	// Short circuit if there are no args
	if len(q.indices) == 0 {
		buf.WriteString(q.query)
		return
	}

	var last int
	argNum := 1

	for i, idx := range q.indices {
		buf.WriteString(q.query[last:idx])
		last = idx + 2

		if enc, ok := args[i].(fast.StringEncoder); ok {
			enc.EncodeString(buf)
		} else {
			buf.WriteByte('$')
			buf.WriteInt(argNum)
			argNum++
		}
	}

	buf.WriteString(q.query[last:])

	return
}

// func

// argNumber returns the next argument to evaluate, which is either the value of the passed-in
// argNum or the value of the bracketed integer that begins format[i:]. It also returns
// the new value of i, that is, the index of the next byte of the format to process.
func argNumber(argNum int, format string, i int) (newArgNum, newi int, found bool) {
	if len(format) <= i || format[i] != '[' {
		return argNum, i, false
	}

	index, wid, ok := parseArgNumber(format[i:])

	if ok && 0 <= index {
		return index, i + wid, true
	}

	return argNum, i + wid, ok
}

// parseArgNumber returns the value of the bracketed number, minus 1
// (explicit argument numbers are one-indexed but we want zero-indexed).
// The opening bracket is known to be present at format[0].
// The returned values are the index, the number of bytes to consume
// up to the closing paren, if present, and whether the number parsed
// ok. The bytes to consume will be 1 if no closing paren is present.
func parseArgNumber(format string) (index int, wid int, ok bool) {
	// There must be at least 3 bytes: [n].
	if len(format) < 3 {
		return 0, 1, false
	}

	// Find closing bracket.
	for i := 1; i < len(format); i++ {
		if format[i] == ']' {
			width, ok, newi := parsenum(format, 1, i)
			if !ok || newi != i {
				return 0, i + 1, false
			}
			return width - 1, i + 1, true // arg numbers are one-indexed and skip paren.
		}
	}
	return 0, 1, false
}

// tooLarge reports whether the magnitude of the integer is
// too large to be used as a formatting width or precision.
func tooLarge(x int) bool {
	const max int = 1e6
	return x > max || x < -max
}

// parsenum converts ASCII to integer.  num is 0 (and isnum is false) if no number present.
func parsenum(s string, start, end int) (num int, isnum bool, newi int) {
	if start >= end {
		return 0, false, end
	}
	for newi = start; newi < end && '0' <= s[newi] && s[newi] <= '9'; newi++ {
		if tooLarge(num) {
			return 0, false, end // Overflow; crazy long number most likely.
		}
		num = num*10 + int(s[newi]-'0')
		isnum = true
	}
	return
}
