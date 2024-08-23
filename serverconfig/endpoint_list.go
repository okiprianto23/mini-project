package serverconfig

import (
	"main-xyz/endpoint"
	"main-xyz/endpoint/consumer"
	"main-xyz/endpoint/login"
	"main-xyz/endpoint/transaction"
	"main-xyz/endpoint/user"
)

func (s *ServerAttribute) InitEndpoint() {
	endp := endpoint.NewEndpoint()

	//set endpoint
	endp.AddEndpoint(login.NewLoginEndpoint(s.Validator.HttpController, s.listService.Login))
	endp.AddEndpoint(user.NewUserEndpoint(s.Validator.HttpController, s.listService.User))
	endp.AddEndpoint(consumer.NewConsumerEndpoint(s.Validator.HttpController, s.listService.Consumer))
	endp.AddEndpoint(transaction.NewTransactionEndpoint(s.Validator.HttpController, s.listService.Transaction))

	endp.ServeEndpoint()
}
