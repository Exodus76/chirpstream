package main

import (
	"chirpstream/internal/database"
	"chirpstream/internal/userrelationship"
	"chirpstream/pkg/config"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	cfg, err := config.Init(".")
	if err != nil {
		log.Fatalf("Cant parse config file %v\n", err)
	}

	dbSession, err := database.InitScyllaDb(cfg.Scylladb.Host, cfg.Scylladb.Username, cfg.Scylladb.Password, cfg.Scylladb.RelationKeyspace)
	if err != nil {
		log.Fatalf("Failed to initialize ScyllaDB: %v\n", err)
	}
	defer dbSession.Close()

	//run database migration
	err = database.MigrateRelationship(&dbSession, cfg.Scylladb.RelationKeyspace)
	if err != nil {
		log.Fatalf("Failed to run database migration: %v\n", err)
	}

	repo := userrelationship.NewRepo(&dbSession)
	service := userrelationship.NewService(repo)
	handler := userrelationship.NewHandler(service)

	mux := httprouter.New()
	handler.RegisterRoutes(mux)

	log.Println("Server started on port :3222")
	err = http.ListenAndServe("localhost:3222", mux)
	if err != nil {
		log.Fatalf("Error starting server %v\n", err)
	}

}
