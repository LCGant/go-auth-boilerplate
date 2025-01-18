package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/LCGant/go-auth-boilerplate/config"
	"github.com/LCGant/go-auth-boilerplate/controllers"
)

func main() {
	// Connect to the database
	config.ConnectDB()

	// Initialize Gin
	r := gin.Default()

	// Base directory for your static files (HTML, CSS, JS, images, etc.)
	staticDir := "public"

	// 1) Serve static files (CSS, JS, images) at /static
	//    Example: /static/css/style.css => ./public/css/style.css
	r.Static("/static", filepath.Join(staticDir))

	// 2) Define GET routes for specific HTML pages

	// / => Redirect to /login
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/login")
	})

	// /login => login.html
	r.GET("/login", func(c *gin.Context) {
		filePath := filepath.Join(staticDir, "pages", "login.html")
		c.File(filePath)
	})

	// /register => register.html
	r.GET("/register", func(c *gin.Context) {
		filePath := filepath.Join(staticDir, "pages", "register.html")
		c.File(filePath)
	})

	// /reset-password => reset-password.htm
	r.GET("/reset-password", func(c *gin.Context) {
		filePath := filepath.Join(staticDir, "pages", "reset-password.html")
		c.File(filePath)
	})

	// 3) API routes (POST and GET), avoiding wildcard conflicts
	r.POST("/register", controllers.RegisterHandler(config.DB))
	r.POST("/login", controllers.LoginHandler(config.DB))
	r.POST("/reset-password", controllers.ResetPasswordHandler(config.DB))
	r.POST("/forgot-password", controllers.ForgotPasswordHandler(config.DB))
	r.GET("/verify-token", controllers.VerifyTokenHandler(config.DB))
	r.GET("/verify-email", controllers.VerifyEmailHandler(config.DB))
	r.GET("/verify-email-token", controllers.VerifyEmailTokenHandler(config.DB))

	// 4) If no route matches, serve 404.html
	r.NoRoute(func(c *gin.Context) {
		filePath := filepath.Join(staticDir, "pages", "404.html")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// If 404.html doesn't exist, return a JSON error
			c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
			return
		}
		c.File(filePath)
	})

	// 5) Start the server on port 8080
	fmt.Println("Server running at http://127.0.0.1:8080/")
	r.Run(":8080")
}
