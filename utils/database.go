package utils

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"main-xyz/config"

	_ "github.com/go-sql-driver/mysql"
)

func DBAddParam() *DBAddressParam {
	return &DBAddressParam{
		maxOpenConnection: 50,
		maxIdleConnection: 10,
	}
}

type DBAddressParam struct {
	address           string
	defaultSchema     string
	maxOpenConnection int
	maxIdleConnection int
}

// sql string username:password@tcp(localhost:3306)/dbname
func (d *DBAddressParam) SetAddress(username, password, host, port, dbname string) *DBAddressParam {
	d.address = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname)
	return d
}

func (d *DBAddressParam) AddParseTime() *DBAddressParam {
	d.address += "?parseTime=true"
	return d
}

func (d *DBAddressParam) MaxOpenConnection(con int) *DBAddressParam {
	d.maxOpenConnection = con
	return d
}

func (d *DBAddressParam) MaxIdleConnection(con int) *DBAddressParam {
	d.maxIdleConnection = con
	return d
}

type DBInfo struct {
	instance      *sql.DB
	driver        string
	connectionStr string
}

var instance *sql.DB

func GetDBConnection(
	logger *config.LoggerCustom,
	param *DBAddressParam,
) *sql.DB {
	_dbInfo := DBInfo{
		instance:      nil,
		driver:        "mysql",
		connectionStr: param.address,
	}

	_db, _err := getInstance(_dbInfo, logger)
	if _err != nil {
		logger.Logger.Fatal("Error Find When Connect to Database", zap.Error(_err))
	}

	_db.SetMaxOpenConns(param.maxOpenConnection)
	_db.SetMaxIdleConns(param.maxIdleConnection)

	return _db

}

func getInstance(connInfo DBInfo, logger *config.LoggerCustom) (*sql.DB, error) {
	var _errOpen error

	instance, _errOpen = sql.Open(connInfo.driver, connInfo.connectionStr)

	if _errOpen != nil {
		logger.Logger.Error(fmt.Sprintf("Connect failed to DB %v", connInfo), zap.Error(_errOpen))
		instance = nil
	}

	return instance, _errOpen
}
