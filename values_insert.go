package pg

func (db *DB) Insert(table Identifier, conflictColumns ...string) *Values {
	v := db.valPool.Acquire()
	v.table = table
	v.conflictColumns = append(v.conflictColumns, conflictColumns...)

	return v
}
