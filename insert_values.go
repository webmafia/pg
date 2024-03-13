package pg

import (
	"context"
	"slices"
)

func (db *DB) InsertValues(ctx context.Context, table Identifier, vals *Values, conflictColumns ...string) (count int64, err error) {
	if vals.Empty() {
		return
	}

	inst := db.instPool.Acquire()
	defer db.instPool.Release(inst)

	db.insertValuesQuery(inst, table, vals, conflictColumns)

	cmd, err := db.exec(ctx, inst.buf.String(), inst.args...)

	if err != nil {
		err = instError(err, inst)
	} else {
		count = cmd.RowsAffected()
	}

	vals.reset()

	return
}

func (db *DB) insertValuesQuery(inst *inst, table Identifier, vals *Values, conflictColumns []string) {
	inst.buf.WriteString("INSERT INTO ")
	table.EncodeString(inst.buf)
	inst.buf.WriteString(" (")
	writeIdentifiers(inst.buf, vals.columns)
	inst.buf.WriteString(") VALUES(")

	for i := range vals.values {
		if i != 0 {
			inst.buf.WriteByte(',')
		}

		writeQueryArg(inst.buf, &inst.args, vals.values[i])
	}

	inst.buf.WriteByte(')')

	if len(conflictColumns) > 0 {
		inst.buf.WriteByte('\n')
		inst.buf.WriteString("ON CONFLICT (")
		writeIdentifiers(inst.buf, conflictColumns)
		inst.buf.WriteString(") ")

		var written bool

		for i := range vals.columns {
			if slices.Contains(conflictColumns, vals.columns[i]) {
				continue
			}

			if written {
				inst.buf.WriteString(", ")
			} else {
				inst.buf.WriteString("DO UPDATE SET ")
				written = true
			}

			writeIdentifier(inst.buf, vals.columns[i])
			inst.buf.WriteString(" = EXCLUDED.")
			writeIdentifier(inst.buf, vals.columns[i])
		}

		if !written {
			inst.buf.WriteString("DO NOTHING")
		}
	}
}
