package main

import (
	"log"
	"os"

	"go_book_api/api"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database
	api.InitDB()
	r := gin.Default()

	r.POST("/token", api.GenerateJWT)

	// Define routes
	protected := r.Group("/", api.JWTAuthMiddleware())
	{
		protected.POST("/books", api.CreateBook)
		protected.GET("/books", api.GetBooks)
		protected.GET("/books/:id", api.GetBookByID)
		protected.PUT("/books/:id", api.UpdateBook)
		protected.DELETE("/books/:id", api.DeleteBook)
	}
	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server: ", err)
	} else {
		log.Println("Server running on port " + port)
	}
}
