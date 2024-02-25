package main

import (
	"github.com/askaroe/social-media-api/pkg/social-media/model"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Message string `json:"message"`
		UserId  string `json:"userId"`
		PostId  string `json:"postId"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
	}

	comment := &model.Comment{
		Message: input.Message,
		UserId:  input.UserId,
		PostId:  input.PostId,
	}

	err = app.models.Comments.Insert(comment)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJson(w, http.StatusCreated, comment)
}

func (app *application) getCommentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["commentId"]

	id, err := strconv.Atoi(param)

	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid comment Id")
		return
	}

	comment, err := app.models.Comments.Get(id)

	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}
	app.respondWithJson(w, http.StatusOK, comment)
}

func (app *application) updateCommentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["commentId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	comment, err := app.models.Comments.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		Message *string `json:"message"`
	}

	err = app.readJSON(w, r, &input)

	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Message != nil {
		comment.Message = *input.Message
	}

	err = app.models.Comments.Update(comment)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
	}

	app.respondWithJson(w, http.StatusOK, comment)
}

func (app *application) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	params := vars["commentId"]

	id, err := strconv.Atoi(params)

	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Comment Id")
		return
	}

	err = app.models.Comments.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
	}

	app.respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}
