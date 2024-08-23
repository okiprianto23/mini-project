package repository

import "database/sql"

type CreditLimitModel struct {
	ID                  sql.NullInt64
	UUID                sql.NullString
	ConsumerID          sql.NullInt64
	MonthlyInstallments sql.NullFloat64
	InterestRate        sql.NullFloat64
	Tenor               sql.NullInt64
	Limit               sql.NullFloat64
	RemainingLimit      sql.NullFloat64
	DefaultCreatedUpdated
}
