package account

import (
	"fmt"
	"sportlink/api/domain/account"
)

type Dto struct {
	EntityId string `dynamodbav:"EntityId"`
	Id       string `dynamodbav:"Id"`
	Email    string `dynamodbav:"Email"`
	Nickname string `dynamodbav:"Nickname"`
	Picture  string `dynamodbav:"Picture,omitempty"`
}

func From(entity account.Entity) (Dto, error) {
	if entity.Email == "" {
		return Dto{}, fmt.Errorf("email could not be empty")
	}

	return Dto{
		EntityId: "Entity#Account",
		Id:       account.GenerateAccountID(entity.Email),
		Email:    entity.Email,
		Nickname: entity.Nickname,
		Picture:  entity.Picture,
	}, nil
}

func (d *Dto) ToDomain() account.Entity {
	return account.Entity{
		ID:       d.Id,
		Email:    d.Email,
		Nickname: d.Nickname,
		Picture:  d.Picture,
	}
}
