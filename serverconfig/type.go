package serverconfig

import (
	"database/sql"
	"github.com/go-redis/redis/v7"
	"main-xyz/config"
	"main-xyz/dao"
	"main-xyz/error/bundles"
	"main-xyz/router"
	"main-xyz/server"
	"main-xyz/service/consumer"
	"main-xyz/service/login"
	"main-xyz/service/transaction"
	"main-xyz/service/user"
	"main-xyz/token"
	"main-xyz/tx_helper"
	"main-xyz/utils/validator/basic_validator"
	"main-xyz/utils/validator/multipart_validator"
	"main-xyz/utils/validator/tag_validator"
)

type ServerAttribute struct {
	config   config.Configuration
	logger   *config.LoggerCustom
	Bundles  bundles.Bundles
	TXHelper tx_helper.TXHelper

	redis        *redis.Client
	DBConnection *sql.DB

	listDAO     ListDAO
	listService ListService

	MultipartReader server.MultipartReader
	TokenLocal      token.LocalToken
	Validator       validator
}

type validator struct {
	basicValidator     basic_validator.BasicValidator
	tagValidator       tag_validator.TagValidator
	tokenValidator     token.UserJWTValidator
	multipartValidator multipart_validator.MultipartValidator
	HttpController     *router.HTTPController
}

type ListService struct {
	Login       login.LoginService
	User        user.UserService
	Consumer    consumer.ConsumerService
	Transaction transaction.TransactionService
}

type ListDAO struct {
	UserDAO        dao.UserDAO
	ClientTokenDAO dao.ClientTokenDAO
	ConsumerDAO    dao.ConsumerDAO
	CreditLimitDAO dao.CreditLimitDAO
	TransactionDAO dao.TransactionDAO
}
