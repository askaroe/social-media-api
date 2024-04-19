package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/askaroe/social-media-api/pkg/jsonlog"
	"github.com/askaroe/social-media-api/pkg/social-media/model"
	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	models model.Models
	logger *jsonlog.Logger
}

func main() {
	fmt.Println("Started server")
	var cfg config
	flag.StringVar(&cfg.port, "port", ":8081", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:admin@localhost/social_media?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	// Connect to DB
	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	app := &application{
		config: cfg,
		models: model.NewModels(db),
	}

	app.run()
}

func (app *application) run() {
	fmt.Println("Running")
	r := mux.NewRouter()

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

	log.Printf("Starting server on %s\n", app.config.port)
	err := http.ListenAndServe(app.config.port, r)
	log.Fatal(err)
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
