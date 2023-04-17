package helpers

import (
	"database/sql"
)

type ProcessItem struct {
	Host         sql.NullString `db:"HOST"`
	User         sql.NullString `db:"USER"`
	Command      sql.NullString `db:"COMMAND"`
	State        sql.NullString `db:"STATE"`
	DB           sql.NullString `db:"DB"`
	Info         sql.NullString `db:"INFO"`
	ID           int64          `db:"ID"`
	Time         int64          `db:"TIME"`
	TimeMs       int64          `db:"TIME_MS"`
	RowsSent     int64          `db:"ROWS_SENT"`
	RowsExamined int64          `db:"ROWS_EXAMINED"`
}
