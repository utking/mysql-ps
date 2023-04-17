package helpers

import (
	"database/sql"
)

type ProcessItem struct {
	Host         string         `db:"HOST"`
	User         string         `db:"USER"`
	Command      string         `db:"COMMAND"`
	State        string         `db:"STATE"`
	DB           sql.NullString `db:"DB"`
	Info         sql.NullString `db:"INFO"`
	ID           int64          `db:"ID"`
	Time         int64          `db:"TIME"`
	TimeMs       int64          `db:"TIME_MS"`
	RowsSent     int64          `db:"ROWS_SENT"`
	RowsExamined int64          `db:"ROWS_EXAMINED"`
}
