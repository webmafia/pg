package pg

import (
	"context"
	"slices"

	"github.com/webmafia/fast"
)

type Values struct {
	conflictColumns []string
	columns         []string
	values          []any
	db              *DB
	table           Identifier
	update          bool
}

func (r *Values) Reset() {
	r.columns = r.columns[:0]
	r.values = r.values[:0]
}

func (r *Values) Value(column string, value any) *Values {
	r.columns = append(r.columns, column)
	r.values = append(r.values, value)

	return r
}

func (r *Values) Exec(ctx context.Context) (err error) {
	inst := r.db.instPool.Acquire()
	defer r.db.instPool.Release(inst)

	r.EncodeString(inst.buf)
	_, err = r.db.db.Exec(ctx, inst.buf.String(), r.values...)

	if err != nil {
		err = errorFromCopy(err, inst.buf, r.values)
	}

	r.Reset()

	return
}

func (r *Values) ExecAndReturn(ctx context.Context, column string, bind any) (err error) {
	inst := r.db.instPool.Acquire()
	defer r.db.instPool.Release(inst)

	r.EncodeString(inst.buf)

	inst.buf.WriteByte('\n')
	inst.buf.WriteString("RETURNING ")
	writeIdentifier(inst.buf, column)

	row := r.db.db.QueryRow(ctx, inst.buf.String(), r.values...)
	err = row.Scan(bind)

	if err != nil {
		err = errorFromCopy(err, inst.buf, r.values)
	}

	r.Reset()

	return
}

var _ fast.StringEncoder = (*Values)(nil)

func (r *Values) EncodeString(b *fast.StringBuffer) {
	b.WriteString("INSERT INTO ")
	r.table.EncodeString(b)
	b.WriteString(" (")
	writeIdentifiers(b, r.columns)
	b.WriteString(") VALUES(")

	for i := range r.values {
		if i != 0 {
			b.WriteByte(',')
		}

		b.WriteByte('$')
		b.WriteInt(i + 1)
	}

	b.WriteByte(')')

	if len(r.conflictColumns) > 0 {
		b.WriteByte('\n')
		b.WriteString("ON CONFLICT (")
		writeIdentifiers(b, r.conflictColumns)
		b.WriteString(") ")

		var written bool

		for i := range r.columns {
			if slices.Contains(r.conflictColumns, r.columns[i]) {
				continue
			}

			if written {
				b.WriteString(", ")
			} else {
				b.WriteString("DO UPDATE SET ")
				written = true
			}

			writeIdentifier(b, r.columns[i])
			b.WriteString(" = EXCLUDED.")
			writeIdentifier(b, r.columns[i])
		}

		if !written {
			b.WriteString("DO NOTHING")
		}
	}
}
