# Database Seeder

Script untuk populate database PostgreSQL dengan initial data.

## 🌱 Data yang Di-seed

### Users
1. **Admin User** - `admin@example.com`
2. **Test User** - `user@example.com`

### Tasks
- Beberapa contoh task untuk admin dan user.

## 🚀 Cara Menjalankan

### 1. Pastikan PostgreSQL Running
Pastikan konfigurasi di `.env` sudah benar untuk koneksi PostgreSQL.

### 2. Jalankan Seeder
```bash
# Dari root project
go run cmd/seeder/main.go
```

### 3. Output
```
🌱 Starting database seeder...
INFO[2026-02-21T16:22:15+07:00] Using PostgreSQL for seeding
🌱 Seeding completed!
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

## ⚙️ Features

- ✅ **Duplicate Check** - Skip users yang sudah ada
- ✅ **Password Hashing** - Semua password di-hash dengan bcrypt
- ✅ **UUID Generation** - Auto-generate UUID untuk setiap user
- ✅ **Timestamps** - CreatedAt & UpdatedAt otomatis
- ✅ **Error Handling** - Graceful error handling

## 🔐 Default Passwords

**⚠️ IMPORTANT:** Ganti password ini di production!

- Admin: `admin123`
- User: `user123`
