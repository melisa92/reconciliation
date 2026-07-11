package model

type ReconciliationSummary struct {
	TotalTrxRecords                    int                        `json:"total_trx_records"`
	TotalProceesRecords                int                        `json:"total_processed_records"`
	TotalMatchedTrxRecords             int                        `json:"total_matched_trx_records"`
	TotalUnmatchedTransactionRecords   int                        `json:"total_unmatched_trx"`
	TotalUnmatchedBankStatementRecords int                        `json:"total_unmatched_bank_statement"`
	ListUnmatchTrx                     []Transaction              `json:"list_unmatch_trx"`
	ListUnmatchBankStatement           map[string][]BankStatement `json:"list_unmatch_bank_statement"`
	TotalDiscrepancies                 float64                    `json:"total_descrepancies"`
}
