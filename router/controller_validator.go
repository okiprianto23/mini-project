package router

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"main-xyz/constanta"
	context2 "main-xyz/context"
	internalCtx "main-xyz/context"
	error2 "main-xyz/error"
	"main-xyz/token"
)

type ControllerValidator struct {
	fixedToken     string
	tokenValidator token.UserJWTValidator
	redis          *redis.Client
}

func (cv ControllerValidator) WhitelistValidator(ctx context.Context, header map[string]string) error {
	return nil
}

func (cv ControllerValidator) UserAccessValidator(
	ctx context.Context,
	header map[string]string,
) error {
	return cv.getTokenOnRedis(ctx, header)
}

func (cv ControllerValidator) getTokenOnRedis(
	ctx context.Context,
	header map[string]string,
) error {

	token := header[constanta.AuthorizationHeaderConstanta]

	isNotChecked, err := cv.tokenValidator.IsTokenUnidentified(
		ctx,
		token,
	)

	if err != nil {
		return err
	}

	if isNotChecked {
		err = cv.tokenValidator.ValidateJWTToken(
			ctx,
			token,
		)

		if err != nil {
			return err
		}

	}

	redisResult, err := cv.redis.Get(
		token,
	).Result()

	if err != nil {
		if err == redis.Nil {
			return error2.ErrUnauthorized
		}

		return err
	}

	return cv.setAuthDataToContext(ctx, redisResult)
}

func (cv ControllerValidator) setAuthDataToContext(
	ctx context.Context,
	result string,
) error {

	var tokenModel context2.AuthAccessTokenModel

	err := json.Unmarshal([]byte(result), &tokenModel)

	if err != nil {
		return error2.ErrUnauthorized
	}

	_ctx, valid := ctx.Value(constanta.ApplicationContextConstanta).(*internalCtx.ContextModel)
	if !valid {
		_ctx = internalCtx.NewContextModel()
	}

	_ctx.AuthAccessTokenModel = tokenModel

	if err != nil {
		return error2.ErrUnauthorized
	}

	_ctx.AuthAccessTokenModel.ResourceUserID = tokenModel.ResourceUserID
	_ctx.AuthAccessTokenModel.SignatureKey = tokenModel.SignatureKey
	_ctx.AuthAccessTokenModel.AliasName = tokenModel.AliasName
	_ctx.AuthAccessTokenModel.Locale = tokenModel.Locale
	_ctx.AuthAccessTokenModel.ClientID = tokenModel.ClientID

	ctx = context.WithValue(ctx, constanta.ApplicationContextConstanta, _ctx)

	return nil
}
