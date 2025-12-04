package pg

import "github.com/webmafia/fast"

var _ StringEncoder = ChainedIdentifier{Identifier(""), Identifier("")}

type ChainedIdentifier [2]StringEncoder

// EncodeString implements StringEncoder.
func (t ChainedIdentifier) EncodeString(b *fast.StringBuffer) {
	switch t[0].(type) {
	case ChainedIdentifier:
		b.WriteByte('(')
		t[0].EncodeString(b)
		b.WriteByte(')')
	default:
		t[0].EncodeString(b)
	}

	b.WriteByte('.')
	t[1].EncodeString(b)
}

func (t ChainedIdentifier) Col(col string) ChainedIdentifier {
	return ChainedIdentifier{t, Identifier(col)}
}

func (t ChainedIdentifier) Alias(col string) Alias {
	return Alias{t[1], Identifier(col)}
}

func (t ChainedIdentifier) GetJSONValue(field string) ChainedJSONField {
	return ChainedJSONField{t, JSONField(field)}
}
