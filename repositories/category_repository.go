package repositories

import (
	"errors"

	"pizza-app/database"
	"pizza-app/models"
)

var ErrCategoryNotFound = errors.New("category not found")

func CreateCategory(name string) (models.Category, error) {
	var cat models.Category
	cat.Name = name
	err := database.DB.QueryRow(
		`INSERT INTO categories (name) VALUES ($1) RETURNING id`, name,
	).Scan(&cat.ID)
	if err != nil {
		return models.Category{}, err
	}
	return cat, nil
}

func GetAllCategories() ([]models.Category, error) {
	rows, err := database.DB.Query(`SELECT id, name FROM categories ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []models.Category{}
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func UpdateCategory(id int, name string) (models.Category, error) {
	res, err := database.DB.Exec(`UPDATE categories SET name = $1 WHERE id = $2`, name, id)
	if err != nil {
		return models.Category{}, err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return models.Category{}, ErrCategoryNotFound
	}
	return models.Category{ID: id, Name: name}, nil
}

// DeleteCategory removes a category. Pizzas referencing it aren't
// deleted — the migration's ON DELETE SET NULL means they simply
// become uncategorized instead of disappearing from the menu.
func DeleteCategory(id int) error {
	res, err := database.DB.Exec(`DELETE FROM categories WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrCategoryNotFound
	}
	return nil
}