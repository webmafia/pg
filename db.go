package pg

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/webmafia/fast"
)

type DB struct {
	db      *pgxpool.Pool
	bufPool *fast.Pool[fast.StringBuffer]
	argPool *fast.Pool[[]any]
}

func NewDB(db *pgxpool.Pool) *DB {
	return &DB{
		db: db,
		bufPool: fast.NewPool[fast.StringBuffer](func(sb *fast.StringBuffer) {
			sb.Grow(256)
		}, func(sb *fast.StringBuffer) {
			sb.Reset()
		}),
		argPool: fast.NewPool[[]any](func(a *[]any) {
			*a = make([]any, 0, 5)
		}, func(a *[]any) {
			*a = (*a)[:0]
		}),
	}
}
