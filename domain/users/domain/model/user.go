package model

import "time"

type User struct {
	UserID             string    `json:"userId,required"`
	Email              string    `json:"user_email,required"`
	Name               string    `json:"user_name,required"`
	UserIdentifier     string    `json:"user_identifier,required"`
	UserPassword       string    `json:"user_password,required"`
	UserTypeIdentifier string    `json:"user_type_identifier,required"`
	DateCreated        time.Time `json:"-"`
	UserModify         int       `json:"product_user_modify,omitempty"`
	DateModify         time.Time `json:"-"`
}
