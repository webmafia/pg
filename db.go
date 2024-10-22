package pg

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/webmafia/fast"
)

type DB struct {
	db      *pgxpool.Pool
	valPool *fast.Pool[Values]
}

func NewDB(db *pgxpool.Pool) *DB {
	return &DB{
		db: db,
		valPool: fast.NewPool(func(v *Values) {
			v.columns = make([]string, 0, 8)
			v.values = make([]any, 0, 8)
		}, func(v *Values) {
			v.reset()
		}),
	}
}

func (db *DB) Close() {
	db.db.Close()
}
