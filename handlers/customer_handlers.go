package handlers

import (
	"net/http"

	"pizza-app/middleware"
	"pizza-app/models"
	"pizza-app/repositories"

	"github.com/gin-gonic/gin"
)

// POST /api/customers/signup
func CustomerSignupHandler(c *gin.Context) {
	var req models.CustomerSignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, phone, and a password of at least 6 characters are required"})
		return
	}

	customer, err := repositories.CreateCustomer(req.Name, req.Phone, req.Password)
	if err != nil {
		if err == repositories.ErrCustomerExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create account"})
		return
	}

	token, expiresAt, err := middleware.GenerateCustomerToken(customer.ID, customer.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "account created but failed to start session — please log in"})
		return
	}

	c.JSON(http.StatusCreated, models.CustomerLoginResponse{
		Token:     token,
		Name:      customer.Name,
		ExpiresAt: expiresAt,
	})
}

// POST /api/customers/login
func CustomerLoginHandler(c *gin.Context) {
	var req models.CustomerLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone and password are required"})
		return
	}

	customer, err := repositories.VerifyCustomerCredentials(req.Phone, req.Password)
	if err != nil {
		// Same generic message either way — never reveal whether the
		// phone number exists or the password was wrong.
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid phone number or password"})
		return
	}

	token, expiresAt, err := middleware.GenerateCustomerToken(customer.ID, customer.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
		return
	}

	c.JSON(http.StatusOK, models.CustomerLoginResponse{
		Token:     token,
		Name:      customer.Name,
		ExpiresAt: expiresAt,
	})
}

// GET /api/customers/orders — behind RequireCustomerAuth
func GetMyOrdersHandler(c *gin.Context) {
	customerID, _ := c.Get("customer_id")
	id, ok := customerID.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "please log in to view your orders"})
		return
	}

	history, err := repositories.GetOrdersByCustomerID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load order history"})
		return
	}
	c.JSON(http.StatusOK, history)
}