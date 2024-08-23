package dao

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"main-xyz/config"
	"main-xyz/repository"
)

func NewClientToken(db *sql.DB, logger *config.LoggerCustom) ClientTokenDAO {
	return ClientTokenDAO{
		db:        db,
		tableName: "client_token",
		logger:    logger,
	}
}

type ClientTokenDAO struct {
	db        *sql.DB
	tableName string
	logger    *config.LoggerCustom
}

func (ct ClientTokenDAO) InsertClientToken(userParam repository.ClientTokenModel) (err error) {

	query := fmt.Sprintf(`INSERT INTO %s 
    	(uuid_id, user_id, token, expires_at, created_by, created_client) 
	VALUES 
		(UUID(), ?, ?, ?, ?, ?)`, ct.tableName)

	param := []interface{}{
		userParam.UserID.Int64, userParam.Token.String, userParam.ExpiredAt.Time, userParam.CreatedBy.Int64,
		userParam.CreatedClient.String,
	}
	stmt, errorS := ct.db.Prepare(query)
	if errorS != nil {
		err = errorS
		ct.logger.Logger.Error("Error to prepare query", zap.Error(errorS))
		return
	}
	defer stmt.Close()

	_, errorS = stmt.Exec(param...)
	if errorS != nil {
		err = errorS
		ct.logger.Logger.Error("Error to exec query", zap.Error(errorS))
		return
	}

	return
}

func (ct ClientTokenDAO) CheckExistClientToken(model repository.ClientTokenModel) (resultID sql.NullInt64, err error) {
	query := fmt.Sprintf(`SELECT id FROM %s WHERE user_id = ? AND token = ?`, ct.tableName)

	param := []interface{}{model.UserID.Int64, model.Token.String}

	results := ct.db.QueryRow(query, param...)

	err = results.Scan(&resultID)
	if err != nil && err != sql.ErrNoRows {
		ct.logger.Logger.Error("Error to get query", zap.Error(err))
		return
	}

	err = nil
	return
}
