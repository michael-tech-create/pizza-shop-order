package repositories

import (
	"database/sql"
	"errors"

	"pizza-app/database"
	"pizza-app/models"

	"golang.org/x/crypto/bcrypt"
)

var ErrCustomerExists = errors.New("an account with this phone number already exists")
var ErrInvalidCustomerCredentials = errors.New("invalid phone number or password")

// CreateCustomer hashes the password and inserts a new customer account.
func CreateCustomer(name, phone, plainPassword string) (models.Customer, error) {
	// Check for an existing account first so we can return a clear,
	// specific error rather than a raw unique-constraint DB error.
	var existingID int
	err := database.DB.QueryRow(`SELECT id FROM customers WHERE phone = $1`, phone).Scan(&existingID)
	if err == nil {
		return models.Customer{}, ErrCustomerExists
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return models.Customer{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return models.Customer{}, err
	}

	var customer models.Customer
	customer.Name = name
	customer.Phone = phone
	err = database.DB.QueryRow(
		`INSERT INTO customers (name, phone, password_hash) VALUES ($1, $2, $3) RETURNING id, created_at`,
		name, phone, string(hash),
	).Scan(&customer.ID, &customer.CreatedAt)
	if err != nil {
		return models.Customer{}, err
	}
	return customer, nil
}

func GetCustomerByPhone(phone string) (models.Customer, error) {
	var c models.Customer
	err := database.DB.QueryRow(
		`SELECT id, name, phone, password_hash, created_at FROM customers WHERE phone = $1`, phone,
	).Scan(&c.ID, &c.Name, &c.Phone, &c.PasswordHash, &c.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Customer{}, ErrInvalidCustomerCredentials
		}
		return models.Customer{}, err
	}
	return c, nil
}

// VerifyCustomerCredentials mirrors the admin login's timing-safe pattern —
// a dummy bcrypt comparison runs even when the phone number isn't found,
// so a nonexistent account takes about as long to reject as a wrong
// password, avoiding a phone-number enumeration side channel.
func VerifyCustomerCredentials(phone, password string) (models.Customer, error) {
	customer, err := GetCustomerByPhone(phone)
	if err != nil {
		_ = bcrypt.CompareHashAndPassword(
			[]byte("$2a$10$7EqJtq98hPqEX7fNZaFWoOhi5L4Lhgk1Xb6/U6V5/9b3o3o4k1n0K"),
			[]byte(password),
		)
		return models.Customer{}, ErrInvalidCustomerCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(customer.PasswordHash), []byte(password)); err != nil {
		return models.Customer{}, ErrInvalidCustomerCredentials
	}

	return customer, nil
}

// GetOrdersByCustomerID returns order history for a single logged-in
// customer — scoped by customer_id so one customer can never see
// another's orders.
func GetOrdersByCustomerID(customerID int) ([]models.CustomerOrderHistoryItem, error) {
	query := `
		SELECT o.id, p.name, oi.size, oi.quantity, o.total_cost, o.status, o.payment_status
		FROM orders o
		JOIN order_items oi ON o.id = oi.order_id
		JOIN pizzas p ON oi.pizza_id = p.id
		WHERE o.customer_id = $1
		ORDER BY o.id DESC
	`
	rows, err := database.DB.Query(query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	history := []models.CustomerOrderHistoryItem{}
	for rows.Next() {
		var h models.CustomerOrderHistoryItem
		if err := rows.Scan(&h.OrderID, &h.PizzaName, &h.Size, &h.Quantity, &h.TotalCost, &h.Status, &h.PaymentStatus); err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	return history, nil
}