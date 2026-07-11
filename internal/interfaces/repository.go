package interfaces

import (
	"context"

	"github.com/melisa92/reconciliation/internal/model"
)

type BankStatementInterface interface {
	GetBankStatement(ctx context.Context) ([]*model.BankStatement, error)
}

type TransactionStatementInterface interface {
	GetTransaction(ctx context.Context) ([]*model.Transaction, error)
}
