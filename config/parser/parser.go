package parser

import "github.com/caarlos0/env/v9"

func Load(cfg any) error {
	opts := env.Options{
		TagName: "config",
	}

	return env.ParseWithOptions(cfg, opts)
}
