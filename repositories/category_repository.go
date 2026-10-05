package repositories

import (
	"errors"

	"pizza-app/database"
	"pizza-app/models"
)

var ErrCategoryNotFound = errors.New("category not found")

func CreateCategory(name string, imageURL string) (models.Category, error) {
	var cat models.Category
	cat.Name = name
	cat.ImageURL = imageURL
	err := database.DB.QueryRow(
		`INSERT INTO categories (name, image_url) VALUES ($1, $2) RETURNING id`, name, imageURL,
	).Scan(&cat.ID)
	if err != nil {
		return models.Category{}, err
	}
	return cat, nil
}

// GetAllCategories includes a live pizza_count per category (via LEFT JOIN
// so empty categories still show up with count 0) — this is what makes
// the pizza↔category relationship visible from the Categories tab, even
// though the actual assignment still happens from each pizza's own Edit
// form (a category can hold many pizzas, so managing it the other way
// around — one dropdown per pizza — stays the simpler mental model).
func GetAllCategories() ([]models.Category, error) {
	query := `
		SELECT c.id, c.name, COALESCE(c.image_url, ''), COUNT(p.id) AS pizza_count
		FROM categories c
		LEFT JOIN pizzas p ON p.category_id = c.id
		GROUP BY c.id, c.name, COALESCE(c.image_url, '')
		ORDER BY c.name
	`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []models.Category{}
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.ImageURL, &c.PizzaCount); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func UpdateCategory(id int, name string, imageURL string) (models.Category, error) {
	res, err := database.DB.Exec(`UPDATE categories SET name = $1, image_url = $2 WHERE id = $3`, name, imageURL, id)
	if err != nil {
		return models.Category{}, err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return models.Category{}, ErrCategoryNotFound
	}
	return models.Category{ID: id, Name: name, ImageURL: imageURL}, nil
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