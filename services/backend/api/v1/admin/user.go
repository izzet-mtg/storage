package admin

import (
	"crypto/sha512"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/izzet-mtg/storage/services/backend/db"
)

type User struct {
	Name     string
	Email    string
	Password string
}

func CreateUser(p *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		queries := db.New(p)

		var u User
		if err := c.ShouldBindJSON(&u); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user structure"})
			return
		}

		s := sha512.Sum512([]byte(u.Password))
		hp, err := bcrypt.GenerateFromPassword(s[:], 10)
		if err != nil {
			log.Printf("[Error] cannot create user")
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		cu, err := queries.CreateUser(c, db.CreateUserParams{
			Name:     u.Name,
			Email:    u.Email,
			Password: hex.EncodeToString(hp),
		})
		if err != nil {
			log.Printf("[Error] cannot create user")
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		log.Printf("%v\n", cu)

		c.JSON(http.StatusOK, gin.H{"message": "user created"})
	}
}
