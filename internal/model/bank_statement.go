package model

type BankStatement struct {
	UniqueID string  `json:"uniqueID"`
	BankName string  `json:"bankName"`
	Amount   float64 `json:"amount"` // for debit will be nagative
	Date     string  `json:"date"`
}
