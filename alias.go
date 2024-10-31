package pg

import "github.com/webmafia/fast"

var _ StringEncoder = Alias{Identifier(""), Identifier("")}

type Alias [2]StringEncoder

// EncodeString implements StringEncoder.
func (t Alias) EncodeString(b *fast.StringBuffer) {
	t[0].EncodeString(b)
	b.WriteString(" AS ")
	t[1].EncodeString(b)
}

func (t Alias) Col(col string) ChainedIdentifier {
	return ChainedIdentifier{t[1], Identifier(col)}
}
