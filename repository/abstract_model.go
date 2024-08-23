package repository

import "database/sql"

type DefaultCreatedUpdated struct {
	UpdatedBy     sql.NullInt64
	UpdatedAt     sql.NullTime
	UpdatedClient sql.NullString
	CreatedBy     sql.NullInt64
	CreatedAt     sql.NullTime
	CreatedClient sql.NullString
	CreatedByName sql.NullString
	UpdatedByName sql.NullString
	Deleted       sql.NullBool
}
