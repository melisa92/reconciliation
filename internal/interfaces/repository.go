package interfaces

import (
	"context"
	"time"

	"github.com/melisa92/reconciliation/internal/model"
)

type BankStatementInterface interface {
	GetBankStatement(ctx context.Context, startDate time.Time, endDate time.Time) ([]model.Transaction, error)
}

type TransactionStatementInterface interface {
	GetTransaction(ctx context.Context, startDate time.Time, endDate time.Time) ([]model.Transaction, error)
}
