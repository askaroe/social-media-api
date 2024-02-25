package main

import (
	"github.com/askaroe/social-media-api/pkg/social-media/model"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Image   string `json:"image"`
		Caption string `json:"caption"`
		UserId  string `json:"userId"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
	}

	post := &model.Post{
		Image:   input.Image,
		Caption: input.Caption,
		UserId:  input.UserId,
	}

	err = app.models.Posts.Insert(post)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJson(w, http.StatusCreated, post)
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["postId"]

	id, err := strconv.Atoi(param)

	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid post Id")
		return
	}

	post, err := app.models.Posts.Get(id)

	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}
	app.respondWithJson(w, http.StatusOK, post)
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["postId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	post, err := app.models.Posts.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		Image   *string `json:"image"`
		Caption *string `json:"caption"`
	}

	err = app.readJSON(w, r, &input)

	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Image != nil {
		post.Image = *input.Image
	}
	if input.Caption != nil {
		post.Caption = *input.Caption
	}

	err = app.models.Posts.Update(post)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
	}

	app.respondWithJson(w, http.StatusOK, post)
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	params := vars["postId"]

	id, err := strconv.Atoi(params)

	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Post Id")
		return
	}

	err = app.models.Posts.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
	}

	app.respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}
