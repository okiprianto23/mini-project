package token

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"main-xyz/constanta"
	internalCtx "main-xyz/context"
	error2 "main-xyz/error"
	"main-xyz/utils/text"
	"main-xyz/utils/validator/basic_validator"
	"strings"
	"time"
)

type userJWTValidator struct {
	config         JWTValidatorConfig
	client         *redis.Client
	basicValidator basic_validator.BasicValidator
}

func (u userJWTValidator) ValidateJWTToken(
	ctx context.Context,
	jwtTokenStr string,
) (
	err error,
) {
	if jwtTokenStr == "" {
		err = error2.ErrUnauthorized
		return
	} else {

		var payload PayloadJWTToken
		payload, err = u.ValidateTokenWithoutCheckSignature(
			jwtTokenStr,
			u.config.ResourceID,
		)

		if err != nil {
			return
		}

		_ctx, valid := ctx.Value(constanta.ApplicationContextConstanta).(*internalCtx.ContextModel)
		if !valid {
			_ctx = internalCtx.NewContextModel()
		}

		_ctx.AuthAccessTokenModel.ClientID = payload.ClientID
		_ctx.AuthAccessTokenModel.Locale = payload.Locale

		ctx = context.WithValue(ctx, constanta.ApplicationContextConstanta, _ctx)

		return u.checkTokenInRedis(
			ctx,
			jwtTokenStr,
		)
	}
}

func (u userJWTValidator) ValidateTokenWithoutCheckSignature(
	jwtTokenStr string,
	resourceID string,
) (
	jwtToken PayloadJWTToken,
	err error,
) {
	jwtToken, err = u.ConvertJWTToPayload(jwtTokenStr)
	if err != nil {
		return
	}

	if time.Now().Unix() > jwtToken.ExpiresAt.Unix() {
		err = error2.ErrExpiredToken
		return
	}

	if resourceID != "" {
		if !u.basicValidator.CheckIsResourceIDExist(jwtToken.Resource, resourceID) {
			err = error2.ErrJWTForbiddenByResource
			return
		}
	}

	return
}

func (u userJWTValidator) ConvertJWTToPayload(
	jwtTokenStr string,
) (
	jwtToken PayloadJWTToken,
	err error,
) {
	splitJWT := strings.Split(jwtTokenStr, ".")
	if len(splitJWT) == 3 {
		payload := splitJWT[1]

		byteData, errs := text.Base64decoder(payload)
		if errs != nil {
			err = errs
			return
		}
		_ = json.Unmarshal(byteData, &jwtToken)
	} else {
		err = error2.ErrUnauthorized
	}

	return
}

func (u userJWTValidator) checkTokenInRedis(
	ctx context.Context,
	token string,
) (
	err error,
) {
	notChecked, err := u.IsTokenUnidentified(ctx, token)

	if err != nil {
		return
	}

	//jika tidak ada lempat 401
	if notChecked {
		err = error2.ErrUnauthorized
		return
	}

	return nil
}

func (u userJWTValidator) IsTokenUnidentified(
	_ context.Context,
	token string,
) (bool, error) {

	redisResult, err := u.client.Get(
		token,
	).Result()

	if err != nil && err != redis.Nil {
		return false, err
	}

	if redisResult == constanta.INVALID_TOKEN_REDIS_VALUE {
		return false, error2.ErrUnauthorized
	}

	return redisResult == "", nil
}

func (u userJWTValidator) LogoutToken(
	ctx context.Context,
	jwtToken string,
) (
	err error,
) {

	isChecked, err := u.IsTokenUnidentified(ctx, jwtToken)
	if err != nil {
		return
	}

	if isChecked {
		return error2.ErrUnauthorized
	}

	ttl, err := u.client.TTL(jwtToken).Result()

	if err != nil {
		return
	}

	_, err = u.client.Set(
		jwtToken,
		constanta.INVALID_TOKEN_REDIS_VALUE,
		ttl,
	).Result()

	if err != nil {
		return
	}

	return
}
