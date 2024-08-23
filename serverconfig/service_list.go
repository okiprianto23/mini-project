package serverconfig

import (
	"main-xyz/service/consumer"
	"main-xyz/service/login"
	"main-xyz/service/transaction"
	"main-xyz/service/user"
)

func (s *ServerAttribute) InitService() {
	result := ListService{}

	result.Login = login.NewLoginService(
		s.logger,
		s.listDAO.UserDAO,
		s.listDAO.ClientTokenDAO,
		s.redis,
		s.TokenLocal,
	)

	result.User = user.NewUserService(
		s.logger,
		s.listDAO.UserDAO,
		s.Bundles,
	)

	result.Consumer = consumer.NewConsumerService(
		s.logger,
		s.listDAO.ConsumerDAO,
		s.listDAO.UserDAO,
		s.listDAO.CreditLimitDAO,
		s.Bundles,
		s.TXHelper,
	)

	result.Transaction = transaction.NewTransactionService(
		s.logger,
		s.listDAO.TransactionDAO,
		s.listDAO.ConsumerDAO,
		s.listDAO.CreditLimitDAO,
		s.TXHelper,
		s.redis,
		s.Bundles,
	)

	s.listService = result
}
