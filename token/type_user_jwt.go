package token

import (
	"context"
	"github.com/go-redis/redis/v7"
	"main-xyz/utils/validator/basic_validator"
	"time"
)

type JWTValidatorConfig struct {
	ResourceID  string
	UserKey     string
	InternalKey string
	ClientID    string
	UserID      int64
	Version     string
	Duration    time.Duration
	FixedToken  string
}

type UserJWTValidator interface {
	ValidateJWTToken(
		ctx context.Context,
		jwtTokenStr string,
	) (
		err error,
	)

	IsTokenUnidentified(
		ctx context.Context,
		token string,
	) (
		bool,
		error,
	)

	LogoutToken(
		ctx context.Context,
		jwtToken string,
	) (
		err error,
	)
}

func NewJWTTokenValidator(
	config JWTValidatorConfig,
	client *redis.Client,
	validator basic_validator.BasicValidator,
) UserJWTValidator {
	return userJWTValidator{
		config:         config,
		client:         client,
		basicValidator: validator,
	}
}
