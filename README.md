# Arithmetic Calculator API

This is the backend service for the application. It handles user operations, calculates results, and communicates with the database.

Ensure you have the following installed:

- **Go** (version 1.22 or above)
- **PostgreSQL** (or any database of your choice)
- **Docker** (optional, for containerization)

## Setup and Installation

1. **Clone the Repository**

    ```bash
    git clone https://github.com/ezequielfornes/arithmetic-calculator-go.git
    ```

2. **Environment Variables**

   Create a `.env` file in the root directory with the following:

    ```bash
    DATABASE_URL= dbConnection
    ```

   Alternatively, you can set the environment variables in your system or cloud environment settings (e.g., Render, Heroku).

3. **Install Dependencies**

    ```bash
    go mod tidy
    ```

4. **Database Setup**

   If you are using PostgreSQL:

  - Create a PostgreSQL database:

    ```sql
    CREATE DATABASE arithmetic_db;
    ```
   Make sure your `DATABASE_URL` in the `.env` file is pointing to your database correctly.

5. **Run the Application**

    ```bash
    go run main.go
    ```

   The backend will be running at `http://localhost:8080`.

6. **Testing**

   Run tests using:

    ```bash
    go test ./...
    ```

7. **Docker Setup (Optional)**

   To run the backend in Docker, ensure Docker is installed, and create a `Dockerfile` and a `docker-compose.yml`:

   **Dockerfile**

    ```dockerfile
    FROM golang:1.22

    WORKDIR /app
    COPY . .

    RUN go mod tidy
    RUN go build -o main .

    EXPOSE 8080
    CMD ["./main"]
    ```

   **docker-compose.yml**

    ```yaml
    version: "3"
    services:
      app:
        build: .
        ports:
          - "8080:8080"
        environment:
          - DATABASE_URL=your_postgres_connection_string
          - SENDGRID_API_KEY=your_sendgrid_api_key
      db:
        image: postgres:latest
        environment:
          POSTGRES_USER: postgres
          POSTGRES_PASSWORD: your_password
          POSTGRES_DB: your_database_name
        ports:
          - "5432:5432"
    ```

   Run with Docker:

    ```bash
    docker-compose up --build
    ```

## API Endpoints

- **POST /operation**

  Calculates the result of an operation and stores it.

  **Body Example**

  ```json
  {
    "user_id": 1,
    "operation": "addition",
    "amount": 10
  }
## Technologies
  - Go (Gin framework)
  - PostgreSQL (for database)

## Types of Operation
- "addition"
- "subtraction"
- "multiplication"
- "division"
- "square_root"
- "random_string"

## cURL Examples
## Login Test User

curl -X POST http://localhost:8080/api/v1/auth/login \
-H "Content-Type: application/json" \
-d '{
"username": "testuser1@example.com",
"password": "password123"
}'

Expected Response

{
"token": "your_jwt_token_here"
}

Perform Operation

curl -X POST http://localhost:8080/api/v1/operation \
-H "Content-Type: application/json" \
-H "Authorization: Bearer your_jwt_token_here" \
-d '{
"type": "addition",
"amount": 50.0
}'
