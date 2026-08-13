// Package model defines the product domain entity exposed over the API.
package model

import "time"

// Product represents a product exposed through the API (request and response body).
type Product struct {
	ProductID          int64     `json:"product_id"`
	ProductName        string    `json:"product_name"`
	ProductAmount      int       `json:"product_amount"`
	ProductPrice       float64   `json:"product_price"`
	ProductUserCreated int       `json:"product_user_created"`
	ProductDateCreated time.Time `json:"-"`
	ProductUserModify  int       `json:"product_user_modify"`
	ProductDateModify  time.Time `json:"-"`
}
