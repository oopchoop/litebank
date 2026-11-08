package auth

import (
	"time"
)

type NewUser struct {
	PhoneNumber  string
	PasswordHash string
	FirstName    string
	LastName     string
	Patronymic   *string
}

type RegisterInput struct {
	PhoneNumber string
	Password    string
	FirstName   string
	LastName    string
	Patronymic  string
}

type AuthInput struct {
	PhoneNumber string
	Password    string
}

type User struct {
	ID           int64     `db:"id" json:"id"`
	PhoneNumber  string    `db:"phone_number" json:"phone_number"`
	PasswordHash string    `db:"password_hash" json:"-"`
	FirstName    string    `db:"first_name" json:"first_name"`
	LastName     string    `db:"last_name" json:"last_name"`
	Patronymic   *string   `db:"patronymic" json:"patronymic"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}
