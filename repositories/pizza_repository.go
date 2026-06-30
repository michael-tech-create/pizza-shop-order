package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"pizza-app/database"
	"pizza-app/models"
)


// PIZZA CRUD FUNCTIONS


func CreatePizza(pizza models.Pizza) error {
	query := `
	INSERT INTO pizzas (name, price, description)
	VALUES ($1, $2, $3)
	`
	_, err := database.DB.Exec(query, pizza.Name, pizza.Price, pizza.Description)
	return err
}

func GetAllPizzas() ([]models.Pizza, error) {
	query := `SELECT id, name, price, description FROM pizzas`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pizzas []models.Pizza
	for rows.Next() {
		var pizza models.Pizza
		err := rows.Scan(&pizza.ID, &pizza.Name, &pizza.Price, &pizza.Description)
		if err != nil {
			return nil, err
		}

		images, err := GetPizzaImages(pizza.ID)
		if err == nil {
			pizza.Images = images
		}
		pizzas = append(pizzas, pizza)
	}
	return pizzas, nil
}

func GetPizzaByID(id int) (models.Pizza, error) {
	query := `SELECT id, name, price, description FROM pizzas WHERE id = $1`
	var pizza models.Pizza

	err := database.DB.QueryRow(query, id).Scan(&pizza.ID, &pizza.Name, &pizza.Price, &pizza.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Pizza{}, fmt.Errorf("pizza with id %d not found", id)
		}
		return models.Pizza{}, err
	}
	return pizza, nil
}

func UpdatePizza(id int, pizza models.Pizza) (models.Pizza, error) {
	query := `
	UPDATE pizzas 
	SET name = $1, price = $2, description = $3
	WHERE id = $4
	`
	_, err := database.DB.Exec(query, pizza.Name, pizza.Price, pizza.Description, id)
	if err != nil {
		return models.Pizza{}, err
	}
	pizza.ID = id
	return pizza, nil
}

func DeletePizza(id int) (models.Pizza, error) {
	pizza, err := GetPizzaByID(id)
	if err != nil {
		return models.Pizza{}, err
	}

	query := `DELETE FROM pizzas WHERE id = $1`
	_, err = database.DB.Exec(query, id)
	if err != nil {
		return models.Pizza{}, err
	}
	return pizza, nil
}

func SearchPizza(queryStr string) ([]models.Pizza, error) {
	query := `
		SELECT id, name, price, description 
		FROM pizzas
		WHERE LOWER(name) LIKE LOWER($1) OR LOWER(description) LIKE LOWER($1)
	`
	rows, err := database.DB.Query(query, "%"+queryStr+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pizzas []models.Pizza
	for rows.Next() {
		var p models.Pizza
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Description)
		if err != nil {
			return nil, err
		}

		images, err := GetPizzaImages(p.ID)
		if err == nil {
			p.Images = images
		}
		pizzas = append(pizzas, p)
	}
	return pizzas, nil
}


// PIZZA IMAGES RELATIONSHIP FUNCTIONS


