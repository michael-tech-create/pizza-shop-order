package models

import "time"

type Customer struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Phone        string    `json:"phone"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

// CustomerSignupRequest is the payload for creating a new customer account.
type CustomerSignupRequest struct {
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// CustomerLoginRequest is the payload for logging into an existing account.
type CustomerLoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// CustomerLoginResponse is returned on successful signup or login.
type CustomerLoginResponse struct {
	Token     string `json:"token"`
	Name      string `json:"name"`
	ExpiresAt int64  `json:"expires_at"`
}

// CustomerAuthClaims is encoded into a customer's session token. Kept
// separate from the admin's AuthClaims so a customer token can never be
// mistaken for (or reused as) an admin token, even though both use the
// same signing mechanism under the hood.
type CustomerAuthClaims struct {
	CustomerID int    `json:"customer_id"`
	Phone      string `json:"phone"`
	Exp        int64  `json:"exp"`
	Iat        int64  `json:"iat"`
}

// CustomerOrderHistoryItem is one row in a logged-in customer's order
// history — similar to OrderResponse but scoped to what a customer
// should see about their own orders (no other customers' data).
type CustomerOrderHistoryItem struct {
	OrderID       int    `json:"order_id"`
	PizzaName     string `json:"pizza_name"`
	Size          string `json:"size"`
	Quantity      int    `json:"quantity"`
	TotalCost     int    `json:"total_cost"`
	Status        string `json:"status"`
	PaymentStatus string `json:"payment_status"`
	Reviewed      bool   `json:"reviewed"`
}
