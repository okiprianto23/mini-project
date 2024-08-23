package serverconfig

import (
	"fmt"
	"main-xyz/config"
	"main-xyz/constanta"
	error2 "main-xyz/error"
	"main-xyz/error/bundles"
	"main-xyz/router"
	"main-xyz/server"
	"main-xyz/token"
	"main-xyz/tx_helper"
	"main-xyz/utils"
	"main-xyz/utils/validator/basic_validator"
	"main-xyz/utils/validator/multipart_validator"
	"main-xyz/utils/validator/tag_validator"
	"strconv"
)

func NewServerAttribute(cfg config.Configuration, log *config.LoggerCustom) ServerAttribute {
	return ServerAttribute{
		config: cfg,
		logger: log,
	}
}

func (s *ServerAttribute) Init() (err error) {

	//Create Connection DB
	mysqlCon := s.config.Mysql
	s.DBConnection = utils.GetDBConnection(s.logger,
		utils.DBAddParam().
			SetAddress(mysqlCon.Username, mysqlCon.Password, mysqlCon.Host, strconv.Itoa(mysqlCon.Port), mysqlCon.DBName).
			AddParseTime().
			MaxOpenConnection(mysqlCon.MaxOpenConnection).
			MaxIdleConnection(mysqlCon.MaxIdleConnection),
	)

	//Redis Conenction
	s.redis = utils.ConnectRedis(
		utils.NewRedisParam(
			s.config.Redis.Host,
			s.config.Redis.Port,
		).
			DB(s.config.Redis.DB).
			MaxRetries(s.config.Redis.MaxRetries).
			Password(s.config.Redis.Password).
			Username(s.config.Redis.Username),
	)

	//set bundles
	s.Bundles, err = bundles.NewBundles("i18n", "en-US", s.logger)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Set for tx Helper
	s.TXHelper = tx_helper.NewTXHelper(s.DBConnection)

	//set formator error
	errFormator := error2.NewErrorFormator(s.Bundles).DefaultInternalCode("E-5-GL-SRV-001").
		DefaultLanguage("en-US").
		Version(s.config.Server.Version)

	//set jwtConfig
	jwtConfig := token.JWTValidatorConfig{
		ResourceID: s.config.Server.ResourceID,
		UserKey:    s.config.Token.UserKey,
		ClientID:   s.config.Credentials.ClientID,
		UserID:     s.config.Credentials.UserID,
		Version:    s.config.Server.Version,
		Duration:   s.config.Token.Duration,
		FixedToken: s.config.Token.FixedInternalToken,
	}

	//set multipartReader
	s.MultipartReader = server.NewMultipartReader(
		50000,
		s.config.File.Directory,
		s.logger,
	).SetFileNameAliasingFunction(server.UUIDAliasingFunction)

	//set validator
	s.Validator.tagValidator = tag_validator.NewTagValidator()
	s.Validator.basicValidator = basic_validator.BasicValidator{}
	s.Validator.multipartValidator = multipart_validator.NewMultipartValidator(
		s.MultipartReader,
		s.Validator.basicValidator,
		s.Validator.tagValidator,
	)

	//token
	s.Validator.tokenValidator = token.NewJWTTokenValidator(
		jwtConfig,
		s.redis,
		s.Validator.basicValidator,
	)

	s.Validator.HttpController = router.NewHTTPController()
	s.Validator.HttpController.Version(s.config.Server.Version)
	s.Validator.HttpController.BasicValidator(s.Validator.basicValidator)
	s.Validator.HttpController.TagValidator(s.Validator.tagValidator)
	s.Validator.HttpController.MultipartValidator(s.Validator.multipartValidator)
	s.Validator.HttpController.MultipartReader(s.MultipartReader)
	s.Validator.HttpController.Formator(errFormator)
	s.Validator.HttpController.Redis(s.redis)
	s.Validator.HttpController.TokenValidator(s.Validator.tokenValidator)

	//set for token local generate
	s.TokenLocal = token.NewLocalToken("main", s.config.Token.UserKey, s.config.Credentials.ClientID, constanta.Default1DayExpired)

	//init dao
	s.InitDAO()

	//init service
	s.InitService()

	return
}
