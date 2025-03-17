## Overview

This Task Manager application is built using Golang, PostgreSQL, and GORM ORM.

## Features

- User Registration and Authentication
- JWT-based Authentication
- Task Management (Create, Read, Update, Delete)
- Pagination support for listing tasks
- Following SOLID principles and ACID transactions
- Dockerized for easy deployment

## Project Structure

```
task-manager/
├── dto/
│   ├── request/        # Request DTOs
│   └── response/       # Response DTOs
├── handler/            # HTTP handlers
├── middleware/         # Middleware components
├── model/              # Data models
├── repository/         # Database operations
├── service/            # Business logic
├── util/               # Utility functions
├── .env                # Environment variables
├── Dockerfile          # Docker build configuration
├── docker-compose.yml  # Docker compose file
├── go.mod              # Go modules file
├── go.sum              # Go modules checksums
└── main.go             # Application entry point
```

## Prerequisites

- Docker
- Docker Compose
- Golang 1.21+

## Installation

### 1. Clone the repository

```sh
git clone https://github.com/user020603/postgresql.git
cd postgresql/task-manager
```

### 2. Set up the environment variables

Create a `.env` file in the root directory and add the following environment variables:

```
# Database Configuration
DB_HOST
DB_PORT
DB_USER
DB_PASSWORD
DB_NAME

# JWT Configuration
JWT_SECRET
JWT_EXPIRATION

# Server Configuration
PORT
```

### 3. Build and run the Docker containers

```sh
docker-compose up --build
```

### 4. Access the application

The application should now be running on `http://localhost:8080`.

## API Endpoints

### Authentication

- **Register**
  - `POST /api/v1/auth/register`
  - Request Body:
    ```json
    {
      "username": "testuser",
      "email": "test@example.com",
      "password": "password123"
    }
    ```

- **Login**
  - `POST /api/v1/auth/login`
  - Request Body:
    ```json
    {
      "username": "testuser",
      "password": "password123"
    }
    ```

- **Profile**
  - `GET /api/v1/auth/profile`
  - Headers: `Authorization: Bearer <JWT_TOKEN>`

### Tasks

- **Create Task**
  - `POST /api/v1/tasks`
  - Headers: `Authorization: Bearer <JWT_TOKEN>`
  - Request Body:
    ```json
    {
      "title": "Complete project",
      "description": "Finish implementing the task manager",
      "assigned_to_id": 1,
      "due_date": "2023-12-31T23:59:59Z"
    }
    ```

- **Get All Tasks**
  - `GET /api/v1/tasks`
  - Headers: `Authorization: Bearer <JWT_TOKEN>`
  - Query Parameters:
    - `page`: Page number (default: 1)
    - `page_size`: Number of tasks per page (default: 10)

- **Get My Tasks**
  - `GET /api/v1/tasks/my`
  - Headers: `Authorization: Bearer <JWT_TOKEN>`
  - Query Parameters:
    - `page`: Page number (default: 1)
    - `page_size`: Number of tasks per page (default: 10)

- **Get Task by ID**
  - `GET /api/v1/tasks/:id`
  - Headers: `Authorization: Bearer <JWT_TOKEN>`

- **Update Task**
  - `PUT /api/v1/tasks/:id`
  - Headers: `Authorization: Bearer <JWT_TOKEN>`
  - Request Body:
    ```json
    {
      "title": "Updated title",
      "description": "Updated description",
      "status": "IN_PROGRESS",
      "assigned_to_id": 2,
      "due_date": "2024-01-31T23:59:59Z"
    }
    ```

- **Delete Task**
  - `DELETE /api/v1/tasks/:id`
  - Headers: `Authorization: Bearer <JWT_TOKEN>`

## Testing the API

You can use tools like `curl` or Postman to test the API endpoints. Here are some example `curl` commands:

1. Register a user:
```sh
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'
```

2. Login:
```sh
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

3. Create a task (Use the token from login response):
```sh
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Complete project","description":"Finish implementing the task manager","assigned_to_id":1,"due_date":"2023-12-31T23:59:59Z"}'
```

4. Get all tasks:
```sh
curl -X GET http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer YOUR_TOKEN"
```

5. Get user profile:
```sh
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a pull request

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Acknowledgments

- [GORM](https://gorm.io/)
- [Gin Gonic](https://github.com/gin-gonic/gin)
- [Docker](https://www.docker.com/)
- [PostgreSQL](https://www.postgresql.org/)