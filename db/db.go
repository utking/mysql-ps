package db

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/utking/mysql-ps/helpers"
)

type Querier interface {
	Select(dest interface{}, query string, args ...interface{}) error
}

type DBStore struct {
	Db Querier
}

func ConnectDB(user, password, dsn string) (*DBStore, error) {
	conn, err := sqlx.Connect(
		"mysql",
		fmt.Sprintf("%s:%s@%s/sys", user, password, dsn),
	)
	if err != nil {
		return nil, err
	}

	return &DBStore{Db: conn}, nil
}

func (s *DBStore) Close() error {
	type closer interface{ Close() error }
	if c, ok := s.Db.(closer); ok {
		return c.Close()
	}
	return nil
}

func (s *DBStore) GetProcessList(filters []string, databases []interface{}) ([]helpers.ProcessItem, error) {
	list := []helpers.ProcessItem{}
	filterBuilder := strings.Builder{}

	filterBuilder.WriteString(
		"SELECT * FROM INFORMATION_SCHEMA.PROCESSLIST WHERE COMMAND != 'Sleep'",
	)
	filterBuilder.WriteString(
		" AND USER != 'system user'",
	)

	if len(databases) > 0 {
		placeholders := make([]string, len(databases))
		for i := range databases {
			placeholders[i] = "?"
		}
		filterBuilder.WriteString(" AND DB IN (" + strings.Join(placeholders, ",") + ")")
	}

	for _, filter := range filters {
		filterBuilder.WriteString(fmt.Sprintf(" AND %s", filter))
	}

	filterBuilder.WriteString(" ORDER BY time DESC")

	if err := s.Db.Select(
		&list,
		filterBuilder.String(),
		databases...,
	); err != nil {
		return list, err
	}

	return list, nil
}
