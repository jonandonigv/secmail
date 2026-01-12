package main

import (
	"log"
	"os"
	"secmail/internal/auth"
	"secmail/internal/database"
	"secmail/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Database DSN from environment variable
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	db, err := database.InitDB(dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	r := gin.Default()

	// Auth routes
	// Public routes
	r.POST("/register", func(c *gin.Context) {
		auth.Register(c, db)
	})
	r.POST("/login", func(c *gin.Context) {
		auth.Login(c, db)
	})

	// Protected routes
	emails := r.Group("/emails")
	emails.Use(auth.JWTMiddleware())
	{
		emails.POST("/send", func(c *gin.Context) {
			handlers.SendEmail(c, db)
		})
		emails.GET("/inbox", func(c *gin.Context) {
			handlers.GetInbox(c, db)
		})
	}

	log.Println("Server starting on :8080")
	r.Run(":8080")
}
