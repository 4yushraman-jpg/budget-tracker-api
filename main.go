package main

import (
	"budget-tracker/database"
	"budget-tracker/handlers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	database.InitDB()

	http.HandleFunc("POST /budget", handlers.AddTransactionHandler)
	http.HandleFunc("GET /budget", handlers.GetSummaryHandler)

	fmt.Println("Server starting on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
