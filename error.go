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
