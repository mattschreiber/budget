package models

import (
	"fmt"
	"sync"
	"time"
)

// AllBudgetEntries ...
func AllBudgetEntries(startDate, endDate time.Time) ([]Model, error) {
	// now := time.Now()
	rows, err := db.Query(`SELECT b.id, b.credit, b.debit, b.trans_date, s.store_name, c.category_name, b.store_id, b.category_id, pt.payment_name
	FROM budget as b join store as s on b.store_id = s.id join category as c on b.category_id = c.id
	left join payment_type as pt on b.payment_type_id = pt.id
    WHERE b.trans_date  BETWEEN $1 AND $2 ORDER BY b.trans_date DESC, id DESC`, startDate, endDate)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var budgetEntries []Model
	for rows.Next() {
		var budgetRow Model
		// err := rows.Scan(&budgetRow.Id, &budgetRow.Credit, &budgetRow.Debit, &budgetRow.Trans_date, &budgetRow.Store_name, &budgetRow.Category_name, &budgetRow.Store_id, &budgetRow.Category_id)
		err := rows.Scan(&budgetRow.Id, &budgetRow.Credit, &budgetRow.Debit, &budgetRow.Trans_date, &budgetRow.St.Store_name, &budgetRow.Cat.Category_name, &budgetRow.St.Id, &budgetRow.Cat.Id, &budgetRow.Pt.Payment_name)
		if err != nil {
			fmt.Println("error scaninng", err)
			return nil, err
		}
		budgetEntries = append(budgetEntries, budgetRow)
	}
	if err = rows.Err(); err != nil {
		fmt.Println("error scanning a row", err)
		return nil, err
	}
	return budgetEntries, nil
}

// BudgetTotal ...
func BudgetTotal(t time.Time) (balance int, err error) {

	rows, err := db.Query("SELECT SUM(credit - debit) as balance FROM budget where trans_date <= $1", t)
	if err != nil {
		return -1, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&balance)
		if err != nil {
			return -1, err
		}
	}
	if err = rows.Err(); err != nil {
		return -1, err
	}

	return balance, nil
}

// CreateBudgetEntry ...
func CreateBudgetEntry(budget Model) (id int, err error) {
	// only care about date so set time to 0
	err = db.QueryRow("INSERT INTO budget (credit, debit, trans_date, store_id, category_id, payment_type_id) VALUES($1, $2, $3, $4, $5, $6)RETURNING id",
		budget.Credit, budget.Debit, budget.Trans_date.In(getEst()), budget.St.Id, budget.Cat.Id, budget.Pt.Id).Scan(&id)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	return id, nil
}

// DeleteBudgetEntry ...
func DeleteBudgetEntry(id string) (count int64, err error) {
	deleteEntryStmt := "DELETE FROM budget where id = $1"
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

// UpdateBudgetEntry method accepts a ledgerEntry as input and returns an integer for number of records updated
func UpdateBudgetEntry(budgetEntry Model) (count int64, err error) {

	updateStmt := fmt.Sprintf(`UPDATE budget SET debit = $1, credit = $2 WHERE id = $3`)
	res, err := db.Exec(updateStmt, budgetEntry.Debit, budgetEntry.Credit, budgetEntry.Id)
	if err != nil {
		return -1, err
	}

	count, err = res.RowsAffected()
	if err != nil {
		return -1, err
	}
	return count, nil

}

// AutoPayBudgetEntries ....
func AutoPayBudgetEntries(today time.Time) ([]Model, error) {
	// now := time.Now()
	rows, err := db.Query(`SELECT b.id, b.credit, b.debit, b.store_id, b.category_id, b.payment_type_id
    FROM budget as b join store as s on b.store_id = s.id
    WHERE trans_date = $1
    AND s.auto_pay = true`, today)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var budgetEntries []Model
	for rows.Next() {
		var budgetRow Model
		// err := rows.Scan(&budgetRow.Id, &budgetRow.Credit, &budgetRow.Debit, &budgetRow.Trans_date, &budgetRow.Store_name, &budgetRow.Category_name, &budgetRow.Store_id, &budgetRow.Category_id)
		err := rows.Scan(&budgetRow.Id, &budgetRow.Credit, &budgetRow.Debit, &budgetRow.St.Id, &budgetRow.Cat.Id, &budgetRow.Pt.Id)
		if err != nil {
			fmt.Println("error scaninng", err)
			return nil, err
		}
		budgetEntries = append(budgetEntries, budgetRow)
	}
	if err = rows.Err(); err != nil {
		fmt.Println("error scanning a row", err)
		return nil, err
	}
	return budgetEntries, nil
}

// GetBudgetBalance ...
func GetBudgetBalance(startDate time.Time, endDate time.Time, wg *sync.WaitGroup, total *TotalAmounts) {
	var balance int
	defer wg.Done()
	err := db.QueryRow("SELECT SUM(credit-debit) as balance FROM budget WHERE trans_date BETWEEN $1 AND $2",
		startDate, endDate).Scan(&balance)
	if err != nil {
		// fmt.Println(err)
		total.BudgetAmount = 0
	}
	total.BudgetAmount = balance
}
