# API Documentation: Kartala Task Management System

Sistem manajemen tugas dengan autentikasi JWT dan arsitektur modular.

## Base URL
`http://localhost:8080/api/v1`

## Auth Endpoints

### 1. Register User
Mendaftarkan akun baru.
- **URL**: `/auth/register`
- **Method**: `POST`
- **Body**:
```json
{
  "name": "Jane Doe",
  "username": "janedoe",
  "email": "jane@example.com",
  "password": "password123",
  "confirm_password": "password123"
}
```
- **Responses**:
  - `201 Created`: Registrasi berhasil.
  - `400 Bad Request`: Validasi gagal atau internal error.
  - `409 Conflict`: Email/username sudah terpakai.

### 2. Login
Mendapatkan token akses.
- **URL**: `/auth/login`
- **Method**: `POST`
- **Body**:
```json
{
  "identifier": "janedoe", // bisa email atau username
  "password": "password123",
  "remember_me": false
}
```
- **Responses**:
  - `200 OK`: Mengembalikan `access_token` dan `refresh_token`.
  - `401 Unauthorized`: Kredensial salah.

---

## Task Endpoints (Authenticated)
Semua endpoint tugas memerlukan header: `Authorization: Bearer <access_token>`

### 3. Get All Tasks
- **URL**: `/tasks`
- **Method**: `GET`
- **Query Params**: `page`, `limit`, `title`, `status`
- **Response**: List tugas milik user yang sedang aktif dengan metadata pagination.

### 4. Create Task
- **URL**: `/tasks`
- **Method**: `POST`
- **Body**:
```json
{
  "title": "Belajar Go",
  "description": "Mempelajari Goroutines dan Channels",
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
