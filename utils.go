package pg

import "github.com/webmafia/fast"

//go:inline
func writeIdentifier(b *fast.StringBuffer, str string) {
	b.WriteByte('"')
	b.WriteString(str)
	b.WriteByte('"')
}

//go:inline
func writeAnyIdentifier(b *fast.StringBuffer, str any) {
	switch v := str.(type) {
	case StringEncoder:
		v.EncodeString(b)
	case string:
		writeIdentifier(b, v)
	}
}

func writeIdentifiers(b *fast.StringBuffer, strs []string) {
	for i := range strs {
		if i != 0 {
			b.WriteByte(',')
		}

		writeIdentifier(b, strs[i])
	}
}

//go:inline
func writeQueryArg(b *fast.StringBuffer, args *[]any, val any) {
	*args = append(*args, val)
	b.WriteByte('$')
	b.WriteInt(len(*args))
}

func writeAny(b *fast.StringBuffer, args *[]any, val any) {
	switch v := val.(type) {

	case StringEncoder:
		v.EncodeString(b)

	case QueryEncoder:
		v.EncodeQuery(b, args)

	default:
		writeQueryArg(b, args, val)

	}
}
