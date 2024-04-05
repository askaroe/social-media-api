package main

import (
	"net/http"
	"strconv"

	"github.com/askaroe/social-media-api/pkg/social-media/model"
	"github.com/askaroe/social-media-api/pkg/social-media/validator"
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
		Age          string `json:"age"`
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
		Age:          input.Age,
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJson(w, http.StatusCreated, user)
}

func (app *application) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string
		model.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{
		"id", "name", "username",
		"-id", "-name", "-username",
	}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	users, err := app.models.Users.GetAll(input.Name, input.Filters)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"users": users}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
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
		Age          *string `json:"age"`
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
	if input.Age != nil {
		user.Age = *input.Age
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
