package user

type User struct {
	ID           int64  `json:"id" db:"id"`
	Email        string `json:"email" db:"email"`
	Username     string `json:"username" db:"username"`
	PasswordHash string `json:"-" db:"password_hash"`
}
