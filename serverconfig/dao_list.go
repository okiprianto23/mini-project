package serverconfig

import (
	"main-xyz/dao"
)

func (s *ServerAttribute) InitDAO() {
	result := ListDAO{}

	result.UserDAO = dao.NewRepoUser(s.DBConnection, s.logger)
	result.ClientTokenDAO = dao.NewClientToken(s.DBConnection, s.logger)
	result.ConsumerDAO = dao.NewConsumerDAO(s.DBConnection, s.logger)
	result.CreditLimitDAO = dao.NewCreditLimitDAO(s.DBConnection, s.logger)
	result.TransactionDAO = dao.NewTransactionDAO(s.DBConnection, s.logger)

	s.listDAO = result
}
