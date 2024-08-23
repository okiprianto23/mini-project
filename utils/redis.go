package utils

import (
	"fmt"
	"github.com/go-redis/redis/v7"
)

func NewRedisParam(
	host string,
	port ...int,
) *redisParam {
	return &redisParam{
		host: host,
		port: port,
	}
}

type redisParam struct {
	host       string
	port       []int
	db         int
	password   string
	username   string
	maxRetries int
}

func (r *redisParam) DB(db int) *redisParam {
	r.db = db
	return r
}

func (r *redisParam) Password(password string) *redisParam {
	r.password = password
	return r
}

func (r *redisParam) Username(username string) *redisParam {
	r.username = username
	return r
}

func (r *redisParam) MaxRetries(max int) *redisParam {
	r.maxRetries = max
	return r
}

func ConnectRedis(
	param *redisParam,
) *redis.Client {

	if len(param.port) == 0 {
		param.port = append(param.port, 6379)
	}

	opts := &redis.Options{
		Addr: fmt.Sprintf("%s:%d", param.host, param.port[0]),
		DB:   0,
	}

	if param.db > 0 {
		opts.DB = param.db
	}

	if param.password != "" {
		opts.Password = param.password
	}

	if param.username != "" {
		opts.Username = param.username
	}

	if param.maxRetries > 0 {
		opts.MaxRetries = param.maxRetries
	}

	return redis.NewClient(opts)
}
