package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (db *DB) Transaction(ctx context.Context, readOnly ...bool) (tx Tx, err error) {
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

	tx = Tx{
		db:   db,
		conn: conn,
	}

	return
}

type Tx struct {
	db     *DB
	conn   *pgxpool.Conn
	closed bool
}

func (tx *Tx) Commit(ctx context.Context) (err error) {
	return tx.close(ctx, "COMMIT")
}

func (tx *Tx) Rollback(ctx context.Context) (err error) {
	return tx.close(ctx, "ROLLBACK")
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

// Query sends a query to the server and returns a Rows to read the results. Only errors encountered sending the query
// and initializing Rows will be returned. Err() on the returned Rows must be checked after the Rows is closed to
// determine if the query executed successfully.
//
// The returned Rows must be closed before the connection can be used again. It is safe to attempt to read from the
// returned Rows even if an error is returned. The error will be the available in rows.Err() after rows are closed. It
// is allowed to ignore the error returned from Query and handle it in Rows.
//
// It is possible for a call of FieldDescriptions on the returned Rows to return nil even if the Query call did not
// return an error.
//
// It is possible for a query to return one or more rows before encountering an error. In most cases the rows should be
// collected before processing rather than processed while receiving each row. This avoids the possibility of the
// application processing rows from a query that the server rejected. The CollectRows function is useful here.
func (tx *Tx) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	inst := tx.db.instPool.Acquire()
	defer tx.db.instPool.Release(inst)

	encodeQuery(inst.buf, query, args, &inst.args)

	return tx.conn.Query(ctx, query, inst.args...)
}

// QueryRow is a convenience wrapper over Query. Any error that occurs while querying is deferred
// until calling Scan on the returned Row. That Row will error with ErrNoRows if no rows are returned.
func (tx *Tx) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	inst := tx.db.instPool.Acquire()
	defer tx.db.instPool.Release(inst)

	encodeQuery(inst.buf, query, args, &inst.args)

	return tx.conn.QueryRow(ctx, query, inst.args...)
}

// Exec executes sql. sql can be either a prepared statement name or an SQL string. arguments should be referenced
// positionally from the sql string as $1, $2, etc.
func (tx *Tx) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	inst := tx.db.instPool.Acquire()
	defer tx.db.instPool.Release(inst)

	encodeQuery(inst.buf, query, args, &inst.args)

	return tx.conn.Exec(ctx, query, inst.args...)
}
