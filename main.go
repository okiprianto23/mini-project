package main

import (
	migrate "github.com/rubenv/sql-migrate"
	"go.uber.org/zap"
	"main-xyz/config"
	"main-xyz/router"
	"main-xyz/serverconfig"
	"main-xyz/utils"
)

func main() {
	appconfig := config.AppConfig
	serverAttribute := serverconfig.NewServerAttribute(appconfig, config.Logger)
	err := serverAttribute.Init()

	defer func() {
		if serverAttribute.DBConnection != nil {
			serverAttribute.DBConnection.Close()
		}
	}()

	//Migrate SQL
	_, err = utils.NewDBMigration().DbMigrateMysql(serverAttribute.DBConnection, "migrations/global_migrations", migrate.Up)
	if err != nil {
		config.Logger.Logger.Fatal("Error to migration database", zap.Error(err))
		return
	}

	//For API
	router.InitHttpService(serverAttribute.Validator.HttpController)

	serverAttribute.InitEndpoint()

	router.StartService(appconfig, config.Logger)
}
