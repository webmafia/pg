package pg

import "github.com/webmafia/fast"

var _ fast.StringEncoder = Table("")

type Table string

// EncodeString implements fast.StringEncoder.
func (t Table) EncodeString(b *fast.StringBuffer) {
	b.WriteString(string(t))
}
