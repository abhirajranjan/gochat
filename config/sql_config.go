package config

import "time"

type SqlConfig struct {
	DSN        string
	SqlTimeout time.Duration
}
