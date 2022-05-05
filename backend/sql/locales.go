package sql

import (
	"fmt"
)

func GetLocales() (map[string]string, error) {
	locales, err := DB.Query("SELECT * FROM locales")
	if err != nil {
		return nil, fmt.Errorf("GetLocales: %s", err)
	}

	var result = make(map[string]string)
	for locales.Next() {
		var locale string
		var value string
		err = locales.Scan(&locale, &value)
		if err != nil {
			return nil, fmt.Errorf("GetLocales: %s", err)
		}
	}
	return result, nil
}
