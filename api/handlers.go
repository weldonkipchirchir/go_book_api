package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

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

func GetBooks(c *gin.Context) {
	var books []Book
	if err := DB.Find(&books).Error; err != nil {
		RespondJSON(c, http.StatusInternalServerError, "Failed to retrieve books", nil)
		return
	}
	RespondJSON(c, http.StatusOK, "Books retrieved successfully", books)
}

func GetBookByID(c *gin.Context) {
	id := c.Param("id")
	var book Book
	if err := DB.First(&book, id).Error; err != nil {
		RespondJSON(c, http.StatusNotFound, "Book not found", nil)
		return
	}
	RespondJSON(c, http.StatusOK, "Book retrieved successfully", book)
}

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
