package repositories

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"pizza-app/database"
)

const paystackBaseURL = "https://api.paystack.co"

var ErrOrderNotFound = errors.New("order not found")

// paystackHTTPClient has an explicit timeout — without one, a slow or
// unresponsive Paystack request would hang this handler (and the
// connection holding it) indefinitely instead of failing cleanly.
var paystackHTTPClient = &http.Client{Timeout: 15 * time.Second}

// --- Paystack API response shapes (internal — only the fields we use) ---

type paystackInitializeResponse struct {
	Status bool `json:"status"`
	Data   struct {
		AuthorizationURL string `json:"authorization_url"`
		Reference        string `json:"reference"`
	} `json:"data"`
}

type paystackVerifyResponse struct {
	Status bool `json:"status"`
	Data   struct {
		Status    string `json:"status"`
		Reference string `json:"reference"`
		Amount    int    `json:"amount"`
	} `json:"data"`
}

// paystackSecretKey is read on each call (not at package init) so the
// binary doesn't crash on startup if it's briefly unset.
func paystackSecretKey() string {
	return os.Getenv("PAYSTACK_SECRET_KEY")
}

// InitializePaystackTransaction asks Paystack to open a transaction for the
// given order and returns the hosted checkout URL to send the customer to.
// The amount always comes from our own DB — never from client input — so a
// tampered request can't pay less than the real order total.
func InitializePaystackTransaction(orderID int, email string) (authURL string, reference string, err error) {
	secret := paystackSecretKey()
	if secret == "" {
		return "", "", errors.New("PAYSTACK_SECRET_KEY is not set")
	}

	var totalCost int
	err = database.DB.QueryRow(`SELECT total_cost FROM orders WHERE id = $1`, orderID).Scan(&totalCost)
	if err != nil {
		return "", "", ErrOrderNotFound
	}

	// Paystack expects the amount in kobo (smallest unit, 1 Naira = 100 kobo).
	// If total_cost is already stored in kobo in your schema, remove the *100.
	amountKobo := totalCost * 100

	reqBody, _ := json.Marshal(map[string]interface{}{
		"email":        email,
		"amount":       amountKobo,
		"reference":    fmt.Sprintf("order_%d_%d", orderID, totalCost),
		"callback_url": os.Getenv("PAYSTACK_CALLBACK_URL"),
		"metadata": map[string]interface{}{
			"order_id": orderID,
		},
	})

	req, _ := http.NewRequest("POST", paystackBaseURL+"/transaction/initialize", bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+secret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := paystackHTTPClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var parsed paystackInitializeResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", "", err
	}
	if !parsed.Status {
		return "", "", errors.New("paystack rejected the initialize request")
	}

	// Save the reference against the order so the webhook/verify call can
	// find its way back to this order later.
	if _, err := database.DB.Exec(`UPDATE orders SET payment_reference = $1 WHERE id = $2`, parsed.Data.Reference, orderID); err != nil {
		return "", "", err
	}

	return parsed.Data.AuthorizationURL, parsed.Data.Reference, nil
}

// VerifyPaystackTransaction calls Paystack's Verify endpoint directly. Used
// by a callback page as a fallback in case the webhook hasn't arrived yet.
func VerifyPaystackTransaction(reference string) (status string, orderID int, err error) {
	secret := paystackSecretKey()
	if secret == "" {
		return "", 0, errors.New("PAYSTACK_SECRET_KEY is not set")
	}

	req, _ := http.NewRequest("GET", paystackBaseURL+"/transaction/verify/"+reference, nil)
	req.Header.Set("Authorization", "Bearer "+secret)

	resp, err := paystackHTTPClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var parsed paystackVerifyResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", 0, err
	}
	if !parsed.Status {
		return "", 0, errors.New("paystack verify request failed")
	}

	paymentStatus := "pending"
	switch parsed.Data.Status {
	case "success":
		paymentStatus = "paid"
	case "failed", "abandoned":
		paymentStatus = "failed"
	}

	var id int
	err = database.DB.QueryRow(
		`UPDATE orders SET payment_status = $1 WHERE payment_reference = $2 RETURNING id`,
		paymentStatus, reference,
	).Scan(&id)
	if err != nil {
		return "", 0, err
	}

	return paymentStatus, id, nil
}

// VerifyWebhookSignature checks Paystack's X-Paystack-Signature header, an
// HMAC-SHA512 of the raw request body using your secret key. This is what
// proves a webhook request genuinely came from Paystack and not an
// attacker hitting your endpoint directly to fake a "payment successful".
func VerifyWebhookSignature(rawBody []byte, signature string) bool {
	secret := paystackSecretKey()
	if secret == "" {
		return false
	}
	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write(rawBody)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}

// MarkOrderPaidByReference is called once the webhook signature is verified
// and the event confirms a successful charge.
func MarkOrderPaidByReference(reference string) error {
	_, err := database.DB.Exec(`UPDATE orders SET payment_status = 'paid' WHERE payment_reference = $1`, reference)
	return err
}

// MarkOrderFailedByReference handles failed/abandoned webhook events.
func MarkOrderFailedByReference(reference string) error {
	_, err := database.DB.Exec(`UPDATE orders SET payment_status = 'failed' WHERE payment_reference = $1`, reference)
	return err
}
