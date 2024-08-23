package repository

import "database/sql"

type UserAdmin struct {
	ID       sql.NullInt64
	Username sql.NullString
	Password sql.NullString

	UserID       sql.NullInt64
	ClientID     sql.NullString
	RoleName     sql.NullString
	Permissions  sql.NullString
	GroupName    sql.NullString
	DataScope    sql.NullString
	IPWhitelist  sql.NullString
	SignatureKey sql.NullString
	Locale       sql.NullString
	AliasName    sql.NullString
	ClientAlias  sql.NullString
	DefaultCreatedUpdated
}

type UserInformation struct {
	ID             sql.NullInt64
	AliasName      sql.NullString
	Username       sql.NullString
	Locale         sql.NullString
	ResourceUserID sql.NullInt64
	ClientID       sql.NullString
	Email          sql.NullString
	AuthUserID     sql.NullInt64
}
