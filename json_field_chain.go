package pg

import "github.com/webmafia/fast"

var _ StringEncoder = ChainedJSONField{Identifier(""), JSONField("")}

type ChainedJSONField [2]StringEncoder

// EncodeString implements StringEncoder.
func (t ChainedJSONField) EncodeString(b *fast.StringBuffer) {
	t[0].EncodeString(b)
	b.Write([]byte{'-', '>'})
	t[1].EncodeString(b)
}
