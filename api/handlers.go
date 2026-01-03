/*
Package api provides RESTful handlers and middleware for managing a library of books.
It includes functionality for creating, reading, updating, and deleting books, as well as generating JWT tokens for authentication.
*/
package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the database connection using environment variables.
// It loads the database configuration from a .env file and migrates the Book schema.
func InitDB() {
	// if err := godotenv.Load(); err != nil {
	// 	log.Println("No .env file found, using system environment variables")
	// }

	// dsn := os.Getenv("DATABASE_URL")
	// if dsn == "" {
	// 	log.Fatal("DATABASE_URL is not set")
	// }

	var err error
	DB, err = gorm.Open(postgres.Open("postgresql://neondb_owner:npg_FtXzVU5KHDh4@ep-divine-voice-adf4bkqc-pooler.c-2.us-east-1.aws.neon.tech/neondb?sslmode=require&channel_binding=requires"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")

	if err := DB.AutoMigrate(&Book{}); err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}
}

// CreateBook handles the creation of a new book in the database.
// It expects a JSON payload with book details and responds with the created book.
func CreateBook(c *gin.Context) {
	var book Book

	//bind the request body to book struct
	if err := c.ShouldBindJSON(&book); err != nil {
		RespondJSON(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}
	if err := DB.Create(&book).Error; err != nil {
		RespondJSON(c, http.StatusInternalServerError, "Failed to create book", nil)
		return
	}
	RespondJSON(c, http.StatusCreated, "Book created successfully", book)
}

// GetBooks retrieves all books from the database.
// It responds with a list of books.
func GetBooks(c *gin.Context) {
	var books []Book
	if err := DB.Find(&books).Error; err != nil {
		RespondJSON(c, http.StatusInternalServerError, "Failed to retrieve books", nil)
		return
	}
	RespondJSON(c, http.StatusOK, "Books retrieved successfully", books)
}

// GetBookByID retrieves a book by its ID from the database.
// It responds with the book details if found.
func GetBookByID(c *gin.Context) {
	id := c.Param("id")
	var book Book
	if err := DB.First(&book, id).Error; err != nil {
		RespondJSON(c, http.StatusNotFound, "Book not found", nil)
		return
	}
	RespondJSON(c, http.StatusOK, "Book retrieved successfully", book)
}

// UpdateBook updates an existing book in the database.
// It expects a JSON payload with updated book details and responds with the updated book.
func UpdateBook(c *gin.Context) {
	id := c.Param("id")
	var book Book
	if err := DB.First(&book, id).Error; err != nil {
		RespondJSON(c, http.StatusNotFound, "Book not found", nil)
		return
	}
	if err := c.ShouldBindJSON(&book); err != nil {
		RespondJSON(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}
	if err := DB.Save(&book).Error; err != nil {
		RespondJSON(c, http.StatusInternalServerError, "Failed to update book", nil)
		return
	}
	RespondJSON(c, http.StatusOK, "Book updated successfully", book)
}

// DeleteBook deletes a book by its ID from the database.
// It responds with a success message if the deletion is successful.
func DeleteBook(c *gin.Context) {
	id := c.Param("id")
	var book Book
	if err := DB.First(&book, id).Error; err != nil {
		RespondJSON(c, http.StatusNotFound, "Book not found", nil)
		return
	}

	if err := DB.Delete(&book).Error; err != nil {
		RespondJSON(c, http.StatusInternalServerError, "Failed to delete book", nil)
		return
	}
	RespondJSON(c, http.StatusOK, "Book deleted successfully", nil)
}

// GenerateJWT generates a JWT token for authenticated users.
// It expects a JSON payload with username and password, and responds with the token if credentials are valid.
func GenerateJWT(c *gin.Context) {
	var loginRequest loginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		RespondJSON(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if loginRequest.Username != "admin" || loginRequest.Password != "password" {
		RespondJSON(c, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": loginRequest.Username,
		"exp":      jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		RespondJSON(c, http.StatusInternalServerError, "Failed to generate token", nil)
		return
	}

	RespondJSON(c, http.StatusOK, "Token generated successfully", gin.H{"token": tokenString})
}
