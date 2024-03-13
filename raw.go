package pg

import "github.com/webmafia/fast"

func Raw(s string, args ...any) QueryEncoder {
	return Cond(func(buf *fast.StringBuffer, queryArgs *[]any) {
		if len(args) > 0 {
			encodeQuery(buf, s, args, queryArgs)
		} else {
			buf.WriteString(s)
		}

	})
}
