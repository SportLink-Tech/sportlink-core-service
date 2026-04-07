package account

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

// Entity represents an account in the domain
type Entity struct {
	ID        string // DynamoDB sort key: EMAIL#<email>
	AccountID string // Public-facing short ID (ULID)
	Email     string
	Nickname  string
	FirstName string
	LastName  string
	Picture   string
}

// NewAccount creates a new account entity with an ID generated from the email
func NewAccount(email, nickname string) Entity {
	return Entity{
		ID:        GenerateAccountID(email),
		AccountID: GenerateULID(),
		Email:     email,
		Nickname:  nickname,
	}
}

// NewGoogleAccount creates an account from Google OAuth data
func NewGoogleAccount(email, givenName, familyName, picture string) Entity {
	return Entity{
		ID:        GenerateAccountID(email),
		AccountID: GenerateULID(),
		Email:     email,
		Nickname:  "",
		FirstName: givenName,
		LastName:  familyName,
		Picture:   picture,
	}
}

// GenerateAccountID generates the DynamoDB sort key based on the email
// Format: EMAIL#{{email}}
func GenerateAccountID(email string) string {
	return fmt.Sprintf("EMAIL#%s", email)
}

func GenerateULID() string {
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}
