package pg

import "github.com/webmafia/fast"

//go:inline
func op(col any, op string, val any) QueryEncoder {
	if val == nil {
		if op == " = " {
			return Cond(func(buf *fast.StringBuffer, queryArgs *[]any) {
				writeAnyIdentifier(buf, col)
				buf.WriteString(" IS NULL")
			})
		} else if op == " != " {
			return Cond(func(buf *fast.StringBuffer, queryArgs *[]any) {
				writeAnyIdentifier(buf, col)
				buf.WriteString(" IS NOT NULL")
			})
		}
	}

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

func Any(col any, val any) QueryEncoder {
	return Cond(func(buf *fast.StringBuffer, queryArgs *[]any) {
		writeAnyIdentifier(buf, col)
		buf.WriteString(" = ANY (")

		switch v := val.(type) {

		case StringEncoder:
			v.EncodeString(buf)

		case QueryEncoder:
			v.EncodeQuery(buf, queryArgs)

		default:
			writeQueryArg(buf, queryArgs, val)
		}

		buf.WriteByte(')')
	})
}

func All(col any, val any) QueryEncoder {
	return Cond(func(buf *fast.StringBuffer, queryArgs *[]any) {
		writeAnyIdentifier(buf, col)
		buf.WriteString(" = ALL (")

		switch v := val.(type) {

		case StringEncoder:
			v.EncodeString(buf)

		case QueryEncoder:
			v.EncodeQuery(buf, queryArgs)

		default:
			writeQueryArg(buf, queryArgs, val)
		}

		buf.WriteByte(')')
	})
}

// Prefix LIKE match. Optionally search on suffix as well.
func Like(col any, val string, alsoSuffix ...bool) QueryEncoder {
	if len(alsoSuffix) > 0 && alsoSuffix[0] {
		val = "%" + val + "%"
	} else {
		val += "%"
	}

	return op(col, " LIKE ", val)
}

// Negative prefix LIKE match. Optionally search on suffix as well.
func NotLike(col any, val string, alsoSuffix ...bool) QueryEncoder {
	if len(alsoSuffix) > 0 && alsoSuffix[0] {
		val = "%" + val + "%"
	} else {
		val += "%"
	}

	return op(col, " NOT LIKE ", val)
}
