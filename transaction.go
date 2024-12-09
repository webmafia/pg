package pg

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrReleasedTransaction = errors.New("tried to operate on a released transaction")

const savepointName = "pseudo_tx"

func (db *DB) Transaction(ctx context.Context, readOnly ...bool) (tx *Tx, err error) {
	tx = db.txPool.Acquire()
	tx.ctx = ctx

	if parent, ok := ctx.(*Tx); ok {
		tx.conn = parent.conn
		tx.sp = true
	} else if tx.conn, err = db.db.Acquire(ctx); err != nil {
		tx.release()
		return
	}

	var ro bool

	if len(readOnly) > 0 {
		ro = readOnly[0]
	}

	if err = tx.begin(ctx, ro); err != nil {
		tx.release()
	}

	return
}

var _ context.Context = (*Tx)(nil)

type Tx struct {
	ctx    context.Context
	db     *DB
	conn   *pgxpool.Conn
	sp     bool
	closed bool
}

func (tx *Tx) begin(ctx context.Context, readOnly bool) (err error) {
	if tx.sp {
		_, err = tx.conn.Exec(ctx, "SAVEPOINT "+savepointName)
	} else if readOnly {
		_, err = tx.conn.Exec(ctx, "BEGIN READ ONLY")
	} else {
		_, err = tx.conn.Exec(ctx, "BEGIN")
	}

	return
}

func (tx *Tx) Commit(ctx ...context.Context) (err error) {
	if tx.closed {
		return ErrReleasedTransaction
	}

	c := tx.ctx

	if len(ctx) > 0 {
		c = ctx[0]
	}

	if tx.sp {
		return tx.close(c, "RELEASE SAVEPOINT "+savepointName)
	}

	return tx.close(c, "COMMIT")
}

func (tx *Tx) Release(ctx ...context.Context) (err error) {
	defer tx.release()

	c := tx.ctx

	if len(ctx) > 0 {
		c = ctx[0]
	}

	if tx.sp {
		return tx.close(c, "ROLLBACK TO SAVEPOINT "+savepointName, "RELEASE SAVEPOINT "+savepointName)
	}

	return tx.close(c, "ROLLBACK")
}

func (tx *Tx) close(ctx context.Context, query ...string) (err error) {
	if tx.closed {
		return
	}

	for i := range query {
		if _, err = tx.conn.Exec(ctx, query[i]); err != nil {
			return
		}
	}

	tx.closed = true
	return
}

func (tx *Tx) release() {
	tx.db.txPool.Release(tx)
}

// Deadline implements context.Context.
func (tx *Tx) Deadline() (deadline time.Time, ok bool) {
	return tx.ctx.Deadline()
}

// Done implements context.Context.
func (tx *Tx) Done() <-chan struct{} {
	return tx.ctx.Done()
}

// Err implements context.Context.
func (tx *Tx) Err() error {
	return tx.ctx.Err()
}

// Value implements context.Context.
func (tx *Tx) Value(key any) any {
	return tx.ctx.Value(key)
}
