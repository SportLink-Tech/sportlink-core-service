package helper

import (
	"context"
	"testing"

	accountuc "sportlink/api/application/account/usecases"
	"sportlink/api/domain/account"
)

// AccountBuilder builds and persists an account.Entity for e2e tests.
// All fields have sensible defaults; override only what matters for the test.
type AccountBuilder struct {
	t         *testing.T
	repo      account.Repository
	email     string
	nickname  string
	firstName string
	lastName  string
	picture   string
}

// NewAccountBuilder returns a builder with sensible defaults.
// Email is randomized to avoid collisions between test runs.
func NewAccountBuilder(t *testing.T, repo account.Repository) *AccountBuilder {
	t.Helper()
	return &AccountBuilder{
		t:         t,
		repo:      repo,
		email:     randomEmail(),
		nickname:  "johndoe",
		firstName: "John",
		lastName:  "Doe",
	}
}

// randomEmail returns a unique email address with a ULID-based local part,
// giving a collision probability negligible enough for test use.
func randomEmail() string {
	return account.GenerateULID() + "@test.com"
}

func (b *AccountBuilder) WithEmail(email string) *AccountBuilder {
	b.email = email
	return b
}

func (b *AccountBuilder) WithNickname(nickname string) *AccountBuilder {
	b.nickname = nickname
	return b
}

func (b *AccountBuilder) WithFirstName(firstName string) *AccountBuilder {
	b.firstName = firstName
	return b
}

func (b *AccountBuilder) WithLastName(lastName string) *AccountBuilder {
	b.lastName = lastName
	return b
}

func (b *AccountBuilder) WithPicture(picture string) *AccountBuilder {
	b.picture = picture
	return b
}

// Build creates the entity, saves it via the repository and returns it.
// It calls t.Fatal on any error.
func (b *AccountBuilder) Build(ctx context.Context) *account.Entity {
	b.t.Helper()

	entity := account.Entity{
		ID:        account.GenerateAccountID(b.email),
		AccountID: account.GenerateULID(),
		Email:     b.email,
		Nickname:  b.nickname,
		FirstName: b.firstName,
		LastName:  b.lastName,
		Picture:   b.picture,
	}

	uc := accountuc.NewCreateAccountUC(b.repo, account.NewValidator())
	result, err := uc.Invoke(ctx, entity)
	if err != nil {
		b.t.Fatalf("AccountBuilder: failed to save account: %v", err)
	}

	return result
}
