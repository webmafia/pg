package pg

import (
	"strings"

	"github.com/webmafia/fast"
)

func encodeQuery(buf *fast.StringBuffer, format string, args []any) {
	var cursor, argIdx int
	argNum := 1

	for {
		if argIdx >= len(args) {
			break
		}

		i := strings.IndexByte(format[cursor:], '%')

		if i < 0 {
			break
		}

		idx := cursor + i
		i = idx + 1
		buf.WriteString(format[cursor:idx])
		cursor = idx + 2

		// Double % means an escaped %
		if format[i] == '%' {
			buf.WriteByte('%')
			continue
		}

		c := format[i]

		if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') {
			arg := args[argIdx]

			if enc, ok := arg.(fast.StringEncoder); ok {
				enc.EncodeString(buf)
			} else {
				buf.WriteByte('$')
				buf.WriteInt(argNum)
				argNum++
			}

			argIdx++
		}
	}

	buf.WriteString(format[cursor:])
}
