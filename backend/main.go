package main

import (
	"log"
	"net/http"

	"sezzle-calculator/backend/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /add", handlers.Add)
	mux.HandleFunc("POST /subtract", handlers.Subtract)
	mux.HandleFunc("POST /multiply", handlers.Multiply)
	mux.HandleFunc("POST /divide", handlers.Divide)
	mux.HandleFunc("POST /sqrt", handlers.Sqrt)
	mux.HandleFunc("POST /power", handlers.Power)
	mux.HandleFunc("POST /percentage", handlers.Percentage)

	handler := handlers.CORSMiddleware(mux)

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
