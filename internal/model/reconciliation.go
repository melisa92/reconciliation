package model

type ReconciliationSummary struct {
	TotalMatchedTrx    int     `json:"total_matched_trx"`
	TotalUnmatchedTrx  int     `json:"total_unmatched_trx"`
	TotalDiscrepancies float64 `json:"total_descrepancies"`
}
