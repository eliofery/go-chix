package model

const UserTableName = "users"

type User struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Age          int    `json:"age"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}
