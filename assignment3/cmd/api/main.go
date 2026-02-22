package main

import (
	"log"
	"net/http"

	"golang/internal/app"
	httpHandler "golang/internal/delivery/http"
	"golang/internal/repository/_postgres"

	"github.com/gorilla/mux"
)

func main() {

	dsn := "postgres://robert@localhost/assignment3?sslmode=disable"
	db, err := _postgres.NewPostgresDB(dsn)
	if err != nil {
		log.Fatal(err)
	}

	app := app.NewApp(db)

	r := mux.NewRouter()

	r.Use(httpHandler.LoggingMiddleware)
	r.Use(httpHandler.AuthMiddleware)

	r.HandleFunc("/health", HealthCheck).Methods("GET")
	r.HandleFunc("/users", app.Handler.GetUsers).Methods("GET")
	r.HandleFunc("/users/{id}", app.Handler.GetUserByID).Methods("GET")
	r.HandleFunc("/users", app.Handler.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", app.Handler.UpdateUser).Methods("PATCH")
	r.HandleFunc("/users/{id}", app.Handler.DeleteUser).Methods("DELETE")

	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", r)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
