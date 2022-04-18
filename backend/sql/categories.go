package sql

type Category struct {
	Name string `db:"category_name"`
	Icon string `db:"icon"`
}

// GetCategories returns all categories
func GetCategories() ([]Category, error) {
	var categories []Category
	row, err := DB.Query("SELECT * FROM categories ORDER BY category_name")
	if err != nil {
		return nil, err
	}

	for row.Next() {
		var category Category
		err = row.Scan(&category.Name, &category.Icon)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	HandleSQLErrors(row)

	return categories, nil
}
