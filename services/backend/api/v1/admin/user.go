package adminapi

import (
	"crypto/sha512"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	"github.com/izzet-mtg/storage/services/backend/db"
)

func generateHashPassword(p string) (string, error) {
	s := sha512.Sum512([]byte(p))
	hp, err := bcrypt.GenerateFromPassword(s[:], 10)
	if err != nil {
		return "", err
	}
	return string(hp), nil
}

type User struct {
	Name     string `validate:"required"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
	IsAdmin  bool
}

func CreateUser(p *pgxpool.Pool, rc *redis.Client, exp time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		queries := db.New(p)

		var u User
		if err := c.ShouldBindJSON(&u); err != nil {
			log.Println("[Info] reject user request with invalid structure")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user structure"})
			return
		}

		validate := validator.New(validator.WithRequiredStructEnabled())
		if err := validate.Struct(u); err != nil {
			log.Println("[Info] reject user request with invalid structure")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user structure"})
			return
		}

		hp, err := generateHashPassword(u.Password)
		if err != nil {
			log.Printf("[Error] cannot create user")
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		_, err = queries.CreateUser(c, db.CreateUserParams{
			Name:     u.Name,
			Email:    u.Email,
			Password: hp,
			IsAdmin:  u.IsAdmin,
		})
		if err != nil {
			log.Printf("[Error] cannot create user")
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "user created"})
	}
}
