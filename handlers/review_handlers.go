package handlers

import (
	"net/http"

	"pizza-app/models"
	"pizza-app/repositories"

	"github.com/gin-gonic/gin"
)

// POST /api/reviews — behind RequireCustomerAuth
func CreateReviewHandler(c *gin.Context) {
	customerIDVal, _ := c.Get("customer_id")
	customerID, ok := customerIDVal.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "please log in to leave a review"})
		return
	}

	var req models.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "a rating (1-5) and order_id are required"})
		return
	}

	review, err := repositories.CreateReview(customerID, req.OrderID, req.Rating, req.Comment)
	if err != nil {
		switch err {
		case repositories.ErrOrderNotYours, repositories.ErrOrderNotEligibleForReview,
			repositories.ErrAlreadyReviewed, repositories.ErrOrderNotFoundForReview:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save review"})
		}
		return
	}

	c.JSON(http.StatusCreated, review)
}

// GET /api/reviews — public, powers the homepage reviews section
func GetPublicReviewsHandler(c *gin.Context) {
	reviews, err := repositories.GetPublicReviews(12)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load reviews"})
		return
	}

	stats, err := repositories.GetReviewStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load review stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reviews": reviews,
		"stats":   stats,
	})
}