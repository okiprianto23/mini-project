package consumer

import (
	"main-xyz/constanta"
	"main-xyz/router"
	"main-xyz/service/consumer"
	"net/http"
)

func NewConsumerEndpoint(
	httpValidator *router.HTTPController,
	srv consumer.ConsumerService,
) *consumerEndpoint {
	return &consumerEndpoint{
		httpValidator: httpValidator,
		srv:           srv,
	}
}

type consumerEndpoint struct {
	httpValidator *router.HTTPController
	srv           consumer.ConsumerService
}

func (e consumerEndpoint) RegisterEndpoint() {

	e.httpValidator.HandleFunc(
		router.NewHandleFuncParam(
			constanta.PathMain+"/consumer/insert",
			e.httpValidator.WrapService(
				router.NewWarpServiceParam(
					e.srv,
					e.srv.InsertConsumer,
					e.httpValidator.UserAccessValidator,
				).
					Menu("insert").
					Multipart(),
			),
			http.MethodPost,
		),
	)

	e.httpValidator.HandleFunc(
		router.NewHandleFuncParam(
			constanta.PathMain+"/consumer/list",
			e.httpValidator.WrapService(
				router.NewWarpServiceParam(
					e.srv,
					e.srv.GetListConsumer,
					e.httpValidator.UserAccessValidator,
				),
			),
			http.MethodGet,
		),
	)

	e.httpValidator.HandleFunc(
		router.NewHandleFuncParam(
			constanta.PathMain+"/consumer/credit/{ID}",
			e.httpValidator.WrapService(
				router.NewWarpServiceParam(
					e.srv,
					e.srv.GetListCreditByConsumerID,
					e.httpValidator.UserAccessValidator,
				).PathParams("ID"),
			),
			http.MethodGet,
		),
	)

}
