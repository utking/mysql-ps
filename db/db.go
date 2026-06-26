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

type Filter struct {
	Column   string
	Operator string
	Value    string
}

var allowedColumns = map[string]bool{
	"ID": true, "USER": true, "HOST": true, "DB": true,
	"COMMAND": true, "TIME": true, "STATE": true, "INFO": true,
}

type DBStore struct {
	db Querier
}

func ConnectDB(user, password, dsn string) (*DBStore, error) {
	conn, err := sqlx.Connect(
		"mysql",
		fmt.Sprintf("%s:%s@%s/sys", user, password, dsn),
	)
	if err != nil {
		return nil, err
	}

	return &DBStore{db: conn}, nil
}

func (s *DBStore) Close() error {
	type closer interface{ Close() error }
	if c, ok := s.db.(closer); ok {
		return c.Close()
	}
	return nil
}

func (s *DBStore) GetProcessList(filters []Filter, databases []interface{}) ([]helpers.ProcessItem, error) {
	list := []helpers.ProcessItem{}
	query := strings.Builder{}
	args := []interface{}{}

	query.WriteString(
		"SELECT * FROM INFORMATION_SCHEMA.PROCESSLIST WHERE COMMAND != 'Sleep'",
	)
	query.WriteString(
		" AND USER != 'system user'",
	)

	if len(databases) > 0 {
		placeholders := make([]string, len(databases))
		for i := range databases {
			placeholders[i] = "?"
		}
		query.WriteString(" AND DB IN (" + strings.Join(placeholders, ",") + ")")
		args = append(args, databases...)
	}

	for _, f := range filters {
		if !allowedColumns[strings.ToUpper(f.Column)] {
			continue
		}
		query.WriteString(fmt.Sprintf(" AND %s %s ?", f.Column, f.Operator))
		args = append(args, f.Value)
	}

	query.WriteString(" ORDER BY time DESC")

	if err := s.db.Select(&list, query.String(), args...); err != nil {
		return list, err
	}

	return list, nil
}
