package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"os"

	"github.com/melisa92/reconciliation/internal/model"
	"github.com/melisa92/reconciliation/internal/repository"
	"github.com/melisa92/reconciliation/internal/usecase"
)

func main() {
	flagGenerateCSV := flag.Bool("generate-csv", false, "Generate dummy CSV files")
	flagRunTest := flag.Bool("run-test", false, "Run reconciliation test scenarios")

	flag.Parse()

	switch {
	case *flagGenerateCSV:
		fmt.Println("-- Start Generate CSV --")
		generateCsvSample()
		fmt.Println("-- Finish Generate CSV --")
	case *flagRunTest:
		fmt.Print("-- Start Running TestCase --\n\n")
		initTestCase()
		fmt.Println("\n-- Finish Running TestCase --")
	default:
		fmt.Println("Please specify a command:")
		fmt.Println("  --generate-csv")
		fmt.Println("  --run-test")
	}
}

func initTestCase() {
	transactionFilePath := "data/transaction.csv"
	bankStatementFilePath := []string{
		"data/bankBCA.csv",
		"data/bankMandiri.csv",
	}

	// Init Repo
	bankStatementRepo := repository.NewBankStatementRepo(bankStatementFilePath)
	transactionRepo := repository.NewTransactionRepo(transactionFilePath)

	// Init usecase
	reconciliationUc := usecase.NewTransactionUsecase(bankStatementRepo, transactionRepo)

	// Run Test Case
	ctx := context.Background()
	summary, err := reconciliationUc.ReconciliationProcess(ctx, "2026-01-01", "2026-01-10")
	if err != nil {
		fmt.Println("got error when doing reconciliation, errmsg:", err.Error())
		return
	}
	printSummary(summary)
}

func printSummary(summary *model.ReconciliationSummary) {
	fmt.Printf("Total Data Transaction from CSV = %d\n", summary.TotalTrxRecords)
	fmt.Printf("Total Data Transaction Proceed = %d\n", summary.TotalProceesRecords)
	fmt.Printf("Total Match Transaction = %d\n", summary.TotalMatchedTrxRecords)
	fmt.Printf("Total Unmatch Transaction = %d, here's the breakdown:\n", summary.TotalUnmatchedTransactionRecords+summary.TotalUnmatchedBankStatementRecords)
	fmt.Printf("- Unmatch Transaction (not found in bank statement) = %d\n", summary.TotalUnmatchedTransactionRecords)
	if summary.TotalUnmatchedTransactionRecords > 0 {
		fmt.Println("- Unmatch Transaction (list data): ")
		for _, v := range summary.ListUnmatchTrx {
			fmt.Printf(">> TrxID: %s, TrxType: %s, Amount: %.2f, DateTime: %s\n", v.TrxID, model.TrxTypeToStr(v.Type), v.Amount, v.TransactionTime.Format("2006-01-02T15:04:05Z"))
		}
	}

	fmt.Printf("\n- Unmatch Bank Statement (not found in transaction) = %d\n", summary.TotalUnmatchedBankStatementRecords)
	if summary.TotalUnmatchedBankStatementRecords > 0 {
		fmt.Println("- Breakdown Unmatch Bank Statement (group by Bank Name): ")
		for k, v := range summary.ListUnmatchBankStatement {
			fmt.Printf("[Bank %s]\n", k)
			for i := range v {
				fmt.Printf(">> UniqueID: %s, Amount: %.2f, Date: %s\n", v[i].UniqueID, v[i].Amount, v[i].Date)
			}
		}
	}

	fmt.Println("\nTotal Discrepancies for Matched Transactions = ", summary.TotalDiscrepancies)
}

type DummyCSVTransaction struct {
	TrxID           string  `json:"trxID"`
	Amount          float64 `json:"amount"`
	Type            string  `json:"type"`
	TransactionTime string  `json:"transactionTime"`
}

func generateCsvSample() {
	transactions := []DummyCSVTransaction{
		{
			TrxID:           "TRX001",
			Amount:          100.00,
			Type:            "CREDIT",
			TransactionTime: "2026-01-01T10:00:00Z",
		},
		{
			TrxID:           "TRX002",
			Amount:          200.00,
			Type:            "DEBIT",
			TransactionTime: "2026-01-02T10:30:00Z",
		},
		{
			TrxID:           "TRX003",
			Amount:          300.00,
			Type:            "CREDIT",
			TransactionTime: "2026-01-03T14:15:00Z",
		},
		{
			TrxID:           "TRX004",
			Amount:          400.00,
			Type:            "CREDIT",
			TransactionTime: "2026-01-04T16:45:00Z",
		},
		{
			TrxID:           "TRX005",
			Amount:          500.00,
			Type:            "DEBIT",
			TransactionTime: "2026-01-05T08:45:00Z",
		},
	}
	GenerateTransactionCSV("data/transaction.csv", transactions)

	bank1 := []model.BankStatement{
		{
			UniqueID: "BCA001",
			Amount:   100.00,
			Date:     "2026-01-01",
			BankName: "BCA",
		},
		{
			UniqueID: "BCA002",
			Amount:   -205.00,
			Date:     "2026-01-02",
			BankName: "BCA",
		},
		{
			UniqueID: "BCA003",
			Amount:   305.00,
			Date:     "2026-01-03",
			BankName: "BCA",
		},
		{
			UniqueID: "BCA004",
			Amount:   999.00,
			Date:     "2026-01-10",
			BankName: "BCA",
		},
	}
	GenerateBankStatementCSV("data/bankBCA.csv", bank1)

	bank2 := []model.BankStatement{
		{
			UniqueID: "Mandiri001",
			Amount:   398.00,
			Date:     "2026-01-04",
			BankName: "Mandiri",
		},
		{
			UniqueID: "Mandiri002",
			Amount:   -500.00,
			Date:     "2026-01-05",
			BankName: "Mandiri",
		},
		{
			UniqueID: "Mandiri003",
			Amount:   300.00,
			Date:     "2026-01-03",
			BankName: "Mandiri",
		},
		{
			UniqueID: "Mandiri004",
			Amount:   400.00,
			Date:     "2026-01-09",
			BankName: "Mandiri",
		},
	}
	GenerateBankStatementCSV("data/bankMandiri.csv", bank2)

}

func GenerateTransactionCSV(filepath string, data []DummyCSVTransaction) {
	file, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{
		"trxID",
		"amount",
		"type",
		"transactionTime",
	})
	for _, v := range data {
		writer.Write([]string{
			v.TrxID,
			fmt.Sprint(v.Amount),
			v.Type,
			v.TransactionTime,
		})
	}
}

func GenerateBankStatementCSV(filepath string, data []model.BankStatement) {
	file, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{
		"uniqueID",
		"amount",
		"date",
		"bankName",
	})
	for _, v := range data {
		writer.Write([]string{
			v.UniqueID,
			fmt.Sprint(v.Amount),
			v.Date,
			v.BankName,
		})
	}
}
