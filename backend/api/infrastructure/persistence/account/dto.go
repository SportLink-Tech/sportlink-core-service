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
	Password string `dynamodbav:"Password"`
}

func From(entity account.Entity) (Dto, error) {
	if entity.Email == "" {
		return Dto{}, fmt.Errorf("email could not be empty")
	}

	hashedPassword, err := entity.GetHashedPassword()
	if err != nil {
		return Dto{}, fmt.Errorf("error hashing password: %w", err)
	}

	return Dto{
		EntityId: "Entity#Account",
		Id:       account.GenerateAccountID(entity.Email),
		Email:    entity.Email,
		Nickname: entity.Nickname,
		Password: hashedPassword,
	}, nil
}

func (d *Dto) ToDomain() account.Entity {
	return account.Entity{
		ID:       d.Id,
		Email:    d.Email,
		Nickname: d.Nickname,
		Password: d.Password,
	}
}
