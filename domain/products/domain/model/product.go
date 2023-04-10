package model

import "time"

//Product type struct for database anime
type Product struct {
	ProductID          string    `json:"product_id,required"`
	ProductName        string    `json:"product_name,omitempty"`
	ProductAmount      int       `json:"product_amount,omitempty"`
	ProductPrice 	   float64   `json:"product_price,omitempty"`
	ProductUserCreated int       `json:"product_user_created,omitempty"`
	ProductDateCreated time.Time `json:"-"`
	ProductUserModify  int       `json:"product_user_modify,omitempty"`
	ProductDateModify  time.Time `json:"-"`
}
