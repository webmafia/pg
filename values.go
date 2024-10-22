package pg

import "github.com/webmafia/fast"

var (
	_ QueryEncoder = (*Values)(nil)
	_ QueryEncoder = cols{}
	_ QueryEncoder = vals{}
)

//go:inline
func (db *DB) AcquireValues() *Values {
	return db.valPool.Acquire()
}

//go:inline
func (db *DB) ReleaseValues(v *Values) {
	db.valPool.Release(v)
}

type Values struct {
	columns []string
	values  []any
}

func (r *Values) reset() {
	clear(r.columns)
	clear(r.values)
	r.columns = r.columns[:0]
	r.values = r.values[:0]
}

func (r *Values) Value(column string, value any) *Values {
	r.columns = append(r.columns, column)
	r.values = append(r.values, value)

	return r
}

func (r *Values) Len() int {
	return len(r.columns)
}

func (r *Values) Empty() bool {
	return len(r.columns) == 0
}

// EncodeQuery implements QueryEncoder.
func (r *Values) EncodeQuery(buf *fast.StringBuffer, queryArgs *[]any) {
	for i := range r.columns {
		if i != 0 {
			buf.WriteString(", ")
		}

		writeIdentifier(buf, r.columns[i])
		buf.WriteString(" = ")
		writeQueryArg(buf, queryArgs, r.values[i])
	}
}

func (r *Values) colEncoder() QueryEncoder {
	return cols{v: r}
}

type cols struct {
	v *Values
}

// EncodeQuery implements QueryEncoder.
func (c cols) EncodeQuery(buf *fast.StringBuffer, _ *[]any) {
	writeIdentifiers(buf, c.v.columns)
}

func (r *Values) valEncoder() QueryEncoder {
	return vals{v: r}
}

type vals struct {
	v *Values
}

// EncodeQuery implements QueryEncoder.
func (c vals) EncodeQuery(buf *fast.StringBuffer, args *[]any) {
	for i := range c.v.values {
		if i != 0 {
			buf.WriteByte(',')
		}

		writeQueryArg(buf, args, c.v.values[i])
	}
}
