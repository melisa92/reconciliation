package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/melisa92/reconciliation/internal/interfaces"
	"github.com/melisa92/reconciliation/internal/model"
)

type ReconciliationUsecase struct {
	BankStatementRepo interfaces.BankStatementInterface
	TransactionRepo   interfaces.TransactionStatementInterface
}

func NewTransactionUsecase(
	bankStatementRepo interfaces.BankStatementInterface,
	transactionRepo interfaces.TransactionStatementInterface,
) *ReconciliationUsecase {
	return &ReconciliationUsecase{
		BankStatementRepo: bankStatementRepo,
		TransactionRepo:   transactionRepo,
	}
}

func (u *ReconciliationUsecase) ReconciliationProcess(ctx context.Context, startDate, endDate string) (*model.ReconciliationSummary, error) {
	transactionData, err := u.TransactionRepo.GetTransaction(ctx)
	if err != nil {
		return nil, err
	}

	bankStatementData, err := u.BankStatementRepo.GetBankStatement(ctx)
	if err != nil {
		return nil, err
	}

	var totalMatch, totalTransactionProceed, totalUnmatchBankStatement int
	var totalDiscrepancies float64
	var unmatchTransaction []model.Transaction
	unmatchBankStatement := make(map[string][]model.BankStatement)

	mapBankStatement := make(map[string][]model.BankStatement)

	startTime, err := time.Parse("2006-01-02T15:04:05Z", fmt.Sprint(startDate, "T00:00:00Z"))
	if err != nil {
		return nil, errors.New("start time is not valid (expected format: 2006-01-02)")
	}
	endTime, err := time.Parse("2006-01-02T15:04:05Z", fmt.Sprint(endDate, "T23:59:59Z"))
	if err != nil {
		return nil, errors.New("end time is not valid (expected format: 2006-01-02)")
	}
	if startTime.After(endTime) {
		return nil, errors.New("start time cannot be grater than end time")
	}

	// make map for bankstatement data
	for _, v := range bankStatementData {
		statementTime, _ := time.Parse("2006-01-02", v.Date)
		if statementTime.Before(startTime) || statementTime.After(endTime) {
			continue
		}

		trxType := model.DefineTrxType("CREDIT")
		if v.Amount < 0 {
			trxType = model.DefineTrxType("DEBIT")
		}
		key := fmt.Sprintf("%d~%s", trxType, v.Date)
		mapBankStatement[key] = append(mapBankStatement[key], *v)
	}

	for _, v := range transactionData {
		// skip if transaction type is undefined
		if v.Type == model.Undefined || v.TransactionTime.Before(startTime) || v.TransactionTime.After(endTime) {
			continue
		}

		key := fmt.Sprintf("%d~%s", v.Type, v.TransactionTime.Format("2006-01-02"))
		if dataMapBank, ok := mapBankStatement[key]; ok {
			if len(dataMapBank) == 0 {
				unmatchTransaction = append(unmatchTransaction, *v)
				totalTransactionProceed++
				continue
			}
			var totalDiscr float64
			var idxMinDicr int
			for k2 := range dataMapBank {
				d := math.Abs(v.Amount - math.Abs(dataMapBank[k2].Amount))
				if k2 == 0 {
					totalDiscr = d
					continue
				}
				if d < totalDiscr {
					totalDiscr = d
					idxMinDicr = k2
				}
			}
			totalMatch++
			totalDiscrepancies += totalDiscr

			if idxMinDicr >= 0 {
				dataMapBank = append(dataMapBank[:idxMinDicr], dataMapBank[idxMinDicr+1:]...)
			}
			mapBankStatement[key] = dataMapBank
		} else {
			unmatchTransaction = append(unmatchTransaction, *v)
		}
		totalTransactionProceed++
	}

	for _, v := range mapBankStatement {
		if len(v) == 0 {
			continue
		}
		for _, b := range v {
			totalUnmatchBankStatement++
			unmatchBankStatement[b.BankName] = append(unmatchBankStatement[b.BankName], b)
		}
	}
	return &model.ReconciliationSummary{
		TotalTrxRecords:                    len(transactionData),
		TotalProceesRecords:                totalTransactionProceed,
		TotalMatchedTrxRecords:             totalMatch,
		TotalUnmatchedTransactionRecords:   len(unmatchTransaction),
		TotalUnmatchedBankStatementRecords: totalUnmatchBankStatement,
		ListUnmatchTrx:                     unmatchTransaction,
		ListUnmatchBankStatement:           unmatchBankStatement,
		TotalDiscrepancies:                 totalDiscrepancies,
	}, nil
}
