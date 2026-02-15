package handlers

import (
	"budget-tracker/database"
	"budget-tracker/models"
	"encoding/json"
	"net/http"
)

var ValidTypes = map[string]bool{
	"income":  true,
	"expense": true,
}

var ValidCategories = map[string]bool{
	"food":          true,
	"rent":          true,
	"salary":        true,
	"entertainment": true,
	"utilities":     true,
}

func AddTransactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	var req models.AddTransactionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "Amount must be more than zero", http.StatusBadRequest)
		return
	}

	if !ValidTypes[req.Type] || !ValidCategories[req.Category] {
		http.Error(w, "Invalid type or category or both", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO transactions (amount, type, category) VALUES (?, ?, ?)`
	_, err = database.DB.Exec(query, req.Amount, req.Type, req.Category)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Transaction added",
	})
}

func GetSummaryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	query := `SELECT type, SUM(amount) FROM transactions GROUP BY type`
	rows, err := database.DB.Query(query)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var resp models.GetSummaryResponse
	for rows.Next() {
		var t string
		var amount float64

		if err := rows.Scan(&t, &amount); err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			return
		}

		switch t {
		case "income":
			resp.TotalIncome = amount
		case "expense":
			resp.TotalExpense = amount
		}
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Database iteration error", http.StatusInternalServerError)
		return
	}

	resp.NetBalance = resp.TotalIncome - resp.TotalExpense

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
