package pg

import "github.com/webmafia/fast"

var _ fast.StringEncoder = Identifier("")

type Identifier string

// EncodeString implements fast.StringEncoder.
func (t Identifier) EncodeString(b *fast.StringBuffer) {
	b.WriteByte('"')
	b.WriteString(string(t))
	b.WriteByte('"')
}
