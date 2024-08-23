package repository

import "database/sql"

type ClientTokenModel struct {
	ID            sql.NullInt64
	ClientID      sql.NullString
	UserID        sql.NullInt64
	Token         sql.NullString
	ExpiredAt     sql.NullTime
	CreatedAt     sql.NullTime
	IsActive      sql.NullBool
	CreatedBy     sql.NullInt64
	CreatedClient sql.NullString
	Deleted       sql.NullBool
}
