package pg

import "context"

func (db *DB) Delete(ctx context.Context, table Identifier, cond QueryEncoder) (count int64, err error) {
	cmd, err := db.Exec(ctx, "DELETE FROM %T WHERE %T", table, cond)

	if err != nil {
		return
	}

	return cmd.RowsAffected(), nil
}
