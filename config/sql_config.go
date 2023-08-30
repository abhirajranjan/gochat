package config

import "time"

type SqlConfig struct {
	DSN        string        `config:"SQL_DSN"`
	SqlTimeout time.Duration `config:"SQL_TIMEOUT"`
}
