package model

import (
	"database/sql"
	"log"
)

type User struct {
	Id           string `json:"id"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
	ProfilePhoto string `json:"profilePhoto"`
	Name         string `json:"name"`
	Username     string `json:"username"`
	Description  string `json:"description"`
	Email        string `json:"email"`
	Password     string `json:"password"`
}
type UserModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}
