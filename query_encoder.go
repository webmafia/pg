package pg

import (
	"strings"

	"github.com/webmafia/fast"
)

type QueryEncoder interface {
	EncodeQuery(buf *fast.StringBuffer, queryArgs *[]any)
}

func encodeQuery(buf *fast.StringBuffer, format string, args []any, queryArgs *[]any) {
	var cursor int

	for {
		if len(args) == 0 {
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
			writeAny(buf, queryArgs, args[0])

			args = args[1:]
		}
	}

	buf.WriteString(format[cursor:])
}
