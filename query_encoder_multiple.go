package pg

import "github.com/webmafia/fast"

func Multiple(enc ...QueryEncoder) EncodeQuery {
	return func(buf *fast.StringBuffer, queryArgs *[]any) {
		for i := range enc {
			if i != 0 {
				buf.WriteString(", ")
			}

			enc[i].EncodeQuery(buf, queryArgs)
		}
	}
}
