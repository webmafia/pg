package pg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/webmafia/fast"
)

type QueryEncoder interface {
	EncodeQuery(buf *fast.StringBuffer, queryArgs *[]any)
}

var _ QueryEncoder = (EncodeQuery)(nil)

type EncodeQuery func(buf *fast.StringBuffer, queryArgs *[]any)

func (fn EncodeQuery) EncodeQuery(buf *fast.StringBuffer, queryArgs *[]any) {
	fn(buf, queryArgs)
}

var queryEncoderNoop = EncodeQuery(func(_ *fast.StringBuffer, _ *[]any) {})

func encodeQuery(buf *fast.StringBuffer, format string, args []any, queryArgs *[]any) (err error) {
	var cursor int
	var argNum int

	for {
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

		if format[i] == '[' {
			end := strings.IndexByte(format[i:], ']')

			if end < 0 {
				return errors.New("missing ']'")
			}

			num, err := strconv.Atoi(format[i+1 : i+end])

			if err != nil {
				return err
			}

			if num < 1 {
				return errors.New("argument number must be at least 1")
			}

			argNum = num
			cursor += end + 1
			i += end + 1
		} else {
			argNum++
		}

		if argNum > len(args) {
			return fmt.Errorf("argument number %d does not exist", argNum)
		}

		c := format[i]

		if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') {
			writeAny(buf, queryArgs, args[argNum-1])
		} else {
			return fmt.Errorf("unsupported placeholder '%%%s'", string(c))
		}
	}

	buf.WriteString(format[cursor:])
	return
}
