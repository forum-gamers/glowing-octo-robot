package wallet

import "time"

type Wallet struct {
	Id        string    `db:"id"`
	UserId    string    `db:"userId"`
	Balance   float64   `db:"balance"`
	Coin      float64   `db:"coin"`
	CreatedAt time.Time `db:"createdAt"`
	UpdatedAt time.Time `db:"updatedAt"`
}
