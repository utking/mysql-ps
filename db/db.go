package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/utking/mysql-ps/helpers"
)

var (
	db  *sqlx.DB
	err error
)

func ConnectDB(user, password, dsn string) error {
	if db, err = sqlx.Connect("mysql", fmt.Sprintf("%s:%s@%s/sys", user, password, dsn)); err != nil {
		return err
	}

	return nil
}

func GetProcessList() ([]helpers.ProcessItem, error) {
	list := []helpers.ProcessItem{}

	if err := db.Select(&list, "SELECT * FROM INFORMATION_SCHEMA.PROCESSLIST WHERE DB != 'sys' AND COMMAND != 'Sleep' AND USER != 'system user' ORDER BY time DESC"); err != nil {
		return list, err
	}

	return list, nil
}
