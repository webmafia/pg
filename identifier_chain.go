package pg

import "github.com/webmafia/fast"

var _ StringEncoder = ChainedIdentifier{Identifier(""), Identifier("")}

type ChainedIdentifier [2]StringEncoder

// EncodeString implements StringEncoder.
func (t ChainedIdentifier) EncodeString(b *fast.StringBuffer) {
	t[0].EncodeString(b)
	b.WriteByte('.')
	t[1].EncodeString(b)
}

func (t ChainedIdentifier) Col(col string) ChainedIdentifier {
	return ChainedIdentifier{t, Identifier(col)}
}
