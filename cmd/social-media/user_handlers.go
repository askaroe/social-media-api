package main

import (
	"net/http"
	"strconv"

	"github.com/askaroe/social-media-api/pkg/social-media/model"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ProfilePhoto string `json:"profilePhoto"`
		Name         string `json:"name"`
		Username     string `json:"username"`
		Description  string `json:"description"`
		Email        string `json:"email"`
		Password     string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 8)

	user := &model.User{
		ProfilePhoto: input.ProfilePhoto,
		Name:         input.Name,
		Username:     input.Username,
		Description:  input.Description,
		Email:        input.Email,
		Password:     string(hashed),
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJson(w, http.StatusCreated, user)
}

func (app *application) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := app.models.Users.GetAll()
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}
	app.respondWithJson(w, http.StatusOK, users)
}

func (app *application) getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["userId"]

	id, err := strconv.Atoi(param)

	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid user Id")
		return
	}

	user, err := app.models.Users.GetById(id)

	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}
	app.respondWithJson(w, http.StatusOK, user)
}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["userId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid menu ID")
		return
	}

	user, err := app.models.Users.GetById(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		ProfilePhoto *string `json:"profilePhoto"`
		Name         *string `json:"name"`
		Username     *string `json:"username"`
		Description  *string `json:"description"`
		Email        *string `json:"email"`
		Password     *string `json:"password"`
	}

	err = app.readJSON(w, r, &input)

	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.ProfilePhoto != nil {
		user.ProfilePhoto = *input.ProfilePhoto
	}
	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.Username != nil {
		user.Username = *input.Username
	}
	if input.Description != nil {
		user.Description = *input.Description
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Password != nil {
		user.Password = *input.Password
	}

	err = app.models.Users.Update(user)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
	}

	app.respondWithJson(w, http.StatusOK, user)
}

func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	params := vars["userId"]

	id, err := strconv.Atoi(params)

	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid User Id")
		return
	}

	err = app.models.Users.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
	}

	app.respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}
