package pg

import "github.com/webmafia/fast"

//go:inline
func op(col any, op string, val any) QueryEncoder {
	return Cond(func(buf *fast.StringBuffer, queryArgs *[]any) {
		writeAnyIdentifier(buf, col)
		buf.WriteString(op)
		writeAny(buf, queryArgs, val)
	})
}

func Eq(col any, val any) QueryEncoder {
	return op(col, " = ", val)
}

func NotEq(col any, val any) QueryEncoder {
	return op(col, " != ", val)
}

func Gt(col any, val any) QueryEncoder {
	return op(col, " > ", val)
}

func Gte(col any, val any) QueryEncoder {
	return op(col, " >= ", val)
}

func Lt(col any, val any) QueryEncoder {
	return op(col, " < ", val)
}

func Lte(col any, val any) QueryEncoder {
	return op(col, " <= ", val)
}

func In(col any, val any) QueryEncoder {
	return Cond(func(buf *fast.StringBuffer, queryArgs *[]any) {
		writeAnyIdentifier(buf, col)

		switch v := val.(type) {

		case StringEncoder:
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
	})
}

func NotIn(col any, val any) QueryEncoder {
	return Cond(func(buf *fast.StringBuffer, queryArgs *[]any) {
		writeAnyIdentifier(buf, col)

		switch v := val.(type) {

		case StringEncoder:
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
