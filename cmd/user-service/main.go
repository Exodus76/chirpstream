package main

import (
	"chirpstream/internal/user"
	"chirpstream/pkg/config"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("Cant parse config file %v\n", err)
	}

	DBPoolConf(cfg.Databases.Users)

	// --- Migration stuff ---

	migrationDir := "internal/data/migrations"

	fmt.Println("Running db migrations...")

	db, err := GetStdlibDB()
	if err != nil {
		log.Fatalf("Failed to get standard DB connection for migration %v\n", err)
	}

	//get the handle and then defer closing it
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("WARN: Error closing migration DB handle: %v\n", err)
		}
	}(db)

	error := goose.Up(db, migrationDir)
	if error != nil {
		log.Fatalf("Failed to run migrations %v\n", error)
	}

	// --- repository stuff ---
	repo := user.NewRepo(pool)
	service := user.NewService(repo)
	handler := user.NewHandler(service)

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

var pool *pgxpool.Pool

func DBPoolConf(dconf config.DBConfig) {

	config, err := pgxpool.ParseConfig(dconf.DBConnstring())

	pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Database Connection pool initialized successfully")
}

// GetStdlibDB returns a standard *sql.DB connection pool, needed for goose.
// Imp: close this DB handle when done.
func GetStdlibDB() (*sql.DB, error) {
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

func CloseDB() {
	if pool != nil {
		pool.Close()
		fmt.Println("Database connection pool closed.")
	}
}
