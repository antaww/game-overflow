package sql

import "sort"

type TagListItem struct {
	Name  string `json:"name" db:"tag_name"`
	Count int
}

// GetTrendingTags returns the trending tags with number of usages, sorted by count, limited by the limit
func GetTrendingTags(limit int) ([]TagListItem, error) {
	var tags []TagListItem
	rows, err := DB.Query("SELECT tag_name, Count(tag_name) as count FROM have GROUP BY tag_name HAVING count ORDER BY count DESC LIMIT ?;", limit)
	if err != nil {
		return tags, err
	}

	for rows.Next() {
		var tag TagListItem
		err = rows.Scan(&tag.Name, &tag.Count)
		if err != nil {
			return tags, err
		}
		tags = append(tags, tag)
	}

	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Count > tags[j].Count
	})

	return tags, nil
}
