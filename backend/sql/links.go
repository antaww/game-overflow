package sql

import "fmt"

type Link struct {
	IdUser int64  `db:"id_user"`
	Name   string `db:"name"`
	Link   string `db:"link"`
}

// GetLinks returns all links of a user
func GetLinks(idUser int64) ([]Link, error) {
	rows, err := DB.Query("SELECT * FROM links WHERE id_user = ?", idUser)
	if err != nil {
		return nil, err
	}

	var links []Link
	for rows.Next() {
		var link Link
		err = rows.Scan(&link.IdUser, &link.Name, &link.Link)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	return links, nil
}

// SetLinks sets all links of a user, overwriting all previous links if any, a user has a maximum of 5 links
func SetLinks(idUser int64, links []Link) error {
	_, err := DB.Exec("DELETE FROM links WHERE id_user = ?", idUser)
	if err != nil {
		return err
	}

	request := "INSERT INTO links (id_user, name, link) VALUES "
	for i, link := range links {
		if i != 0 {
			request += ", "
		}
		request += fmt.Sprintf("(%d, '%s', '%s')", idUser, link.Name, link.Link)
	}

	_, err = DB.Exec(request)
	if err != nil {
		return err
	}

	return nil
}
