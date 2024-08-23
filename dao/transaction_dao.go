package dao

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"main-xyz/config"
	"main-xyz/context"
	"main-xyz/repository"
)

func NewTransactionDAO(
	db *sql.DB,
	logger *config.LoggerCustom,
) TransactionDAO {
	return TransactionDAO{
		db:        db,
		logger:    logger,
		tableName: "transaction",
	}
}

type TransactionDAO struct {
	db        *sql.DB
	logger    *config.LoggerCustom
	tableName string
}

func (t TransactionDAO) InsertTransaction(_ *context.ContextModel, tx *sql.Tx, data repository.TransactionModel) (result int64, err error) {
	query := fmt.Sprintf(`INSERT INTO %s (uuid_id, consumer_id, limit_id, contract_number, otr, admin_fee, installment_amount, 
                interest_amount, asset_name, transaction_date, created_by, created_client, updated_by, updated_client) 
		VALUES (UUID(), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) `, t.tableName)

	param := []interface{}{data.ConsumerID.Int64, data.LimitID.Int64, data.ContractNumber.String, data.OTR.Float64, data.AdminFee.Float64,
		data.InstallmentAmount.Float64, data.InterestAmount.Float64, data.AssetName.String, data.TransactionDate.Time,
		data.CreatedBy.Int64, data.CreatedClient.String, data.UpdatedBy.Int64, data.UpdatedClient.String,
	}

	var (
		stmt     *sql.Stmt
		resultDB sql.Result
	)

	stmt, err = tx.Prepare(query)
	if err != nil {
		t.logger.Logger.Error("Error to prepare query", zap.Error(err))
		return
	}
	defer stmt.Close()

	resultDB, err = stmt.Exec(param...)
	if err != nil {
		t.logger.Logger.Error("Error to exec query", zap.Error(err))
		return
	}

	result, _ = resultDB.LastInsertId()
	return
}
