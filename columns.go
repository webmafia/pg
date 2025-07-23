package pg

import "github.com/webmafia/fast"

var _ QueryEncoder = columns{}

type columns []StringEncoder

// EncodeQuery implements QueryEncoder.
func (c columns) EncodeQuery(buf *fast.StringBuffer, _ *[]any) {
	for i := range c {
		if i != 0 {
			buf.WriteByte(',')
		}

		c[i].EncodeString(buf)
	}
}

// Merges multiple StringEncoders to a single QueryEncoder.
func Columns(cols []StringEncoder) QueryEncoder {
	return columns(cols)
}

func Col(s string) StringEncoder {
	return EncodeString(func(buf *fast.StringBuffer) {
		buf.WriteString(s)
	})
}
