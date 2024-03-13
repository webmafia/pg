package pg

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
