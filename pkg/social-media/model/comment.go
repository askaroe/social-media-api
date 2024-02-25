package model

import (
	"database/sql"
	"log"
)

type Comment struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Message   string `json:"text"`
	UserId    string `json:"userId"`
	PostId    string `json:"PostId"`
}

type CommentModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}
