package models
import "time"

type Pizza struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Price       int          `json:"price"`
	Description string       `json:"description"`
	Images      []PizzaImage `json:"images,omitempty"`
}

type PizzaImage struct {
	ID       int    `json:"id"`
	PizzaID  int    `json:"pizza_id"`
	ImageURL string `json:"image_url"`
}

type DashboardStats struct {
	TotalOrders     int `json:"total_orders"`
	PendingOrders   int `json:"pending_orders"`
	DeliveredOrders int `json:"delivered_orders"`
	CancelledOrders int `json:"cancelled_orders"`
	Revenue         int `json:"revenue"`
}

type BestSellingPizza struct {
	PizzaName string `json:"pizza_name"`
	Sold      int    `json:"sold"`
}

type Order struct {
ID           int         `json:"id"`
	CustomerName string      `json:"customer_name"`
	Phone        string      `json:"phone"`
	Address      string      `json:"address"`
	TotalCost    int         `json:"total_cost"`
	Status       string      `json:"status"`
	Items        []OrderItem `json:"items,omitempty"`
	
	
}

type PizzaDetailResponse struct {
	Pizza  Pizza        `json:"pizza"`
	Images []PizzaImage `json:"images"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}

type OrderResponse struct {
	OrderID      int    `json:"order_id"`
	CustomerName string `json:"customer_name"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	PizzaName    string `json:"pizza_name"`
	Quantity     int    `json:"quantity"`
	TotalCost    int    `json:"total_cost"`
	Status       string `json:"status"`
}

type OrderItem struct {
ID      int `json:"id,omitempty"` // DB primary key
	OrderID  int `json:"order_id,omitempty"`
	PizzaID  int `json:"pizza_id" binding:"required,gt=0"`
	Quantity int `json:"quantity" binding:"required,gt=0"`
	SubTotal int `json:"sub_total,omitempty"`
}

type CreateOrderRequest struct {
CustomerName string      `json:"customer_name" binding:"required"`
	Phone        string      `json:"phone" binding:"required"`
	Address      string      `json:"address" binding:"required"`
	Items        []OrderItem `json:"items" binding:"required,min=1,dive"`
}

type Admin struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}
 
// LoginRequest is the payload sent from login.html.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
 
// LoginResponse is returned on successful authentication.
type LoginResponse struct {
	Token     string `json:"token"`
	Username  string `json:"username"`
	ExpiresAt int64  `json:"expires_at"` // unix seconds, lets the frontend show "session expires in..."
}
 
// AuthClaims is encoded into the JWT payload.
type AuthClaims struct {
	AdminID  int    `json:"admin_id"`
	Username string `json:"username"`
	Exp      int64  `json:"exp"`
	Iat      int64  `json:"iat"`
}
 

var Pizzas = []Pizza{
	{
		ID:          1,
		Name:        "Pepperoni",
		Price:       12,
		Description: "Spicy pepperoni with extra cheese",
	},
	{
		ID:          2,
		Name:        "BBQ Chicken",
		Price:       12,
		Description: "Original smoky BBQ flavor",
	},
}

var NextPizzaID = 3