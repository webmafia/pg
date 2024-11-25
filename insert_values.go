package pg

import (
	"context"
	"slices"

	"github.com/webmafia/fast"
)

type InsertOptions struct {
	OnConflict   func(vals *Values) EncodeQuery
	ReturnColumn string
	ReturnDst    any // A pointer that ReturnColumn should be scanned to
}

func (db *DB) InsertValues(ctx context.Context, table Identifier, vals *Values, options ...InsertOptions) (count int64, err error) {
	if vals.Empty() {
		return
	}

	var onConflict EncodeQuery
	var returning QueryEncoder

	if len(options) > 0 {
		if options[0].OnConflict != nil {
			onConflict = options[0].OnConflict(vals)
		}

		if options[0].ReturnColumn != "" {
			returning = Raw("RETURNING %v", Identifier(options[0].ReturnColumn))
		}
	}

	if len(options) > 0 && options[0].ReturnDst != nil {
		row := db.QueryRow(ctx,
			"INSERT INTO %T (%T) VALUES(%T) %T %T",
			table,
			vals.colEncoder(),
			vals.valEncoder(),
			onConflict,
			returning,
		)

		if err = row.Scan(options[0].ReturnDst); err != nil {
			return
		}

		count = 1
	} else {
		cmd, err := db.Exec(ctx,
			"INSERT INTO %T (%T) VALUES(%T) %T %T",
			table,
			vals.colEncoder(),
			vals.valEncoder(),
			onConflict,
			returning,
		)

		if err != nil {
			return 0, err
		}

		count = cmd.RowsAffected()
	}

	vals.reset()

	return
}

func DoUpdate(numConflictingColumns int, ignoreColumns ...string) func(vals *Values) EncodeQuery {
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
