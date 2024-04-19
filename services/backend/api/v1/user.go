package api

import (
	"context"
	"crypto/sha512"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	"github.com/izzet-mtg/storage/services/backend/db"
	"github.com/izzet-mtg/storage/services/backend/user"
)

func generateHashPassword(p string) (string, error) {
	s := sha512.Sum512([]byte(p))
	hp, err := bcrypt.GenerateFromPassword(s[:], 10)
	if err != nil {
		return "", err
	}
	return string(hp), nil
}

func isSamePassword(p, hp string) bool {
	s := sha512.Sum512([]byte(p))
	err := bcrypt.CompareHashAndPassword([]byte(hp), s[:])
	return err == nil
}

type LoginingUser struct {
	NameOrEmail string `validate:"required"`
	RawPassword string `validate:"required" json:"password"`
}

var ErrNoSuchUser = errors.New("no such user")
var ErrInvalidLoginUser = errors.New("challenged login by invalid user")

func login(ctx context.Context, p *pgxpool.Pool, rc *redis.Client, lu LoginingUser, exp time.Duration) (string, error) {
	queries := db.New(p)

	u, err := queries.GetUser(ctx, lu.NameOrEmail)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return "", ErrNoSuchUser
	case err != nil:
		return "", err
	}
	if !isSamePassword(lu.RawPassword, u.Password) {
		return "", ErrInvalidLoginUser
	}

	return user.Login(ctx, rc, u.ID, exp)
}

type User struct {
	Name     string `validate:"required"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
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
		})
		if err != nil {
			log.Printf("[Error] cannot create user")
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		lu := LoginingUser{
			NameOrEmail: u.Name,
			RawPassword: u.Password,
		}
		si, err := login(c, p, rc, lu, exp)
		if err != nil {
			log.Println("[Error] cannot loging user")
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot login user"})
			return
		}
		c.Header("Authorization", fmt.Sprintf("Bearer %s", si))
		c.JSON(http.StatusOK, gin.H{"message": "user created"})
	}
}

func Login(p *pgxpool.Pool, rc *redis.Client, exp time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		var lu LoginingUser
		if err := c.ShouldBindJSON(&lu); err != nil {
			log.Println("[Info] reject user request with invalid structure")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user structure"})
			return
		}
		validate := validator.New(validator.WithRequiredStructEnabled())
		if err := validate.Struct(lu); err != nil {
			log.Println("[Info] reject user request with invalid structure")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user structure"})
			return
		}

		si, err := login(c, p, rc, lu, exp)
		if err != nil {
			log.Println("[Error] cannot loging user")
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot login user"})
			return
		}
		c.Header("Authorization", fmt.Sprintf("Bearer %s", si))
		c.JSON(http.StatusOK, gin.H{"message": "logined"})
	}
}

var bearerTokenRegexp = regexp.MustCompile("^\\s*Bearer\\s*")

func Logout(rc *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := bearerTokenRegexp.ReplaceAllString(c.GetHeader("Authorization"), "")
		fmt.Printf("Authorization = %v\n", t)
		if err := user.Logout(c, rc, t); err != nil {
			log.Println("[Error] cannot delete session id")
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "logout"})
	}
}
