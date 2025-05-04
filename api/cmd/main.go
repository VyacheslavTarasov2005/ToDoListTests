package main

import (
	"HITS_ToDoList_Tests/internal/infrastructure/db"
	"log"
)

func main() {
	dbConn, err := db.NewPostgresConnection("localhost", "postgres", "123456", "ToDoDb",
		"5432")
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}

	if err = db.Migrate(dbConn); err != nil {
		log.Fatalf("Failed to migrate db: %v", err)
	}

	log.Println("Application started")
}
