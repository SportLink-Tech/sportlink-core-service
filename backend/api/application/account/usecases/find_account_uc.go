package usecases

import (
	"context"
	"fmt"
	"strings"

	"sportlink/api/domain/account"
)

// FindAccountInput selects an account by public AccountID (ULID) or by email address.
// Exactly one of AccountID or Email must be non-empty (after trim).
type FindAccountInput struct {
	AccountID string // Public ULID-based ID
	Email     string
}

type FindAccountUC struct {
	repository account.Repository
}

func NewFindAccountUC(repository account.Repository) *FindAccountUC {
	return &FindAccountUC{repository: repository}
}

func (f *FindAccountUC) Invoke(ctx context.Context, input FindAccountInput) (*[]account.Entity, error) {
	accountID := strings.TrimSpace(input.AccountID)
	email := strings.TrimSpace(input.Email)

	if accountID != "" && email != "" {
		return nil, fmt.Errorf("only one of account id or email must be provided")
	}
	if accountID == "" && email == "" {
		return nil, fmt.Errorf("account id or email is required")
	}

	var query account.DomainQuery
	if accountID != "" {
		query.AccountIDs = []string{accountID}
	} else {
		query.Emails = []string{email}
	}

	result, err := f.repository.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
