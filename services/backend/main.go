package main

import (
	"context"
	"log"
	//"reflect"

	"github.com/jackc/pgx/v5/pgxpool"
	//"github.com/jackc/pgx/v5/pgtype"

	//"github.com/izzet-mtg/storage/services/backend/db"
)

func run(pgConStr string) error {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, pgConStr)
	if err != nil {
		return err
	}
	defer pool.Close()

	return nil
}

func main() {
	if err := run(""); err != nil {
		log.Fatal(err)
	}
}
