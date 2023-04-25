package model

import "time"

type User struct {
	UserID             string    `json:"userId"`
	Email              string    `json:"user_email"`
	Name               string    `json:"user_name"`
	UserIdentifier     string    `json:"user_identifier"`
	UserPassword       string    `json:"user_password"`
	UserTypeIdentifier string    `json:"user_type_identifier"`
	DateCreated        time.Time `json:"-"`
	UserModify         int       `json:"product_user_modify"`
	DateModify         time.Time `json:"-"`
}
