package account

import (
	"fmt"
	"sportlink/api/domain/account"
)

type Dto struct {
	EntityId  string `dynamodbav:"EntityId"`
	Id        string `dynamodbav:"Id"`
	AccountId string `dynamodbav:"AccountId,omitempty"`
	Email     string `dynamodbav:"Email"`
	Nickname  string `dynamodbav:"Nickname,omitempty"`
	FirstName string `dynamodbav:"FirstName,omitempty"`
	LastName  string `dynamodbav:"LastName,omitempty"`
	Picture   string `dynamodbav:"Picture,omitempty"`
}

func From(entity account.Entity) (Dto, error) {
	if entity.Email == "" {
		return Dto{}, fmt.Errorf("email could not be empty")
	}

	return Dto{
		EntityId:  "Entity#Account",
		Id:        account.GenerateAccountID(entity.Email),
		AccountId: entity.AccountID,
		Email:     entity.Email,
		Nickname:  entity.Nickname,
		FirstName: entity.FirstName,
		LastName:  entity.LastName,
		Picture:   entity.Picture,
	}, nil
}

func (d *Dto) ToDomain() account.Entity {
	return account.Entity{
		ID:        d.Id,
		AccountID: d.AccountId,
		Email:     d.Email,
		Nickname:  d.Nickname,
		FirstName: d.FirstName,
		LastName:  d.LastName,
		Picture:   d.Picture,
	}
}
