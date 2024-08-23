package repository

import "database/sql"

type ConsumerModel struct {
	ID          sql.NullInt64
	UUID        sql.NullString
	UserID      sql.NullInt64
	NIK         sql.NullString
	FullName    sql.NullString
	LegalName   sql.NullString
	BirthPlace  sql.NullString
	BirthDate   sql.NullTime
	Salary      sql.NullFloat64
	KTPPhoto    sql.NullString
	SelfiePhoto sql.NullString
	DefaultCreatedUpdated
}
