package pg

import "github.com/webmafia/fast"

type inst struct {
	buf  *fast.StringBuffer
	args []any
}
