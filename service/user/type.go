package user

import (
	"main-xyz/context"
	"main-xyz/router"
)

type UserService interface {
	GetDTO() interface{}
	GetMultipartDTO() interface{}
	GetListUser(ctx *context.ContextModel, param router.URLParam, dto interface{}) (headers map[string]string, output interface{}, err error)
}
