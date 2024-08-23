package dao

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"main-xyz/config"
	"main-xyz/context"
	"main-xyz/repository"
)

func NewCreditLimitDAO(
	db *sql.DB,
	logger *config.LoggerCustom,
) CreditLimitDAO {
	return CreditLimitDAO{
		db:        db,
		logger:    logger,
		tableName: "credit_limit",
	}
}

type CreditLimitDAO struct {
	db        *sql.DB
	logger    *config.LoggerCustom
	tableName string
}

func (c CreditLimitDAO) InsertCreditLimit(_ *context.ContextModel, tx *sql.Tx, model repository.CreditLimitModel) (result int64, err error) {
	query := fmt.Sprintf(`INSERT INTO %s 
    	(uuid_id, consumer_id, monthly_installments, interest_rate, tenor, limit_amount, remaining_limit_amount,
    	 created_by, created_client, updated_by, updated_client) 
	VALUES 
		(UUID(), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, c.tableName)

	var (
		stmt     *sql.Stmt
		resultDB sql.Result
	)

	params := []interface{}{
		model.ConsumerID.Int64, model.MonthlyInstallments.Float64, model.InterestRate.Float64, model.Tenor.Int64,
		model.Limit.Float64, model.RemainingLimit.Float64, model.CreatedBy.Int64, model.CreatedClient.String,
		model.UpdatedBy.Int64, model.UpdatedClient.String,
	}

	stmt, err = tx.Prepare(query)
	if err != nil {
		c.logger.Logger.Error("Error to prepare credit limit", zap.Error(err))
		return
	}
	defer stmt.Close()

	resultDB, err = stmt.Exec(params...)
	if err != nil {
		c.logger.Logger.Error("Error to exec query credit limit", zap.Error(err))
		return
	}

	result, _ = resultDB.LastInsertId()
	return
}

func (c CreditLimitDAO) GetCreditLimitByConsumerID(_ *context.ContextModel, consumerID int64) (result []interface{}, err error) {
	query := fmt.Sprintf(`SELECT 
		limit_id, uuid_id, consumer_id, tenor, limit_amount, created_by, 
		created_client, created_at, updated_by, updated_client, updated_at,
		monthly_installments, interest_rate
	FROM %s 
	WHERE consumer_id = ? `, c.tableName)

	rows, errS := c.db.Query(query, consumerID)
	if errS != nil {
		err = errS
		c.logger.Logger.Error("Error to query credit limit", zap.Error(err))
		return
	}

	defer rows.Close()

	for rows.Next() {
		temp := repository.CreditLimitModel{}
		err = rows.Scan(
			&temp.ID,
			&temp.UUID,
			&temp.ConsumerID,
			&temp.Tenor,
			&temp.Limit,
			&temp.CreatedBy,
			&temp.CreatedClient,
			&temp.CreatedAt,
			&temp.UpdatedBy,
			&temp.UpdatedClient,
			&temp.UpdatedAt,
			&temp.MonthlyInstallments,
			&temp.InterestRate,
		)
		if err != nil {
			c.logger.Logger.Error("Error to scan credit limit", zap.Error(err))
			return
		}

		result = append(result, temp)
	}

	return
}

func (c CreditLimitDAO) UpdateCreditLimitRemaining(_ *context.ContextModel, tx *sql.Tx, model repository.CreditLimitModel) error {
	query := fmt.Sprintf(`UPDATE %s SET remaining_limit_amount = ? WHERE limit_id = ? AND consumer_id = ?`, c.tableName)

	_, err := tx.Exec(query, model.RemainingLimit.Float64, model.ID.Int64, model.ConsumerID.Int64)
	if err != nil {
		c.logger.Logger.Error("Error to exec query", zap.Error(err))
		return err
	}

	return nil
}
