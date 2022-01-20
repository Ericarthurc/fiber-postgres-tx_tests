package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

var DBPool *pgxpool.Pool

func DbConnect() {
	var err error
	DBPool, err = pgxpool.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("PGX_USER"), os.Getenv("PGX_PASSWORD"), os.Getenv("PGX_HOST"), os.Getenv("PGX_PORT"), os.Getenv("PGX_DATABASE")))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	_, err = DBPool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS ITEMS (
		id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
		product TEXT NOT NULL,
		manufacturer TEXT NOT NULL,
		device_type TEXT NOT NULL,
		serial TEXT NOT NULL,
		condition TEXT NOT NULL,
		year TEXT NOT NULL
	)`)
	if err != nil {
		fmt.Println(err.Error())
	}

	// _, err = DBPool.Exec(context.Background(), `CREATE EXTENSION pg_trgm`)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// _, err = DBPool.Exec(context.Background(), `ALTER TABLE ITEMS ADD COLUMN tsv tsvector`)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// _, err = DBPool.Exec(context.Background(), `UPDATE items SET tsv =
	// 	setweight(to_tsvector(product), 'A') ||
	// 	setweight(to_tsvector(serial), 'B') ||
	// 	setweight(to_tsvector(year), 'C') ||
	// 	setweight(to_tsvector(condition), 'D')`)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// _, err = DBPool.Exec(context.Background(), `CREATE INDEX ix_items_tsv ON items USING GIN(tsv)`)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// _, err = DBPool.Exec(context.Background(), `CREATE TYPE user_roles AS ENUM ('admin', 'maintainer', 'guest')`)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	_, err = DBPool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS USERS (
		id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		role TEXT NOT NULL,
		tokens TEXT[] NOT NULL
	)`)
	if err != nil {
		fmt.Println(err.Error())
	}

}
