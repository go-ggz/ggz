package model

import (
	"github.com/go-xorm/xorm"
	"github.com/go-xorm/xorm/migrate"
)

// Migrations for db migrate
var migrations = []*migrate.Migration{
	{
		ID: "201709201400",
		Migrate: func(tx *xorm.Engine) error {
			return tx.Sync2(&User{})
		},
		Rollback: func(tx *xorm.Engine) error {
			return tx.DropTables(&User{})
		},
	},
	{
		ID: "201711181402",
		Migrate: func(tx *xorm.Engine) error {
			// drop column and ignore error.
			tx.Exec("ALTER TABLE user DROP COLUMN passwd")
			tx.Exec("ALTER TABLE user DROP COLUMN user_name")

			return nil
		},
		Rollback: func(tx *xorm.Engine) error {
			return nil
		},
	},
}
