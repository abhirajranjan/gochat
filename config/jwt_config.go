package config

import "time"

type JwtConfig struct {
	Key    string        `config:"JWT_KEY"`
	Expiry time.Duration `config:"JWT_EXPIRY"`
}
