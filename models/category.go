package models

import (
	"fmt"
)

type Category struct {
	Id            int    `json:"id"`
	Category_name string `json:"category_name"`
}

// A function that returns a slice containing every category fromt
func GetAllCategories() ([]Category, error) {
	rows, err := db.Query("SELECT c.id, c.category_name FROM category as c order by c.category_name")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.Id, &category.Category_name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil

}
