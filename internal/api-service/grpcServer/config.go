package grpcServer

import "time"

type GrpcConfig struct {
	Addr                string        `mapstructure:"addr"`
	ConnectionTimeout   time.Duration `mapstructure:"connectionTimeout"`
	FunctionCallTimeout time.Duration `mapstructure:"functioncCallTimeout"`
}
