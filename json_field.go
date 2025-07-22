package pg

import "github.com/webmafia/fast"

var _ StringEncoder = JSONField("")

type JSONField string

// EncodeString implements StringEncoder.
func (t JSONField) EncodeString(b *fast.StringBuffer) {
	b.WriteByte('\'')
	b.WriteString(string(t))
	b.WriteByte('\'')
}
