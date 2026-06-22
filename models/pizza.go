package models

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
	ID        int `json:"id"`
	PizzaId   int `json:"pizza_id"`
	Quantity  int `json:"quantity"`
	TotalCost int `json:"total_cost"`
}

type PizzaDetailResponse struct {
	Pizza  Pizza        `json:"pizza"`
	Images []PizzaImage `json:"images"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}

type OrderResponse struct {
	OrderID   int    `json:"order_id"`
	PizzaName string `json:"pizza_name"`
	Quantity  int    `json:"quantity"`
	TotalCost int    `json:"total_cost"`
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
