package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type CopyFromSource = pgx.CopyFromSource

// Insert data (usually a very big batch) into a table via Postgres' COPY command. Any error will abort the whole batch.
func (db *DB) CopyFrom(ctx context.Context, tableName Identifier, columnNames []string, rowSrc CopyFromSource) (int64, error) {
	if tx, ok := ctx.(*Tx); ok {
		return tx.conn.CopyFrom(ctx, pgx.Identifier{string(tableName)}, columnNames, rowSrc)
	}

	return db.db.CopyFrom(ctx, pgx.Identifier{string(tableName)}, columnNames, rowSrc)
}
