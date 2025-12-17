package account

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Entity represents an account in the domain
type Entity struct {
	ID       string
	Email    string
	Nickname string
	Password string
}

// NewAccount creates a new account entity with an ID generated from the email
func NewAccount(email, nickname, password string) Entity {
	return Entity{
		ID:       GenerateAccountID(email),
		Email:    email,
		Nickname: nickname,
		Password: password,
	}
}

// GenerateAccountID generates an account ID based on the email
// Format: EMAIL#{{email}}
func GenerateAccountID(email string) string {
	return fmt.Sprintf("EMAIL#%s", email)
}

// GetHashedPassword returns the hashed version of the password using bcrypt
func (e Entity) GetHashedPassword() (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}
