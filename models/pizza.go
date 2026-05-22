package models

type Pizza struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
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
