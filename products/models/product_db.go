package models

import "time"

// ProductDB is the persistence entity for the public.products table.
// It carries only db tags; API-facing fields live on Product.
type ProductDB struct {
	ProductID          int64     `db:"product_id"`
	ProductName        string    `db:"product_name"`
	ProductAmount      int       `db:"product_amount"`
	ProductPrice       float64   `db:"product_price"`
	ProductUserCreated int       `db:"product_user_created"`
	ProductDateCreated time.Time `db:"product_date_created"`
	ProductUserModify  int       `db:"product_user_modify"`
	ProductDateModify  time.Time `db:"product_date_modify"`
}
