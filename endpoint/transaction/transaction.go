package transaction

import (
	"main-xyz/constanta"
	"main-xyz/router"
	"main-xyz/service/transaction"
	"net/http"
)

func NewTransactionEndpoint(
	httpValidator *router.HTTPController,
	srv transaction.TransactionService,
) *transactionEndpoint {
	return &transactionEndpoint{
		httpValidator: httpValidator,
		srv:           srv,
	}
}

type transactionEndpoint struct {
	httpValidator *router.HTTPController
	srv           transaction.TransactionService
}

func (e transactionEndpoint) RegisterEndpoint() {

	e.httpValidator.HandleFunc(
		router.NewHandleFuncParam(
			constanta.PathMain+"/transaction",
			e.httpValidator.WrapService(
				router.NewWarpServiceParam(
					e.srv,
					e.srv.TransactionServiceInsert,
					e.httpValidator.UserAccessValidator,
				).
					Menu("insert"),
			),
			http.MethodPost,
		),
	)

}
