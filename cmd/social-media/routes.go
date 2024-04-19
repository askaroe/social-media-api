package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// routes is our main application's router.
func (app *application) routes() http.Handler {
	r := mux.NewRouter()
	// Convert the app.notFoundResponse helper to a http.Handler using the http.HandlerFunc()
	// adapter, and then set it as the custom error handler for 404 Not Found responses.
	r.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)

	// Convert app.methodNotAllowedResponse helper to a http.Handler and set it as the custom
	// error handler for 405 Method Not Allowed responses
	r.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	// healthcheck
	r.HandleFunc("/api/v1/healthcheck", app.healthcheckHandler).Methods("GET")

	v1 := r.PathPrefix("/api/v1").Subrouter()

	// Users
	v1.HandleFunc("/register", app.createUserHandler).Methods("POST")
	v1.HandleFunc("/users", app.getAllUsersHandler).Methods("GET")
	v1.HandleFunc("/users/{userId:[0-9]+}", app.getUserByIdHandler).Methods("GET")
	v1.HandleFunc("/users/{userId:[0-9]+}", app.updateUserHandler).Methods("PUT")
	v1.HandleFunc("/users/{userId:[0-9]+}", app.deleteUserHandler).Methods("DELETE")

	// Posts
	v1.HandleFunc("/posts", app.createPostHandler).Methods("POST")
	v1.HandleFunc("/posts/{postId:[0-9]+}", app.getPostHandler).Methods("GET")
	v1.HandleFunc("/posts/{postId:[0-9]+}", app.updatePostHandler).Methods("PUT")
	v1.HandleFunc("/posts/{postId:[0-9]+}", app.deletePostHandler).Methods("DELETE")

	// Comments
	v1.HandleFunc("/comments", app.createCommentHandler).Methods("POST")
	v1.HandleFunc("/comments/{commentId:[0-9]+}", app.getCommentHandler).Methods("GET")
	v1.HandleFunc("/comments/{commentId:[0-9]+}", app.updateCommentHandler).Methods("PUT")
	v1.HandleFunc("/comments/{commentId:[0-9]+}", app.deleteCommentHandler).Methods("DELETE")

	// Members
	v1.HandleFunc("/members", app.registerMemberHandler).Methods("POST")
	v1.HandleFunc("/members/activated", app.activateMemberHandler).Methods("PUT")
	v1.HandleFunc("/tokens/authentication", app.createAuthenticationTokenHandler).Methods("POST")

	// Wrap the router with the panic recovery middleware and rate limit middleware.
	return app.authenticate(r)
}
