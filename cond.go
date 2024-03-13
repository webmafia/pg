package pg

import "github.com/webmafia/fast"

type Cond func(buf *fast.StringBuffer, queryArgs *[]any)

func (c Cond) EncodeQuery(buf *fast.StringBuffer, queryArgs *[]any) {
	c(buf, queryArgs)
}
