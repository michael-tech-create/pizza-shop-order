package repositories

import (
	"database/sql"
	"pizza-app/database"
	"pizza-app/models"
	"fmt"
	"errors"
)

func CreatePizza(pizza models.Pizza) error {
	query :=  `
	INSERT INTO pizzas (name, price, description)
	VALUES ($1, $2, $3)
	`

	_, err := database.DB.Exec(
		query,
		pizza.Name,
		pizza.Price,
		pizza.Description,
	)

	return err
}

func GetAllPizzas() ([]models.Pizza, error) {
	query := `
	SELECT id, name, price, description
	FROM pizzas
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var pizzas []models.Pizza
	

	for rows.Next() {
		var pizza models.Pizza


		err := rows.Scan(
			&pizza.ID,
			&pizza.Name,
			&pizza.Price,
			&pizza.Description,
		)

	if err != nil {
		return nil, err
	}

	pizzas = append(pizzas, pizza)
	}

	return pizzas, nil
}

func GetPizzaByID(id int) (models.Pizza, error) {
	query := `
	SELECT id, name, price, description
	FROM pizzas

	WHERE id = $1
	`

	var pizza models.Pizza

err := database.DB.QueryRow(query, id).Scan(
	&pizza.ID,
	&pizza.Name,
	&pizza.Price,
	&pizza.Description,
)

if err != nil {
	if errors.Is(err, sql.ErrNoRows) {
	 return models.Pizza{}, fmt.Errorf("pizza with id %d not found",id)

	}
	return models.Pizza{},err
}

return pizza, nil
}


func DeletePizza(id int) (models.Pizza, error) {
	pizza, err := GetPizzaByID(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Pizza{}, fmt.Errorf("pizza not found", id)
		}
		return models.Pizza{}, err
	}

	query := `
	DELETE FROM pizzas
	WHERE id = $1`

	_, err = database.DB.Exec(query, id)
		if err != nil {
			return models.Pizza{}, err
		}

// rowsAffected, err := result.RowsAffected()

// if err != nil {
// 	return nil, err
// }

// if rowsAffected == 0 {
// 	return "no record found matching"
// }

return pizza, nil
}

func CreateOrder(pizzaID int, quantity int) (models.Order, error) {

    pizza, err := GetPizzaByID(pizzaID)
    if err != nil {
        return models.Order{}, err
    }

    total := pizza.Price * quantity

    order := models.Order{
        PizzaId: pizzaID,
        Quantity: quantity,
        TotalCost: total,
    }

    savedOrder, err := SaveOrder(order)
    if err != nil {
        return models.Order{}, err
    }

    return savedOrder, nil
}

func SaveOrder(order models.Order) (models.Order, error) {
	query := `INSERT INTO orders (pizza_id, quantity, total_cost)
	VALUES ($1, $2, $3)
	RETURNING id`

	err := database.DB.QueryRow(
		query,
		order.PizzaId,
		order.Quantity,
		order.TotalCost,
	).Scan(&order.ID)

	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}


func GetAllOrders() ([]models.Order, error) {

	query := `
	SELECT id, pizza_id, quantity, total_cost
	FROM orders
	ORDER BY id DESC
	`

	rows, err := database.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var orders []models.Order

	for rows.Next() {
		var order models.Order

		err := rows.Scan(
			&order.ID,
			&order.PizzaId,
			&order.Quantity,
			&order.TotalCost,
		)

		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}
 return orders, nil 
}

func GetAllOrdersWithPizzaName() ([]models.OrderResponse, error) {

    query := `
        SELECT 
            orders.id,
            pizzas.name,
            orders.quantity,
            orders.total_cost
        FROM orders
        JOIN pizzas ON orders.pizza_id = pizzas.id
        ORDER BY orders.id DESC
    `

    rows, err := database.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var orders []models.OrderResponse

    for rows.Next() {

        var order models.OrderResponse

        err := rows.Scan(
            &order.OrderID,
            &order.PizzaName,
            &order.Quantity,
            &order.TotalCost,
        )

        if err != nil {
            return nil, err
        }

        orders = append(orders, order)
    }

    return orders, nil
}

func UpdateOrderStatus(id int, status string) (error) {
	query := `UPDATE orders
	SET status = $1
	WHERE id = $2`

	_, err := database.DB.Exec(
		query,
		status,
		id,
	)

	switch status {
	case "Pending", "Preparing", "Delivered", "Cancelled":
		//valid
	default:
		return fmt.Errorf("invalid status")
	}

	return err
}


func UpdatePizza(id int, pizza models.Pizza) (models.Pizza, error) {

	query := `
	UPDATE pizzas
	SET name = $1,
	    price = $2,
	    description = $3
	WHERE id = $4
	RETURNING id, name, price, description
	`

	var updated models.Pizza

	err := database.DB.QueryRow(
		query,
		pizza.Name,
		pizza.Price,
		pizza.Description,
		id,
	).Scan(
		&updated.ID,
		&updated.Name,
		&updated.Price,
		&updated.Description,
	)

	if err != nil {
		return models.Pizza{}, err
	}

	return updated, nil
}