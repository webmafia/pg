package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/webmafia/fast"
)

type DB struct {
	db      *pgxpool.Pool
	valPool *fast.Pool[Values]
}

func New(ctx context.Context, connString string, alterConfig ...func(*pgxpool.Config)) (db *DB, err error) {
	config, err := pgxpool.ParseConfig(connString)

	if err != nil {
		return
	}

	if len(alterConfig) > 0 && alterConfig[0] != nil {
		alterConfig[0](config)
	}

	oldBeforeConnect := config.BeforeConnect
	oldAfterRelease := config.AfterRelease
	oldBeforeClose := config.BeforeClose

	config.BeforeConnect = func(ctx context.Context, cc *pgx.ConnConfig) error {

		if oldBeforeConnect != nil {
			if err := oldBeforeConnect(ctx, cc); err != nil {
				return err
			}
		}

		cc.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe
		cc.StatementCacheCapacity = 0

		return nil
	}

	config.AfterRelease = func(c *pgx.Conn) bool {
		resetConnMem(c)

		if oldAfterRelease != nil {
			return oldAfterRelease(c)
		}

		return true
	}

	config.BeforeClose = func(c *pgx.Conn) {
		purgeConnMem(c)

		if oldBeforeClose != nil {
			oldBeforeClose(c)
		}
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		return
	}

	return newDB(pool), nil
}

func newDB(db *pgxpool.Pool) *DB {
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
