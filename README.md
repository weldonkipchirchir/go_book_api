# Go Book API

A RESTful API for managing books built with Go, Gin, and GORM.

## Features

- Create, read, update, and delete books
- PostgreSQL database integration
- RESTful API endpoints
- JSON responses
- Unit tests with SQLite in-memory database

## Prerequisites

- Go 1.23.2 or later
- PostgreSQL database

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/weldonkipchirchir/go_book_api.git
   cd go_book_api
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Set up your PostgreSQL database and update the connection string in `api/handlers.go` if necessary.

## Usage

1. Run the server:

   ```bash
   go run cmd/main.go
   ```

   The server will start on port 8000 by default. You can change the port by setting the `PORT` environment variable:

   ```bash
   PORT=8080 go run cmd/main.go
   ```

## API Endpoints

### Create a Book

- **POST** `/books`
- Request Body:
  ```json
  {
    "title": "Book Title",
    "author": "Author Name",
    "year": 2023
  }
  ```
- Response:
  ```json
  {
    "status": 201,
    "message": "Book created successfully",
    "data": {
      "id": 1,
      "title": "Book Title",
      "author": "Author Name",
      "year": 2023
    }
  }
  ```

### Get All Books

- **GET** `/books`
- Response:
  ```json
  {
    "status": 200,
    "message": "Books retrieved successfully",
    "data": [
      {
        "id": 1,
        "title": "Book Title",
        "author": "Author Name",
        "year": 2023
      }
    ]
  }
  ```

### Get a Book by ID

- **GET** `/books/:id`
- Response:
  ```json
  {
    "status": 200,
    "message": "Book retrieved successfully",
    "data": {
      "id": 1,
      "title": "Book Title",
      "author": "Author Name",
      "year": 2023
    }
  }
  ```

### Update a Book

- **PUT** `/books/:id`
- Request Body:
  ```json
  {
    "title": "Updated Title",
    "author": "Updated Author",
    "year": 2024
  }
  ```
- Response:
  ```json
  {
    "status": 200,
    "message": "Book updated successfully",
    "data": {
      "id": 1,
      "title": "Updated Title",
      "author": "Updated Author",
      "year": 2024
    }
  }
  ```

### Delete a Book

- **DELETE** `/books/:id`
- Response:
  ```json
  {
    "status": 200,
    "message": "Book deleted successfully"
  }
  ```

## Testing

Run the unit tests:

```bash
go test ./tests/
```

## Project Structure

```
go_book_api/
├── api/
│   ├── handlers.go    # API handlers and database initialization
│   └── model.go       # Data models and response structures
├── cmd/
│   └── main.go        # Application entry point
├── tests/
│   └── main_test.go   # Unit tests
├── go.mod
├── go.sum
└── README.md
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
