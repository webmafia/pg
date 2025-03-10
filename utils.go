package pg

import (
	"strings"

	"github.com/webmafia/fast"
)

var closedChannel chan struct{}

func init() {
	closedChannel = make(chan struct{})
	close(closedChannel)
}

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
		dot := strings.IndexByte(v, '.')

		if dot >= 0 {
			writeIdentifier(b, v[:dot])
			b.WriteByte('.')
			v = v[dot+1:]
		}

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
		if v != nil {
			v.EncodeString(b)
		}

	case QueryEncoder:
		if v != nil {
			v.EncodeQuery(b, args)
		}

	default:
		writeQueryArg(b, args, val)

	}
}
