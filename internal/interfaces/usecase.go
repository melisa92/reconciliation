package interfaces

import (
	"context"

	"github.com/melisa92/reconciliation/internal/model"
)

type ReconciliationInterface interface {
	ReconciliationProcess(ctx context.Context, startDate, endDate string) (*model.ReconciliationSummary, error)
}
