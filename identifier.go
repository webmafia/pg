package pg

import "github.com/webmafia/fast"

var _ StringEncoder = Identifier("")

type Identifier string

// EncodeString implements StringEncoder.
func (t Identifier) EncodeString(b *fast.StringBuffer) {
	b.WriteByte('"')
	b.WriteString(string(t))
	b.WriteByte('"')
}
