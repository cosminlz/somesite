package model

import (
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// UserID is the type identifier for a User
type UserID string

// User is the struct representing a user
type User struct {
	ID           UserID     `json:"id,omitempty" db:"user_id"`
	Email        *string    `json:"email,omitempty" db:"email"`
	PasswordHash *[]byte    `json:"-" db:"password_hash"`
	CreatedAt    *time.Time `json:"-" db:"created_at"`
	DeletedAt    *time.Time `json:"-" db:"deleted_at"`
}

//Verify all fields before create or update
func (u *User) Verify() error {
	if u.Email == nil || len(*u.Email) == 0 {
		return errors.New("Email is required")
	}
	return nil
}

// SetPassword updates a user's password
func (u *User) SetPassword(password string) error {

	hash, err := HashPassword(password)
	if err != nil {
		return err
	}

	u.PasswordHash = &hash

	return nil
}

// CheckPassword checks if the input password is the correct one
func (u *User) CheckPassword(password string) error {
	if u.PasswordHash != nil && len(*u.PasswordHash) == 0 {
		return errors.New("password not set")
	}

	return bcrypt.CompareHashAndPassword(*u.PasswordHash, []byte(password))
}

// HashPassword hashes a user's raw password
func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
