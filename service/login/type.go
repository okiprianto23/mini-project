package login

import (
	"main-xyz/context"
	"main-xyz/router"
)

type LoginService interface {
	GetDTO() interface{}
	GetMultipartDTO() interface{}
	LoginService(*context.ContextModel, router.URLParam, interface{}) (map[string]string, interface{}, error)
}
