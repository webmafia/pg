package pg

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/webmafia/fast"
)

type DB struct {
	db       *pgxpool.Pool
	instPool *fast.Pool[inst]
	valPool  *fast.Pool[Values]
}

type inst struct {
	buf  *fast.StringBuffer
	args []any
}

func NewDB(pool *pgxpool.Pool) *DB {
	db := &DB{
		db: pool,
		instPool: fast.NewPool[inst](func(i *inst) {
			i.buf = fast.NewStringBuffer(256)
			i.args = make([]any, 0, 5)
		}, func(i *inst) {
			i.buf.Reset()
			i.args = i.args[:0]
		}),
	}

	db.valPool = fast.NewPool[Values](func(v *Values) {
		v.conflictColumns = make([]string, 0, 5)
		v.columns = make([]string, 0, 5)
		v.values = make([]any, 0, 5)
		v.db = db
	}, func(v *Values) {
		v.conflictColumns = v.conflictColumns[:0]
		v.columns = v.columns[:0]
		v.values = v.values[:0]
		v.update = false
	})

	return db
}
