package config

import (
	"context"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	DynamoDbCfg DynamoDbCfg
	AuthCfg     AuthCfg
}

type DynamoDbCfg struct {
	Region string `env:"AWS_REGION,default=us-west-2"`
	Url    string `env:"DYNAMODB_URL,default=http://localhost:4566"`
}

type AuthCfg struct {
	GoogleClientID string `env:"GOOGLE_CLIENT_ID,required"`
	JWTSecret      string `env:"JWT_SECRET,required"`
}

func LoadConfig(ctx context.Context) (*Config, error) {
	var cfg Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
