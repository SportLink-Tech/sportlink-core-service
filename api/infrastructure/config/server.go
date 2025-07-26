package config

import (
	"context"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	DynamoDbCfg DynamoDbCfg
}
type DynamoDbCfg struct {
	Region string `env:"AWS_REGION,default=us-west-2"`
	Url    string `env:"DYNAMODB_URL,default=http://localhost:4566"`
}

func LoadConfig(ctx context.Context) (*Config, error) {
	var cfg Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
