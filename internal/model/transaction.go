package model

import "time"

// 1. Define a custom type for the enum
type TrxType int

// 2. Declare the values using iota within a const block
const (
	Undefined TrxType = iota // 0
	Debit                    // 1
	Credit                   // 2
)

type Transaction struct {
	TrxID           string    `json:"trxID"`
	Amount          float64   `json:"amount"`
	Type            TrxType   `json:"type"`
	TransactionTime time.Time `json:"transactionTime"`
}
