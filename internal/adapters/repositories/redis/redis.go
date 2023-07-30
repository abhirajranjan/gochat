package redis

import "gochat/logger"

type redisRepo struct {
	logger logger.ILogger
}

func NewRedisRepository(l logger.ILogger) *redisRepo {
	return &redisRepo{logger: l}
}
