package model

import (
	"context"
	"database/sql"
	"log"
	"time"
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

func (u UserModel) Insert(user *User) error {
	query := `
			INSERT INTO users (profilePhoto, name, username, description, email, password)
			VALUES($1, $2, $3, $4, $5, $6)
			RETURNING id, createdAt, updatedAt
			`
	args := []interface{}{user.ProfilePhoto, user.Name, user.Username, user.Description, user.Email, user.Password}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	return u.DB.QueryRowContext(ctx, query, args...).Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt)
}

func (u UserModel) Get(id int) (*User, error) {
	query := `
		SELECT id, createdAt, updatedAt, profilePhoto, name, username, description, email, password
		FROM users
		WHERE id = $1
		`
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := u.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt, &user.ProfilePhoto,
		&user.Name, &user.Username, &user.Description, &user.Email, &user.Password)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u UserModel) Update(user *User) error {
	query := `
		UPDATE users
		SET profilePhoto = $1, name = $2, username = $3, description = $4, email = $5, password = $6
		WHERE id = $7
		RETURNING updatedAt
		`

	args := []interface{}{user.ProfilePhoto, user.Name, user.Username, user.Description, user.Email, user.Password, user.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return u.DB.QueryRowContext(ctx, query, args...).Scan(&user.UpdatedAt)
}

func (u UserModel) Delete(id int) error {
	query := `
		DELETE FROM users
		WHERE id = $1
		`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	_, err := u.DB.ExecContext(ctx, query, id)

	return err
}