func GetPizzaImages(pizzaID int) ([]models.PizzaImage, error) {
	query := `SELECT id, pizza_id, image_url FROM pizza_images WHERE pizza_id = $1`
	rows, err := database.DB.Query(query, pizzaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []models.PizzaImage
	for rows.Next() {
		var img models.PizzaImage
		if err := rows.Scan(&img.ID, &img.PizzaID, &img.ImageURL); err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, nil
}

func SavePizzaImage(pizzaID int, imageURL string) error {
	query := `INSERT INTO pizza_images (pizza_id, image_url) VALUES ($1, $2)`
	_, err := database.DB.Exec(query, pizzaID, imageURL)
	return err
}

// ORDERS LOGIC INTERFACES


// CreateOrder inserts a new order and returns the full record including the
// DB-assigned id and the default status ("pending" set by column DEFAULT).
// ORDERS LOGIC INTERFACES


func CreateOrder(customerName string, phone string, address string, items []models.OrderItem) (models.Order, error) {
    tx, err := database.DB.Begin()
    if err != nil {
        return models.Order{}, err
    }
    defer tx.Rollback() // Rollback if any step fails

    // 1. Insert main order
    var orderID int
    err = tx.QueryRow(`INSERT INTO orders (customer_name, phone, address, total_cost, status) 
                       VALUES ($1, $2, $3, 0, 'pending') RETURNING id`, 
                       customerName, phone, address).Scan(&orderID)
    if err != nil {
        return models.Order{}, err
    }

    // 2. Insert items and sum the total
    var totalCost int
    for _, item := range items {
        var price int
        err := tx.QueryRow("SELECT price FROM pizzas WHERE id = $1", item.PizzaID).Scan(&price)
        if err != nil {
            return models.Order{}, err
        }
        
        subTotal := price * item.Quantity
        totalCost += subTotal
        
        _, err = tx.Exec("INSERT INTO order_items (order_id, pizza_id, quantity, sub_total) VALUES ($1, $2, $3, $4)", 
                          orderID, item.PizzaID, item.Quantity, subTotal)
        if err != nil {
            return models.Order{}, err
        }
    }

    // 3. Update the final total cost
    _, err = tx.Exec("UPDATE orders SET total_cost = $1 WHERE id = $2", totalCost, orderID)
    if err != nil {
        return models.Order{}, err
    }

    return models.Order{ID: orderID, TotalCost: totalCost}, tx.Commit()
}

// OrderResponse including status — required by the frontend order table.
func GetAllOrdersWithPizzaName() ([]models.OrderResponse, error) {
	// This joins orders -> order_items -> pizzas
	query := `
		SELECT o.id, o.customer_name, p.name, oi.quantity, o.total_cost, o.status
		FROM orders o
		JOIN order_items oi ON o.id = oi.order_id
		JOIN pizzas p ON oi.pizza_id = p.id
		ORDER BY o.id DESC
	`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.OrderResponse
	for rows.Next() {
		var o models.OrderResponse
		err := rows.Scan(&o.OrderID, &o.CustomerName, &o.PizzaName, &o.Quantity, &o.TotalCost, &o.Status)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

// UpdateOrderStatus accepts: pending | preparing | delivered | cancelled

func UpdateOrderStatus(id int, status string) error {
	query := `UPDATE orders SET status = $1 WHERE id = $2`
	_, err := database.DB.Exec(query, status, id)
	return err
}


// DASHBOARD METRICS AND STATISTICS


func GetDashboardStats() (models.DashboardStats, error) {
	var stats models.DashboardStats

	revenueQuery := `SELECT COALESCE(SUM(total_cost), 0) FROM orders WHERE status = 'delivered'`
	err := database.DB.QueryRow(revenueQuery).Scan(&stats.Revenue)
	if err != nil {
		return stats, err
	}

	countQuery := `
		SELECT 
			COUNT(*),
			COUNT(*) FILTER (WHERE status = 'pending'),
			COUNT(*) FILTER (WHERE status = 'delivered'),
			COUNT(*) FILTER (WHERE status = 'cancelled')
		FROM orders
	`
	err = database.DB.QueryRow(countQuery).Scan(
		&stats.TotalOrders,
		&stats.PendingOrders,
		&stats.DeliveredOrders,
		&stats.CancelledOrders,
	)
	if err != nil {
		return stats, err
	}
	return stats, nil
}

func GetBestSellingPizza() (models.BestSellingPizza, error) {
	var best models.BestSellingPizza
	query := `
		SELECT p.name, COALESCE(SUM(o.quantity), 0) as total_sold
		FROM orders o
		JOIN pizzas p ON o.pizza_id = p.id
		GROUP BY p.name
		ORDER BY total_sold DESC
		LIMIT 1
	`
	err := database.DB.QueryRow(query).Scan(&best.PizzaName, &best.Sold)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.BestSellingPizza{PizzaName: "No sales recorded", Sold: 0}, nil
		}
		return best, err
	}
	return best, nil
}