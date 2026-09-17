package main

import (
	"log"
	"net/http"
	"os"

	"sezzle-calculator/backend/handlers"
)

func resolvePort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return "8080"
}

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

	addr := ":" + resolvePort()
	log.Println("listening on", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}
