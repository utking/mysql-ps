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
	if db, err = sqlx.Connect(
		"mysql",
		fmt.Sprintf("%s:%s@%s/sys", user, password, dsn),
	); err != nil {
		return err
	}

	return nil
}

func GetProcessList(filters []string, databases []interface{}) ([]helpers.ProcessItem, error) {
	list := []helpers.ProcessItem{}
	filterBuilder := strings.Builder{}

	filterBuilder.WriteString(
		"SELECT * FROM INFORMATION_SCHEMA.PROCESSLIST WHERE COMMAND != 'Sleep'",
	)
	filterBuilder.WriteString(
		" AND USER != 'system user'",
	)

	if len(databases) > 0 {
		filterBuilder.WriteString(" AND DB IN (")
		// every database must be added to string as ? placeholder
		// and separated from the previous with a comma.
		// this is done to avoid SQL injection
		for i := range databases {
			filterBuilder.WriteString("?")
			if i < len(databases)-1 {
				filterBuilder.WriteString(",")
			}
		}
		filterBuilder.WriteString(")")
	}

	for _, filter := range filters {
		filterBuilder.WriteString(fmt.Sprintf(" AND %s", filter))
	}

	filterBuilder.WriteString(" ORDER BY time DESC")

	if err := db.Select(
		&list,
		filterBuilder.String(),
		databases...,
	); err != nil {
		return list, err
	}

	return list, nil
}
