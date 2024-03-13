package pg

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func (db *DB) Transaction(ctx context.Context, readOnly ...bool) (tx *Tx, err error) {
	conn, err := db.db.Acquire(ctx)

	if err != nil {
		return
	}

	if len(readOnly) > 0 && readOnly[0] {
		if _, err = conn.Exec(ctx, "BEGIN READ ONLY"); err != nil {
			return
		}
	} else {
		if _, err = conn.Exec(ctx, "BEGIN"); err != nil {
			return
		}
	}

	tx = &Tx{
		ctx:  ctx,
		conn: conn,
	}

	return
}

var _ context.Context = (*Tx)(nil)

type Tx struct {
	ctx    context.Context
	conn   *pgxpool.Conn
	closed bool
}

func (tx *Tx) Commit(ctx ...context.Context) (err error) {
	if len(ctx) > 0 {
		return tx.close(ctx[0], "COMMIT")
	}

	return tx.close(tx.ctx, "COMMIT")
}

func (tx *Tx) Rollback(ctx ...context.Context) (err error) {
	if len(ctx) > 0 {
		return tx.close(ctx[0], "ROLLBACK")
	}

	return tx.close(tx.ctx, "ROLLBACK")
}

func (tx *Tx) close(ctx context.Context, query string) (err error) {
	if tx.closed {
		return
	}

	_, err = tx.conn.Exec(ctx, query)
	tx.conn.Release()
	tx.closed = true
	return
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
