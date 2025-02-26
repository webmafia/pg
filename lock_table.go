package pg

import (
	"context"
	"errors"

	"github.com/webmafia/fast"
)

var _ StringEncoder = LockMode(0)

type LockMode byte

const (
	SharedLock    LockMode = iota // Others can read, but not write
	ExclusiveLock                 // Other can neither read nor write
)

// EncodeString implements fast.StringEncoder.
func (lm LockMode) EncodeString(b *fast.StringBuffer) {
	if lm == ExclusiveLock {
		b.WriteString("EXCLUSIVE MODE")
	} else {
		b.WriteString("SHARE MODE")
	}
}

// Locks a table during the remainder of the transaction. Returns an error if not within a transaction.
// Lock mode can be either SharedLock or ExclusiveLock.
func (db *DB) LockTable(ctx context.Context, table Identifier, lockMode ...LockMode) (err error) {
	if _, ok := ctx.(*Tx); !ok {
		return errors.New("cannot lock table outside transaction")
	}

	var mode LockMode = SharedLock

	if lockMode != nil {
		mode = lockMode[0]
	}

	_, err = db.Exec(ctx, "LOCK TABLE %t IN %s", mode)
	return
}
