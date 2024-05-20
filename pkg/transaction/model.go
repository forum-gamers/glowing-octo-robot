package transaction

import "time"

type Transaction struct {
	Id              string              `db:"id"`
	UserId          string              `db:"userId"`
	Amount          float64             `db:"amount"`
	Type            string              `db:"type"`
	Currency        TransactionCurrency `db:"currency"`
	Status          TransactionStatus   `db:"status"`
	TransactionDate time.Time           `db:"transactionDate"`
	Description     string              `db:"description"`
	Discount        float64             `db:"discount"`
	Detail          string              `db:"detail"`
	Signature       string              `db:"signature"`
	ItemId          string              `db:"itemId"`
	CreatedAt       time.Time           `db:"createdAt"`
	UpdatedAt       time.Time           `db:"updatedAt"`
}
