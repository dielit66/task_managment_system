package entities

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash" db:"password_hash"`
}

type AuthUserDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
