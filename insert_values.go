package pg

import (
	"context"
	"slices"

	"github.com/webmafia/fast"
)

func (db *DB) InsertValues(ctx context.Context, table Identifier, vals *Values, conflictColumns ...string) (count int64, err error) {
	if vals.Empty() {
		return
	}

	cmd, err := db.Exec(ctx,
		"INSERT INTO %T (%T) VALUES(%T) %T",
		table,
		vals.colEncoder(),
		vals.valEncoder(),
		conflictingColumns(vals, conflictColumns),
	)

	if err == nil {
		vals.reset()
		count = cmd.RowsAffected()
	}

	return
}

func conflictingColumns(vals *Values, cols []string) EncodeQuery {
	return func(buf *fast.StringBuffer, queryArgs *[]any) {
		if len(cols) == 0 {
			return
		}

		buf.WriteString("ON CONFLICT (")
		writeIdentifier(buf, cols[0])
		buf.WriteString(") ")

		var written bool

		for i := range vals.columns {
			if slices.Contains(cols, vals.columns[i]) {
				continue
			}

			if written {
				buf.WriteString(", ")
			} else {
				buf.WriteString("DO UPDATE SET ")
				written = true
			}

			writeIdentifier(buf, vals.columns[i])
			buf.WriteString(" = EXCLUDED.")
			writeIdentifier(buf, vals.columns[i])
		}

		if !written {
			buf.WriteString("DO NOTHING")
		}
	}
}
