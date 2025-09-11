package main

import (
	"chirpstream/internal/chirps"
	"chirpstream/pkg/config"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/julienschmidt/httprouter"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg, err := config.Init(".")
	if err != nil {
		log.Fatalf("Cant parse config file %v\n", err)
	}

	pool, err := NewDBPool(cfg.Databases.Chirps)
	if err != nil {
		log.Fatalf("Failed to create database pool: %v\n", err)
	}

	defer CloseDB(pool)

	// --- Migration stuff ---
	err = gooseMigrations(pool)
	if err != nil {
		log.Fatalf("Failed to run migration: %v\n", err)
	}

	// --- repository stuff ---
	repo := chirps.NewRepo(pool)
	service := chirps.NewService(repo)
	handler := chirps.NewHandler(service)

	mux := httprouter.New()
	mux.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})

	handler.RegisterRoutes(mux)

	log.Println("Server started on port :3220")
	errr := http.ListenAndServe("localhost:3220", mux)
	if errr != nil {
		log.Fatalf("Error starting server %v\n", errr)
	}
}

func NewDBPool(conf config.DBConfig) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(conf.DBConnstring())
	if err != nil {
		log.Fatalf("Failed to create initial pgxpool config: %v\n", err)
		os.Exit(1)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Database Connection pool initialized successfully")

	return pool, err
}

// GetStdlibDB returns a standard *sql.DB connection pool, needed for goose.
// Imp: close this DB handle when done.
func GetStdlibDB(pool *pgxpool.Pool) (*sql.DB, error) {
	if pool == nil {
		return nil, fmt.Errorf("database pool not initialized")
	}
	// OpenDBFromPool creates a standard sql.DB wrapper around the pgxpool
	stdlibDB := stdlib.OpenDBFromPool(pool)

	// Verify the connection
	if err := stdlibDB.Ping(); err != nil {
		return nil, fmt.Errorf("unable to ping database via stdlib: %w", err)
	}
	return stdlibDB, nil
}

func CloseDB(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
		fmt.Println("Database connection pool closed.")
	}
}

func gooseMigrations(pool *pgxpool.Pool) error {

	migrationDir := "internal/data/migrations"

	fmt.Println("Running db migrations...")

	db, err := GetStdlibDB(pool)
	if err != nil {
		return fmt.Errorf("Failed to get standard DB connection for migration %v\n", err)
	}

	//get the handle and then defer closing it
	defer func(db *sql.DB) error {
		err := db.Close()
		if err != nil {
			return fmt.Errorf("WARN: Error closing migration DB handle: %v\n", err)
		}

		return err
	}(db)

	error := goose.Up(db, migrationDir)
	if error != nil {
		return fmt.Errorf("Failed to run migrations %v\n", error)
	}

	return err
}
