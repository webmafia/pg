package pg

import (
	"context"
)

func (db *DB) UpdateValues(ctx context.Context, table Identifier, vals *Values, cond QueryEncoder) (count int64, err error) {
	if vals.Empty() {
		return
	}

	inst := db.instPool.Acquire()
	defer db.instPool.Release(inst)

	inst.buf.WriteString("UPDATE ")
	table.EncodeString(inst.buf)
	inst.buf.WriteString(" SET ")

	for i := range vals.columns {
		if i != 0 {
			inst.buf.WriteString(", ")
		}

		writeIdentifier(inst.buf, vals.columns[i])
		inst.buf.WriteString(" = ")
		writeQueryArg(inst.buf, &inst.args, vals.values[i])
	}

	inst.buf.WriteString(" WHERE ")
	cond.EncodeQuery(inst.buf, &inst.args)

	cmd, err := db.exec(ctx, inst.buf.String(), inst.args...)

	if err != nil {
		err = instError(err, inst)
	} else {
		count = cmd.RowsAffected()
	}

	vals.reset()

	return
}
