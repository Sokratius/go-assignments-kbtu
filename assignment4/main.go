package main

import (
	"log"
	"net/http"

	"assignment4/db"
	"assignment4/handlers"
	"assignment4/repository"
)

func main() {
	dbConn, err := db.Connect()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	repo := repository.NewRepository(dbConn)
	handler := handlers.NewUserHandler(repo)

	http.HandleFunc("/users", handler.GetUsers)
	http.HandleFunc("/users/common-friends", handler.GetCommonFriends)

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
