package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"pizza-app/database"
	"pizza-app/models"
	"strings"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// PIZZA CRUD FUNCTIONS

func CreatePizza(pizza models.Pizza) error {
	query := `
	INSERT INTO pizzas (name, price_small, price_medium, price_large, description, category_id)
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := database.DB.Exec(query, pizza.Name, pizza.PriceSmall, pizza.PriceMedium, pizza.PriceLarge, pizza.Description, pizza.CategoryID)
	return err
}

func GetAllPizzas() ([]models.Pizza, error) {
	query := `
		SELECT p.id, p.name, p.price_small, p.price_medium, p.price_large, p.description,
		       p.category_id, COALESCE(c.name, '') AS category_name
		FROM pizzas p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.is_active
	`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}

	var pizzas []models.Pizza
	var ids []int
	for rows.Next() {
		var pizza models.Pizza
		err := rows.Scan(&pizza.ID, &pizza.Name, &pizza.PriceSmall, &pizza.PriceMedium, &pizza.PriceLarge, &pizza.Description,
			&pizza.CategoryID, &pizza.CategoryName)
		if err != nil {
			rows.Close()
			return nil, err
		}
		pizzas = append(pizzas, pizza)
		ids = append(ids, pizza.ID)
	}
	rows.Close()

	// One extra round trip for ALL pizzas' images combined, instead of
	// one round trip PER pizza (was N+1 — a menu of 10 pizzas meant 11
	// total queries; now it's always exactly 2, regardless of menu size).
	imagesByPizza, err := getImagesForPizzaIDs(ids)
	if err != nil {
		return nil, err
	}
	for i := range pizzas {
		pizzas[i].Images = imagesByPizza[pizzas[i].ID]
	}

	return pizzas, nil
}

func GetPizzaByID(id int) (models.Pizza, error) {
	query := `
		SELECT p.id, p.name, p.price_small, p.price_medium, p.price_large, p.description,
		       p.category_id, COALESCE(c.name, '') AS category_name
		FROM pizzas p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1 AND p.is_active
	`
	var pizza models.Pizza

	err := database.DB.QueryRow(query, id).Scan(&pizza.ID, &pizza.Name, &pizza.PriceSmall, &pizza.PriceMedium, &pizza.PriceLarge, &pizza.Description,
		&pizza.CategoryID, &pizza.CategoryName)
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
	SET name = $1, price_small = $2, price_medium = $3, price_large = $4, description = $5, category_id = $6
	WHERE id = $7
	`
	_, err := database.DB.Exec(query, pizza.Name, pizza.PriceSmall, pizza.PriceMedium, pizza.PriceLarge, pizza.Description, pizza.CategoryID, id)
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

	var orderCount int
	if err := database.DB.QueryRow(`SELECT COUNT(*) FROM order_items WHERE pizza_id = $1`, id).Scan(&orderCount); err != nil {
		return models.Pizza{}, err
	}

	// Past orders store pizza_id and look up the name with a join. Removing the
	// row would either fail the foreign key or wipe those line items, so a pizza
	// that has been ordered is taken off the menu instead of deleted.
	if orderCount > 0 {
		_, err = database.DB.Exec(`UPDATE pizzas SET is_active = false WHERE id = $1`, id)
		if err != nil {
			return models.Pizza{}, err
		}
		return pizza, nil
	}

	_, err = database.DB.Exec(`DELETE FROM pizzas WHERE id = $1`, id)
	if err != nil {
		return models.Pizza{}, err
	}
	return pizza, nil
}

