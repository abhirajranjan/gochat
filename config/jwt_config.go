package config

import "time"

type JwtConfig struct {
	Key    string
	Expiry time.Duration
}
