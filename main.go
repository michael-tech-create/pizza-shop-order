package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"pizza-app/database"
	"pizza-app/middleware"
	"pizza-app/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Embed only the files the browser needs. node_modules and scratch
// scripts under frontend/ stay out of the binary.
//
//go:embed frontend/index.html frontend/login.html frontend/admin.html frontend/activity.html frontend/pizza.html frontend/checkout.html frontend/my-orders.html frontend/userlogin.html frontend/payment-callback.html frontend/css frontend/js frontend/public
var frontendFS embed.FS

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on system environment variables")
	}

	middleware.RequireConfiguredSecret()
	database.ConnectDataBase()

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	// Hosts such as Render sit behind a proxy. Trust none by default so
	// client IP headers cannot be spoofed; set TRUSTED_PROXIES to a
	// comma-separated list when you need real client IPs.
	if proxies := os.Getenv("TRUSTED_PROXIES"); proxies != "" {
		parts := splitCSV(proxies)
		if err := router.SetTrustedProxies(parts); err != nil {
			log.Fatalf("invalid TRUSTED_PROXIES: %v", err)
		}
	} else if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatalf("failed to configure proxies: %v", err)
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	frontendContent, err := fs.Sub(frontendFS, "frontend")
	if err != nil {
		log.Fatalf("failed to create frontend sub-filesystem: %v", err)
	}

	router.StaticFS("/css", http.FS(mustSubFS(frontendContent, "css")))
	router.StaticFS("/js", http.FS(mustSubFS(frontendContent, "js")))
	router.StaticFS("/public", http.FS(mustSubFS(frontendContent, "public")))

	if err := os.MkdirAll("uploads", 0o755); err != nil {
		log.Fatalf("failed to create uploads directory: %v", err)
	}
	router.Static("/uploads", "./uploads")

	pages := []string{
		"index.html",
		"login.html",
		"admin.html",
		"activity.html",
		"pizza.html",
		"checkout.html",
		"my-orders.html",
		"userlogin.html",
		"payment-callback.html",
	}
	for _, name := range pages {
		handler := serveFile(frontendFS, "frontend/"+name)
		router.GET("/"+name, handler)
		if name == "index.html" {
			router.GET("/", handler)
			continue
		}
		router.GET("/"+strings.TrimSuffix(name, ".html"), handler)
	}

	router.GET("/health", func(c *gin.Context) {
		if err := database.DB.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "db unreachable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	routes.SetupRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("server starting on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}

func allowedOrigins() []string {
	origins := []string{
		"http://localhost:5500",
		"http://127.0.0.1:5500",
		"http://localhost:8080",
		"https://pizza-shop-order-ten.vercel.app",
	}
	if extra := os.Getenv("ALLOWED_ORIGINS"); extra != "" {
		origins = append(origins, splitCSV(extra)...)
	}
	return origins
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func serveFile(efs embed.FS, path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := efs.ReadFile(path)
		if err != nil {
			c.String(http.StatusNotFound, "404 page not found")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	}
}

func mustSubFS(sys fs.FS, dir string) fs.FS {
	sub, err := fs.Sub(sys, dir)
	if err != nil {
		log.Fatalf("missing embedded directory %s: %v", dir, err)
	}
	return sub
}
