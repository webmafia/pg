package pg

import (
	"context"
	"strconv"

	"github.com/cespare/xxhash/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/webmafia/fast"
	"github.com/webmafia/pg/internal/lru"
)

const connDataKey = "github.com/webmafia/pg"

type connectionMemory struct {
	buf   *fast.StringBuffer
	args  []any
	hash  xxhash.Digest
	stmts *lru.Cache[uint64, *pgconn.StatementDescription]
}

func (c *connectionMemory) reset() {

	// Don't keep too large buffer
	if c.buf.Cap() > 1024 {
		c.buf = fast.NewStringBuffer(256)
	} else {
		c.buf.Reset()
	}

	clear(c.args)
	c.args = c.args[:0]
	c.hash.Reset()
}

func getConnMem(conn *pgx.Conn) *connectionMemory {
	m := conn.PgConn().CustomData()

	if c, ok := m[connDataKey].(*connectionMemory); ok {
		return c
	}

	c := &connectionMemory{
		buf:  fast.NewStringBuffer(256),
		args: make([]any, 0, 4),
		stmts: lru.NewCache[uint64, *pgconn.StatementDescription](64, func(_ uint64, stmt *pgconn.StatementDescription) {
			conn.Deallocate(context.Background(), stmt.Name)
		}),
	}

	c.reset()

	m[connDataKey] = c
	return c
}

func (c *connectionMemory) stmt(ctx context.Context, conn *pgx.Conn, format string, args []any) (stmt *pgconn.StatementDescription, err error) {
	if err = encodeQuery(c.buf, format, args, &c.args); err != nil {
		return nil, err
	}

	if _, err = c.hash.Write(c.buf.Bytes()); err != nil {
		return
	}

	return c.stmts.GetOrSet(c.hash.Sum64(), func(hash uint64) (stmt *pgconn.StatementDescription, err error) {
		stmt, err = conn.Prepare(ctx, "stmt_"+strconv.FormatUint(hash, 36), c.buf.String())

		if err == nil {
			stmt.SQL = ""
		}

		return
	})
}
