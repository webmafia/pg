package pg

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/webmafia/fast"
)

type DB struct {
	db       *pgxpool.Pool
	instPool *fast.Pool[inst]
}

type inst struct {
	buf  *fast.StringBuffer
	args []any
}

func NewDB(db *pgxpool.Pool) *DB {
	return &DB{
		db: db,
		instPool: fast.NewPool[inst](func(i *inst) {
			i.buf = fast.NewStringBuffer(256)
			i.args = make([]any, 0, 5)
		}, func(i *inst) {
			i.buf.Reset()
			i.args = i.args[:0]
		}),
	}
}
