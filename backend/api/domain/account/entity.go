package account

import "fmt"

// Entity represents an account in the domain
type Entity struct {
	ID       string
	Email    string
	Nickname string
	Picture  string
}

// NewAccount creates a new account entity with an ID generated from the email
func NewAccount(email, nickname string) Entity {
	return Entity{
		ID:       GenerateAccountID(email),
		Email:    email,
		Nickname: nickname,
	}
}

// NewGoogleAccount creates an account from Google OAuth data
func NewGoogleAccount(email, name, picture string) Entity {
	return Entity{
		ID:       GenerateAccountID(email),
		Email:    email,
		Nickname: name,
		Picture:  picture,
	}
}

// GenerateAccountID generates an account ID based on the email
// Format: EMAIL#{{email}}
func GenerateAccountID(email string) string {
	return fmt.Sprintf("EMAIL#%s", email)
}
