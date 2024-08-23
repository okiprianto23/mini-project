package login

import (
	"database/sql"
	"go.uber.org/zap"
	"main-xyz/constanta"
	"main-xyz/context"
	error2 "main-xyz/error"
	"main-xyz/repository"
	"main-xyz/router"
	"main-xyz/utils"
	"time"
)

func (l loginService) LoginService(ctx *context.ContextModel, _ router.URLParam, dto interface{}) (header map[string]string, output interface{}, err error) {
	var (
		result      repository.UserAdmin
		redisAccess context.RedisAuthAccessTokenModel
	)

	//set expired token 1 hari
	expiredToken := time.Now().Add(24 * time.Hour)
	expiredTokenDrt := constanta.Default1DayExpired

	dtoIn := l.parseDTO(dto)

	//get user in db by username
	result, err = l.userDAO.CheckIsUserAdminByUsername(ctx, repository.UserAdmin{
		Username: sql.NullString{String: dtoIn.Username},
	})

	if err != nil {
		return
	}

	if result.ClientID.String == "" {
		err = error2.ErrUnauthorized
		return
	}

	// Check Password
	isValid := utils.CheckPasswordHash(dtoIn.Password, result.Password.String)
	if !isValid {
		err = error2.ErrUnauthorized
		return
	}

	// Generate token
	dtoIn.Token = l.localToken.GenerateLocalToken(result.ClientID.String, result.UserID.Int64, "", constanta.LocaleHitAPI)

	//create redis model
	redisAccess, err = l.createJSONForRedisAdmin(dto, result, expiredToken)
	if err != nil {
		return
	}

	//SET TO REDIS
	err = l.redisClient.Set(dtoIn.Token, utils.StructToJSON(redisAccess), expiredTokenDrt).Err()
	if err != nil {
		l.logger.Logger.Error("Error to set token to redis", zap.Error(err))
		return
	}

	output, err = l.generateAccessOutLogin(ctx, result.ClientID.String, dtoIn.Token)
	if err != nil {
		return
	}
	return
}

func (l loginService) createJSONForRedisAdmin(in interface{}, userModel repository.UserAdmin, expired time.Time) (redisAccess context.RedisAuthAccessTokenModel, err error) {
	var (
		dtoIn         = l.parseDTO(in)
		clientTokenID sql.NullInt64
	)

	redisAccess = context.RedisAuthAccessTokenModel{
		ResourceUserID: userModel.UserID.Int64,
		SignatureKey:   userModel.SignatureKey.String,
		Locale:         userModel.Locale.String,
		AliasName:      userModel.AliasName.String,
		ClientAlias:    userModel.ClientAlias.String,
	}

	//set to client token
	clientTokenModel := repository.ClientTokenModel{
		UserID:        sql.NullInt64{Int64: userModel.UserID.Int64},
		Token:         sql.NullString{String: dtoIn.Token},
		ExpiredAt:     sql.NullTime{Time: expired},
		CreatedBy:     sql.NullInt64{Int64: userModel.UserID.Int64},
		CreatedClient: sql.NullString{String: userModel.AliasName.String},
	}

	clientTokenID, err = l.clientTokenDAO.CheckExistClientToken(clientTokenModel)
	if err != nil {
		return
	}

	if clientTokenID.Int64 == 0 {
		err = l.clientTokenDAO.InsertClientToken(clientTokenModel)
		if err != nil {
			return
		}
	}

	return
}
