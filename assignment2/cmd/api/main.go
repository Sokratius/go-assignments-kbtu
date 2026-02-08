package main

import (
	"log"
	"net/http"

	"assignment2/internal/handlers"
	"assignment2/internal/middleware"
)

func main() {
	taskHandler := handlers.NewTaskHandler()

	mux := http.NewServeMux()

	mux.Handle("/tasks", middleware.APIKey(
		middleware.Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				taskHandler.GetTasks(w, r)
			case http.MethodPost:
				taskHandler.CreateTask(w, r)
			case http.MethodPatch:
				taskHandler.UpdateTask(w, r)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		})),
	))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
