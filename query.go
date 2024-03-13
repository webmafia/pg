package pg

import (
	"strings"
	"unsafe"

	"github.com/webmafia/fast"
)

var _ QueryEncoder = Query("")

type Query string

// EncodeQuery implements QueryEncoder.
func (q Query) EncodeQuery(buf *fast.StringBuffer, args []any, queryArgs *[]any) {
	var cursor int
	argNum := 1
	format := *(*string)(unsafe.Pointer(&q))

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
			if enc, ok := args[0].(fast.StringEncoder); ok {
				enc.EncodeString(buf)
			} else {
				buf.WriteByte('$')
				buf.WriteInt(argNum)
				*queryArgs = append(*queryArgs, args[0])
				argNum++
			}

			args = args[1:]
		}
	}

	buf.WriteString(format[cursor:])
}

func encodeQuery(buf *fast.StringBuffer, format string, args []any, queryArgs *[]any) {
	var cursor int
	argNum := 1

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
			if enc, ok := args[0].(fast.StringEncoder); ok {
				enc.EncodeString(buf)
			} else {
				buf.WriteByte('$')
				buf.WriteInt(argNum)
				*queryArgs = append(*queryArgs, args[0])
				argNum++
			}

			args = args[1:]
		}
	}

	buf.WriteString(format[cursor:])
}
