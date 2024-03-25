package pg

import (
	"github.com/webmafia/fast"
)

type MultiAnd interface {
	QueryEncoder
	And(QueryEncoder) MultiAnd
}

type MultiOr interface {
	QueryEncoder
	Or(QueryEncoder) MultiOr
}

func And(ops ...QueryEncoder) MultiAnd {
	return &multi{
		ops: ops,
		del: " AND ",
	}
}

func Or(ops ...QueryEncoder) MultiOr {
	return &multi{
		ops: ops,
		del: " OR ",
	}
}

var (
	_ MultiAnd = (*multi)(nil)
	_ MultiOr  = (*multi)(nil)
)

type multi struct {
	ops []QueryEncoder
	del string
}

// And implements MultiAnd.
func (m *multi) And(v QueryEncoder) MultiAnd {
	m.ops = append(m.ops, v)
	return m
}

// Or implements MultiOr.
func (m *multi) Or(v QueryEncoder) MultiOr {
	m.ops = append(m.ops, v)
	return m
}

// EncodeQuery implements MultiAnd.
func (m *multi) EncodeQuery(buf *fast.StringBuffer, queryArgs *[]any) {
	if len(m.ops) == 0 {
		buf.WriteString("1 = 1")
		return
	}

	buf.WriteByte('(')

	for i := range m.ops {
		if i != 0 {
			buf.WriteString(m.del)
		}

		m.ops[i].EncodeQuery(buf, queryArgs)
	}

	buf.WriteByte(')')
}
