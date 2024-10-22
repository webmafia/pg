package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/webmafia/pg"
)

func example(ctx context.Context) (err error) {
	config, err := pgxpool.ParseConfig("postgresql://postgres:postgres@localhost/postgres?sslmode=disable")

	if err != nil {
		return
	}

	config.BeforeConnect = func(ctx context.Context, cc *pgx.ConnConfig) error {
		cc.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe
		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		return
	}

	db := pg.NewDB(pool)
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
