package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"pizza-app/models"
	"pizza-app/repositories"

	"github.com/gin-gonic/gin"
)

// POST /api/payments/initialize
// Called by the frontend right after an order is created. Returns a
// Paystack-hosted checkout URL to redirect the customer to.
func InitializePaymentHandler(c *gin.Context) {
	var req models.InitializePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_id and a valid email are required"})
		return
	}

	authURL, reference, err := repositories.InitializePaystackTransaction(req.OrderID, req.Email)
	if err != nil {
		log.Println("payment init error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not start payment"})
		return
	}

	c.JSON(http.StatusOK, models.InitializePaymentResponse{
		AuthorizationURL: authURL,
		Reference:        reference,
	})
}

// GET /api/payments/verify/:reference
// A backup check the frontend's callback page can call after Paystack
// redirects the customer back, in case the webhook hasn't landed yet.
func VerifyPaymentHandler(c *gin.Context) {
	reference := c.Param("reference")

	status, orderID, err := repositories.VerifyPaystackTransaction(reference)
	if err != nil {
		log.Println("payment verify error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not verify payment"})
		return
	}

	c.JSON(http.StatusOK, models.VerifyPaymentResponse{
		OrderID:       orderID,
		PaymentStatus: status,
		Reference:     reference,
	})
}

// POST /api/payments/webhook
// Paystack calls this server-to-server whenever a transaction's status
// changes. This is the source of truth for payment status — more
// reliable than the redirect, since the customer's browser might close
// before the redirect completes.
func PaystackWebhookHandler(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read request body"})
		return
	}

	signature := c.GetHeader("X-Paystack-Signature")
	if !repositories.VerifyWebhookSignature(body, signature) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
		return
	}

	var event models.PaystackWebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	switch event.Event {
	case "charge.success":
		if err := repositories.MarkOrderPaidByReference(event.Data.Reference); err != nil {
			log.Println("failed to mark order paid:", err)
		}
	default:
		if event.Data.Status == "failed" {
			if err := repositories.MarkOrderFailedByReference(event.Data.Reference); err != nil {
				log.Println("failed to mark order failed:", err)
			}
		}
	}

	// Respond 200 quickly regardless — Paystack retries on non-200s, and
	// we've already handled (or intentionally ignored) the event.
	c.JSON(http.StatusOK, gin.H{"received": true})
}
