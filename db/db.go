package db

import (
	"fmt"
	"strings"

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

func GetProcessList(filters []string) ([]helpers.ProcessItem, error) {
	list := []helpers.ProcessItem{}
	filterBuilder := strings.Builder{}

	filterBuilder.WriteString("SELECT * FROM INFORMATION_SCHEMA.PROCESSLIST WHERE COMMAND != 'Sleep'")
	filterBuilder.WriteString(" AND USER != 'system user'")

	for _, filter := range filters {
		filterBuilder.WriteString(fmt.Sprintf(" AND %s", filter))
	}

	filterBuilder.WriteString(" ORDER BY time DESC")

	if err := db.Select(&list, filterBuilder.String()); err != nil {
		return list, err
	}

	return list, nil
}
