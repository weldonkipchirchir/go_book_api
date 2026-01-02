package tests

import (
	"bytes"
	"encoding/json"
	"go_book_api/api"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	dsn := ":memory:"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}
	db.AutoMigrate(&api.Book{})
	return db
}

func addBook() api.Book {
	book := api.Book{Title: "Go Programming", Author: "John Doe", Year: 2023}
	api.DB.Create(&book)
	return book
}

func TestCreateBook(t *testing.T) {
	api.DB = setupTestDB()
	router := gin.Default()
	router.POST("/books", api.CreateBook)

	book := api.Book{
		Title:  "Go Programming",
		Author: "John Doe",
		Year:   2023,
	}

	jsonValue, _ := json.Marshal(book)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	type APIResponse struct {
		Message string   `json:"message"`
		Data    api.Book `json:"data"`
	}

	var resp APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Data.Title != book.Title ||
		resp.Data.Author != book.Author ||
		resp.Data.Year != book.Year {
		t.Errorf("Expected %+v, got %+v", book, resp.Data)
	}
}

func TestGetBooks(t *testing.T) {
	api.DB = setupTestDB()
	router := gin.Default()
	router.GET("/books", api.GetBooks)

	addedBook := addBook()

	req, _ := http.NewRequest("GET", "/books", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	type APIResponse struct {
		Message string     `json:"message"`
		Data    []api.Book `json:"data"`
	}

	var resp APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("Expected 1 book, got %d", len(resp.Data))
	}
	if resp.Data[0].ID != addedBook.ID {
		t.Errorf("Expected book ID %d, got %d", addedBook.ID, resp.Data[0].ID)
	}

}

func TestGetBookByID(t *testing.T) {
	api.DB = setupTestDB()
	router := gin.Default()
	router.GET("/books/:id", api.GetBookByID)
	addedBook := addBook()

	req, _ := http.NewRequest("GET", "/books/"+strconv.Itoa(int(addedBook.ID)), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	type APIResponse struct {
		Message string   `json:"message"`
		Data    api.Book `json:"data"`
	}
	var resp APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Data.ID != addedBook.ID {
		t.Errorf("Expected book ID %d, got %d", addedBook.ID, resp.Data.ID)
	}
}

func TestUpdateBook(t *testing.T) {
	api.DB = setupTestDB()
	router := gin.Default()
	router.PUT("/books/:id", api.UpdateBook)
	addedBook := addBook()

	updatedBook := api.Book{
		Title:  "Advanced Go Programming",
		Author: "Jane Smith",
		Year:   2024,
	}

	jsonValue, _ := json.Marshal(updatedBook)
	req, _ := http.NewRequest("PUT", "/books/"+strconv.Itoa(int(addedBook.ID)), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	type APIResponse struct {
		Message string   `json:"message"`
		Data    api.Book `json:"data"`
	}

	var resp APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Data.Title != updatedBook.Title ||
		resp.Data.Author != updatedBook.Author ||
		resp.Data.Year != updatedBook.Year {
		t.Errorf("Expected %+v, got %+v", updatedBook, resp.Data)
	}

}

func TestDeleteBook(t *testing.T) {
	api.DB = setupTestDB()
	router := gin.Default()
	router.DELETE("/books/:id", api.DeleteBook)
	addedBook := addBook()

	req, _ := http.NewRequest("DELETE", "/books/"+strconv.Itoa(int(addedBook.ID)), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp api.JsonResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Message != "Book deleted successfully" {
		t.Errorf("Expected message 'Book deleted successfully', got '%s'", resp.Message)
	}

	var book api.Book
	result := api.DB.First(&book, addedBook.ID)
	if result.Error == nil {
		t.Errorf("Expected book to be deleted, but it still exists")
	}
}
