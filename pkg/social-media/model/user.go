package model

import (
	"context"
	"database/sql"
	"fmt"
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

func (u UserModel) GetAll(username string, age string, filters Filters) ([]*User, Metadata, error) {
	if age == "" {
		age = "0"
	}
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, createdAt, updatedAt, profilePhoto, name, username, description, email, password, age
		FROM users
		WHERE (LOWER(username) = LOWER($1) OR $1 = '')
		AND (age >= $2 OR $2 = 0)
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4
	`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{username, age, filters.limit(), filters.offset()}

	rows, err := u.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	users := []*User{}

	for rows.Next() {
		var user User

		err := rows.Scan(
			&totalRecords,
			&user.Id,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.ProfilePhoto,
			&user.Name,
			&user.Username,
			&user.Description,
			&user.Email,
			&user.Password,
			&user.Age,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return users, metadata, nil // Return users and nil error
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
