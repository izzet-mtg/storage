package main

import (
	"context"
	"log"
	"os"
	//"reflect"

	"github.com/gin-gonic/gin"
	adminv1 "github.com/izzet-mtg/storage/services/backend/api/v1/admin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func run(pgConStr string) error {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, pgConStr)
	if err != nil {
		return err
	}
	defer pool.Close()

	r := gin.Default()
	v1 := r.Group("v1")
	v1.POST("admin/user", adminv1.CreateUser(pool))
	r.Run()

	return nil
}

func main() {
	if err := run(os.Getenv("DB_URI")); err != nil {
		log.Fatal(err)
	}
}
