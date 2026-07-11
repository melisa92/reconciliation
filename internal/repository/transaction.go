package repository

import (
	"context"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/melisa92/reconciliation/internal/model"
)

type TransactiontRepo struct {
	FilePath string
}

func NewTransactionRepo(filepath string) *TransactiontRepo {
	return &TransactiontRepo{
		FilePath: filepath,
	}
}

func (r *TransactiontRepo) GetTransaction(ctx context.Context, startTime, endTime time.Time) ([]*model.Transaction, error) {
	file, err := os.Open(r.FilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	dt := []*model.Transaction{}
	reader := csv.NewReader(file)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		if row[0] == "trxID" {
			continue
		}

		amountStr := row[1]
		trxTimeStr := row[3]

		amount, _ := strconv.ParseFloat(amountStr, 64)
		trxTime, _ := time.ParseInLocation("2006-01-02T15:04:05Z", trxTimeStr, time.Local)

		dt = append(dt, &model.Transaction{
			TrxID:           row[0],
			Type:            model.DefineTrxType(row[2]),
			Amount:          amount,
			TransactionTime: trxTime,
		})
	}

	return dt, nil
}
