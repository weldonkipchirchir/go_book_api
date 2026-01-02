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

	// Define routes
	r.POST("/books", api.CreateBook)
	r.GET("/books", api.GetBooks)
	r.GET("/books/:id", api.GetBookByID)
	r.PUT("/books/:id", api.UpdateBook)
	r.DELETE("/books/:id", api.DeleteBook)
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
