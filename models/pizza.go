package models

type Pizza struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
}

type Order struct{
	ID int `json:"id"`
	PizzaId int `json:"pizza_id"`
	Quantity int `json:"quantity"`
	TotalCost int `json:"total_cost"`
}


type UpdateStatusRequest struct {
	Status string `json:"status"`
}
// blue print 

type OrderResponse struct {
	OrderID int `json:"order_id"`
	PizzaName string `json:"pizza_name"`
	Quantity int `json:"quantity"`
	TotalCost int `json:total_cost`
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
