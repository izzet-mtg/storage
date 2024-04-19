package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	apiv1 "github.com/izzet-mtg/storage/services/backend/api/v1"
	adminapiv1 "github.com/izzet-mtg/storage/services/backend/api/v1/admin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func run(conStr string, redisURL string, exp time.Duration) error {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, conStr)
	if err != nil {
		return err
	}
	defer pool.Close()

	ropts, err := redis.ParseURL(redisURL)
	if err != nil {
		return err
	}
	rc := redis.NewClient(ropts)

	r := gin.Default()
	v1 := r.Group("v1")
	v1.POST("login", apiv1.Login(pool, rc, exp))
	v1.DELETE("logout", apiv1.Logout(rc))
	adminv1 := v1.Group("admin")
	adminv1.POST("user", adminapiv1.CreateUser(pool, rc, exp))
	r.Run()

	return nil
}

func main() {
	exp, err := time.ParseDuration("24h")
	if err != nil {
		log.Fatal(err)
	}
	if err := run(os.Getenv("DB_URI"), os.Getenv("REDIS_URL"), exp); err != nil {
		log.Fatal(err)
	}
}
