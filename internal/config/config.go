package config

import (
	"context"
	"log/slog"

	"github.com/valkey-io/valkey-go"
)

type Config struct {
	Client      valkey.Client
	Logger      *slog.Logger
	Port        string
	ServiceName string
}

func NewConfig(ctx context.Context, env *Env, logger *slog.Logger) (*Config, error) {
	client, clientErr := NewValkey(ctx, env.ValkeyURL)
	if clientErr != nil {
		return nil, clientErr
	}

	return &Config{
		Client:      client,
		Logger:      logger,
		Port:        env.Port,
		ServiceName: env.ServiceName,
	}, nil
}

func (c *Config) Close() {
	if c.Client != nil {
		c.Client.Close()
	}
}
