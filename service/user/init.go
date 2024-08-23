package user

import (
	"main-xyz/config"
	"main-xyz/dao"
	"main-xyz/dto/in"
	"main-xyz/error/bundles"
)

func NewUserService(
	logger *config.LoggerCustom,
	userDAO dao.UserDAO,
	bundles bundles.Bundles,
) UserService {
	return userService{
		logger:  logger,
		userDAO: userDAO,
		bundles: bundles,
	}
}

type userService struct {
	logger  *config.LoggerCustom
	userDAO dao.UserDAO
	bundles bundles.Bundles
}

func (u userService) GetDTO() interface{} {
	return &in.UserRequest{}
}

func (u userService) GetMultipartDTO() interface{} {
	return nil
}
