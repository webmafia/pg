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

	if _, err = db.Exec(ctx, `INSERT INTO %T ("name") VALUES(%T)`, pg.Identifier("test"), "a"); err != nil {
		return
	}

	if _, err = db.Exec(ctx, `INSERT INTO %T ("name") VALUES(%T)`, pg.Identifier("test"), "b"); err != nil {
		return
	}

	if _, err = db.Exec(ctx, `INSERT INTO %T ("name") VALUES(%T)`, pg.Identifier("test"), "c"); err != nil {
		return
	}

	log.Println("done")

	return
}
