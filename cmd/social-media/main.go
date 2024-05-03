package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/askaroe/social-media-api/pkg/jsonlog"
	"github.com/askaroe/social-media-api/pkg/social-media/model"
	"github.com/askaroe/social-media-api/pkg/social-media/model/filler"
	"github.com/askaroe/social-media-api/pkg/vcs"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/peterbourgon/ff/v3"

	_ "github.com/lib/pq"
)

var (
	version = vcs.Version()
)

type config struct {
	port       int
	env        string
	fill       bool
	migrations string
	db         struct {
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
	// var cfg config
	// flag.StringVar(&cfg.port, "port", ":8081", "API server port")
	// flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	// flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:admin@localhost/social_media?sslmode=disable", "PostgreSQL DSN")
	// flag.Parse()
	fs := flag.NewFlagSet("demo-app", flag.ContinueOnError)

	var (
		cfg        config
		fill       = fs.Bool("fill", false, "Fill database with dummy data")
		migrations = fs.String("migrations", "", "Path to migration files folder. If not provided, migrations do not applied")
		port       = fs.Int("port", 8060, "API Server port")
		env        = fs.String("env", "development", "Environment (development|staging|production)")
		dbDsn      = fs.String("dsn", "postgres://postgres:admin@db:5432/social_media?sslmode=disable", "PostgreSQL DSN")
	)

	// Init logger
	logger := jsonlog.NewLogger(os.Stdout, jsonlog.LevelInfo)

	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVars()); err != nil {
		logger.PrintFatal(err, nil)
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}

	cfg.port = *port
	cfg.env = *env
	cfg.fill = *fill
	cfg.db.dsn = *dbDsn
	cfg.migrations = *migrations

	logger.PrintInfo("starting application with configuration", map[string]string{
		"port":       fmt.Sprintf("%d", cfg.port),
		"fill":       fmt.Sprintf("%t", cfg.fill),
		"env":        cfg.env,
		"db":         cfg.db.dsn,
		"migrations": cfg.migrations,
	})

	// Connect to DB
	db, err := openDB(cfg)
	if err != nil {
		logger.PrintError(err, nil)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.PrintFatal(err, nil)
		}
	}()

	app := &application{
		config: cfg,
		models: model.NewModels(db),
		logger: logger,
	}

	if cfg.fill {
		err = filler.PopulateDatabase(app.models)
		if err != nil {
			logger.PrintFatal(err, nil)
			return
		}
	}

	if err := app.serve(); err != nil {
		logger.PrintFatal(err, nil)
	}

}

func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config // struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// https://github.com/golang-migrate/migrate?tab=readme-ov-file#use-in-your-go-project
	if cfg.migrations != "" {
		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			return nil, err
		}
		m, err := migrate.NewWithDatabaseInstance(
			cfg.migrations,
			"postgres", driver)
		if err != nil {
			return nil, err
		}
		m.Up()
	}

	return db, nil
}
