# API Documentation: Kartala Task Management System

Task management system with JWT authentication and modular architecture.

## Base URL
`http://localhost:8080/api/v1`

## Auth Endpoints

### 1. Register User
Register a new account.
- **URL**: `/auth/register`
- **Method**: `POST`
- **Body**:
```json
{
  "username": "janedoe",
  "email": "jane@example.com",
  "password": "password123"
}
```
- **Responses**:
  - `201 Created`: Registration successful.
  - `400 Bad Request`: Validation failure or invalid request.
  - `409 Conflict`: Email or username already taken.

### 2. Login
Authenticate and get access tokens.
- **URL**: `/auth/login`
- **Method**: `POST`
- **Body**:
```json
{
  "email": "jane@example.com",
  "password": "password123"
}
```
- **Responses**:
  - `200 OK`: Returns `access_token` and `refresh_token`.
  - `401 Unauthorized`: Invalid email or password.

---

## Task Endpoints (Authenticated)
All task endpoints require header: `Authorization: Bearer <access_token>`

### 3. Get All Tasks
- **URL**: `/tasks`
- **Method**: `GET`
- **Query Params**: `page`, `limit`, `title`, `status`
- **Response**: List of tasks belonging to the active user with pagination metadata.

### 4. Create Task
- **URL**: `/tasks`
- **Method**: `POST`
- **Body**:
```json
{
  "title": "Learn Go",
  "description": "Learning Goroutines and Channels",
  "status": "todo"
}
```

### 5. Update Task
- **URL**: `/tasks/:id`
- **Method**: `PUT`
- **Body**: (Optional fields) `title`, `description`, `status`

### 6. Delete Task
- **URL**: `/tasks/:id`
- **Method**: `DELETE`
