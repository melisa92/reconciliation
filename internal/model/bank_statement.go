package model

import "time"

type BankStatement struct {
	UniqueID string    `json`
	Amount   float64   `json:"amount"` // for debit will be nagative
	Date     time.Time `json:"date"`
}
