package helpers

import "os"

const (
	SQLLogName = "./queries.sql"
	sqlLogPerm = 0o644
)

func WriteSQLLog(logLine string, appendLine bool) error {
	var (
		f   *os.File
		err error
	)

	if appendLine {
		f, err = os.OpenFile(SQLLogName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, sqlLogPerm)
	} else {
		f, err = os.OpenFile(SQLLogName, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, sqlLogPerm)
	}

	if err == nil {
		defer f.Close()

		_, _ = f.WriteString(logLine)
	}

	return err
}
