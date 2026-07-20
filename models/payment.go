package models

// is created, to start a Paystack transaction for that order's total.
type InitializePaymentRequest struct {
	OrderID int    `json:"order_id" binding:"required,gt=0"`
	Email   string `json:"email" binding:"required,email"`
}

// InitializePaymentResponse hands back the Paystack-hosted checkout URL
// for the frontend to redirect the customer to.
type InitializePaymentResponse struct {
	AuthorizationURL string `json:"authorization_url"`
	Reference        string `json:"reference"`
}

// VerifyPaymentResponse is returned when the frontend's callback page
// checks payment status directly (a fallback to the webhook).
type VerifyPaymentResponse struct {
	OrderID       int    `json:"order_id"`
	PaymentStatus string `json:"payment_status"` // paid | failed | pending
	Reference     string `json:"reference"`
}

// PaystackWebhookEvent is the payload Paystack POSTs to your webhook URL
// whenever a transaction's status changes. Only the fields we use are
// declared here — Paystack sends more, but we don't need them.
type PaystackWebhookEvent struct {
	Event string `json:"event"` // e.g. "charge.success"
	Data  struct {
		Status    string `json:"status"`
		Reference string `json:"reference"`
	} `json:"data"`
}