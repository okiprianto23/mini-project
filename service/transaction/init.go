package transaction

import (
	"github.com/go-redis/redis/v7"
	"main-xyz/config"
	"main-xyz/dao"
	"main-xyz/dto/in"
	"main-xyz/error/bundles"
	"main-xyz/tx_helper"
)

func NewTransactionService(
	logger *config.LoggerCustom,
	transactionDAO dao.TransactionDAO,
	consumerDAO dao.ConsumerDAO,
	creditLimitDAO dao.CreditLimitDAO,
	txHelper tx_helper.TXHelper,
	redis *redis.Client,
	bundles bundles.Bundles,
) TransactionService {
	service := transactionService{
		logger:         logger,
		transactionDAO: transactionDAO,
		consumerDAO:    consumerDAO,
		creditLimitDAO: creditLimitDAO,
		txHelper:       txHelper,
		redis:          redis,
		bundles:        bundles,
	}

	return &service
}

type transactionService struct {
	logger         *config.LoggerCustom
	transactionDAO dao.TransactionDAO
	consumerDAO    dao.ConsumerDAO
	creditLimitDAO dao.CreditLimitDAO
	txHelper       tx_helper.TXHelper
	redis          *redis.Client
	bundles        bundles.Bundles
}

func (t *transactionService) GetDTO() interface{} {
	return &in.TransactionRequest{}
}

func (t *transactionService) GetMultipartDTO() interface{} {
	return nil
}

func (t *transactionService) parseDTO(ts interface{}) *in.TransactionRequest {
	return ts.(*in.TransactionRequest)
}
