// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import ()

type User struct {
	ID       int64
	Name     string
	Password string
	Email    string
	IsAdmin  bool
}
