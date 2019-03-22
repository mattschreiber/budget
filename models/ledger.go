package models

import (
	"fmt"
	"time"

	// "bytes"
	"database/sql"
	"sync"

	"github.com/mattschreiber/budget/email"
	"github.com/mattschreiber/budget/utilities"
)

type autoPay struct {
	count    int
	numTimes int
}

func AllLedgerEntries(startDate, endDate time.Time) ([]Model, error) {
	// now := time.Now()
	// before := time.Date(1900, 01, 15, 0, 0, 0, 0, time.UTC)
	rows, err := db.Query(`SELECT l.id, l.credit, l.debit, l.trans_date, s.store_name, c.category_name, l.store_id, l.category_id
    FROM ledger as l join store as s on l.store_id = s.id join category as c on l.category_id = c.id
    WHERE trans_date  BETWEEN $1 AND $2 ORDER BY l.trans_date DESC, id DESC`, startDate, endDate)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var ledgerEntries []Model
	for rows.Next() {
		var ledgerRow Model
		err := rows.Scan(&ledgerRow.Id, &ledgerRow.Credit, &ledgerRow.Debit, &ledgerRow.Trans_date, &ledgerRow.St.Store_name,
			&ledgerRow.Cat.Category_name, &ledgerRow.St.Id, &ledgerRow.Cat.Id)
		if err != nil {
			return nil, err
		}
		ledgerEntries = append(ledgerEntries, ledgerRow)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ledgerEntries, nil
}

// func GetLedgerBalance(startDate time.Time, endDate time.Time, c chan Balance) {
//   var balance int
//   err := db.QueryRow("SELECT sum(credit-debit) as balance from ledger WHERE trans_date BETWEEN $1 AND $2", startDate, endDate).Scan(&balance)
//   if err != nil {
//     // fmt.Println(err)
//     c <- Balance{0, err}
//   }
//   c <- Balance{balance, nil}
// }

func CreateLedgerEntry(ledger Model) (id int, err error) {
	// only care about date so set time to 0
	err = db.QueryRow("INSERT INTO ledger (credit, debit, trans_date, store_id, category_id) VALUES($1, $2, $3, $4, $5)RETURNING id",
		ledger.Credit, ledger.Debit, ledger.Trans_date.In(getEst()), ledger.St.Id, ledger.Cat.Id).Scan(&id)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	return id, nil
}

func DeleteLedgerEntry(id string) (count int64, err error) {
	deleteEntryStmt := "DELETE FROM ledger where id = $1"
	res, err := db.Exec(deleteEntryStmt, id)
	if err != nil {
		return -1, err
	}
	count, err = res.RowsAffected()
	if err != nil {
		return -1, err
	}
	return count, nil
}

// this funcion checks for budget entries on a given day that are marked auto_pay for their store.
// The method makes sure there isn't already an entry for the current month, creates the entry and then sends an
// email with confirmation of the new entries
func AutoPay() {

	var defaultDate = utilities.InitDates()

	// var emailBody bytes.Buffer

	var countNewEntries int

	// find all budgetEntries with a trans_date of today
	budgetEntries, err := AutoPayBudgetEntries(defaultDate.Today)
	if err != nil {
		fmt.Println(err)
		return
	}
	// if there are entries for today
	if budgetEntries != nil {
		//loop through todays entries, check if a corresponding entry exists on the ledger for current monthly
		// and if no ledger entry, then create one
		for _, entry := range budgetEntries {
			// autoPay type used to determine how many times an entry should show up for a given month
			var autopay autoPay
			err = db.QueryRow(`SELECT COUNT(*), store.auto_pay_num FROM ledger JOIN store ON ledger.store_id = store.id
        WHERE extract(month from trans_date) = $1 AND extract(year from trans_date) = $2
        AND ledger.store_id = $3 AND ledger.category_id = $4
        GROUP BY store.auto_pay_num`, defaultDate.CurrentMonth, defaultDate.CurrentYear, entry.St.Id, entry.Cat.Id).Scan(&autopay.count, &autopay.numTimes)
			if err != nil {
				if err == sql.ErrNoRows {
					entry.Trans_date = defaultDate.Today
					_, err = CreateLedgerEntry(entry)
					if err != nil {
						fmt.Println(err)
						return
					}
					countNewEntries += 1
				} else {
					fmt.Println(err)
					return
				}
			}
			if autopay.count < autopay.numTimes {
				entry.Trans_date = defaultDate.Today
				_, err = CreateLedgerEntry(entry)
				if err != nil {
					fmt.Println(err)
					return
				}
				countNewEntries += 1
			}
		}
	}
	// send email
	mail := email.Mail{}
	mail.SenderId = "matt.schreiber01@gmail.com"
	mail.ToIds = []string{"matt.schreiber01@gmail.com"}
	mail.Subject = "New Ledger Entries"
	mail.Body = fmt.Sprintf("Created %d new entries at %s", countNewEntries, defaultDate.Today)
	mail.SendMail()
}

func GetLedgerBalance(startDate time.Time, endDate time.Time, wg *sync.WaitGroup, total *TotalAmounts) {
	var balance int
	defer wg.Done()
	err := db.QueryRow("SELECT sum(credit-debit) as balance from ledger WHERE trans_date BETWEEN $1 AND $2", startDate, endDate).Scan(&balance)
	if err != nil {
		// fmt.Println(err)
		total.LedgerAmount = 0
	}
	total.LedgerAmount = balance
}
