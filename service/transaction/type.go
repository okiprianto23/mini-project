package transaction

import (
	"main-xyz/context"
	"main-xyz/router"
)

type TransactionService interface {
	GetDTO() interface{}
	GetMultipartDTO() interface{}

	TransactionServiceInsert(ctx *context.ContextModel, _ router.URLParam, dto interface{}) (map[string]string, interface{}, error)
}
