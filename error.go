package pg

import (
	"fmt"

	"github.com/webmafia/fast"
)

var (
	_ error        = Error{}
	_ fmt.Stringer = Error{}
)

type Error struct {
	err   error
	query string
	args  []any
}

func errorFromCopy(err error, buf *fast.StringBuffer, args []any) Error {
	e := Error{
		err:   err,
		query: string(buf.Bytes()),
		args:  make([]any, len(args)),
	}

	copy(e.args, args)

	return e
}

func (err Error) Error() string {
	return err.err.Error()
}

func (err Error) String() string {
	return err.err.Error()
}

func (err Error) Query() string {
	return err.query
}

func (err Error) Args() []any {
	return err.args
}
