package models

import "time"

type Transaction struct {
	ID       int       `json:"id"`
	Amount   float64   `json:"amount"`
	Type     string    `json:"type"`
	Category string    `json:"category"`
	Date     time.Time `json:"date"`
}

type AddTransactionRequest struct {
	Amount   float64 `json:"amount"`
	Type     string  `json:"type"`
	Category string  `json:"category"`
}

type GetSummaryResponse struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	NetBalance   float64 `json:"net_balance"`
}
