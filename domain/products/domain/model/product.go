package model

import "time"

// Product type struct for database anime.
type Product struct {
	ProductID          string    `json:"product_id"`
	ProductName        string    `json:"product_name"`
	ProductAmount      int       `json:"product_amount"`
	ProductPrice       float64   `json:"product_price"`
	ProductUserCreated int       `json:"product_user_created"`
	ProductDateCreated time.Time `json:"-"`
	ProductUserModify  int       `json:"product_user_modify"`
	ProductDateModify  time.Time `json:"-"`
}
