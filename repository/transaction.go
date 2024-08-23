package repository

import "database/sql"

type TransactionModel struct {
	ID                sql.NullInt64
	UUID              sql.NullString
	ConsumerID        sql.NullInt64
	LimitID           sql.NullInt64
	ContractNumber    sql.NullString
	OTR               sql.NullFloat64
	AdminFee          sql.NullFloat64
	InstallmentAmount sql.NullFloat64
	InterestAmount    sql.NullFloat64
	AssetName         sql.NullString
	TransactionDate   sql.NullTime
	DefaultCreatedUpdated
}
