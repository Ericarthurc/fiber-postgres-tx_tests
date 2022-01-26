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

	_, err = DBPool.Exec(context.Background(), `SET timezone = 'America/Los_Angeles'`)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	_, err = DBPool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS SERVICES (
		id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
		time TIMESTAMPTZ NOT NULL,
		seats INTEGER NOT NULL
	)`)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	_, err = DBPool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS USERS (
		id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		userseats INTEGER NOT NULL CHECK (userseats >= 1 AND userseats <= 10),
		servicetime TIMESTAMPTZ NOT NULL,
		serviceid TEXT NOT NULL
	)`)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

}
