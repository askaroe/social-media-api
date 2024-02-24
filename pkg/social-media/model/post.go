package model

import (
	"database/sql"
	"log"
)

type Post struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Image     string `json:"image"`
	Caption   string `json:"caption"`
	UserId    string `json:"UserId"`
}
type PostModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}
