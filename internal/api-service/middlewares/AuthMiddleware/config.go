package AuthMiddleware

import "time"

type AuthConf struct {
	Realm           string        `mapstructure:"relm"`
	Key             []byte        `mapstructure:"key"`
	IdentityKey     string        `mapstructure:"identityKey"`
	TokenLookup     string        `mapstructure:"tokenLookup"`
	TimeoutDuration time.Duration `mapstructure:"timeoutDuration"`
	MaxRefresh      time.Duration `mapstructure:"maxRefresh"`
	TokenHeadName   string        `mapstructure:"tokenHeadName"`
}
