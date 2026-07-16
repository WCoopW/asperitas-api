package user

import "time"

// UserSchema maps 1:1 to the users table.
type UserSchema struct {
	ID           string    `db:"id"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}
