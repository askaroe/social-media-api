package model

import (
	"context"
	"database/sql"
	"log"
	"time"
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

func (c CommentModel) Insert(comment *Comment) error {
	query := `
			INSERT INTO comments (message, userId, postId)
			VALUES($1, $2, $3)
			RETURNING id, createdAt, updatedAt
			`
	args := []interface{}{comment.Message, comment.UserId, comment.PostId}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	return c.DB.QueryRowContext(ctx, query, args...).Scan(&comment.Id, &comment.CreatedAt, &comment.UpdatedAt)
}

func (c CommentModel) Get(id int) (*Comment, error) {
	query := `
		SELECT id, createdAt, updatedAt, message, userId, postId
		FROM comments
		WHERE id = $1
		`
	var comment Comment
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := c.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&comment.Id, &comment.CreatedAt, &comment.UpdatedAt, &comment.Message,
		&comment.UserId, &comment.PostId)

	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (c CommentModel) Update(comment *Comment) error {
	query := `
		UPDATE comments
		SET message = $1
		WHERE id = $2
		RETURNING updatedAt
		`

	args := []interface{}{comment.Message, comment.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.DB.QueryRowContext(ctx, query, args...).Scan(&comment.UpdatedAt)
}

func (c CommentModel) Delete(id int) error {
	query := `
		DELETE FROM comments
		WHERE id = $1
		`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	_, err := c.DB.ExecContext(ctx, query, id)

	return err
}
