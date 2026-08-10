package repositories

import (
	"database/sql"
	"errors"
	"strings"

	"pizza-app/database"
	"pizza-app/models"

	"github.com/lib/pq"
)

var ErrOrderNotEligibleForReview = errors.New("this order isn't eligible for a review yet — it must be delivered first")
var ErrOrderNotYours = errors.New("this order doesn't belong to your account")
var ErrAlreadyReviewed = errors.New("you've already reviewed this order")
var ErrOrderNotFoundForReview = errors.New("order not found")

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}
	return false
}

// CreateReview is what makes every review a "verified order" review — it
// checks the order actually belongs to this customer AND has reached
// "delivered" status before allowing a review to be attached to it.
func CreateReview(customerID, orderID, rating int, comment string) (models.Review, error) {
	var orderCustomerID sql.NullInt64
	var status string
	err := database.DB.QueryRow(`SELECT customer_id, status FROM orders WHERE id = $1`, orderID).
		Scan(&orderCustomerID, &status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Review{}, ErrOrderNotFoundForReview
		}
		return models.Review{}, err
	}

	if !orderCustomerID.Valid || int(orderCustomerID.Int64) != customerID {
		return models.Review{}, ErrOrderNotYours
	}
	if status != "delivered" {
		return models.Review{}, ErrOrderNotEligibleForReview
	}

	var review models.Review
	err = database.DB.QueryRow(
		`INSERT INTO reviews (order_id, customer_id, rating, comment) 
         VALUES ($1, $2, $3, $4) RETURNING id, created_at`,
		orderID, customerID, rating, comment,
	).Scan(&review.ID, &review.CreatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return models.Review{}, ErrAlreadyReviewed
		}
		return models.Review{}, err
	}

	review.OrderID = orderID
	review.CustomerID = customerID
	review.Rating = rating
	review.Comment = comment
	return review, nil
}

// abbreviateName turns "Chidinma Adaeze" into "Chidinma A." — real full
// names are never exposed on the public homepage.
func abbreviateName(fullName string) string {
	parts := strings.Fields(fullName)
	if len(parts) == 0 {
		return "Anonymous"
	}
	if len(parts) == 1 {
		return parts[0]
	}
	last := parts[len(parts)-1]
	return parts[0] + " " + strings.ToUpper(last[:1]) + "."
}

// GetPublicReviews returns the most recent reviews for the homepage. The
// pizza name uses a scalar subquery rather than a JOIN so an order with
// multiple pizzas doesn't produce duplicate review rows.
func GetPublicReviews(limit int) ([]models.PublicReview, error) {
	query := `
		SELECT c.name, r.rating, r.comment, r.created_at,
		       (SELECT p.name FROM order_items oi
		        JOIN pizzas p ON oi.pizza_id = p.id
		        WHERE oi.order_id = r.order_id LIMIT 1) AS pizza_name
		FROM reviews r
		JOIN customers c ON r.customer_id = c.id
		ORDER BY r.created_at DESC
		LIMIT $1
	`
	rows, err := database.DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews := []models.PublicReview{}
	for rows.Next() {
		var fullName string
		var pr models.PublicReview
		if err := rows.Scan(&fullName, &pr.Rating, &pr.Comment, &pr.CreatedAt, &pr.PizzaName); err != nil {
			return nil, err
		}
		pr.CustomerDisplayName = abbreviateName(fullName)
		reviews = append(reviews, pr)
	}
	return reviews, nil
}

// GetReviewStats computes real numbers — replaces the old hardcoded
// "98% would order again" style placeholders on the homepage.
func GetReviewStats() (models.ReviewStats, error) {
	var stats models.ReviewStats
	err := database.DB.QueryRow(`SELECT COALESCE(AVG(rating), 0), COUNT(*) FROM reviews`).
		Scan(&stats.AverageRating, &stats.TotalReviews)
	return stats, err
}