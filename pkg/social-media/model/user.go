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
	Age          string `json:"age"`
}
type UserModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (u UserModel) Insert(user *User) error {
	query := `
			INSERT INTO users (profilePhoto, name, username, description, email, password, age)
			VALUES($1, $2, $3, $4, $5, $6, $7)
			RETURNING id, createdAt, updatedAt
			`
	args := []interface{}{user.ProfilePhoto, user.Name, user.Username, user.Description, user.Email, user.Password, user.Age}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	return u.DB.QueryRowContext(ctx, query, args...).Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt)
}

func (u UserModel) GetAll(name string, filters Filters) ([]*User, error) {
	query := `
		SELECT id, createdAt, updatedAt, profilePhoto, name, username, description, email, password, age
		FROM users
		ORDER BY id
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := u.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []*User{}

	for rows.Next() {
		var user User

		err := rows.Scan(
			&user.Id,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Name,
			&user.ProfilePhoto,
			&user.Username,
			&user.Description,
			&user.Email,
			&user.Password,
			&user.Age,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, err
}

func (u UserModel) GetById(id int) (*User, error) {
	query := `
		SELECT id, createdAt, updatedAt, profilePhoto, name, username, description, email, password, age
		FROM users
		WHERE id = $1
		`
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := u.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt, &user.ProfilePhoto,
		&user.Name, &user.Username, &user.Description, &user.Email, &user.Password, &user.Age)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u UserModel) Update(user *User) error {
	query := `
		UPDATE users
		SET profilePhoto = $1, name = $2, username = $3, description = $4, email = $5, password = $6, age = $7
		WHERE id = $8
		RETURNING updatedAt
		`

	args := []interface{}{user.ProfilePhoto, user.Name, user.Username, user.Description, user.Email, user.Password, user.Age, user.Id}
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
