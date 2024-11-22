package pg

import (
	"context"
	"slices"

	"github.com/webmafia/fast"
)

type InsertOption func(vals *Values) EncodeQuery

func (db *DB) InsertValues(ctx context.Context, table Identifier, vals *Values, options ...InsertOption) (count int64, err error) {
	if vals.Empty() {
		return
	}

	var upsert EncodeQuery

	if len(options) > 0 && options[0] != nil {
		upsert = options[0](vals)
	}

	cmd, err := db.Exec(ctx,
		"INSERT INTO %T (%T) VALUES(%T) %T",
		table,
		vals.colEncoder(),
		vals.valEncoder(),
		upsert,
	)

	if err == nil {
		vals.reset()
		count = cmd.RowsAffected()
	}

	return
}

func Upsert(numConflictingColumns int, ignoreColumns ...string) InsertOption {
	return func(vals *Values) EncodeQuery {
		return func(buf *fast.StringBuffer, queryArgs *[]any) {
			if len(ignoreColumns) == 0 {
				return
			}

			if numConflictingColumns <= 0 || numConflictingColumns > len(ignoreColumns) {
				numConflictingColumns = len(ignoreColumns)
			}

			buf.WriteString("ON CONFLICT (")
			writeIdentifiers(buf, ignoreColumns[:numConflictingColumns])
			buf.WriteString(") ")

			var written bool

			for i := range vals.columns {
				if slices.Contains(ignoreColumns, vals.columns[i]) {
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
}
