package login

import (
	"github.com/go-redis/redis/v7"
	"main-xyz/config"
	"main-xyz/context"
	"main-xyz/dao"
	in2 "main-xyz/dto/in"
	"main-xyz/dto/out"
	"main-xyz/repository"
	"main-xyz/token"
)

func NewLoginService(
	logger *config.LoggerCustom,
	userDAO dao.UserDAO,
	clientTokenDAO dao.ClientTokenDAO,
	redisClient *redis.Client,
	local token.LocalToken,
) LoginService {
	return loginService{
		logger:         logger,
		userDAO:        userDAO,
		clientTokenDAO: clientTokenDAO,
		redisClient:    redisClient,
		localToken:     local,
	}
}

type loginService struct {
	logger         *config.LoggerCustom
	userDAO        dao.UserDAO
	clientTokenDAO dao.ClientTokenDAO
	redisClient    *redis.Client
	localToken     token.LocalToken
}

func (l loginService) GetDTO() interface{} {
	return &in2.LoginRequest{}
}

func (l loginService) GetMultipartDTO() interface{} {
	return nil
}

func (l loginService) parseDTO(in interface{}) *in2.LoginRequest {
	return in.(*in2.LoginRequest)
}

func (l loginService) generateAccessOutLogin(ctx *context.ContextModel, clientID, token string) (resultS *out.LoginResponse, err error) {
	var (
		resultUser repository.UserInformation
	)

	resultUser, err = l.userDAO.GetInformationUser(ctx, clientID)
	if err != nil {
		return
	}

	resultS = &out.LoginResponse{
		UserID:    resultUser.ResourceUserID.Int64,
		Token:     token,
		Locale:    resultUser.Locale.String,
		Username:  resultUser.Username.String,
		AliasName: resultUser.AliasName.String,
	}

	return
}
