package pg

import "context"

func (db *DB) Delete(ctx context.Context, table Identifier, cond QueryEncoder) (count int64, err error) {
	inst := db.instPool.Acquire()
	defer db.instPool.Release(inst)

	db.deleteQuery(inst, table, cond)

	cmd, err := db.exec(ctx, inst.buf.String(), inst.args...)

	if err != nil {
		err = instError(err, inst)
	} else {
		count = cmd.RowsAffected()
	}

	return
}

func (db *DB) deleteQuery(inst *inst, table Identifier, cond QueryEncoder) {
	inst.buf.WriteString("DELETE FROM ")
	table.EncodeString(inst.buf)
	inst.buf.WriteString(" WHERE ")
	cond.EncodeQuery(inst.buf, &inst.args)
}
