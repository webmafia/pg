package pg

import "github.com/webmafia/fast"

func Order(column string, order ...string) QueryEncoder {
	return orderBy{
		column: column,
		desc:   len(order) > 0 && (order[0] == "desc" || order[0] == "DESC"),
	}
}

type orderBy struct {
	column string
	desc   bool
}

// EncodeString implements fast.StringEncoder.
func (o orderBy) EncodeQuery(buf *fast.StringBuffer, queryArgs *[]any) {
	writeQueryArg(buf, queryArgs, o.column)

	if o.desc {
		buf.WriteString(" DESC")
	}
}
