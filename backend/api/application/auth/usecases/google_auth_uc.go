package usecases

import (
	"context"
	"fmt"
	"sportlink/api/application/auth/service"
	"sportlink/api/domain/account"
)

type GoogleAuthResult struct {
	JWTToken  string
	AccountID string
}

type GoogleAuthUC struct {
	googleVerifier service.GoogleTokenVerifier
	accountRepo    account.Repository
	jwtService     service.JWTService
}

func NewGoogleAuthUC(
	googleVerifier service.GoogleTokenVerifier,
	accountRepo account.Repository,
	jwtService service.JWTService,
) *GoogleAuthUC {
	return &GoogleAuthUC{
		googleVerifier: googleVerifier,
		accountRepo:    accountRepo,
		jwtService:     jwtService,
	}
}

func (uc *GoogleAuthUC) Invoke(ctx context.Context, idToken string) (*GoogleAuthResult, error) {
	tokenInfo, err := uc.googleVerifier.Verify(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("invalid Google token: %w", err)
	}

	accounts, err := uc.accountRepo.Find(ctx, account.DomainQuery{Emails: []string{tokenInfo.Email}})
	if err != nil {
		return nil, fmt.Errorf("error finding account: %w", err)
	}

	var accountID string
	if len(accounts) == 0 {
		newAccount := account.NewGoogleAccount(tokenInfo.Email, tokenInfo.Name, tokenInfo.Picture)
		if err := uc.accountRepo.Save(ctx, newAccount); err != nil {
			return nil, fmt.Errorf("error creating account: %w", err)
		}
		accountID = newAccount.ID
	} else {
		accountID = accounts[0].ID
	}

	jwtToken, err := uc.jwtService.Generate(accountID)
	if err != nil {
		return nil, fmt.Errorf("error generating token: %w", err)
	}

	return &GoogleAuthResult{JWTToken: jwtToken, AccountID: accountID}, nil
}
