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

func (r *TransactiontRepo) GetTransaction(ctx context.Context) ([]*model.Transaction, error) {
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

		trxTime, _ := time.ParseInLocation("2006-01-02T15:04:05Z", row[3], time.Local)

		amountStr := row[1]
		amount, _ := strconv.ParseFloat(amountStr, 64)

		dt = append(dt, &model.Transaction{
			TrxID:           row[0],
			Type:            model.DefineTrxType(row[2]),
			Amount:          amount,
			TransactionTime: trxTime,
		})
	}

	return dt, nil
}
