package pg

import "github.com/webmafia/fast"

type StringEncoder = fast.StringEncoder

var _ StringEncoder = (EncodeString)(nil)

type EncodeString func(buf *fast.StringBuffer)

func (fn EncodeString) EncodeString(buf *fast.StringBuffer) {
	fn(buf)
}
