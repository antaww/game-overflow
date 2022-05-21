package sql

import (
	"fmt"
)

var locales map[string]string

// GetLocales returns the locales map, first getting from SQL then memoizing it
func GetLocales() (map[string]string, error) {
	if locales != nil {
		return locales, nil
	}

	rows, err := DB.Query("SELECT * FROM locales")
	if err != nil {
		return nil, fmt.Errorf("GetLocales: %s", err)
	}

	var result = make(map[string]string)
	for rows.Next() {
		var locale string
		var value string
		err = rows.Scan(&locale, &value)
		if err != nil {
			return nil, fmt.Errorf("GetLocales: %s", err)
		}
		result[locale] = value
	}

	locales = result
	return result, nil
}
