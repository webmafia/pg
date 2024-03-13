package pg

import "github.com/webmafia/fast"

type QueryEncoder interface {
	EncodeQuery(buf *fast.StringBuffer, args []any, queryArgs *[]any)
}
