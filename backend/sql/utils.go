package sql

import (
	"database/sql"
	"main/utils"
)

// HandleSQLErrors closes the connection to the database and logs the error if any
func HandleSQLErrors(rows *sql.Rows) {
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			utils.SQLError(err)
		}
	}(rows)

	defer func(rows *sql.Rows) {
		err := rows.Err()
		if err != nil {
			utils.SQLError(err)
		}
	}(rows)
}

// Results parse the first Row of results, and put the values into dest parameters
func Results(rows *sql.Rows, dest ...interface{}) error {
	if rows.Next() {
		err := rows.Scan(dest...)
		if err != nil {
			return err
		}
	}
	return nil
}

func contains(tags []Tags, tag Tags) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}
