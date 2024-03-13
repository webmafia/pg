package pg

import "github.com/webmafia/fast"

//go:inline
func op(col, op string, val any) QueryEncoder {
	return Cond(func(buf *fast.StringBuffer, queryArgs *[]any) {
		writeIdentifier(buf, col)
		buf.WriteString(op)
		writeAny(buf, queryArgs, val)
	})
}

func Eq(col string, val any) QueryEncoder {
	return op(col, " = ", val)
}

func NotEq(col string, val any) QueryEncoder {
	return op(col, " != ", val)
}

func Gt(col string, val any) QueryEncoder {
	return op(col, " > ", val)
}

func Gte(col string, val any) QueryEncoder {
	return op(col, " >= ", val)
}

func Lt(col string, val any) QueryEncoder {
	return op(col, " < ", val)
}

func Lte(col string, val any) QueryEncoder {
	return op(col, " <= ", val)
}

func In(col string, val any) QueryEncoder {
	return Cond(func(buf *fast.StringBuffer, queryArgs *[]any) {
		writeIdentifier(buf, col)

		switch v := val.(type) {

		case fast.StringEncoder:
			buf.WriteString(" IN (")
			v.EncodeString(buf)
			buf.WriteByte(')')

		case QueryEncoder:
			buf.WriteString(" IN (")
			v.EncodeQuery(buf, queryArgs)
			buf.WriteByte(')')

		default:
			buf.WriteString(" = ANY (")
			writeQueryArg(buf, queryArgs, val)
			buf.WriteByte(')')
		}

		buf.WriteString(" ")
		writeQueryArg(buf, queryArgs, val)
	})
}

func NotIn(col string, val any) QueryEncoder {
	return Cond(func(buf *fast.StringBuffer, queryArgs *[]any) {
		writeIdentifier(buf, col)

		switch v := val.(type) {

		case fast.StringEncoder:
			buf.WriteString(" NOT IN (")
			v.EncodeString(buf)
			buf.WriteByte(')')

		case QueryEncoder:
			buf.WriteString(" NOT IN (")
			v.EncodeQuery(buf, queryArgs)
			buf.WriteByte(')')

		default:
			buf.WriteString(" != ANY (")
			writeQueryArg(buf, queryArgs, val)
			buf.WriteByte(')')
		}

		buf.WriteString(" ")
		writeQueryArg(buf, queryArgs, val)
	})
}
