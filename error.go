package pg

import (
	"fmt"
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

func instError(err error, inst *inst) Error {
	e := Error{
		err:   err,
		query: string(inst.buf.Bytes()),
		args:  make([]any, len(inst.args)),
	}

	copy(e.args, inst.args)

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
