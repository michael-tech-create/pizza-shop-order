package models

import "time"

type Category struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ImageURL   string `json:"image_url,omitempty"`
	PizzaCount int    `json:"pizza_count"`
}

type CreateCategoryRequest struct {
	Name     string `json:"name" binding:"required"`
	ImageURL string `json:"image_url"`
}

type Review struct {
	ID         int       `json:"id"`
	OrderID    int       `json:"order_id"`
	CustomerID int       `json:"customer_id"`
	Rating     int       `json:"rating"`
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
}

// CreateReviewRequest is submitted by a logged-in customer for one of
// their own orders — eligibility (must be delivered, must be theirs,
// must not already have a review) is enforced server-side.
type CreateReviewRequest struct {
	OrderID int    `json:"order_id" binding:"required,gt=0"`
	Rating  int    `json:"rating" binding:"required,gte=1,lte=5"`
	Comment string `json:"comment"`
}

// PublicReview is what's rendered on the homepage. CustomerDisplayName is
// abbreviated (first name + last initial) for privacy — real customer
// full names are never exposed publicly.
type PublicReview struct {
	CustomerDisplayName string    `json:"customer_display_name"`
	Rating              int       `json:"rating"`
	Comment             string    `json:"comment"`
	PizzaName           string    `json:"pizza_name"`
	CreatedAt           time.Time `json:"created_at"`
}

// ReviewStats are real, computed numbers — no fabricated stats like the
// old hardcoded "98% would order again" placeholders.
type ReviewStats struct {
	AverageRating float64 `json:"average_rating"`
	TotalReviews  int     `json:"total_reviews"`
}
