package utils

import (
	"database/sql"
	"errors"
	"github.com/gobuffalo/packr/v2"
	migrate "github.com/rubenv/sql-migrate"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type DBMigration interface {
	DbMigrateMysql(
		db *sql.DB,
		fileMigratePath string,
		direction migrate.MigrationDirection,
	) (
		int,
		error,
	)
}

func NewDBMigration() DBMigration {
	return &dbMigration{}
}

type dbMigration struct {
	resolutionDir string
}

func (dbMigrate *dbMigration) DbMigrateMysql(
	db *sql.DB,
	fileMigratePath string,
	direction migrate.MigrationDirection,
) (
	int,
	error,
) {
	migrations := migrate.PackrMigrationSource{
		Box: packr.New("migrations", fileMigratePath),
	}

	if db != nil {
		n, err := migrate.Exec(db, "mysql", migrations, direction)
		if err != nil {
			return 0, err
		} else {
			if dbMigrate.resolutionDir == "" {
				box := reflect.Indirect(reflect.ValueOf(migrations)).FieldByName("Box")
				resolution := reflect.Indirect(box).Interface().(*packr.Box)
				dir := strings.Replace(resolution.ResolutionDir, "\\", "/", -1)
				splitData := strings.Split(dir, "/")
				dbMigrate.resolutionDir = strings.Join(splitData[0:len(splitData)-1], "/")
			}
			return n, err
		}
	} else {
		return 0, errors.New("null database")
	}
}
