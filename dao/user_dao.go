package dao

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"main-xyz/config"
	"main-xyz/context"
	"main-xyz/repository"
)

func NewRepoUser(db *sql.DB, logger *config.LoggerCustom) UserDAO {
	return UserDAO{
		db:        db,
		tableName: "users",
		logger:    logger,
	}
}

type UserDAO struct {
	db        *sql.DB
	tableName string
	logger    *config.LoggerCustom
}

func (u UserDAO) GetUserByUsername(username string) (*repository.UserAdmin, error) {
	row := u.db.QueryRow("SELECT username, password_hash FROM users WHERE username = ?", username)
	var user repository.UserAdmin
	err := row.Scan(&user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // user tidak ditemukan
		}
		return nil, err
	}
	return &user, nil
}

func (u UserDAO) GetUserByID(id int64) (repository.UserAdmin, error) {
	row := u.db.QueryRow("SELECT id, client_id FROM users WHERE id = ?", id)
	var user repository.UserAdmin
	err := row.Scan(&user.ID, &user.ClientID)
	if err != nil && err != sql.ErrNoRows {
		return user, err
	}
	return user, nil
}

func (u UserDAO) CheckIsUserAdminByUsername(_ *context.ContextModel, userParam repository.UserAdmin) (result repository.UserAdmin, err error) {
	query := fmt.Sprintf(`SELECT 
    		u.id, u.client_id, u.username, u.password_hash, u.signature_key, u.locale, u.alias_name, u.client_alias
		FROM %s u
		WHERE (u.username = ? OR u.email = ?)
	`, u.tableName)

	param := []interface{}{
		userParam.Username.String,
		userParam.Username.String,
	}

	err = u.db.QueryRow(query, param...).Scan(
		&result.UserID, &result.ClientID, &result.Username, &result.Password, &result.SignatureKey,
		&result.Locale, &result.AliasName, &result.ClientAlias,
	)
	if err != nil && err != sql.ErrNoRows {
		u.logger.Logger.Error("Error to get query", zap.Error(err))
		return
	}

	err = nil
	return
}

func (u UserDAO) GetInformationUser(_ *context.ContextModel, clientID string) (result repository.UserInformation, err error) {
	query := fmt.Sprintf(`
		SELECT u.alias_name, u.username, u.id, u.locale, u.client_id
		FROM %s u
		WHERE u.client_id = ?
		`, u.tableName)

	err = u.db.QueryRow(query, clientID).Scan(
		&result.AliasName, &result.Username, &result.ResourceUserID, &result.Locale,
		&result.ClientID,
	)

	if err != nil && err != sql.ErrNoRows {
		u.logger.Logger.Error("Error to query select", zap.Error(err))
		return
	}

	err = nil
	return
}

func (u *UserDAO) GetListUser(ctx *context.ContextModel) (result []interface{}, err error) {
	query := fmt.Sprintf(`SELECT 
		id, username, email, auth_user_id, client_id, locale, alias_name
	FROM %s `, u.tableName)

	var rows *sql.Rows
	rows, err = u.db.Query(query)
	if err != nil {
		u.logger.Logger.Error("Error to select Query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var temp repository.UserInformation

		err = rows.Scan(
			&temp.ID, &temp.Username, &temp.Email, &temp.AuthUserID, &temp.ClientID,
			&temp.Locale, &temp.AliasName,
		)
		if err != nil {
			u.logger.Logger.Error("Error to scan query get list", zap.Error(err))
			return
		}

		result = append(result, temp)
	}

	return
}
