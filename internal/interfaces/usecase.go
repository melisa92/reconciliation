package interfaces

import (
	"context"
	"time"

	"github.com/melisa92/reconciliation/internal/model"
)

type ReconciliationInterface interface {
	ReconciliationProcess(ctx context.Context, startDate time.Time, endDate time.Time) (model.ReconciliationSummary, error)
}
