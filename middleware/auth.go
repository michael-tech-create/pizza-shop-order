package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"pizza-app/models"

	"github.com/gin-gonic/gin"
)

// TOKEN CONFIG


// tokenTTL controls how long a session lasts before re-login is required.
const tokenTTL = 12 * time.Hour

// jwtSecret is read once at package init. In production this MUST come
// from an environment variable — never hardcode it. A fallback is provided
// only so local dev doesn't crash if JWT_SECRET is unset, but this should
// always be overridden via .env in real deployments.
var jwtSecret = func() []byte {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return []byte(s)
	}
	// Dev-only fallback — DO NOT rely on this in production.
	return []byte("dev-insecure-secret-change-me-in-env")
}()



// GenerateToken creates a signed, base64url-encoded token of the form
// "<payload>.<signature>". This is a minimal JWT-style token (HS256)
// built on the standard library so no external dependency is required.
func GenerateToken(adminID int, username string) (string, int64, error) {
	now := time.Now()
	exp := now.Add(tokenTTL).Unix()

	claims := models.AuthClaims{
		AdminID:  adminID,
		Username: username,
		Iat:      now.Unix(),
		Exp:      exp,
	}

	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		return "", 0, err
	}

	payload := base64.RawURLEncoding.EncodeToString(payloadBytes)
	signature := sign(payload)

	token := payload + "." + signature
	return token, exp, nil
}

// sign computes the HMAC-SHA256 signature of the given payload, base64url-encoded.
func sign(payload string) string {
	mac := hmac.New(sha256.New, jwtSecret)
	mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}



var (
	ErrMalformedToken = errors.New("malformed token")
	ErrInvalidSignature = errors.New("invalid token signature")
	ErrTokenExpired     = errors.New("token has expired")
)

// VerifyToken checks the signature and expiry of a token and returns the
// decoded claims if valid.
func VerifyToken(token string) (*models.AuthClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return nil, ErrMalformedToken
	}
	payload, sigGiven := parts[0], parts[1]

	// Constant-time comparison to avoid timing attacks on signature check.
	sigExpected := sign(payload)
	if subtle.ConstantTimeCompare([]byte(sigGiven), []byte(sigExpected)) != 1 {
		return nil, ErrInvalidSignature
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return nil, ErrMalformedToken
	}

	var claims models.AuthClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, ErrMalformedToken
	}

	if time.Now().Unix() > claims.Exp {
		return nil, ErrTokenExpired
	}

	return &claims, nil
}



// RequireAuth protects routes behind a valid Bearer token. Attach to any
// route group that should be admin-only, e.g.:
//
//	admin := router.Group("/api/admin")
//	admin.Use(middleware.RequireAuth())
//	{
//	    admin.GET("/stats", handlers.GetDashboardStatsHandler)
//	}
//
// On success, the verified claims are stored in the Gin context under
// "admin_id" and "admin_username" for handlers to read if needed.
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		
		const prefix = "Bearer "
		if !strings.HasPrefix(header, prefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header must use Bearer scheme",
			})
			return
		}
		token := strings.TrimPrefix(header, prefix)

		claims, err := VerifyToken(token)
		if err != nil {
			status := http.StatusUnauthorized
			msg := "invalid or expired session, please log in again"
			c.AbortWithStatusJSON(status, gin.H{"error": msg})
			return
		}

		c.Set("admin_id", claims.AdminID)
		c.Set("admin_username", claims.Username)
		c.Next()
	}
}