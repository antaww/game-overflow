package sql

import (
	"database/sql"
	"log"
)

func Select(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	defer func(rows *sql.Rows) {
		err := rows.Err()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	return rows, nil
}

func Results(rows *sql.Rows, dest ...interface{}) error {
	for rows.Next() {
		err := rows.Scan(dest...)
		if err != nil {
			return err
		}
	}
	return nil
}