func SearchPizza(queryStr string) ([]models.Pizza, error) {
	query := `
        SELECT p.id, p.name, p.price_small, p.price_medium, p.price_large, p.description,
               p.category_id, COALESCE(c.name, '') AS category_name
        FROM pizzas p
        LEFT JOIN categories c ON p.category_id = c.id
        WHERE p.is_active AND (LOWER(p.name) LIKE LOWER($1) OR LOWER(p.description) LIKE LOWER($1))
    `
	rows, err := database.DB.Query(query, "%"+queryStr+"%")
	if err != nil {
		return nil, err
	}

	pizzas := []models.Pizza{}
	var ids []int

	for rows.Next() {
		var p models.Pizza
		err := rows.Scan(&p.ID, &p.Name, &p.PriceSmall, &p.PriceMedium, &p.PriceLarge, &p.Description,
			&p.CategoryID, &p.CategoryName)
		if err != nil {
			rows.Close()
			return nil, err
		}
		pizzas = append(pizzas, p)
		ids = append(ids, p.ID)
	}
	rows.Close()

	imagesByPizza, err := getImagesForPizzaIDs(ids)
	if err != nil {
		return nil, err
	}
	for i := range pizzas {
		pizzas[i].Images = imagesByPizza[pizzas[i].ID]
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

// getImagesForPizzaIDs fetches images for MANY pizzas in a single query,
// used by GetAllPizzas/SearchPizza to avoid firing one query per pizza
// (which was an N+1 query bug — see comments at those call sites).
func getImagesForPizzaIDs(pizzaIDs []int) (map[int][]models.PizzaImage, error) {
	result := make(map[int][]models.PizzaImage)
	if len(pizzaIDs) == 0 {
		return result, nil
	}

	query := `SELECT id, pizza_id, image_url FROM pizza_images WHERE pizza_id = ANY($1)`
	rows, err := database.DB.Query(query, pq.Array(pizzaIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var img models.PizzaImage
		if err := rows.Scan(&img.ID, &img.PizzaID, &img.ImageURL); err != nil {
			return nil, err
		}
		result[img.PizzaID] = append(result[img.PizzaID], img)
	}
	return result, nil
}

func SavePizzaImage(pizzaID int, imageURL string) error {
	query := `INSERT INTO pizza_images (pizza_id, image_url) VALUES ($1, $2)`
	_, err := database.DB.Exec(query, pizzaID, imageURL)
	return err
}

// ORDERS LOGIC INTERFACES

// ErrPizzaNotFound means an item in the order references a pizza_id that
// doesn't exist in this database — e.g. a stale cart from before a menu
// change, or from a different environment. This is the customer's cart
// being out of date, not a server malfunction, so handlers should turn
// this into a 400, not a 500.
var ErrPizzaNotFound = errors.New("one or more items in your cart no longer exist — please refresh the menu and try again")

func CreateOrder(customerName string, phone string, address string, items []models.OrderItem, customerID *int) (models.Order, error) {
	tx, err := database.DB.Begin()
	if err != nil {
		return models.Order{}, err
	}
	defer tx.Rollback() // Rollback if any step fails

	var orderID int
	err = tx.QueryRow(`INSERT INTO orders (customer_name, phone, address, total_cost, status, customer_id) 
                       VALUES ($1, $2, $3, 0, 'pending', $4) RETURNING id`,
		customerName, phone, address, customerID).Scan(&orderID)
	if err != nil {
		return models.Order{}, err
	}

	// Batch-fetch prices for every DISTINCT pizza referenced, in one
	// round trip — was previously one SELECT per item (N+1, same shape
	// as the images bug fixed earlier).
	idSet := make(map[int]bool)
	var ids []int
	for _, item := range items {
		if !idSet[item.PizzaID] {
			idSet[item.PizzaID] = true
			ids = append(ids, item.PizzaID)
		}
	}

	type sizePrices struct{ small, medium, large int }
	priceByID := make(map[int]sizePrices)

	priceRows, err := tx.Query(`SELECT id, price_small, price_medium, price_large FROM pizzas WHERE is_active AND id = ANY($1)`, pq.Array(ids))
	if err != nil {
		return models.Order{}, err
	}
	for priceRows.Next() {
		var id, s, m, l int
		if err := priceRows.Scan(&id, &s, &m, &l); err != nil {
			priceRows.Close()
			return models.Order{}, err
		}
		priceByID[id] = sizePrices{s, m, l}
	}
	priceRows.Close()

	// Build ONE multi-row INSERT for all order_items instead of N
	// separate INSERT round trips. Still fully parameterized — only the
	// number of placeholder groups is dynamic, never the values
	// themselves, so this is not vulnerable to SQL injection.
	var totalCost int
	valuePlaceholders := make([]string, 0, len(items))
	valueArgs := make([]interface{}, 0, len(items)*5)
	argPos := 1

	for _, item := range items {
		p, ok := priceByID[item.PizzaID]
		if !ok {
			return models.Order{}, ErrPizzaNotFound
		}

		// item.Size is already restricted to small/medium/large by the
		// request binding, so this switch is exhaustive — the default
		// case only guards against an empty string slipping through.
		var price int
		switch item.Size {
		case "small":
			price = p.small
		case "large":
			price = p.large
		default:
			price = p.medium
		}

		subTotal := price * item.Quantity
		totalCost += subTotal

		valuePlaceholders = append(valuePlaceholders,
			fmt.Sprintf("($%d,$%d,$%d,$%d,$%d)", argPos, argPos+1, argPos+2, argPos+3, argPos+4))
		valueArgs = append(valueArgs, orderID, item.PizzaID, item.Size, item.Quantity, subTotal)
		argPos += 5
	}

	insertQuery := "INSERT INTO order_items (order_id, pizza_id, size, quantity, sub_total) VALUES " +
		strings.Join(valuePlaceholders, ",")
	if _, err := tx.Exec(insertQuery, valueArgs...); err != nil {
		return models.Order{}, err
	}

	// Update the final total cost
	_, err = tx.Exec("UPDATE orders SET total_cost = $1 WHERE id = $2", totalCost, orderID)
	if err != nil {
		return models.Order{}, err
	}

	return models.Order{ID: orderID, TotalCost: totalCost}, tx.Commit()
}

func GetAllOrdersWithPizzaName() ([]models.OrderResponse, error) {
	query := `
		SELECT o.id, o.customer_name, o.phone, o.address, p.name, oi.size, oi.quantity, o.total_cost, o.status, o.customer_id
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
		// o.CustomerID is *int — database/sql supports scanning NULL
		// directly into a nil **int destination, no sql.NullInt64 needed.
		err := rows.Scan(&o.OrderID, &o.CustomerName, &o.Phone, &o.Address, &o.PizzaName, &o.Size, &o.Quantity, &o.TotalCost, &o.Status, &o.CustomerID)
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

	// Previously two separate round trips to Supabase (revenue, then
	// counts) — combined into one query cuts this endpoint's network
	// latency roughly in half. The revenue subquery is scoped to
	// 'delivered' orders only, same as before; the COUNT(*) FILTER
	// clauses cover every status in a single pass over the same table.
	query := `
		SELECT
			COALESCE((SELECT SUM(total_cost) FROM orders WHERE status = 'delivered'), 0) AS revenue,
			COUNT(*) AS total_orders,
			COUNT(*) FILTER (WHERE status = 'pending')   AS pending_orders,
			COUNT(*) FILTER (WHERE status = 'delivered') AS delivered_orders,
			COUNT(*) FILTER (WHERE status = 'cancelled') AS cancelled_orders
		FROM orders
	`
	err := database.DB.QueryRow(query).Scan(
		&stats.Revenue,
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
		SELECT p.name, COALESCE(SUM(oi.quantity), 0) as total_sold
		FROM order_items oi
		JOIN pizzas p ON oi.pizza_id = p.id
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

var ErrInvalidCredentials = errors.New("invalid username or password")

// GetAdminByUsername fetches a single admin row by username.
func GetAdminByUsername(username string) (models.Admin, error) {
	query := `SELECT id, username, password_hash, created_at FROM admins WHERE username = $1`
	var admin models.Admin

	err := database.DB.QueryRow(query, username).Scan(
		&admin.ID, &admin.Username, &admin.PasswordHash, &admin.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Admin{}, ErrInvalidCredentials
		}
		return models.Admin{}, err
	}
	return admin, nil
}

// VerifyAdminCredentials looks up the admin and checks the password hash.
// Returns ErrInvalidCredentials for any mismatch — username not found or
// wrong password produce the exact same error and timing characteristics
// are close enough (bcrypt.CompareHashAndPassword always runs) to avoid
// leaking which case occurred via response time.
func VerifyAdminCredentials(username, password string) (models.Admin, error) {
	admin, err := GetAdminByUsername(username)
	if err != nil {
		// Still run a dummy bcrypt comparison so a non-existent username
		// takes roughly the same time as a wrong-password attempt,
		// reducing the username-enumeration timing side-channel.
		_ = bcrypt.CompareHashAndPassword(
			[]byte("$2a$10$7EqJtq98hPqEX7fNZaFWoOhi5L4Lhgk1Xb6/U6V5/9b3o3o4k1n0K"),
			[]byte(password),
		)
		return models.Admin{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		return models.Admin{}, ErrInvalidCredentials
	}

	return admin, nil
}

// CreateAdmin hashes the given password and inserts a new admin row.
// Intended for a one-off seed script or CLI command — not exposed via
// any public HTTP route, since admin self-signup is out of scope for
// a single-restaurant dashboard.
func CreateAdmin(username, plainPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `INSERT INTO admins (username, password_hash) VALUES ($1, $2)`
	_, err = database.DB.Exec(query, username, string(hash))
	return err
}
