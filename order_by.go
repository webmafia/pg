package pg

import "github.com/webmafia/fast"

func Order(column any, order ...string) QueryEncoder {
	return orderBy{
		column: column,
		desc:   len(order) > 0 && (order[0] == "desc" || order[0] == "DESC"),
	}
}

type orderBy struct {
	column any
	desc   bool
}

// EncodeString implements fast.StringEncoder.
func (o orderBy) EncodeQuery(buf *fast.StringBuffer, queryArgs *[]any) {
	writeAnyIdentifier(buf, o.column)

	if o.desc {
		buf.WriteString(" DESC")
	}
}
