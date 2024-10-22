package pg

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	_ pgx.Rows = (*poolRows)(nil)
	_ pgx.Row  = (*poolRow)(nil)
)

type poolRows struct {
	r   pgx.Rows
	c   *pgxpool.Conn
	err error
}

func (rows *poolRows) Close() {
	rows.r.Close()
	if rows.c != nil {
		rows.c.Release()
		rows.c = nil
	}
}

func (rows *poolRows) Err() error {
	if rows.err != nil {
		return rows.err
	}
	return rows.r.Err()
}

func (rows *poolRows) CommandTag() pgconn.CommandTag {
	return rows.r.CommandTag()
}

func (rows *poolRows) FieldDescriptions() []pgconn.FieldDescription {
	return rows.r.FieldDescriptions()
}

func (rows *poolRows) Next() bool {
	if rows.err != nil {
		return false
	}

	n := rows.r.Next()
	if !n {
		rows.Close()
	}
	return n
}

func (rows *poolRows) Scan(dest ...any) error {
	err := rows.r.Scan(dest...)
	if err != nil {
		rows.Close()
	}
	return err
}

func (rows *poolRows) Values() ([]any, error) {
	values, err := rows.r.Values()
	if err != nil {
		rows.Close()
	}
	return values, err
}

func (rows *poolRows) RawValues() [][]byte {
	return rows.r.RawValues()
}

func (rows *poolRows) Conn() *pgx.Conn {
	return rows.r.Conn()
}

type poolRow struct {
	r   pgx.Row
	c   *pgxpool.Conn
	err error
}

func (row *poolRow) Scan(dest ...any) error {
	if row.err != nil {
		return row.err
	}

	panicked := true
	defer func() {
		if panicked && row.c != nil {
			row.c.Release()
		}
	}()
	err := row.r.Scan(dest...)
	panicked = false
	if row.c != nil {
		row.c.Release()
	}
	return err
}
