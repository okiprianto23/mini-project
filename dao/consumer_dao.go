package dao

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"main-xyz/config"
	"main-xyz/context"
	"main-xyz/repository"
	"main-xyz/utils/compare"
)

func NewConsumerDAO(
	db *sql.DB,
	logger *config.LoggerCustom,
) ConsumerDAO {
	return ConsumerDAO{
		logger:    logger,
		db:        db,
		tableName: "consumer",
	}
}

type ConsumerDAO struct {
	logger    *config.LoggerCustom
	db        *sql.DB
	tableName string
}

func (c ConsumerDAO) InsertConsumer(_ *context.ContextModel, tx *sql.Tx, data repository.ConsumerModel) (result int64, err error) {
	query := fmt.Sprintf(`INSERT INTO %s (uuid_id, user_id, nik, full_name, legal_name, birth_place, 
                birth_date, salary, ktp_photo, selfie_photo, created_by, created_client, updated_by, updated_client) 
		VALUES (UUID(), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) `, c.tableName)

	param := []interface{}{data.UserID.Int64, data.NIK.String, data.FullName.String, data.LegalName.String,
		data.BirthPlace.String, data.BirthDate.Time, data.Salary.Float64, data.KTPPhoto.String, data.SelfiePhoto.String,
		data.CreatedBy.Int64, data.CreatedClient.String, data.UpdatedBy.Int64, data.UpdatedClient.String,
	}

	var (
		stmt     *sql.Stmt
		resultDB sql.Result
	)

	stmt, err = tx.Prepare(query)
	if err != nil {
		c.logger.Logger.Error("Error to prepare query", zap.Error(err))
		return
	}
	defer stmt.Close()

	resultDB, err = stmt.Exec(param...)
	if err != nil {
		c.logger.Logger.Error("Error to exec query", zap.Error(err))
		return
	}

	result, _ = resultDB.LastInsertId()
	return
}

func (c ConsumerDAO) CheckNIKConsumer(_ *context.ContextModel, nik string) (result bool, err error) {
	query := fmt.Sprintf(`SELECT nik FROM %s WHERE nik = ? `, c.tableName)

	var nikDb string
	err = c.db.QueryRow(query, nik).Scan(&nikDb)
	if err != nil && err != sql.ErrNoRows {
		c.logger.Logger.Error("Error to query consumer nik", zap.Error(err))
		return
	}

	err = nil
	if !compare.IsEmptyString(nikDb) {
		result = true
		return
	}

	return
}

func (c ConsumerDAO) GetListConsumer(_ *context.ContextModel) (result []interface{}, err error) {
	query := fmt.Sprintf(`SELECT 
		id, uuid_id, user_id, nik, full_name, legal_name, birth_place, birth_date, salary,
		ktp_photo, selfie_photo
	FROM %s `, c.tableName)

	var rows *sql.Rows
	rows, err = c.db.Query(query)
	if err != nil {
		c.logger.Logger.Error("Error to select Query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var temp repository.ConsumerModel

		err = rows.Scan(
			&temp.ID, &temp.UUID, &temp.UserID, &temp.NIK, &temp.FullName, &temp.LegalName, &temp.BirthPlace,
			&temp.BirthDate, &temp.Salary, &temp.KTPPhoto, &temp.SelfiePhoto,
		)
		if err != nil {
			c.logger.Logger.Error("Error to scan query get list", zap.Error(err))
			return
		}

		result = append(result, temp)
	}

	return
}
