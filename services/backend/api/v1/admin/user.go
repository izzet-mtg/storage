package admin

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

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

		cu, err := queries.CreateUser(c, db.CreateUserParams{
			Name:     u.Name,
			Email:    u.Email,
			Password: u.Password,
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
