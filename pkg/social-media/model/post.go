package model

import (
	"context"
	"database/sql"
	"log"
	"time"
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

func (p PostModel) Insert(post *Post) error {
	query := `
			INSERT INTO posts (image, caption, userId)
			VALUES($1, $2, $3)
			RETURNING id, createdAt, updatedAt
			`
	args := []interface{}{post.Image, post.Caption, post.UserId}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	return p.DB.QueryRowContext(ctx, query, args...).Scan(&post.Id, &post.CreatedAt, &post.UpdatedAt)
}

func (p PostModel) Get(id int) (*Post, error) {
	query := `
		SELECT id, createdAt, updatedAt, image, caption, userId
		FROM posts
		WHERE id = $1
		`
	var post Post
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := p.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&post.Id, &post.CreatedAt, &post.UpdatedAt, &post.Image, &post.Caption, &post.UserId)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (p PostModel) Update(post *Post) error {
	query := `
		UPDATE posts
		SET image = $1, caption = $2
		WHERE id = $3
		RETURNING updatedAt
		`

	args := []interface{}{post.Image, post.Caption, post.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return p.DB.QueryRowContext(ctx, query, args...).Scan(&post.UpdatedAt)
}

func (p PostModel) Delete(id int) error {
	query := `
		DELETE FROM posts
		WHERE id = $1
		`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	_, err := p.DB.ExecContext(ctx, query, id)

	return err
}
