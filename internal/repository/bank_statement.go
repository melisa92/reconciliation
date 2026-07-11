package repository

import (
	"context"
	"encoding/csv"
	"io"
	"os"
	"strconv"

	"github.com/melisa92/reconciliation/internal/model"
)

type BankStatementRepo struct {
	FilePath []string
}

func NewBankStatementRepo(filepath []string) *BankStatementRepo {
	return &BankStatementRepo{
		FilePath: filepath,
	}
}

func (r *BankStatementRepo) GetBankStatement(ctx context.Context) ([]*model.BankStatement, error) {
	raw := []*model.BankStatement{}
	for i := range r.FilePath {
		file, err := os.Open(r.FilePath[i])
		if err != nil {
			return nil, err
		}
		defer file.Close()

		reader := csv.NewReader(file)

		for {
			row, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}

			if row[0] == "uniqueID" {
				continue
			}

			amountStr := row[1]
			amount, _ := strconv.ParseFloat(amountStr, 64)

			raw = append(raw, &model.BankStatement{
				UniqueID: row[0],
				BankName: row[3],
				Amount:   amount,
				Date:     row[2],
			})
		}
	}

	return raw, nil
}
