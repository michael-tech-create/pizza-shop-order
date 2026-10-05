package middleware

import (
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"pizza-app/models"

	"github.com/gin-gonic/gin"
)

// customerTokenTTL is longer than the admin's 12h session — customers
// aren't re-authenticating throughout a work shift, so a longer-lived
// "stay logged in" session is more appropriate here.
const customerTokenTTL = 30 * 24 * time.Hour

// GenerateCustomerToken creates a signed session token for a customer
// account. Reuses sign() and signingKey() from auth.go (same package) —
// the signing mechanism is identical, only the claims struct differs,
// which is what keeps a customer token from ever being mistaken for
// (or accepted as) an admin token.
func GenerateCustomerToken(customerID int, phone string) (string, int64, error) {
	now := time.Now()
	exp := now.Add(customerTokenTTL).Unix()

	claims := models.CustomerAuthClaims{
		CustomerID: customerID,
		Phone:      phone,
		Iat:        now.Unix(),
		Exp:        exp,
	}

	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		return "", 0, err
	}

	payload := base64.RawURLEncoding.EncodeToString(payloadBytes)
	signature := sign(payload)

	return payload + "." + signature, exp, nil
}

// VerifyCustomerToken checks signature + expiry and decodes the claims.
func VerifyCustomerToken(token string) (*models.CustomerAuthClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return nil, ErrMalformedToken
	}
	payload, sigGiven := parts[0], parts[1]

	sigExpected := sign(payload)
	if subtle.ConstantTimeCompare([]byte(sigGiven), []byte(sigExpected)) != 1 {
		return nil, ErrInvalidSignature
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return nil, ErrMalformedToken
	}

	var claims models.CustomerAuthClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, ErrMalformedToken
	}

	if time.Now().Unix() > claims.Exp {
		return nil, ErrTokenExpired
	}

	return &claims, nil
}

// RequireCustomerAuth protects routes that need a logged-in customer,
// e.g. GET /api/customers/orders. Mirrors RequireAuth's behavior but
// for the customer claim type, storing "customer_id" in the context.
func RequireCustomerAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		// log.Printf("RequireCostumerAuth: received header = %q", header)
		const prefix = "Bearer "
		if header == "" || !strings.HasPrefix(header, prefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "please log in to view this"})
			return
		}
		token := strings.TrimPrefix(header, prefix)

		claims, err := VerifyCustomerToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired session, please log in again"})
			return
		}

		c.Set("customer_id", claims.CustomerID)
		c.Set("customer_phone", claims.Phone)
		c.Next()
	}
}

// OptionalCustomerID checks for a valid customer Bearer token WITHOUT
// rejecting the request if one isn't present or doesn't verify — this is
// what lets /api/orders stay open to guest checkout while still linking
// the order to an account when the customer happens to be logged in.
// Returns nil when there's no valid logged-in customer.
func OptionalCustomerID(c *gin.Context) *int {
	header := c.GetHeader("Authorization")
	const prefix = "Bearer "
	if header == "" || !strings.HasPrefix(header, prefix) {
		return nil
	}
	token := strings.TrimPrefix(header, prefix)

	claims, err := VerifyCustomerToken(token)
	if err != nil {
		return nil
	}
	id := claims.CustomerID
	return &id
}