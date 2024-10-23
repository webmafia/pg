package main

import (
	"context"
	"log"

	"github.com/webmafia/pg"
)

func example(ctx context.Context) (err error) {
	db, err := pg.New(ctx, "postgresql://postgres:postgres@localhost/postgres?sslmode=disable")

	if err != nil {
		return
	}

	defer db.Close()

	vals := db.AcquireValues()
	defer db.ReleaseValues(vals)

	if _, err = db.InsertValues(ctx, pg.Identifier("test"), vals.Value("name", "x")); err != nil {
		return
	}

	if _, err = db.InsertValues(ctx, pg.Identifier("test"), vals.Value("name", "y")); err != nil {
		return
	}

	if _, err = db.InsertValues(ctx, pg.Identifier("test"), vals.Value("name", "x")); err != nil {
		return
	}

	log.Println("done")

	return
}
