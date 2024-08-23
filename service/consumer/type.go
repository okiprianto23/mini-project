package consumer

import (
	"main-xyz/context"
	"main-xyz/router"
)

type ConsumerService interface {
	GetDTO() interface{}
	GetMultipartDTO() interface{}

	InsertConsumer(ctx *context.ContextModel, param router.URLParam, dto interface{}) (header map[string]string, output interface{}, err error)
	GetListConsumer(ctx *context.ContextModel, _ router.URLParam, dto interface{}) (map[string]string, interface{}, error)
	GetListCreditByConsumerID(ctx *context.ContextModel, param router.URLParam, dto interface{}) (header map[string]string, output interface{}, err error)
}
