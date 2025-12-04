package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

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
func (db *DB) Query(ctx context.Context, query string, args ...any) (rows pgx.Rows, err error) {
	var (
		conn          *pgxpool.Conn
		inTransaction bool
	)

	// If inside a transaction...
	if tx, ok := ctx.(*Tx); ok && tx.conn != nil {

		// ...use the transaction's connection - and don't release it!
		conn = tx.conn
		inTransaction = true
	} else {

		// Otherwise, acquire a connection from the pool - but don't release it!
		if conn, err = db.db.Acquire(ctx); err != nil {
			return
		}
	}

	stmt, mem, err := db.prepare(ctx, conn, query, args)

	if err != nil {
		conn.Release()
		return
	}

	defer mem.reset()

	if rows, err = conn.Query(ctx, stmt.Name, mem.args...); err != nil {
		conn.Release()
		return
	}

	// The connection can't be released until the rows are closed, thus we need to wrap the returned rows
	// and pass along the connection so that it can be released when done. However, if we're inside a transaction,
	// we should NOT wrap it as we don't want to release the connection until the transaction is done.
	if !inTransaction {
		rows = &poolRows{r: rows, c: conn}
	}

	return
}

// QueryRow is a convenience wrapper over Query. Any error that occurs while querying is deferred
// until calling Scan on the returned Row. That Row will error with ErrNoRows if no rows are returned.
func (db *DB) QueryRow(ctx context.Context, query string, args ...any) (row pgx.Row) {
	var (
		conn          *pgxpool.Conn
		inTransaction bool
		err           error
	)

	// If inside a transaction...
	if tx, ok := ctx.(*Tx); ok && tx.conn != nil {

		// ...use the transaction's connection - and don't release it!
		conn = tx.conn
		inTransaction = true
	} else {

		// Otherwise, acquire a connection from the pool - but don't release it!
		if conn, err = db.db.Acquire(ctx); err != nil {
			return &poolRow{err: err}
		}
	}

	stmt, mem, err := db.prepare(ctx, conn, query, args)

	if err != nil {
		conn.Release()
		return &poolRow{err: err}
	}

	defer mem.reset()

	row = conn.QueryRow(ctx, stmt.Name, mem.args...)

	// The connection can't be released until the rows are closed, thus we need to wrap the returned rows
	// and pass along the connection so that it can be released when done. However, if we're inside a transaction,
	// we should NOT wrap it as we don't want to release the connection until the transaction is done.
	if !inTransaction {
		row = &poolRow{r: row, c: conn}
	}

	return
}

// Execute a query that doesn't return any rows. Arguments are passed in printf syntax.
func (db *DB) Exec(ctx context.Context, query string, args ...any) (cmd pgconn.CommandTag, err error) {
	var conn *pgxpool.Conn

	// If inside a transaction...
	if tx, ok := ctx.(*Tx); ok && tx.conn != nil {

		// ...use the transaction's connection - and don't release it!
		conn = tx.conn
	} else {

		// Otherwise, acquire a connection from the pool...
		if conn, err = db.db.Acquire(ctx); err != nil {
			return
		}

		// ...and release it once done.
		defer conn.Release()
	}

	stmt, mem, err := db.prepare(ctx, conn, query, args)

	if err != nil {
		return
	}

	defer mem.reset()

	return conn.Exec(ctx, stmt.Name, mem.args...)
}

func (db *DB) prepare(ctx context.Context, conn *pgxpool.Conn, query string, args []any) (stmt *pgconn.StatementDescription, mem *connectionMemory, err error) {
	c := conn.Conn()
	mem = getConnMem(c)

	if stmt, err = mem.stmt(ctx, c, query, args); err != nil {
		mem.reset()
	}

	return
}
