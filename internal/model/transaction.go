package model

import (
	"strings"
	"time"
)

// 1. Define a custom type for the enum
type TrxType int

// 2. Declare the values using iota within a const block
const (
	Undefined TrxType = iota // 0
	Debit                    // 1
	Credit                   // 2
)

func DefineTrxType(param string) TrxType {
	switch strings.ToUpper(param) {
	case "DEBIT":
		return Debit
	case "CREDIT":
		return Credit
	default:
		return Undefined
	}
}

func TrxTypeToStr(param TrxType) string {
	switch param {
	case Debit:
		return "DEBIT"
	case Credit:
		return "CREDIT"
	default:
		return "UNDEFINED"
	}
}

type Transaction struct {
	TrxID           string    `json:"trxID"`
	Amount          float64   `json:"amount"`
	Type            TrxType   `json:"type"`
	TransactionTime time.Time `json:"transactionTime"`
}
