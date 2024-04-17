package api

import (
	"crypto/sha512"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/izzet-mtg/storage/services/backend/db"
)

type User struct {
	Name     string `validate:"required"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
}

func hashPassword(p string) (string, error) {
	s := sha512.Sum512([]byte(p))
	hp, err := bcrypt.GenerateFromPassword(s[:], 10)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hp), nil
}

func CreateUser(p *pgxpool.Pool) gin.HandlerFunc {
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

		hp, err := hashPassword(u.Password)
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
