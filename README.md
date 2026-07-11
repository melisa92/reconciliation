# Transaction Reconciliation Service

## Overview

This project implements a transaction reconciliation service that compares system transactions against bank statements from multiple banks.

The reconciliation process identifies:

* Matched transactions
* Unmatched system transactions
* Unmatched bank transactions
* Amount discrepancies between matched transactions

---

## Assumptions

### Transaction Matching Logic

Since bank statements do not contain the internal `trxID`, transaction matching is performed using the following rules:

1. Transactions are grouped by:
   * Transaction Date : Transactions are grouped using the date component only (`YYYY-MM-DD`). The time component is ignored during reconciliation.
   * Transaction Type (`DEBIT` or `CREDIT`)
2. A system transaction can only be matched to a single bank statement entry.
3. A bank statement entry can only be matched to a single system transaction.
4. If multiple bank statement candidates exist for a transaction, the candidate with the smallest amount discrepancy is selected.

Example:

System Transaction:
```text
TRX001  Amount=300  Type=CREDIT  Date=2026-01-03
```

Bank Candidates:
```text
ABC001 Amount=305 Date=2026-01-03
DEF001 Amount=300 Date=2026-01-03
```

Result:
```text
TRX001 -> DEF001
```

because:

```text
|300 - 300| = 0
|300 - 305| = 5
```

---

### Amount Discrepancies
A discrepancy is defined as:

```text
absolute(system_amount - absolute(bank_amount))
```

Examples:

```text
System Amount = 200
Bank Amount   = 205

Discrepancy = 5
```

For debit transactions, bank statements may use negative values:

```text
System Transaction:
DEBIT 200

Bank Statement:
-205

Discrepancy = |200 - 205| = 5
```

---

### Unmatched Transactions

A transaction is considered unmatched if:
* A system transaction cannot find any suitable bank statement candidate.
* A bank statement entry is not matched by any system transaction.

---

## Output Summary

The reconciliation process outputs:

* Total number of transactions processed
* Total number of matched transactions
* Total number of unmatched transactions
* List of unmatched system transactions
* List of unmatched bank statements grouped by bank
* Total discrepancies for matched transactions

Example:

```text
Total Data Transaction from CSV = 6
Total Data Transaction Proceed = 5
Total Match Transaction = 5
Total Unmatch Transaction = 1, here's the breakdown:
- Unmatch Transaction (not found in bank statement) = 0
- Unmatch Bank Statement (not found in transaction) = 1
- Breakdown Unmatch Bank Statement (group by Bank Name): 
-- Bank BCA
--- UniqueID: ABC003, Amount: 305.00, Date: 2026-01-03
Total Discrepancies for Matched Transactions =  7
```

---

## Project Structure

```text
cmd/
    main.go

internal/
    interfaces/
    model/
    repository/
    usecase/

data/
    transaction.csv
    bankBCA.csv
    bankMandiri.csv
```

---

## Running the Application

### Run reconciliation test scenario:
1. Update file path under initTestCase() function in main.go
2. Update start date & end date when calling ReconciliationProcess() function
3. Save the file
4. Run command
```bash
go run ./cmd/main.go --run-test
```

### Generate Dummy CSV:
1. Update data inside generateCsvSample() in main.go
2. Run command
```bash
go run ./cmd/main.go --generate-csv
```
